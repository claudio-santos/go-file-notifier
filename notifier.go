package main

import (
	"path/filepath"
	"strings"

	"gopkg.in/toast.v1"
)

// Notifier manages notifications and associated actions
type Notifier struct {
	appName string
}

// NewNotifier creates a new notifier
func NewNotifier(appName string) *Notifier {
	return &Notifier{
		appName: appName,
	}
}

// Notify sends a notification to the user with action buttons
func (n *Notifier) Notify(filePath string, eventType string) {
	fileName := filepath.Base(filePath)

	// Send Windows Toast notification with buttons
	go func() {
		n.sendToastNotification(fileName, filePath, eventType)
	}()
}

// sendToastNotification sends Toast notification with action buttons
func (n *Notifier) sendToastNotification(fileName, filePath, eventType string) {
	// Convert path to file:// URI format
	fileURI := "file:///" + filepath.ToSlash(filePath)
	folderURI := "file:///" + filepath.ToSlash(filepath.Dir(filePath))

	notification := toast.Notification{
		AppID:   n.appName,
		Title:   escapeXML(eventType),
		Message: escapeXML(fileName),
		Actions: []toast.Action{
			{Type: "protocol", Label: "Open File", Arguments: fileURI},
			{Type: "protocol", Label: "Open Folder", Arguments: folderURI},
		},
		Duration: toast.Long,
	}

	notification.Push()
}

// escapeXML escapes special XML characters
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, `'`, "&apos;")
	return s
}
