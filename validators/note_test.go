package validators

import (
	"errors"
	"fmt"
	"testing"
)

func TestIsValidText(t *testing.T) {
	long := "I'm selfish, I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard to handle. But if you can't handle me at my worst, then you sure as hell don' deserve me at my b.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hardimpatient and a little insecure. I make mistakes, I am out of control and at times hard to handle. But if you can't handle me at my worst, then you sure as hell don't deserve me at my best.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard"

	type test struct {
		name     string
		input    string
		expected error
	}

	tests := []test{
		{name: "normal text", input: "keep the room clean", expected: nil},
		{name: "no text", input: "", expected: errors.New("empty note text")},
		{name: "long text", input: long, expected: errors.New("long note text")},
	}
	for _, tc := range tests {
		got := IsValidText(tc.input)
		if fmt.Sprint(got) != fmt.Sprint(tc.expected) {
			t.Log(tc.name)
			t.Fatalf("got: %s   expected:%s", got.Error(), tc.expected.Error())
		}
	}

}

func TestIsValidTitle(t *testing.T) {
	long := "then you sure as hell don't deserve me at my best.I'm selfish, impatient and a little insecure. I make mistakes, I am out of control and at times hard"

	type test struct {
		name     string
		input    string
		expected error
	}

	tests := []test{
		{name: "normal title", input: "keep the room clean", expected: nil},
		{name: "no title", input: "", expected: errors.New("empty note title")},
		{name: "long title", input: long, expected: errors.New("long note title")},
	}
	for _, tc := range tests {
		got := IsValidTitle(tc.input)
		if fmt.Sprint(got) != fmt.Sprint(tc.expected) {
			t.Log(tc.name)
			t.Fatalf("got: %s   expected:%s", got.Error(), tc.expected.Error())
		}
	}

}

func TestIsUUIDValid(t *testing.T) {
	type test struct {
		name     string
		input    string
		expected bool
	}

	tests := []test{
		{name: "Valid uuid", input: "6a3f5046-7d3f-11ed-a1eb-0242ac120002", expected: true},
		{name: "no uuid", input: "", expected: false},
		{name: "invalid uuid", input: "6a3f5046-7d3f-11ed2-a1eb-0242ac12000", expected: false},
	}

	for _, tc := range tests {
		got := IsUUIDValid(tc.input)
		if got != tc.expected {
			t.Log(tc.name)
			t.Fatalf("got: %v     expected: %v", got, tc.expected)
		}
	}

}
