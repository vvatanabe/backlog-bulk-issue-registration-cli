package bbir

import "fmt"

const (
	Name    = "bbir"
	version = "1.0.0"
)

var (
	commit string
	date   string
)

func FmtVersion() string {
	if commit == "" || date == "" {
		return version
	} else {
		return fmt.Sprintf("%s, build %s, date %s", version, commit, date)
	}
}
