package db

type Schedule struct {
	ID         int    `db:"id" json:"id" backup:"-"`
	ProjectID  int    `db:"project_id" json:"project_id" backup:"-"`
	TemplateID int    `db:"template_id" json:"template_id" backup:"-"`
	CronFormat string `db:"cron_format" json:"cron_format"`
	Name       string `db:"name" json:"name"`
	Active     bool   `db:"active" json:"active"`
}

type ScheduleWithTpl struct {
	Schedule
	TemplateName string `db:"tpl_name" json:"tpl_name"`
}
