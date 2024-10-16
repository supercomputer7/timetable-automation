package tasks

import (
	"math/rand"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/ansible-semaphore/semaphore/db_lib"

	"github.com/ansible-semaphore/semaphore/db"
	"github.com/ansible-semaphore/semaphore/db/bolt"
	"github.com/ansible-semaphore/semaphore/util"
)

func TestTaskRunnerRun(t *testing.T) {
	util.Config = &util.ConfigType{
		TmpPath: "/tmp",
	}

	store := bolt.CreateTestStore()

	pool := CreateTaskPool(store)

	go pool.Run()

	var task db.Task

	var err error

	db.StoreSession(store, "", func() {
		task, err = store.CreateTask(db.Task{}, 0)
	})

	if err != nil {
		t.Fatal(err)
	}

	taskRunner := TaskRunner{
		Task: task,
		pool: &pool,
	}
	taskRunner.job = &LocalJob{
		Task:        taskRunner.Task,
		Template:    taskRunner.Template,
		Logger:      &taskRunner,
		App: &db_lib.AnsibleApp{
			Template:   taskRunner.Template,
			Logger:     &taskRunner,
			Playbook: &db_lib.AnsiblePlaybook{
				Logger:     &taskRunner,
				TemplateID: taskRunner.Template.ID,
			},
		},
	}
	taskRunner.run()
}

func TestPopulateDetails(t *testing.T) {
	store := bolt.CreateTestStore()

	proj, err := store.CreateProject(db.Project{})
	if err != nil {
		t.Fatal(err)
	}

	tpl, err := store.CreateTemplate(db.Template{
		Name:          "Test",
		Playbook:      "test.yml",
		ProjectID:     proj.ID,
	})

	if err != nil {
		t.Fatal(err)
	}

	pool := TaskPool{store: store}

	tsk := TaskRunner{
		pool: &pool,
		Task: db.Task{
			TemplateID:  tpl.ID,
			ProjectID:   proj.ID,
		},
	}
	tsk.job = &LocalJob{
		Task:        tsk.Task,
		Template:    tsk.Template,
		Logger:      &tsk,
		App: &db_lib.AnsibleApp{
			Template:   tsk.Template,
			Logger:     &tsk,
			Playbook: &db_lib.AnsiblePlaybook{
				Logger:     &tsk,
				TemplateID: tsk.Template.ID,
			},
		},
	}

	err = tsk.populateDetails()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckTmpDir(t *testing.T) {
	//It should be able to create a random dir in /tmp
	dirName := path.Join(os.TempDir(), util.RandString(rand.Intn(10-4)+4))
	err := checkTmpDir(dirName)
	if err != nil {
		t.Fatal(err)
	}

	//checking again for this directory should return no error, as it exists
	err = checkTmpDir(dirName)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chmod(dirName, os.FileMode(int(0550)))
	if err != nil {
		t.Fatal(err)
	}

	//nolint: vetshadow
	if stat, err := os.Stat(dirName); err != nil {
		t.Fatal(err)
	} else if stat.Mode() != os.FileMode(int(0550)) {
		// File System is not support 0550 mode, skip this test
		return
	}

	err = checkTmpDir(dirName + "/noway")
	if err == nil {
		t.Fatal("You should not be able to write in this folder, causing an error")
	}
	err = os.Remove(dirName)
	if err != nil {
		t.Log(err)
	}
}
