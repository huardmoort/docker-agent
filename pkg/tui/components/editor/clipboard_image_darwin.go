//go:build darwin

package editor

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// readClipboardImage attempts to read an image from the macOS clipboard using
// osascript. If the clipboard contains an image, it is written to a temp file
// and the path is returned. If the clipboard has no image (or osascript fails),
// ("", nil) is returned so the caller can fall through to text-paste behaviour.
func readClipboardImage() (string, error) {
	tmpPath := fmt.Sprintf("%s/clipboard-image-%d.png", os.TempDir(), time.Now().UnixNano())

	script := fmt.Sprintf(`
set imgData to the clipboard as «class PNGf»
set tmpFile to open for access POSIX file %q with write permission
write imgData to tmpFile
close access tmpFile
`, tmpPath)

	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		// No image in clipboard (or osascript unavailable) — silent fallback.
		_ = os.Remove(tmpPath) // clean up any partial file
		return "", nil
	}

	return tmpPath, nil
}
