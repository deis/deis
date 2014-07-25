// Package unit provides unit tests for the Deis project.
package unit

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

const (
	COMMIT_SUBJECT_MAX_LINE_LENGTH int    = 50
	COMMIT_MESSAGE_MAX_LINE_LENGTH int    = 72
	POLICY_DATE                    string = "2014-07-11"
)

type CommitMessage struct {
	Message string
}

func (c CommitMessage) Summary() string {
	return strings.Split(c.Message, "\n")[0]
}

/*
TestCommitsAdhereToCommitStyleGuide does exactly as it says; ensures that all commits
adhere to Deis' commit style guide. It retrieves all commits after a specified policy date
and checks that each commit follows the style guide.

See the style guide for more information on what this function should be checking for:
http://docs.deis.io/en/latest/contributing/standards/#commit-style-guide

Knowledge on retrieving each commit from `git log` using ASCII field/record separators
came from http://blog.lost-theory.org/post/how-to-parse-git-log-output/
*/
func TestCommitsAdhereToCommitStyleGuide(t *testing.T) {
	var validCommitSubject = regexp.MustCompile(`^feat\(.+\): [0-9a-z].+[^.]$|^fix\(.+\): [0-9a-z].+[^.]$|^docs\(.+\): [0-9a-z].+[^.]$|^style\(.+\): [0-9a-z].+[^.]$|^ref\(.+\): [0-9a-z].+[^.]$|^test\(.+\): [0-9a-z].+[^.]$|^chore\(.+\): [0-9a-z].+[^.]$`)

	cmd := exec.Command("git", "log", "--no-merges", "--format=%s%n%b%x1e", "--after", POLICY_DATE)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf(err.Error())
	}
	logs := string(out)
	for _, row := range strings.Split(strings.Trim(logs, "\n\x1e"), "\x1e") {
		commit := CommitMessage{Message: strings.Trim(row, "\n")}
		if !validCommitSubject.MatchString(commit.Summary()) {
			t.Errorf(fmt.Sprintf("not a valid subject: %v", commit.Message))
		}
		if len(commit.Summary()) > COMMIT_SUBJECT_MAX_LINE_LENGTH {
			t.Errorf(fmt.Sprintf(
				"subject cannot be longer than %d characters in length: %s",
				COMMIT_SUBJECT_MAX_LINE_LENGTH,
				commit.Summary()))
		}
		for _, line := range strings.Split(commit.Message, "\n") {
			// a line of source code can be as long as it wants to be
			if strings.HasPrefix(line, "    ") {
				continue
			}
			if len(line) > COMMIT_MESSAGE_MAX_LINE_LENGTH {
				t.Errorf(fmt.Sprintf(
					"lines cannot be longer than %v characters in length: %s",
					COMMIT_MESSAGE_MAX_LINE_LENGTH,
					commit.Summary()))
				break
			}
		}
	}
}
