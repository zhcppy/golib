package version

import (
	"fmt"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var Name string
var Version string
var Commit string

func NewCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   fmt.Sprintf("Show %s version", Name),
		Example: fmt.Sprintf("%s version", Name),
		Run: func(*cobra.Command, []string) {
			fmt.Printf("Version: %s\n", Version)
			fmt.Printf("Git Commit: %s\n", Commit)
			fmt.Printf("Go Version: %s\n", runtime.Version())
			fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
			fmt.Printf("Build Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		},
	}
	return versionCmd
}
