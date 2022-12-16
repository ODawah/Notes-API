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

	type test struct {
		name     string
		input    schemas.Note
		expected *schemas.Note
		err      error
	}

	tests := []test{
		{name: "normal note input", input: schemas.Note{Title: "pick up your son", Text: "go to school and pick him up"}, expected: &schemas.Note{Title: "pick up your son", Text: "go to school and pick him up"}, err: nil},
		{name: "long note title", input: schemas.Note{Title: longTitle, Text: "go to school and pick him up"}, expected: nil, err: errors.New("long note title")},
		{name: "no note title", input: schemas.Note{Title: "", Text: "go to school and pick him up"}, expected: nil, err: errors.New("empty note title")},
		{name: "long note text", input: schemas.Note{Title: "finish the assignment", Text: longText}, expected: nil, err: errors.New("long note text")},
		{name: "no note text", input: schemas.Note{Title: "cook the dinner", Text: ""}, expected: nil, err: errors.New("empty note text")},
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
	CreateNote(schemas.Note{Title: "Test Title", Text: "Test Text"})
	CreateNote(schemas.Note{Title: "finish the tasks", Text: "create the API"})

	type test struct {
		name        string
		searchTitle string
		expected    *schemas.Note
		err         error
	}

	tests := []test{
		{name: "normal note title", searchTitle: "Test Title", expected: &schemas.Note{Title: "test title", Text: "test text"}, err: nil},
		{name: "normal note title with trails", searchTitle: "   finish the tasks", expected: &schemas.Note{Title: "finish the tasks", Text: "create the api"}, err: nil},
		{name: "no note title", searchTitle: "", expected: nil, err: errors.New("empty note title")},
		{name: "long note title", searchTitle: longTitle, expected: nil, err: errors.New("long note title")},
	}

	for _, tc := range tests {
		got, err := FindNoteByTitle(tc.searchTitle)
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
		}
	}

}
