# Go File Notifier

A lightweight file monitoring tool for Windows that watches specified files and sends desktop notifications when changes are detected.

## Features

- File Monitoring - Polls specified files at configurable intervals
- Toast Notifications - Windows 10/11 native notifications with action buttons
- Quick Actions - Open file or folder directly from notification
- Persistent State - Detects changes that occurred while the program was not running
- Lightweight - Simple, minimal console output, no dependencies beyond Go standard library + toast notifications
- Graceful Shutdown - Clean exit on Ctrl+C with state preservation

## Installation

### Prerequisites

- Go 1.21 or later
- Windows 10/11 (for Toast notifications)

### Build from Source

```bash
git clone <repository-url>
cd go-file-notifier
go mod tidy
go build -o go-file-notifier.exe
```

## Usage

### 1. Configure Files to Monitor

Edit `config.json` with the files you want to monitor:

```json
{
  "files": [
    "C:\\Users\\YourName\\Documents\\myfile.txt",
    "C:\\Path\\To\\Another\\file.log"
  ],
  "interval_s": 5
}
```

| Field | Type | Description |
|-------|------|-------------|
| `files` | array | List of absolute file paths to monitor |
| `interval_s` | int | Polling interval in seconds (default: 5) |

### 2. Run the Program

```bash
.\go-file-notifier.exe
```

### 3. Monitor Output

```
Monitoring 1 file(s), interval: 5s
File changed: myfile.txt
```

When a file changes, you'll receive:
- Console notification
- Windows Toast notification with **Open File** and **Open Folder** buttons

### 4. Stop the Program

Press `Ctrl+C` to gracefully shutdown. The program will save its state before exiting.

## Project Structure

```
go-file-notifier/
├── main.go           # Entry point, signal handling
├── config.go         # Configuration loading and validation
├── state.go          # Persistent state management
├── watcher.go        # File monitoring logic (polling)
├── notifier.go       # Windows Toast notifications
├── config.json       # User configuration
├── state.json        # Auto-generated state file
├── go.mod            # Go module dependencies
└── go-file-notifier.exe  # Compiled binary
```

## How It Works

### Startup
1. Load `config.json` with file list and polling interval
2. Load `state.json` (if exists) with previous file states
3. Check for changes that occurred while program was not running
4. Start polling loop

### Monitoring Loop
1. Every `interval_s` seconds, check each file's size and modification time
2. If changed → send Toast notification + update state
3. Debounce: prevent multiple notifications within 1 second
4. Save state after each change

### Shutdown
1. On Ctrl+C, save current state to `state.json`
2. Clean exit

## State File Format

`state.json` is auto-generated and tracks file states:

```json
{
  "files": {
    "C:\\Users\\...\\myfile.txt": {
      "last_modified": "2026-03-16T12:00:00Z",
      "size": 1024
    }
  }
}
```

| Field | Description |
|-------|-------------|
| `last_modified` | File's last modification timestamp |
| `size` | File size in bytes |

## Notifications

When a file change is detected, a Windows Toast notification appears with:

- **Title**: "File Changed"
- **Message**: Filename
- **Buttons**:
  - **Open File** - Opens the file with its default application
  - **Open Folder** - Opens Windows Explorer to the file's folder

## Error Handling

- If a file becomes inaccessible, an error is logged to console
- Monitoring continues for other files
- State is preserved on unexpected shutdown

## Dependencies

| Package | Purpose |
|---------|---------|
| `gopkg.in/toast.v1` | Windows Toast notifications |
| `github.com/nu7hatch/gouuid` | UUID generation (toast dependency) |

## Limitations

- **Windows only** - Toast notifications require Windows 10/11
- **Polling-based** - Uses periodic checks instead of filesystem events (more reliable across editors)
- **No recursive monitoring** - Only monitors specified files, not directories

## Troubleshooting

### Notifications not appearing
- Ensure Windows Toast notifications are enabled for the app
- Check Windows Focus Assist settings (may block notifications)

### File changes not detected
- Increase polling frequency (lower `interval_s` in config)
- Verify file path is correct and accessible

### Build errors
- Ensure Go 1.21+ is installed: `go version`
- Run `go mod tidy` to fetch dependencies
