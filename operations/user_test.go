package operations

import (
	"fmt"
	"testing"

	"github.com/Notes-App/database"
	"github.com/Notes-App/schemas"
	"github.com/Notes-App/validators"
)

func TestCreateUser(t *testing.T) {
	dbErr := database.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer database.CleanUp()
	type test struct {
		name     string
		input    schemas.User
		expected *schemas.User
		err      error
	}

	tests := []test{
		{name: "normal User", input: schemas.User{Email: "test@gmail.com", Password: "ping pong"}, expected: &schemas.User{Email: "test@gmail.com", Password: "ping pong"}, err: nil},
		{name: "empty User", input: schemas.User{}, expected: nil, err: fmt.Errorf("invalid email address")},
		{name: "invalid email", input: schemas.User{Email: "asdafs@.com", Password: "ping pong"}, expected: nil, err: fmt.Errorf("invalid email address")},
		{name: "no password", input: schemas.User{Email: "test@test.com", Password: ""}, expected: nil, err: fmt.Errorf("empty password")},
	}

	for _, tc := range tests {
		t.Log(tc.name)
		_, err := CreateUser(tc.input)
		if fmt.Sprint(err) != fmt.Sprint(tc.err) {
			t.Fatalf("got: %s    expected:%s", err, tc.err)
		}
		var result *schemas.User
		database.DB.First(&result, "email = ?", tc.input.Email)
		if result == nil && tc.err == nil {
			t.Fatalf("got: %v    expected:%v", result, *tc.expected)
		}
		if result != nil && tc.expected != nil {
			if !validators.IsUUIDValid(result.UUID) {
				t.Fatalf("got invalid uuid: %s", result.UUID)
			}
			if result.Email != tc.expected.Email {
				t.Fatalf("got: %s    expected:%s", result.Email, tc.expected.Email)
			}
			if result.Password == tc.expected.Password {
				t.Fatalf("got: %s    expected:%s", result.Password, tc.expected.Password)
			}
		}
	}
}

func TestFindUser(t *testing.T) {
	dbErr := database.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer database.CleanUp()
	_, err := CreateUser(schemas.User{Email: "test@gmail.com", Password: "ping pong"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = CreateUser(schemas.User{Email: "test2@gmail.com", Password: "ping pang"})
	if err != nil {
		t.Fatal(err)
	}

	type test struct {
		name     string
		input    schemas.User
		expected *schemas.User
		err      error
	}

	tests := []test{
		{name: "normal user", input: schemas.User{Email: "test@gmail.com", Password: "ping pong"}, expected: &schemas.User{Email: "test@gmail.com", Password: "ping pong"}, err: nil},
		{name: "empty user", input: schemas.User{}, expected: nil, err: fmt.Errorf("invalid email address")},
		{name: "empty email", input: schemas.User{Password: "ping pong"}, expected: nil, err: fmt.Errorf("invalid email address")},
		{name: "empty password", input: schemas.User{Email: "test@gmail.com"}, expected: nil, err: fmt.Errorf("empty password")},
		{name: "unsigned user", input: schemas.User{Email: "test3@gmail.com", Password: "ping pong"}, expected: nil, err: fmt.Errorf("user not found")},
		{name: "wrong user", input: schemas.User{Email: "test2@gmail.com", Password: "wrong password"}, expected: nil, err: fmt.Errorf("wrong password")},
	}

	for _, tc := range tests {
		t.Log(tc.name)
		got, err := FindUser(tc.input)
		if fmt.Sprint(err) != fmt.Sprint(tc.err) {
			t.Fatalf("got: %s    expected:%s", err, tc.err)
		}
		if got == nil && tc.err == nil {
			t.Fatalf("got: %v    expected:%v", got, tc.expected)
		}
		if got != nil && tc.expected != nil {
			if !validators.IsUUIDValid(got.UUID) {
				t.Fatalf("got invalid uuid: %s", got.UUID)
			}
			if got.Email != tc.expected.Email {
				t.Fatalf("got: %s    expected:%s", got.Email, tc.expected.Email)
			}
			if got.Password == tc.expected.Password {
				t.Fatalf("got: %s    expected:%s", got.Password, tc.expected.Password)
			}
		}
	}

}

func TestFindUserByUUID(t *testing.T) {
	dbErr := database.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer database.CleanUp()
	u1, err := CreateUser(schemas.User{Email: "test@gmail.com", Password: "ping pong"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = CreateUser(schemas.User{Email: "test2@gmail.com", Password: "ping pang"})
	if err != nil {
		t.Fatal(err)
	}

	type test struct {
		name     string
		input    string
		expected *schemas.User
		err      error
	}

	tests := []test{
		{name: "normal user", input: u1.UUID, expected: &schemas.User{Email: "test@gmail.com", Password: "ping pong"}, err: nil},
		{name: "empty uuid", input: "", expected: nil, err: fmt.Errorf("invalid uuid")},
		{name: "wrong uuid", input: "asfadfgsasasf", expected: nil, err: fmt.Errorf("invalid uuid")},
	}

	for _, tc := range tests {
		t.Log(tc.name)
		got, err := FindUserByUUID(tc.input)
		if fmt.Sprint(err) != fmt.Sprint(tc.err) {
			t.Fatalf("got: %s    expected:%s", err, tc.err)
		}
		if got == nil && tc.err == nil {
			t.Fatalf("got: %v    expected:%v", got, tc.expected)
		}
		if got != nil && tc.expected != nil {
			if !validators.IsUUIDValid(got.UUID) {
				t.Fatalf("got invalid uuid: %s", got.UUID)
			}
			if got.Email != tc.expected.Email {
				t.Fatalf("got: %s    expected:%s", got.Email, tc.expected.Email)
			}
			if got.Password == tc.expected.Password {
				t.Fatalf("got: %s    expected:%s", got.Password, tc.expected.Password)
			}
		}
	}

}
