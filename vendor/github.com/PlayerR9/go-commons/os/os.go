package os

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var (
	// clear is a map of the clear functions for each platform.
	clear map[string]func()
)

func init() {
	clear = make(map[string]func())

	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// ClearScreen clears the screen.
//
// Returns:
//   - error: An error if the platform is not supported.
func ClearScreen() error {
	f, ok := clear[runtime.GOOS]
	if !ok {
		return fmt.Errorf("platform %q is yet to be supported by this clear screen function", runtime.GOOS)
	}

	f()

	return nil
}
