# ⚙️ Advanced Process Monitor CLI

An intelligent, cross-platform command-line tool that monitors system processes and provides meaningful desktop notifications for events that matter. Forget the noise; get alerted only for long-running tasks, high resource usage, and suspicious activity.

### Screenshots

### TUI -
![tui](./static/tui.png)

#### Notification -
![notification](./static/notification-2.png)
![notification](./static/notification.png)

---

## ✨ Core Features

- **🧠 Intelligent Notifications**: Only get notified when a task running for **more than 10 seconds** completes. Perfect for monitoring builds, scripts, and long commands.
- **🚨 Advanced Alerts**:
  - **High CPU Usage**: Get an alert if a process consumes high CPU for a sustained period.
  - **Suspicious Name Detection**: Be notified if a process starts with a name commonly associated with malware or miners (e.g., `miner`, `keylog`).
- **🔇 Smart Filtering**: Automatically ignores common, noisy background processes like `Google Chrome Helper`, `mdworker`, and `Code Helper` to keep your notifications relevant.
- **🖥️ Cross-Platform**: Works seamlessly on **macOS** and **Linux**. The core application runs on Windows, but desktop notifications are currently limited to macOS/Linux.
- **💅 Polished Interface**: A clean and beautiful terminal UI built with `Bubble Tea` that displays a live log of important events.

## 🚀 Installation

This tool relies on native command-line notifiers for each OS. Please ensure they are installed first.

#### 1. Install Notification Dependencies

- **On macOS** (using [Homebrew](https://brew.sh/)):

  ```sh
  brew install terminal-notifier
  ```

  _(The tool will prompt you to do this on first run if it's missing.)_

- **On Linux** (Debian/Ubuntu):

  ```sh
  sudo apt-get update && sudo apt-get install libnotify-bin
  ```

  _(For other distributions like Fedora/CentOS, use `sudo yum install libnotify`.)_

#### 2. Install the go-watch-processes

- **From Source:**

  ```bash
  git clone https://github.com/Abhinandan-Khurana/go-watch-processes.git
  cd go-watch-processes
  go build
  ```

- **From Releases:**
  Head over to the [**Releases**](https://github.com/Abhinandan-Khurana/go-watch-processes/releases) page to download a pre-compiled binary for your operating system.

- **Direct Installation using golang:**

```bash
go install -v github.com/Abhinandan-Khurana/go-watch-processes@latest
```

## ▶️ Usage

Simply run the executable from your terminal.

```bash
./go-watch-processes
```

The tool will take over the terminal window to display its event log. Leave it running in the background to monitor your system.

Press `q` or `Ctrl+C` to quit.

**Example TUI Output:**

```

                        ▐     ▌
        ▞▀▌▞▀▖▄▄▖▌  ▌▝▀▖▜▀ ▞▀▖▛▀▖▄▄▖▛▀▖▙▀▖▞▀▖▞▀▖▞▀▖▞▀▘▞▀▘▞▀▖▞▀▘
        ▚▄▌▌ ▌   ▐▐▐ ▞▀▌▐ ▖▌ ▖▌ ▌   ▙▄▘▌  ▌ ▌▌ ▖▛▀ ▝▀▖▝▀▖▛▀ ▝▀▖
        ▗▄▘▝▀     ▘▘ ▝▀▘ ▀ ▝▀ ▘ ▘   ▌  ▘  ▝▀ ▝▀ ▝▀▘▀▀ ▀▀ ▝▀▘▀▀


        By: Abhinandan Khurana | Version: v1.0.0

Monitoring for long-running tasks and suspicious activity... Press 'q' to quit.

[18:24:45] ✔ Task Completed: sleep(PID: 23767, Duration: 11s)
[18:24:58] ✔ Task Completed: screencaptureui (PID: 23883, Duration: 10s)
[18:25:45] ❗ Alert: Suspicious Process Detected AssetCacheLocatorService (PID: 24322)[Matched keyword: 'tor']
[18:26:08] ✔ Task Completed: AddressBookSourceSync (PID: 24240, Duration: 31s)
[18:26:35] ✔ Task Completed: (PID: 24458, Duration: 31s)
[18:26:42] ✔ Task Completed: screencaptureui (PID: 24055, Duration: 1m33s)
[18:26:42] ✔ Task Completed: (PID: 24422, Duration: 44s)
```

## 🔧 Configuration

The list of ignored processes and suspicious keywords are currently hardcoded for simplicity. To customize them, modify the `ignoreList` and `suspiciousKeywords` variables at the top of the `monitor/watcher.go` file and re-build the application.

```go
// monitor/watcher.go

var ignoreList = []string{
 "Google Chrome Helper",
 "mdworker_shared",
 // Add your custom ignored processes here
}

var suspiciousKeywords = []string{"miner", "keylog", "tor", /* ... */}
```
