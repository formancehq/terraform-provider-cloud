package internal

import "fmt"

const (
	Repository  = "formancehq/tf-cloud-provider"
	ServiceName = "terraform-provider-server"
	Version     = "develop"
	BuildDate   = "-"
	Commit      = "-"
)

var App = AppInfo{
	Name:       ServiceName,
	Repository: Repository,
	Version:    Version,
	BuildDate:  BuildDate,
	Commit:     Commit,
}

type AppInfo struct {
	Name       string
	Repository string
	Version    string
	BuildDate  string
	Commit     string
}

func (a AppInfo) String() string {
	return fmt.Sprintf("\n\tName: %s\n\tVersion: %s\n\tBuildDate: %s\n\tCommit: %s\n\tRepository: %s\n\t", a.Name, a.Version, a.BuildDate, a.Commit, a.Repository)
}
