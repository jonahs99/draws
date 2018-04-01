package draws

import (
	"fmt"
	"os/exec"
	"runtime"
)

// source: stack overflow #39320371
// openBrowser opens the specified URL in the default browser of the user.
func openBrowser(url string) error {
	fmt.Printf("OPEN BROWSER: %s\n", url)

	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
