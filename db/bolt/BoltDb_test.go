package bolt

import (
	"fmt"
	"github.com/ansible-semaphore/semaphore/db"
	"reflect"
	"testing"
)

type test1 struct {
	ID             int    `db:"ID"`
	FirstName      string `db:"first_name" json:"firstName"`
	LastName       string `db:"last_name" json:"lastName"`
	Password       string `db:"-" json:"password"`
	PasswordRepeat string `db:"-" json:"passwordRepeat"`
	PasswordHash   string `db:"password" json:"-"`
	Removed        bool   `db:"removed"`
}

func TestMarshalObject_UserWithPwd(t *testing.T) {
	user := db.UserWithPwd{
		Pwd: "123456",
		User: db.User{
			Username: "fiftin",
			Password: "345345234523452345234",
		},
	}

	bytes, err := marshalObject(user)

	if err != nil {
		t.Fatal(fmt.Errorf("function returns error: " + err.Error()))
	}

	str := string(bytes)

	if str != `{"id":0,"created":"0001-01-01T00:00:00Z","username":"fiftin","name":"","email":"","password":"345345234523452345234","admin":false,"external":false,"alert":false}` {
		t.Fatal(fmt.Errorf("incorrect marshalling result"))
	}

	fmt.Println(str)
}

func TestMarshalObject(t *testing.T) {
	test1 := test1{
		FirstName:      "Denis",
		LastName:       "Gukov",
		Password:       "1234556",
		PasswordRepeat: "123456",
		PasswordHash:   "9347502348723",
	}

	bytes, err := marshalObject(test1)

	if err != nil {
		t.Fatal(fmt.Errorf("function returns error: " + err.Error()))
	}

	str := string(bytes)
	if str != `{"ID":0,"first_name":"Denis","last_name":"Gukov","password":"9347502348723","removed":false}` {
		t.Fatal(fmt.Errorf("incorrect marshalling result"))
	}

	fmt.Println(str)
}

func TestUnmarshalObject(t *testing.T) {
	test1 := test1{}
	data := `{
	"first_name": "Denis", 
	"last_name": "Gukov",
	"password": "9347502348723"
}`
	err := unmarshalObject([]byte(data), &test1)
	if err != nil {
		t.Fatal(fmt.Errorf("function returns error: " + err.Error()))
	}
	if test1.FirstName != "Denis" ||
		test1.LastName != "Gukov" ||
		test1.Password != "" ||
		test1.PasswordRepeat != "" ||
		test1.PasswordHash != "9347502348723" {
		t.Fatal(fmt.Errorf("object unmarshalled incorrectly"))
	}
}

func TestGetFieldNameByTag(t *testing.T) {
	f, err := getFieldNameByTagSuffix(reflect.TypeOf(test1{}), "db", "first_name")
	if err != nil {
		t.Fatal(err.Error())
	}

	if f != "FirstName" {
		t.Fatal()
	}
}

func TestGetFieldNameByTag2(t *testing.T) {
	f, err := getFieldNameByTagSuffix(reflect.TypeOf(db.UserWithPwd{}), "db", "id")
	if err != nil {
		t.Fatal(err.Error())
	}
	if f != "ID" {
		t.Fatal()
	}
}

func TestBoltDb_CreateAPIToken(t *testing.T) {
	store := CreateTestStore()

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

	token, err := store.CreateAPIToken(db.APIToken{
		ID:     "f349gyhgqirgysfgsfg34973dsfad",
		UserID: user.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	token2, err := store.GetAPIToken(token.ID)
	if err != nil {
		t.Fatal(err)
	}

	if token2.ID != token.ID {
		t.Fatal()
	}

	tokens, err := store.GetAPITokens(user.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(tokens) != 1 {
		t.Fatal()
	}

	if tokens[0].ID != token.ID {
		t.Fatal()
	}

	err = store.ExpireAPIToken(user.ID, token.ID)
	if err != nil {
		t.Fatal(err)
	}

	token2, err = store.GetAPIToken(token.ID)
	if err != nil {
		t.Fatal(err)
	}

	if !token2.Expired {
		t.Fatal()
	}

	err = store.DeleteAPIToken(user.ID, token.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.GetAPIToken(token.ID)
	if err == nil {
		t.Fatal("Token not deleted")
	}
}
