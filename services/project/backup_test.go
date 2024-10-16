package project

import (
	"github.com/ansible-semaphore/semaphore/db"
	"github.com/ansible-semaphore/semaphore/db/bolt"
	"github.com/ansible-semaphore/semaphore/util"
	"testing"
)

type testItem struct {
	Name string
}

func TestBackupProject(t *testing.T) {
	util.Config = &util.ConfigType{
		TmpPath: "/tmp",
	}

	store := bolt.CreateTestStore()

	proj, err := store.CreateProject(db.Project{
		Name: "Test 123",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.CreateTemplate(db.Template{
		Name:          "Test",
		Playbook:      "test.yml",
		ProjectID:     proj.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	backup, err := GetBackup(proj.ID, store)
	if err != nil {
		t.Fatal(err)
	}

	if backup.Meta.ID != proj.ID {
		t.Fatal("backup meta ID wrong")
	}

	str, err := backup.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	if str != `{"environments":[{"json":"{\"author\": \"Denis\", \"comment\": \"Hello, World!\"}","name":"test"}],"meta":{"alert":false,"max_parallel_tasks":0,"name":"Test 123","type":""},"templates":[{"allow_override_args_in_task":false,"app":"","autorun":false,"name":"Test","playbook":"test.yml","suppress_success_alerts":false,"type":""}],"views":[]}` {
		t.Fatal("Invalid backup content")
	}

	restoredBackup := &BackupFormat{}
	err = restoredBackup.Unmarshal(str)
	if err != nil {
		t.Fatal(err)
	}

	if restoredBackup.Meta.Name != proj.Name {
		t.Fatal("backup meta ID wrong")
	}

	user, err := store.CreateUser(db.UserWithPwd{
		Pwd: "3412341234123",
		User: db.User{
			Username: "test",
			Name:     "Test",
			Email:    "test@example.com",
			Admin:    true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	restoredProj, err := restoredBackup.Restore(user, store)
	if err != nil {
		t.Fatal(err)
	}

	if restoredProj.Name != proj.Name {
		t.Fatal("backup meta ID wrong")
	}

}

func isUnique(items []testItem) bool {
	for i, item := range items {
		for k, other := range items {
			if i == k {
				continue
			}

			if item.Name == other.Name {
				return false
			}
		}
	}

	return true
}

func TestMakeUniqueNames(t *testing.T) {
	items := []testItem{
		{Name: "Project"},
		{Name: "Solution"},
		{Name: "Project"},
		{Name: "Project"},
		{Name: "Project"},
		{Name: "Project"},
	}

	makeUniqueNames(items, func(item *testItem) string {
		return item.Name
	}, func(item *testItem, name string) {
		item.Name = name
	})

	if !isUnique(items) {
		t.Fatal("Not unique names")
	}
}
