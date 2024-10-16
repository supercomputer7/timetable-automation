package project

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ansible-semaphore/semaphore/db"
	"github.com/ansible-semaphore/semaphore/pkg/random"
)

func findNameByID[T db.BackupEntity](ID int, items []T) (*string, error) {
	for _, o := range items {
		if o.GetID() == ID {
			name := o.GetName()
			return &name, nil
		}
	}
	return nil, fmt.Errorf("item %d does not exist", ID)
}
func findEntityByName[T db.BackupEntity](name *string, items []T) *T {
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

func getSchedulesByProject(projectID int, schedules []db.Schedule) []db.Schedule {
	result := make([]db.Schedule, 0)
	for _, o := range schedules {
		if o.ProjectID == projectID {
			result = append(result, o)
		}
	}
	return result
}

func getScheduleByTemplate(templateID int, schedules []db.Schedule) *string {
	for _, o := range schedules {
		if o.TemplateID == templateID {
			return &o.CronFormat
		}
	}
	return nil
}

func getRandomName(name string) string {
	return name + " - " + random.String(10)
}

func makeUniqueNames[T any](items []T, getter func(item *T) string, setter func(item *T, name string)) {
	for i := len(items) - 1; i >= 0; i-- {
		for k, other := range items {
			if k == i {
				break
			}

			name := getter(&items[i])

			if name == getter(&other) {
				randomName := getRandomName(name)
				setter(&items[i], randomName)
				break
			}
		}
	}
}

func (b *BackupDB) makeUniqueNames() {

	makeUniqueNames(b.templates, func(item *db.Template) string {
		return item.Name
	}, func(item *db.Template, name string) {
		item.Name = name
	})

	makeUniqueNames(b.views, func(item *db.View) string {
		return item.Title
	}, func(item *db.View, name string) {
		item.Title = name
	})

}

func (b *BackupDB) load(projectID int, store db.Store) (err error) {

	b.templates, err = store.GetTemplates(projectID, db.TemplateFilter{}, db.RetrieveQueryParams{})
	if err != nil {
		return
	}

	b.views, err = store.GetViews(projectID)
	if err != nil {
		return
	}

	schedules, err := store.GetSchedules()
	if err != nil {
		return
	}

	b.schedules = getSchedulesByProject(projectID, schedules)

	b.meta, err = store.GetProject(projectID)
	if err != nil {
		return
	}

	b.makeUniqueNames()

	return
}

func (b *BackupDB) format() (*BackupFormat, error) {
	views := make([]BackupView, len(b.views))
	for i, o := range b.views {
		views[i] = BackupView{
			o,
		}
	}

	templates := make([]BackupTemplate, len(b.templates))
	for i, o := range b.templates {
		var View *string = nil
		if o.ViewID != nil {
			View, _ = findNameByID[db.View](*o.ViewID, b.views)
		}
		var BuildTemplate *string = nil
		if o.BuildTemplateID != nil {
			BuildTemplate, _ = findNameByID[db.Template](*o.BuildTemplateID, b.templates)
		}

		templates[i] = BackupTemplate{
			Template:      o,
			View:          View,
			BuildTemplate: BuildTemplate,
			Cron:          getScheduleByTemplate(o.ID, b.schedules),
		}
	}

	return &BackupFormat{
		Meta: BackupMeta{
			b.meta,
		},
		Views:              views,
		Templates:          templates,
	}, nil
}

func GetBackup(projectID int, store db.Store) (*BackupFormat, error) {
	backup := BackupDB{}
	if err := backup.load(projectID, store); err != nil {
		return nil, err
	}
	return backup.format()
}

func (b *BackupFormat) Marshal() (res string, err error) {
	data, err := marshalValue(reflect.ValueOf(b))
	if err != nil {
		return
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return
	}

	res = string(bytes)

	return
}

func (b *BackupFormat) Unmarshal(res string) (err error) {
	// Parse the JSON data into a map
	var jsonData interface{}
	if err = json.Unmarshal([]byte(res), &jsonData); err != nil {
		return
	}

	// Start the recursive unmarshaling process
	err = unmarshalValueWithBackupTags(jsonData, reflect.ValueOf(b))
	return
}
