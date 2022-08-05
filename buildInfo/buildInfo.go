package buildInfo

var Version string = "0.1"
var Commit string = "local"

func GetVersion() string {
	return Version + "-" + Commit
}
