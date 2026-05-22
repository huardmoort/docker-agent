//go:build !darwin

package editor

// readClipboardImage is a no-op stub for non-Darwin platforms.
// Image-from-clipboard paste is only supported on macOS.
func readClipboardImage() (string, error) {
	return "", nil
}
