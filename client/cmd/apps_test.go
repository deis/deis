package cmd

import "testing"

func TestPrintLogLinesBadLine(t *testing.T) {
	t.Parallel()

	// Regression test for https://github.com/deis/deis/issues/4420
	logs := `\nDone preparing production files\n\n\u001b[4mRunning \"concat:plugins\" (concat) task\u001b[24m\n`
	if err := printLogs(logs); err != nil {
		t.Fatal(err)
	}

	logs = `\n\n\n`
	if err := printLogs(logs); err != nil {
		t.Fatal(err)
	}
}
