package db

import (
	"fmt"
	"github.com/go-gorp/gorp/v3"
	"time"

	"github.com/ansible-semaphore/semaphore/pkg/task_logger"
	"github.com/ansible-semaphore/semaphore/util"
)

// Task is a model of a task which will be executed by the runner
type Task struct {
	ID         int `db:"id" json:"id"`
	TemplateID int `db:"template_id" json:"template_id" binding:"required"`
	ProjectID  int `db:"project_id" json:"project_id"`

	Status task_logger.TaskStatus `db:"status" json:"status"`

	Debug  bool `db:"debug" json:"debug"`
	DryRun bool `db:"dry_run" json:"dry_run"`
	Diff   bool `db:"diff" json:"diff"`

	// override variables
	Playbook    string  `db:"playbook" json:"playbook"`
	Limit       string  `db:"hosts_limit" json:"limit"`
	Secret      string  `db:"-" json:"secret"`
	Arguments   *string `db:"arguments" json:"arguments"`

	UserID        *int `db:"user_id" json:"user_id"`
	ScheduleID    *int `db:"schedule_id" json:"schedule_id"`

	Created time.Time  `db:"created" json:"created"`
	Start   *time.Time `db:"start" json:"start"`
	End     *time.Time `db:"end" json:"end"`

	Message string `db:"message" json:"message"`
}

func (task *Task) PreInsert(gorp.SqlExecutor) error {
	task.Created = task.Created.UTC()
	return nil
}

func (task *Task) PreUpdate(gorp.SqlExecutor) error {
	if task.Start != nil {
		start := task.Start.UTC()
		task.Start = &start
	}

	if task.End != nil {
		end := task.End.UTC()
		task.End = &end
	}
	return nil
}

func (task *Task) GetUrl() *string {
	if util.Config.WebHost != "" {
		taskUrl := fmt.Sprintf("%s/project/%d/history?t=%d", util.Config.WebHost, task.ProjectID, task.ID)
		return &taskUrl
	}

	return nil
}

func (task *Task) ValidateNewTask(template Template) error {
	return nil
}

// TaskWithTpl is the task data with additional fields
type TaskWithTpl struct {
	Task
	TemplatePlaybook string       `db:"tpl_playbook" json:"tpl_playbook"`
	TemplateAlias    string       `db:"tpl_alias" json:"tpl_alias"`
	TemplateType     TemplateType `db:"tpl_type" json:"tpl_type"`
	TemplateApp      TemplateApp  `db:"tpl_app" json:"tpl_app"`
	UserName         *string      `db:"user_name" json:"user_name"`
	BuildTask        *Task        `db:"-" json:"build_task"`
}

// TaskOutput is the ansible log output from the task
type TaskOutput struct {
	TaskID int       `db:"task_id" json:"task_id"`
	Task   string    `db:"task" json:"task"`
	Time   time.Time `db:"time" json:"time"`
	Output string    `db:"output" json:"output"`
}

type TaskStage struct {
	TaskID        int           `db:"task_id" json:"task_id"`
	Start         *time.Time    `db:"start" json:"start"`
	End           *time.Time    `db:"end" json:"end"`
	StartOutputID *int          `db:"start_output_id" json:"start_output_id"`
	EndOutputID   *int          `db:"end_output_id" json:"end_output_id"`
}
