package pkg

import "testing"

func TestTrimToDashes(t *testing.T) {
	var argtests = []struct {
		input  []string
		output []string
	}{
		{[]string{"foo", "bar", "baz"}, []string{"foo", "bar", "baz"}},
		{[]string{"abc", "def", "--", "ghi"}, []string{"ghi"}},
		{[]string{"abc", "def", "--"}, []string{}},
		{[]string{"--"}, []string{}},
		{[]string{"--", "abc", "def", "ghi"}, []string{"abc", "def", "ghi"}},
		{[]string{"--", "bar", "--", "baz"}, []string{"bar", "--", "baz"}},
		{[]string{"--flagname", "--", "ghi"}, []string{"ghi"}},
		{[]string{"--", "--flagname", "ghi"}, []string{"--flagname", "ghi"}},
	}
	for _, test := range argtests {
		args := TrimToDashes(test.input)
		if len(test.output) != len(args) {
			t.Fatalf("error trimming dashes: expected %s, got %s", test.output, args)
		}
		for i, v := range args {
			if v != test.output[i] {
				t.Fatalf("error trimming dashes: expected %s, got %s", test.output, args)
			}
		}
	}
}
