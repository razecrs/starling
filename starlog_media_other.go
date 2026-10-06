//go:build !windows && !linux

package starling

import "fmt"

type unavailableStarlogMedia struct{}

func newStarlogPlatformMedia() (starlogPlatformMedia, error) {
	return nil, fmt.Errorf("starling: system media detection is unavailable on this platform")
}

func (*unavailableStarlogMedia) read() (StarlogPlayback, error) { return StarlogPlayback{}, nil }
func (*unavailableStarlogMedia) close()                         {}
