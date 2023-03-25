package models

import (
	"errors"
	"fmt"
	"testing"
)

func TestCreateNote(t *testing.T) {
	dbErr := Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer CleanUp()
	longTitle := "then you sure as hell don't deserve me at my best.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard"
	longText := "I'm selfish, I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard to handle. But if you can't handle me at my worst, then you sure as hell don' deserve me at my b.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hardimpatient and a little insecure. I make mistakes, I am out of control and at times hard to handle. But if you can't handle me at my worst, then you sure as hell don't deserve me at my best.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard"

	u, _ := CreateUser(User{Email: "test@test.com", Password: "test password"})

	type test struct {
		name     string
		input    Note
		expected *Note
		err      error
	}

	tests := []test{
		{name: "normal note input", input: Note{Title: "pick up your son", Text: "go to school and pick him up", UserUuid: u.UUID}, expected: &Note{Title: "pick up your son", Text: "go to school and pick him up", UserUuid: u.UUID}, err: nil},
		{name: "long note title", input: Note{Title: longTitle, Text: "go to school and pick him up"}, expected: nil, err: errors.New("long note title")},
		{name: "no note title", input: Note{Title: "", Text: "go to school and pick him up"}, expected: nil, err: errors.New("empty note title")},
		{name: "long note text", input: Note{Title: "finish the assignment", Text: longText}, expected: nil, err: errors.New("long note text")},
		{name: "no note text", input: Note{Title: "cook the dinner", Text: "", UserUuid: u.UUID}, expected: nil, err: errors.New("empty note text")},
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
			t.Fatalf("got: %v    expected:%v", got, tc.err)
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
			if IsUUIDValid(got.UUID) == false {
				t.Log(tc.name)
				t.Fatalf("got invalid uuid: %s", got.UUID)
			}
			if IsUUIDValid(got.UserUuid) == false && tc.expected.UserUuid != got.UserUuid {
				t.Log(tc.name)
				t.Fatalf("got invalid user uuid: %s", got.UUID)
			}
		}

	}

}

