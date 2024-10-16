package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"time"
)

const databaseTimeFormat = "2006-01-02T15:04:05:99Z"

// GetParsedTime returns the timestamp as it will retrieved from the database
// This allows us to create timestamp consistency on return values from create requests
func GetParsedTime(t time.Time) time.Time {
	parsedTime, err := time.Parse(databaseTimeFormat, t.Format(databaseTimeFormat))
	if err != nil {
		log.Error(err)
	}
	return parsedTime
}

func ObjectToJSON(obj interface{}) *string {
	if obj == nil ||
		(reflect.ValueOf(obj).Kind() == reflect.Ptr && reflect.ValueOf(obj).IsNil()) ||
		(reflect.ValueOf(obj).Kind() == reflect.Slice && reflect.ValueOf(obj).IsZero()) {
		return nil
	}
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil
	}
	str := string(bytes)
	return &str
}

type RetrieveQueryParams struct {
	Offset       int
	Count        int
	SortBy       string
	SortInverted bool
	Filter       string
}

type ObjectReferrer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ObjectReferrers struct {
	Templates    []ObjectReferrer `json:"templates"`
	Inventories  []ObjectReferrer `json:"inventories"`
	Repositories []ObjectReferrer `json:"repositories"`
}

type IntegrationReferrers struct {
	IntegrationMatchers      []ObjectReferrer `json:"matchers"`
	IntegrationExtractValues []ObjectReferrer `json:"values"`
}

type IntegrationExtractorChildReferrers struct {
	Integrations []ObjectReferrer `json:"integrations"`
}

// ObjectProps describe database entities.
// It mainly used for NoSQL implementations (currently BoltDB) to preserve same
// data structure of different implementations and easy change it if required.
type ObjectProps struct {
	TableName             string
	Type                  reflect.Type // to which type the table bust be mapped.
	IsGlobal              bool         // doesn't belong to other table, for example to project or user.
	ReferringColumnSuffix string
	PrimaryColumnName     string
	SortableColumns       []string
	DefaultSortingColumn  string
	SortInverted          bool // sort from high to low object ID by default. It is useful for some NoSQL implementations.
}

var ErrNotFound = errors.New("no rows in result set")
var ErrInvalidOperation = errors.New("invalid operation")

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type Store interface {
	// Connect connects to the database.
	// Token parameter used if PermanentConnection returns false.
	// Token used for debugging of session connections.
	Connect(token string)
	Close(token string)

	// PermanentConnection returns true if connection should be kept from start to finish of the app.
	// This mode is suitable for MySQL and Postgres but not for BoltDB.
	// For BoltDB we should reconnect for each request because BoltDB support only one connection at time.
	PermanentConnection() bool

	// IsInitialized indicates is database already initialized, or it is empty.
	// The method is useful for creating required entities in database during first run.
	IsInitialized() (bool, error)

	GetOptions(params RetrieveQueryParams) (map[string]string, error)
	GetOption(key string) (string, error)
	SetOption(key string, value string) error
	DeleteOption(key string) error
	DeleteOptions(filter string) error

	GetUserCount() (int, error)
	GetUsers(params RetrieveQueryParams) ([]User, error)
	CreateUserWithoutPassword(user User) (User, error)
	CreateUser(user UserWithPwd) (User, error)
	DeleteUser(userID int) error

	// UpdateUser updates all fields of the entity except Pwd.
	// Pwd should be present of you want update user password. Empty Pwd ignored.
	UpdateUser(user UserWithPwd) error
	SetUserPassword(userID int, password string) error
	GetUser(userID int) (User, error)
	GetUserByLoginOrEmail(login string, email string) (User, error)

	GetProject(projectID int) (Project, error)
	GetAllProjects() ([]Project, error)
	GetProjects(userID int) ([]Project, error)
	CreateProject(project Project) (Project, error)
	DeleteProject(projectID int) error
	UpdateProject(project Project) error

	GetTemplates(projectID int, filter TemplateFilter, params RetrieveQueryParams) ([]Template, error)
	GetTemplateRefs(projectID int, templateID int) (ObjectReferrers, error)
	CreateTemplate(template Template) (Template, error)
	UpdateTemplate(template Template) error
	GetTemplate(projectID int, templateID int) (Template, error)
	DeleteTemplate(projectID int, templateID int) error

	GetSchedules() ([]Schedule, error)
	GetProjectSchedules(projectID int) ([]ScheduleWithTpl, error)
	GetTemplateSchedules(projectID int, templateID int) ([]Schedule, error)
	CreateSchedule(schedule Schedule) (Schedule, error)
	UpdateSchedule(schedule Schedule) error
	SetScheduleActive(projectID int, scheduleID int, active bool) error
	GetSchedule(projectID int, scheduleID int) (Schedule, error)
	DeleteSchedule(projectID int, scheduleID int) error

	GetAllAdmins() ([]User, error)
	GetProjectUsers(projectID int, params RetrieveQueryParams) ([]UserWithProjectRole, error)
	CreateProjectUser(projectUser ProjectUser) (ProjectUser, error)
	DeleteProjectUser(projectID int, userID int) error
	GetProjectUser(projectID int, userID int) (ProjectUser, error)
	UpdateProjectUser(projectUser ProjectUser) error

	CreateEvent(event Event) (Event, error)
	GetUserEvents(userID int, params RetrieveQueryParams) ([]Event, error)
	GetEvents(projectID int, params RetrieveQueryParams) ([]Event, error)

	GetAPITokens(userID int) ([]APIToken, error)
	CreateAPIToken(token APIToken) (APIToken, error)
	GetAPIToken(tokenID string) (APIToken, error)
	ExpireAPIToken(userID int, tokenID string) error
	DeleteAPIToken(userID int, tokenID string) error

	GetSession(userID int, sessionID int) (Session, error)
	CreateSession(session Session) (Session, error)
	ExpireSession(userID int, sessionID int) error
	TouchSession(userID int, sessionID int) error

	CreateTask(task Task, maxTasks int) (Task, error)
	UpdateTask(task Task) error

	GetTemplateTasks(projectID int, templateID int, params RetrieveQueryParams) ([]TaskWithTpl, error)
	GetProjectTasks(projectID int, params RetrieveQueryParams) ([]TaskWithTpl, error)
	GetTask(projectID int, taskID int) (Task, error)
	DeleteTaskWithOutputs(projectID int, taskID int) error
	GetTaskOutputs(projectID int, taskID int) ([]TaskOutput, error)
	CreateTaskOutput(output TaskOutput) (TaskOutput, error)
	GetTaskStages(projectID int, taskID int) ([]TaskStage, error)
	CreateTaskStage(stage TaskStage) (TaskStage, error)

	GetView(projectID int, viewID int) (View, error)
	GetViews(projectID int) ([]View, error)
	UpdateView(view View) error
	CreateView(view View) (View, error)
	DeleteView(projectID int, viewID int) error
	SetViewPositions(projectID int, viewPositions map[int]int) error

	GetRunner(projectID int, runnerID int) (Runner, error)
	GetRunners(projectID int) ([]Runner, error)
	DeleteRunner(projectID int, runnerID int) error
	GetGlobalRunnerByToken(token string) (Runner, error)
	GetGlobalRunner(runnerID int) (Runner, error)
	GetGlobalRunners(activeOnly bool) ([]Runner, error)
	DeleteGlobalRunner(runnerID int) error
	UpdateRunner(runner Runner) error
	CreateRunner(runner Runner) (Runner, error)
}

