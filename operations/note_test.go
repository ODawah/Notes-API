package operations

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Notes-App/database"
	"github.com/Notes-App/schemas"
	"github.com/Notes-App/validators"
)

func TestCreateNote(t *testing.T) {
	dbErr := database.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer database.CleanUp()
	longTitle := "then you sure as hell don't deserve me at my best.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard"
	longText := "I'm selfish, I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard to handle. But if you can't handle me at my worst, then you sure as hell don' deserve me at my b.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hardimpatient and a little insecure. I make mistakes, I am out of control and at times hard to handle. But if you can't handle me at my worst, then you sure as hell don't deserve me at my best.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard"

	u, _ := CreateUser(schemas.User{Email: "test@test.com", Password: "sgsersdfvsdsfaaa"})

	type test struct {
		name     string
		input    schemas.Note
		expected *schemas.Note
		err      error
	}

	tests := []test{
		{name: "normal note input", input: schemas.Note{Title: "pick up your son", Text: "go to school and pick him up", UserUuid: u.UUID}, expected: &schemas.Note{Title: "pick up your son", Text: "go to school and pick him up", UserUuid: u.UUID}, err: nil},
		{name: "long note title", input: schemas.Note{Title: longTitle, Text: "go to school and pick him up"}, expected: nil, err: errors.New("long note title")},
		{name: "no note title", input: schemas.Note{Title: "", Text: "go to school and pick him up"}, expected: nil, err: errors.New("empty note title")},
		{name: "long note text", input: schemas.Note{Title: "finish the assignment", Text: longText}, expected: nil, err: errors.New("long note text")},
		{name: "no note text", input: schemas.Note{Title: "cook the dinner", Text: "", UserUuid: u.UUID}, expected: nil, err: errors.New("empty note text")},
	}

	for _, tc := range tests {
		t.Log(tc.name)
		got, err := CreateNote(tc.input)
		if fmt.Sprint(err) != fmt.Sprint(tc.err) {
			t.Log(tc.name)
			t.Fatalf("got: %s    expected:%s", err, tc.err)
		}
		if got == nil && tc.err == nil {
			t.Log(tc.name)
			t.Fatalf("got: %s    expected:%s", got, tc.err)
		}
		if got != nil {
			if got.Text != tc.expected.Text {
				t.Log(tc.name)
				t.Fatalf("got: %s    expected:%s", got.Text, tc.expected.Text)
			}
			if got.Title != tc.expected.Title {
				t.Log(tc.name)
				t.Fatalf("got: %s    expected:%s", got.Title, tc.expected.Title)
			}
			if validators.IsUUIDValid(got.UUID) == false {
				t.Log(tc.name)
				t.Fatalf("got invalid uuid: %s", got.UUID)
			}
			if validators.IsUUIDValid(got.UserUuid) == false && tc.expected.UserUuid != got.UserUuid {
				t.Log(tc.name)
				t.Fatalf("got invalid user uuid: %s", got.UUID)
			}
		}

	}

}

func TestFindNoteByTitle(t *testing.T) {
	longTitle := "then you sure as hell don't deserve me at my best.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard"
	dbErr := database.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer database.CleanUp()
	u, _ := CreateUser(schemas.User{Email: "test@test.com", Password: "sgsersdfvsdsfaaa"})
	CreateNote(schemas.Note{Title: "Test Title", Text: "Test Text", UserUuid: u.UUID})
	CreateNote(schemas.Note{Title: "finish the tasks", Text: "create the API", UserUuid: u.UUID})

	type test struct {
		name        string
		searchTitle string
		uuid        string
		expected    *schemas.Note
		err         error
	}

	tests := []test{
		{name: "normal note title", searchTitle: "Test Title", uuid: u.UUID, expected: &schemas.Note{Title: "test title", Text: "test text", UserUuid: u.UUID}, err: nil},
		{name: "normal note title with trails", searchTitle: "   finish the tasks", uuid: u.UUID, expected: &schemas.Note{Title: "finish the tasks", Text: "create the api", UserUuid: u.UUID}, err: nil},
		{name: "normal note title with wrong user uuid", searchTitle: "Test Title", uuid: "123e4567-e89b-12d3-a456-42661417s4000", expected: nil, err: errors.New("invalid uuid")},
		{name: "no note title", searchTitle: "", uuid: u.UUID, expected: nil, err: errors.New("empty note title")},
		{name: "long note title", searchTitle: longTitle, uuid: u.UUID, expected: nil, err: errors.New("long note title")},
	}

	for _, tc := range tests {
		got, err := FindNoteByTitle(tc.searchTitle, tc.uuid)
		if fmt.Sprint(err) != fmt.Sprint(tc.err) {
			t.Log(tc.name)
			t.Fatalf("got: %s    expected:%s", err, tc.err)
		}
		if tc.expected != nil && got != nil {
			if validators.IsUUIDValid(got.UUID) == false {
				t.Log(tc.name)
				t.Fatalf("got invalid uuid: %s", got.UUID)
			}
			if got.Text != tc.expected.Text {
				t.Log(tc.name)
				t.Fatalf("got: %s    expected:%s", got.Text, tc.expected.Text)
			}
			if got.Title != tc.expected.Title {
				t.Log(tc.name)
				t.Fatalf("got: %s    expected:%s", got.Title, tc.expected.Title)
			}
			if got.UserUuid != u.UUID {
				t.Log(tc.name)
				t.Fatalf("got: %s    expected:%s", got.UserUuid, u.UUID)
			}
		}
	}

}

