package notifier

import (
	"fmt"
	"os/exec"
	"runtime"
)

// DesktopNotifier sends desktop notifications using platform-native tools.
type DesktopNotifier struct {
	AppName string
}

// NewDesktopNotifier creates a new DesktopNotifier with the given app name.
func NewDesktopNotifier(appName string) *DesktopNotifier {
	if appName == "" {
		appName = "portwatch"
	}
	return &DesktopNotifier{AppName: appName}
}

// Send dispatches a desktop notification with the given title and message.
func (d *DesktopNotifier) Send(title, message string) error {
	switch runtime.GOOS {
	case "linux":
		return d.sendLinux(title, message)
	case "darwin":
		return d.sendDarwin(title, message)
	case "windows":
		return d.sendWindows(title, message)
	default:
		return fmt.Errorf("desktop notifications not supported on %s", runtime.GOOS)
	}
}

func (d *DesktopNotifier) sendLinux(title, message string) error {
	path, err := exec.LookPath("notify-send")
	if err != nil {
		return fmt.Errorf("notify-send not found: %w", err)
	}
	cmd := exec.Command(path, "-a", d.AppName, title, message)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("notify-send failed: %w, output: %s", err, string(out))
	}
	return nil
}

func (d *DesktopNotifier) sendDarwin(title, message string) error {
	script := fmt.Sprintf(
		`display notification %q with title %q subtitle %q`,
		message, d.AppName, title,
	)
	cmd := exec.Command("osascript", "-e", script)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("osascript failed: %w, output: %s", err, string(out))
	}
	return nil
}

func (d *DesktopNotifier) sendWindows(title, message string) error {
	script := fmt.Sprintf(
		`[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null; `+
			`$template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02); `+
			`$template.SelectSingleNode('//text[@id=1]').InnerText = '%s'; `+
			`$template.SelectSingleNode('//text[@id=2]').InnerText = '%s'; `+
			`[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('%s').Show([Windows.UI.Notifications.ToastNotification]::new($template))`,
		title, message, d.AppName,
	)
	cmd := exec.Command("powershell", "-Command", script)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("powershell toast failed: %w, output: %s", err, string(out))
	}
	return nil
}