var TemplateProps = ObjectProps{
	TableName:             "project__template",
	Type:                  reflect.TypeOf(Template{}),
	PrimaryColumnName:     "id",
	ReferringColumnSuffix: "template_id",
	SortableColumns:       []string{"name"},
	DefaultSortingColumn:  "name",
}

var ScheduleProps = ObjectProps{
	TableName:         "project__schedule",
	Type:              reflect.TypeOf(Schedule{}),
	PrimaryColumnName: "id",
}

var ProjectUserProps = ObjectProps{
	TableName:         "project__user",
	Type:              reflect.TypeOf(ProjectUser{}),
	PrimaryColumnName: "user_id",
}

var ProjectProps = ObjectProps{
	TableName:             "project",
	Type:                  reflect.TypeOf(Project{}),
	PrimaryColumnName:     "id",
	ReferringColumnSuffix: "project_id",
	DefaultSortingColumn:  "name",
	IsGlobal:              true,
}

var UserProps = ObjectProps{
	TableName:         "user",
	Type:              reflect.TypeOf(User{}),
	PrimaryColumnName: "id",
	IsGlobal:          true,
}

var SessionProps = ObjectProps{
	TableName:         "session",
	Type:              reflect.TypeOf(Session{}),
	PrimaryColumnName: "id",
}

var TokenProps = ObjectProps{
	TableName:         "user__token",
	Type:              reflect.TypeOf(APIToken{}),
	PrimaryColumnName: "id",
}

var TaskProps = ObjectProps{
	TableName:         "task",
	Type:              reflect.TypeOf(Task{}),
	PrimaryColumnName: "id",
	IsGlobal:          true,
	SortInverted:      true,
}

var TaskOutputProps = ObjectProps{
	TableName: "task__output",
	Type:      reflect.TypeOf(TaskOutput{}),
}

var TaskStageProps = ObjectProps{
	TableName: "task__stage",
	Type:      reflect.TypeOf(TaskStage{}),
}

var ViewProps = ObjectProps{
	TableName:            "project__view",
	Type:                 reflect.TypeOf(View{}),
	PrimaryColumnName:    "id",
	DefaultSortingColumn: "position",
}

var GlobalRunnerProps = ObjectProps{
	TableName:         "runner",
	Type:              reflect.TypeOf(Runner{}),
	PrimaryColumnName: "id",
	IsGlobal:          true,
}

var OptionProps = ObjectProps{
	TableName:         "option",
	Type:              reflect.TypeOf(Option{}),
	PrimaryColumnName: "key",
	IsGlobal:          true,
}

func (p ObjectProps) GetReferringFieldsFrom(t reflect.Type) (fields []string, err error) {
	n := t.NumField()
	for i := 0; i < n; i++ {
		if !strings.HasSuffix(t.Field(i).Tag.Get("db"), p.ReferringColumnSuffix) {
			continue
		}
		fields = append(fields, t.Field(i).Tag.Get("db"))
	}

	for i := 0; i < n; i++ {
		if t.Field(i).Tag != "" || t.Field(i).Type.Kind() != reflect.Struct {
			continue
		}
		var nested []string
		nested, err = p.GetReferringFieldsFrom(t.Field(i).Type)
		if err != nil {
			return
		}
		fields = append(fields, nested...)
	}

	return
}

func StoreSession(store Store, token string, callback func()) {
	if !store.PermanentConnection() {
		store.Connect(token)
	}

	callback()

	if !store.PermanentConnection() {
		store.Close(token)
	}
}

type MapStringAnyField map[string]interface{}

func (m *MapStringAnyField) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, m)
	case string:
		return json.Unmarshal([]byte(v), m)
	default:
		return errors.New("unsupported type for MapStringAnyField")
	}
}

// Value implements the driver.Valuer interface for MapStringAnyField
func (m MapStringAnyField) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}