func TestFindNoteByTitle(t *testing.T) {
	longTitle := "then you sure as hell don't deserve me at my best.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard"
	dbErr := Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer CleanUp()
	u, _ := CreateUser(User{Email: "test@test.com", Password: "sgsersdfvsdsfaaa"})
	CreateNote(Note{Title: "Test Title", Text: "Test Text", UserUuid: u.UUID})
	CreateNote(Note{Title: "finish the tasks", Text: "create the API", UserUuid: u.UUID})

	type test struct {
		name        string
		searchTitle string
		uuid        string
		expected    *Note
		err         error
	}

	tests := []test{
		{name: "normal note title", searchTitle: "Test Title", uuid: u.UUID, expected: &Note{Title: "test title", Text: "test text", UserUuid: u.UUID}, err: nil},
		{name: "normal note title with trails", searchTitle: "   finish the tasks", uuid: u.UUID, expected: &Note{Title: "finish the tasks", Text: "create the api", UserUuid: u.UUID}, err: nil},
		{name: "normal note title with wrong user uuid", searchTitle: "Test Title", uuid: "123e4567-e89b-12d3-a456-42661417s4000", expected: nil, err: errors.New("invalid uuid")},
		{name: "no note title", searchTitle: "", uuid: u.UUID, expected: nil, err: errors.New("empty note title")},
		{name: "long note title", searchTitle: longTitle, uuid: u.UUID, expected: nil, err: errors.New("long note title")},
	}

	for _, tc := range tests {
		got, err := FindNoteByTitle(tc.searchTitle, tc.uuid)
		if fmt.Sprint(err) != fmt.Sprint(tc.err) {
			t.Log(tc.name)
			t.Fatalf("got: %s  expected:%s", err, tc.err)
		}
		if tc.expected != nil && got != nil {
			if IsUUIDValid(got.UUID) == false {
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
	dbErr := Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer CleanUp()
	u, _ := CreateUser(User{Email: "test@test.com", Password: "sgsersdfvsdsfaaa"})
	n1, _ := CreateNote(Note{Title: "Test Title", Text: "Test Text", UserUuid: u.UUID})
	n2, _ := CreateNote(Note{Title: "finish the tasks", Text: "create the API", UserUuid: u.UUID})

	type test struct {
		name  string
		input string
		uuid  string
		err   error
	}

	tests := []test{
		{name: "normal note uuid", input: n1.UUID, uuid: u.UUID, err: nil},
		{name: "normal note wrong uuid", input: n2.UUID, uuid: "6a3f5046-7d3f-11ed-a1eb-0242ac120002", err: errors.New("user not found")},
		{name: "normal note uuid", input: n2.UUID, uuid: u.UUID, err: nil},
		{name: "no uuid", input: "", uuid: u.UUID, err: errors.New("invalid UUID")},
		{name: "uuid not in db", input: "6a3f5046-7d3f-11ed-a1ed-0242ac130002", uuid: u.UUID, err: errors.New("note not found")},
	}

	for _, tc := range tests {
		err := DeleteNote(tc.input, tc.uuid)
		if fmt.Sprint(err) != fmt.Sprint(tc.err) {
			t.Log(tc.name)
			t.Fatalf("got: %v    expected:%v", err, tc.err)
		}
	}

}

func TestUpdateNote(t *testing.T) {
	dbErr := Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer CleanUp()
	u, _ := CreateUser(User{Email: "test@test.com", Password: "sgsersdfvsdsfaaa"})
	n1, _ := CreateNote(Note{Title: "Test Title", Text: "Test Text", UserUuid: u.UUID})
	n2, _ := CreateNote(Note{Title: "finish the tasks", Text: "create the API", UserUuid: u.UUID})

	type test struct {
		name  string
		uuid  string
		input Note
		err   error
	}

	tests := []test{
		{name: "normal note updates", input: Note{UUID: n1.UUID, Title: "updated title", Text: "updated text", UserUuid: u.UUID}, err: nil},
		{name: "title update only", input: Note{UUID: n1.UUID, Title: "updated title", Text: "updated text", UserUuid: u.UUID}, err: nil},
		{name: "text update only", input: Note{UUID: n2.UUID, Title: "updated title", Text: "updated text", UserUuid: u.UUID}, err: nil},
		{name: "text update with wrong user uuid", input: Note{UUID: n2.UUID, Title: "updated title", Text: "updated text", UserUuid: "6a3f5046-7d3f-11ed-a1eb-0242ac120002"}, err: errors.New("user not found")},
		{name: "no uuid", input: Note{UUID: "", Title: "updated title", Text: "updated text", UserUuid: u.UUID}, err: errors.New("invalid UUID")},
		{name: "empty update", input: Note{UUID: n2.UUID, Title: "", Text: "", UserUuid: u.UUID}, err: nil},
	}

	for _, tc := range tests {
		err := UpdateNote(tc.input)
		if fmt.Sprint(err) != fmt.Sprint(tc.err) {
			t.Log(tc.name)
			t.Fatalf("got: %s    expected:%s", err, tc.err)
		}
	}

}

func TestFindNotes(t *testing.T) {
	dbErr := Connect()
	if dbErr != nil {
		t.Fatal(dbErr)
	}
	defer CleanUp()
	u, _ := CreateUser(User{Email: "test@test.com", Password: "sgsersdfvsdsfaaa"})
	CreateNote(Note{Title: "Test Title", Text: "Test Text", UserUuid: u.UUID})
	CreateNote(Note{Title: "Test Title2", Text: "Test Text2", UserUuid: u.UUID})

	type test struct {
		name     string
		userUuid string
		expected []Note
		err      error
	}

	tests := []test{
		{name: "normal user uuid", userUuid: u.UUID, expected: []Note{{Title: "test title", Text: "test text", UserUuid: u.UUID}, {Title: "test title2", Text: "test text2", UserUuid: u.UUID}}, err: nil},
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
			for i, note := range got {
				if tc.expected != nil && got != nil {
					if IsUUIDValid(note.UUID) == false {
						t.Log(tc.name)
						t.Fatalf("got invalid uuid: %s", note.UUID)
					}
					if note.Text != tc.expected[i].Text {
						t.Log(tc.name)
						t.Fatalf("got: %s    expected:%s", note.Text, tc.expected[i].Text)
					}
					if note.Title != tc.expected[i].Title {
						t.Log(tc.name)
						t.Fatalf("got: %s    expected:%s", note.Title, tc.expected[i].Title)
					}
				}
			}
		}
	}
}
