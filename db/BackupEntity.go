package db

type BackupEntity interface {
	GetID() int
	GetName() string
}

func (e View) GetID() int {
	return e.ID
}

func (e View) GetName() string {
	return e.Title
}

func (e Template) GetID() int {
	return e.ID
}

func (e Template) GetName() string {
	return e.Name
}
