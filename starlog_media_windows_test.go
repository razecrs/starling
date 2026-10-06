//go:build windows

package starling

import (
	"os"
	"testing"
)

func TestWindowsSystemMediaIntegration(t *testing.T) {
	if os.Getenv("STARLING_TEST_SYSTEM_MEDIA") == "" {
		t.Skip("set STARLING_TEST_SYSTEM_MEDIA=1 to query the desktop media session")
	}
	reader, err := newStarlogPlatformMedia()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.close()
	playback, err := reader.read()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("provider=%q title=%q artist=%q playing=%v position=%s duration=%s",
		playback.Provider, playback.Title, playback.Artist, playback.Playing,
		playback.Position, playback.Duration)
}
