package tasks

import (
	"encoding/json"
	"os"

	"github.com/ansible-semaphore/semaphore/db"
	"github.com/ansible-semaphore/semaphore/db_lib"
	"github.com/ansible-semaphore/semaphore/pkg/task_logger"
)

type LocalJob struct {
	// Received constant fields
	Task        db.Task
	Template    db.Template
	Secret      string
	Logger      task_logger.Logger

	App db_lib.LocalApp

	// Internal field
	Process *os.Process
}

func (t *LocalJob) Kill() {
	if t.Process == nil {
		return
	}
	err := t.Process.Kill()
	if err != nil {
		t.Log(err.Error())
	}
}

func (t *LocalJob) Log(msg string) {
	t.Logger.Log(msg)
}

func (t *LocalJob) SetStatus(status task_logger.TaskStatus) {
	t.Logger.SetStatus(status)
}

// func (t *LocalJob) getEnvironmentENV() (res []string, err error) {
// 	environmentVars := make(map[string]string)

// 	if t.Environment.ENV != nil {
// 		err = json.Unmarshal([]byte(*t.Environment.ENV), &environmentVars)
// 		if err != nil {
// 			return
// 		}
// 	}

// 	for key, val := range environmentVars {
// 		res = append(res, fmt.Sprintf("%s=%s", key, val))
// 	}

// 	for _, secret := range t.Environment.Secrets {
// 		if secret.Type != db.EnvironmentSecretEnv {
// 			continue
// 		}
// 		res = append(res, fmt.Sprintf("%s=%s", secret.Name, secret.Secret))
// 	}

// 	return
// }

// nolint: gocyclo
func (t *LocalJob) getShellArgs(username string, incomingVersion *string) (args []string, err error) {
	// extraVars, err := t.getEnvironmentExtraVars(username, incomingVersion)

	if err != nil {
		t.Log(err.Error())
		t.Log("Error getting environment extra vars")
		return
	}

	var templateExtraArgs []string
	if t.Template.Arguments != nil {
		err = json.Unmarshal([]byte(*t.Template.Arguments), &templateExtraArgs)
		if err != nil {
			t.Log("Invalid format of the template extra arguments, must be valid JSON")
			return
		}
	}

	var taskExtraArgs []string
	if t.Template.AllowOverrideArgsInTask && t.Task.Arguments != nil {
		err = json.Unmarshal([]byte(*t.Task.Arguments), &taskExtraArgs)
		if err != nil {
			t.Log("Invalid format of the TaskRunner extra arguments, must be valid JSON")
			return
		}
	}

	// Script to run
	args = append(args, t.Template.Playbook)

	// // Include Environment Secret Vars
	// for _, secret := range t.Environment.Secrets {
	// 	if secret.Type == db.EnvironmentSecretVar {
	// 		args = append(args, fmt.Sprintf("%s=%s", secret.Name, secret.Secret))
	// 	}
	// }

	// Include extra args from template
	args = append(args, templateExtraArgs...)

	// // Include ExtraVars and Survey Vars
	// for name, value := range extraVars {
	// 	if name != "semaphore_vars" {
	// 		args = append(args, fmt.Sprintf("%s=%s", name, value))
	// 	}
	// }

	// Include extra args from task
	args = append(args, taskExtraArgs...)

	return
}

func (t *LocalJob) Run(username string, incomingVersion *string) (err error) {

	t.SetStatus(task_logger.TaskRunningStatus) // It is required for local mode. Don't delete

	var args []string
	var inputs map[string]string

	switch t.Template.App {
	default:
		args, err = t.getShellArgs(username, incomingVersion)
	}

	if err != nil {
		return
	}

	// environmentVariables, err := t.getEnvironmentENV()
	// if err != nil {
		// return
	// }

	// if t.Template.Type != db.TemplateTask {

	// 	environmentVariables = append(environmentVariables, fmt.Sprintf("SEMAPHORE_TASK_TYPE=%s", t.Template.Type))

	// 	if incomingVersion != nil {
	// 		environmentVariables = append(
	// 			environmentVariables,
	// 			fmt.Sprintf("SEMAPHORE_TASK_INCOMING_VERSION=%s", *incomingVersion))
	// 	}

	// 	if t.Template.Type == db.TemplateBuild && t.Task.Version != nil {
	// 		environmentVariables = append(
	// 			environmentVariables,
	// 			fmt.Sprintf("SEMAPHORE_TASK_TARGET_VERSION=%s", *t.Task.Version))
	// 	}
	// }

	environmentVars := []string {}

	return t.App.Run(args, &environmentVars, inputs, func(p *os.Process) {
		t.Process = p
	})
}
