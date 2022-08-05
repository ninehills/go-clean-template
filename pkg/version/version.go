package version

import (
	"fmt"
	"sync"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	branch  = "unknown"
	buildBy = "unknown"
)

type Version struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	Branch  string `json:"branch"`
	BuildBy string `json:"buildBy"`
}

func (v *Version) String() string {
	return fmt.Sprintf("version: %s, commit: %s, date: %s, branch: %s, buildBy: %s", v.Version, v.Commit, v.Date, v.Branch, v.BuildBy)
}

var (
	once sync.Once
	v    *Version
)

func GetVersion() *Version {
	once.Do(func() {
		v = &Version{
			Version: version,
			Commit:  commit,
			Date:    date,
			Branch:  branch,
			BuildBy: buildBy,
		}
	})
	return v
}
