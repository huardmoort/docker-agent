package desktop

import (
	"context"
	"time"

	"github.com/kofalt/go-memoize"
)

// updateInfo is the subset of Docker Desktop's `/update` endpoint payload
// we care about. The endpoint exposes the running app's version and build
// number; the full body also contains appcast and update-status fields
// that are irrelevant here.
type updateInfo struct {
	CurrentVersion string `json:"currentVersion"`
	CurrentBuild   string `json:"currentBuild"`
}

// versionTTL bounds how long a previous lookup result (success or failure)
// is reused. A TTL — rather than the lifetime of the process — means:
//   - if docker-agent starts before Desktop is ready, version detection
//     recovers automatically once Desktop comes up;
//   - if Desktop is upgraded mid-session, the new version is picked up
//     within at most this interval.
//
// 5 minutes is generous enough to keep the per-request overhead negligible
// (the Unix socket call only happens once per window) while still giving
// quick recovery in the failure case.
const versionTTL = 5 * time.Minute

var versionMemoizer = memoize.NewMemoizer(versionTTL, 2*versionTTL)

// GetVersion returns the running Docker Desktop version (e.g. "4.74.0") or
// an empty string if Docker Desktop is not running or the call fails.
//
// The lookup is memoized for [versionTTL] and is bounded by a short
// internal timeout so a stale or missing backend socket cannot stall
// callers on hot paths (it is queried on every outbound built-in tool
// HTTP request). The HTTP call uses [context.Background] rather than a
// caller-supplied context: the result is shared across all callers, so
// we don't want the first caller's deadline or cancellation to bleed
// into other callers' view of the world.
func GetVersion() string {
	v, _, _ := versionMemoizer.Memoize("desktopVersion", func() (any, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		var info updateInfo
		_ = ClientBackend.Get(ctx, "/update", &info)
		return info.CurrentVersion, nil
	})
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}
