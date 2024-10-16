package db_lib

import (
	"github.com/ansible-semaphore/semaphore/db"
	"github.com/ansible-semaphore/semaphore/pkg/task_logger"
)

func CreateApp(template db.Template, logger task_logger.Logger) LocalApp {
		return &ShellApp{
			Template:   template,
			Logger:     logger,
			App:        template.App,
		}
}
