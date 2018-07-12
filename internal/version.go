package internal

import "fmt"

const (
	Name = "backlog-bulk-issue-registration"
	version = "0.9.0"
)

var (
	commit string
	date string
)

func FmtVersion() string {
	if commit == "" || date == "" {
		return version
	} else {
		return fmt.Sprintf("%s, build %s, date %s", version, commit, date)
	}
}