func TestDeleteNote(t *testing.T) {
	dbErr := database.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer database.CleanUp()
	u, _ := CreateUser(schemas.User{Email: "test@test.com", Password: "sgsersdfvsdsfaaa"})
	n1, _ := CreateNote(schemas.Note{Title: "Test Title", Text: "Test Text", UserUuid: u.UUID})
	n2, _ := CreateNote(schemas.Note{Title: "finish the tasks", Text: "create the API", UserUuid: u.UUID})

	type test struct {
		name     string
		input    string
		uuid     string
		expected int
		err      error
	}

	tests := []test{
		{name: "normal note uuid", input: n1.UUID, uuid: u.UUID, expected: 1, err: nil},
		{name: "normal note wrong uuid", input: n2.UUID, uuid: "6a3f5046-7d3f-11ed-a1eb-0242ac120002", expected: 0, err: errors.New("user not found")},
		{name: "normal note uuid", input: n2.UUID, uuid: u.UUID, expected: 1, err: nil},
		{name: "no uuid", input: "", expected: 0, uuid: u.UUID, err: errors.New("invalid UUID")},
		{name: "uuid not in db", input: "6a3f5046-7d3f-11ed-a1eb-0242ac120002", uuid: u.UUID, expected: 0, err: errors.New("note not found")},
	}

	for _, tc := range tests {
		got, err := DeleteNote(tc.input, tc.uuid)
		if fmt.Sprint(err) != fmt.Sprint(tc.err) {
			t.Log(tc.name)
			t.Fatalf("got: %s    expected:%s", err, tc.err)
		}
		if got != tc.expected {
			t.Log(tc.name)
			t.Fatalf("got: %d    expected:%d", got, tc.expected)
		}
	}

}

func TestUpdateNote(t *testing.T) {
	dbErr := database.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer database.CleanUp()
	u, _ := CreateUser(schemas.User{Email: "test@test.com", Password: "sgsersdfvsdsfaaa"})
	n1, _ := CreateNote(schemas.Note{Title: "Test Title", Text: "Test Text", UserUuid: u.UUID})
	n2, _ := CreateNote(schemas.Note{Title: "finish the tasks", Text: "create the API", UserUuid: u.UUID})

	type test struct {
		name     string
		uuid     string
		input    schemas.Note
		expected int
		err      error
	}

	tests := []test{
		{name: "normal note updates", uuid: n1.UUID, input: schemas.Note{Title: "updated title", Text: "updated text", UserUuid: u.UUID}, expected: 1, err: nil},
		{name: "title update only", uuid: n1.UUID, input: schemas.Note{Title: "updated title", Text: "updated text", UserUuid: u.UUID}, expected: 1, err: nil},
		{name: "text update only", uuid: n2.UUID, input: schemas.Note{Title: "updated title", Text: "updated text", UserUuid: u.UUID}, expected: 1, err: nil},
		{name: "text update with wrong user uuid", uuid: n2.UUID, input: schemas.Note{Title: "updated title", Text: "updated text", UserUuid: "6a3f5046-7d3f-11ed-a1eb-0242ac120002"}, expected: 0, err: errors.New("user not found")},
		{name: "no uuid", uuid: "", input: schemas.Note{Title: "updated title", Text: "updated text", UserUuid: u.UUID}, expected: 0, err: errors.New("invalid UUID")},
		{name: "empty update", uuid: n2.UUID, input: schemas.Note{Title: "", Text: "", UserUuid: u.UUID}, expected: 0, err: errors.New("empty updates")},
	}

	for _, tc := range tests {
		got, err := UpdateNote(tc.uuid, tc.input)
		if got != tc.expected {
			t.Log(tc.name)
			t.Fatalf("got: %d    expected:%d", got, tc.expected)
		}
		if fmt.Sprint(err) != fmt.Sprint(tc.err) {
			t.Log(tc.name)
			t.Fatalf("got: %s    expected:%s", err, tc.err)
		}
	}

}

func TestFindNotes(t *testing.T) {
	dbErr := database.Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer database.CleanUp()
	u, _ := CreateUser(schemas.User{Email: "test@test.com", Password: "sgsersdfvsdsfaaa"})
	CreateNote(schemas.Note{Title: "Test Title", Text: "Test Text", UserUuid: u.UUID})

	type test struct {
		name     string
		userUuid string
		expected *schemas.Note
		err      error
	}

	tests := []test{
		{name: "normal user uuid", userUuid: u.UUID, expected: &schemas.Note{Title: "test title", Text: "test text", UserUuid: u.UUID}, err: nil},
		{name: "no user uuid", userUuid: "", expected: nil, err: errors.New("invalid uuid")},
		{name: "wrong user uuid", userUuid: "6a3f5046-7d3f-11ed-a1eb-0242ac120002", expected: nil, err: errors.New("user not found")},
	}

	for _, tc := range tests {
		got, err := FindNotes(tc.userUuid)
		if fmt.Sprint(err) != fmt.Sprint(tc.err) {
			t.Log(tc.name)
			t.Fatalf("got: %s    expected:%s", err, tc.err)
		}
		if got != nil {
			for _, note := range got.Notes {
				if tc.expected != nil && got != nil {
					if validators.IsUUIDValid(note.UUID) == false {
						t.Log(tc.name)
						t.Fatalf("got invalid uuid: %s", note.UUID)
					}
					if note.Text != tc.expected.Text {
						t.Log(tc.name)
						t.Fatalf("got: %s    expected:%s", note.Text, tc.expected.Text)
					}
					if note.Title != tc.expected.Title {
						t.Log(tc.name)
						t.Fatalf("got: %s    expected:%s", note.Title, tc.expected.Title)
					}
				}
			}

		}
	}
}
