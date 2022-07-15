package buildInfo

var Version string = "1.0a"
var Commit string = "local"

func GetVersion() string {
	return Version + "-" + Commit
}
