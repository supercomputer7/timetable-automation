package project

import (
	"fmt"

	"github.com/ansible-semaphore/semaphore/db"
	"github.com/ansible-semaphore/semaphore/services/schedules"
)

func getEntryByName[T BackupEntry](name *string, items []T) *T {
	if name == nil {
		return nil
	}
	for _, o := range items {
		if o.GetName() == *name {
			return &o
		}
	}
	return nil
}

func verifyDuplicate[T BackupEntry](name string, items []T) error {
	n := 0
	for _, o := range items {
		if o.GetName() == name {
			n++
		}
		if n > 2 {
			return fmt.Errorf("%s is duplicate", name)
		}
	}
	return nil
}

func (e BackupView) Verify(backup *BackupFormat) error {
	return verifyDuplicate[BackupView](e.Title, backup.Views)
}

func (e BackupView) Restore(store db.Store, b *BackupDB) error {
	v := e.View
	v.ProjectID = b.meta.ID
	newView, err := store.CreateView(v)
	if err != nil {
		return err
	}
	b.views = append(b.views, newView)
	return nil
}

func (e BackupTemplate) Verify(backup *BackupFormat) error {
	if err := verifyDuplicate[BackupTemplate](e.Name, backup.Templates); err != nil {
		return err
	}

	if e.View != nil && getEntryByName[BackupView](e.View, backup.Views) == nil {
		return fmt.Errorf("view does not exist in views[].name")
	}

	if buildTemplate := getEntryByName[BackupTemplate](e.BuildTemplate, backup.Templates); string(e.Type) == "deploy" && buildTemplate == nil {
		return fmt.Errorf("deploy is build but build_template does not exist in templates[].name")
	}

	if e.Cron != nil {
		if err := schedules.ValidateCronFormat(*e.Cron); err != nil {
			return err
		}
	}

	return nil
}

func (e BackupTemplate) Restore(store db.Store, b *BackupDB) error {
	var BuildTemplateID *int
	if string(e.Type) != "deploy" {
		BuildTemplateID = nil
	} else if k := findEntityByName[db.Template](e.BuildTemplate, b.templates); k == nil {
		BuildTemplateID = nil
	} else {
		BuildTemplateID = &(k.ID)
	}

	var ViewID *int
	if k := findEntityByName[db.View](e.View, b.views); k == nil {
		ViewID = nil
	} else {
		ViewID = &k.ID
	}

	template := e.Template
	template.ProjectID = b.meta.ID
	template.ViewID = ViewID
	template.BuildTemplateID = BuildTemplateID

	newTemplate, err := store.CreateTemplate(template)
	if err != nil {
		return err
	}
	b.templates = append(b.templates, newTemplate)
	if e.Cron != nil {
		_, err := store.CreateSchedule(
			db.Schedule{
				ProjectID:    b.meta.ID,
				TemplateID:   newTemplate.ID,
				CronFormat:   *e.Cron,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (backup *BackupFormat) Verify() error {
	for i, o := range backup.Views {
		if err := o.Verify(backup); err != nil {
			return fmt.Errorf("error at views[%d]: %s", i, err.Error())
		}
	}
	for i, o := range backup.Templates {
		if err := o.Verify(backup); err != nil {
			return fmt.Errorf("error at templates[%d]: %s", i, err.Error())
		}
	}
	return nil
}

func (backup *BackupFormat) Restore(user db.User, store db.Store) (*db.Project, error) {
	var b = BackupDB{}
	project := backup.Meta.Project

	newProject, err := store.CreateProject(project)

	if err != nil {
		return nil, err
	}

	if _, err = store.CreateProjectUser(db.ProjectUser{
		ProjectID: newProject.ID,
		UserID:    user.ID,
		Role:      db.ProjectOwner,
	}); err != nil {
		return nil, err
	}

	b.meta = newProject

	for i, o := range backup.Views {
		if err := o.Restore(store, &b); err != nil {
			return nil, fmt.Errorf("error at views[%d]: %s", i, err.Error())
		}
	}

	deployTemplates := make([]int, 0)
	for i, o := range backup.Templates {
		if string(o.Type) == "deploy" {
			deployTemplates = append(deployTemplates, i)
			continue
		}
		if err := o.Restore(store, &b); err != nil {
			return nil, fmt.Errorf("error at templates[%d]: %s", i, err.Error())
		}
	}

	for _, i := range deployTemplates {
		o := backup.Templates[i]
		if err := o.Restore(store, &b); err != nil {
			return nil, fmt.Errorf("error at templates[%d]: %s", i, err.Error())
		}
	}

	return &newProject, nil
}
