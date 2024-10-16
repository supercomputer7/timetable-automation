package sql

import (
	"database/sql"
	"github.com/ansible-semaphore/semaphore/db"
)

func (d *SqlDb) CreateSchedule(schedule db.Schedule) (newSchedule db.Schedule, err error) {
	insertID, err := d.insert(
		"id",
		"insert into project__schedule (project_id, template_id, cron_format, `name`, `active`)"+
			"values (?, ?, ?, ?, ?, ?)",
		schedule.ProjectID,
		schedule.TemplateID,
		schedule.CronFormat,
		schedule.Name,
		schedule.Active)

	if err != nil {
		return
	}

	newSchedule = schedule
	newSchedule.ID = insertID

	return
}

func (d *SqlDb) UpdateSchedule(schedule db.Schedule) error {
	_, err := d.exec("update project__schedule set "+
		"cron_format=?, "+
		"template_id=?, "+
		"`name`=?, "+
		"`active`=?, "+
		"where project_id=? and id=?",
		schedule.CronFormat,
		schedule.TemplateID,
		schedule.Name,
		schedule.Active,
		schedule.ProjectID,
		schedule.ID)
	return err
}

func (d *SqlDb) GetSchedule(projectID int, scheduleID int) (template db.Schedule, err error) {
	err = d.selectOne(
		&template,
		"select * from project__schedule where project_id=? and id=?",
		projectID,
		scheduleID)

	if err == sql.ErrNoRows {
		err = db.ErrNotFound
	}

	return
}

func (d *SqlDb) DeleteSchedule(projectID int, scheduleID int) error {
	_, err := d.exec("delete from project__schedule where project_id=? and id=?", projectID, scheduleID)
	return err
}

func (d *SqlDb) GetSchedules() (schedules []db.Schedule, err error) {
	_, err = d.selectAll(&schedules, "select * from project__schedule where cron_format != ''")
	return
}

func (d *SqlDb) GetProjectSchedules(projectID int) (schedules []db.ScheduleWithTpl, err error) {
	_, err = d.selectAll(&schedules,
		"SELECT ps.*, pt.name as tpl_name FROM project__schedule ps "+
			"JOIN project__template pt ON pt.id = ps.template_id "+
			"WHERE ps.project_id=?",
		projectID)
	return
}

func (d *SqlDb) GetTemplateSchedules(projectID int, templateID int) (schedules []db.Schedule, err error) {
	_, err = d.selectAll(&schedules,
		"SELECT * FROM project__schedule WHERE project_id=? AND template_id=?",
		projectID,
		templateID)
	return
}

func (d *SqlDb) SetScheduleActive(projectID int, scheduleID int, active bool) error {
	_, err := d.exec("update project__schedule set `active`=? where project_id=? and id=?",
		active,
		projectID,
		scheduleID)
	return err
}
