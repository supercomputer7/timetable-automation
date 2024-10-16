package project

import (
	"github.com/ansible-semaphore/semaphore/db"
)

type BackupDB struct {
	meta         db.Project
	templates    []db.Template
	views        []db.View
	schedules    []db.Schedule
}

type BackupFormat struct {
	Meta               BackupMeta          `backup:"meta"`
	Templates          []BackupTemplate    `backup:"templates"`
	Views              []BackupView        `backup:"views"`
}

type BackupMeta struct {
	db.Project
}

type BackupView struct {
	db.View
}

type BackupTemplate struct {
	db.Template

	BuildTemplate *string               `backup:"build_template"`
	View          *string               `backup:"view"`
	Cron          *string               `backup:"cron"`
}

type BackupEntry interface {
	GetName() string
	Verify(backup *BackupFormat) error
	Restore(store db.Store, b *BackupDB) error
}

func (e BackupView) GetName() string {
	return e.Title
}

func (e BackupTemplate) GetName() string {
	return e.Name
}
