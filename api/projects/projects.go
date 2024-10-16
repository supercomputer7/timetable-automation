package projects

import (
	"net/http"

	"github.com/ansible-semaphore/semaphore/api/helpers"
	"github.com/ansible-semaphore/semaphore/db"
	"github.com/ansible-semaphore/semaphore/util"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/context"
)

// GetProjects returns all projects in this users context
func GetProjects(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*db.User)

	var err error
	var projects []db.Project
	if user.Admin {
		projects, err = helpers.Store(r).GetAllProjects()
	} else {
		projects, err = helpers.Store(r).GetProjects(user.ID)
	}

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, projects)
}

// AddProject adds a new project to the database
func AddProject(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, "user").(*db.User)

	if !user.Admin && !util.Config.NonAdminCanCreateProject {
		log.Warn(user.Username + " is not permitted to edit users")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var bodyWithDemo struct {
		db.Project
		Demo bool `json:"demo"`
	}

	if !helpers.Bind(w, r, &bodyWithDemo) {
		return
	}

	body := bodyWithDemo.Project

	store := helpers.Store(r)

	body, err := store.CreateProject(body)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	_, err = store.CreateProjectUser(db.ProjectUser{ProjectID: body.ID, UserID: user.ID, Role: db.ProjectOwner})
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.EventLog(r, helpers.EventLogCreate, helpers.EventLogItem{
		UserID:      helpers.UserFromContext(r).ID,
		ProjectID:   body.ID,
		ObjectType:  db.EventProject,
		ObjectID:    body.ID,
		Description: "Project created",
	})

	helpers.WriteJSON(w, http.StatusCreated, body)
}
