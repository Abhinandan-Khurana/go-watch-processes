package notifications

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Flag to ensure we only check for dependencies once per session.
var isDependencyChecked = false

// Init checks for required notification tools before the main app starts.
func Init() error {
	if isDependencyChecked {
		return nil
	}

	var err error
	switch runtime.GOOS {
	case "darwin":
		err = checkAndPromptMacOS()
	case "linux":
		err = checkAndPromptLinux()
	case "windows":
		// As requested, notifications are disabled for Windows.
		log.Println("Windows detected. Note: Desktop notifications are not supported on this platform.")
	default:
		err = fmt.Errorf("notifications not supported on this OS: %s", runtime.GOOS)
	}

	isDependencyChecked = true
	return err
}

// Notify sends a desktop notification using the appropriate tool for the OS.
func Notify(title, message string) {
	switch runtime.GOOS {
	case "darwin":
		notifyMacOS(title, message)
	case "linux":
		notifyLinux(title, message)
	case "windows":
		// Do nothing, as requested.
		return
	}
}

// isCommandAvailable checks if a command-line tool is in the user's PATH.
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// --- macOS Implementation ---
func checkAndPromptMacOS() error {
	if isCommandAvailable("terminal-notifier") {
		return nil
	}

	fmt.Println("'terminal-notifier' is required but not found.")
	fmt.Print("Would you like to try and install it via Homebrew? (y/n): ")

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(response)) != "y" {
		fmt.Println("Installation skipped. Notifications will be disabled.")
		return nil
	}

	fmt.Println("Running 'brew install terminal-notifier'...")
	cmd := exec.Command("brew", "install", "terminal-notifier")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install terminal-notifier: %w", err)
	}

	fmt.Println("'terminal-notifier' installed successfully.")
	return nil
}

func notifyMacOS(title, message string) {
	if !isCommandAvailable("terminal-notifier") {
		return // Silently fail if not installed.
	}
	cmd := exec.Command("terminal-notifier", "-title", title, "-message", message)
	cmd.Run()
}

// --- Linux Implementation ---
func checkAndPromptLinux() error {
	if isCommandAvailable("notify-send") {
		return nil
	}
	fmt.Println("Warning: 'notify-send' command not found.")
	fmt.Println("Notifications will be disabled. To enable them, please install the proper package for your distribution.")
	fmt.Println("Example for Debian/Ubuntu: sudo apt-get install libnotify-bin")
	fmt.Println("Example for Fedora/CentOS: sudo yum install libnotify")
	return nil
}

func notifyLinux(title, message string) {
	if !isCommandAvailable("notify-send") {
		return // Silently fail if not installed.
	}
	cmd := exec.Command("notify-send", title, message)
	cmd.Run()
}
