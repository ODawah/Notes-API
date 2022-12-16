package validators

import (
	"testing"
)

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
