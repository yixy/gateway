package version

import "fmt"

var (
	Ver       = "unknown"
	Env       = "unknown"
	BuildTime = "unknown"
)

func Print() {
	fmt.Println("version", Ver)
	fmt.Println("environment", Env)
	fmt.Println("buildtime", BuildTime)
}
