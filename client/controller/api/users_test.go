package api

import (
	"sort"
	"testing"
)

func TestUsersSorted(t *testing.T) {
	users := Users{
		{1, "", false, "Zulu", "", "", "", false, false, ""},
		{2, "", false, "Beta", "", "", "", false, false, ""},
		{3, "", false, "Gamma", "", "", "", false, false, ""},
		{4, "", false, "Alpha", "", "", "", false, false, ""},
	}

	sort.Sort(users)
	expectedUsernames := []string{"Alpha", "Beta", "Gamma", "Zulu"}

	for i, user := range users {
		if expectedUsernames[i] != user.Username {
			t.Errorf("Expected users to be sorted %v, Got %v at index %v", expectedUsernames[i], user.Username, i)
		}
	}
}
