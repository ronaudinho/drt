//go:build integration
// +build integration

package valve

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	t.Run("download dota2 replay file should be success", func(t *testing.T) {
		// this replay url can be found using OpenDota API
		// for docs how to get the replayURL using OpenDota API
		// can be read here: https://docs.opendota.com/#tag/matches%2Fpaths%2F~1matches~1%7Bmatch_id%7D%2Fget
		replayURL := "http://replay153.valve.net/570/7132230434_1635105612.dem.bz2"

		r := NewDefaultReplay()

		ctx := context.Background()

		destination := "/tmp"

		err := r.Download(ctx, replayURL, destination)

		assert.NoError(t, err)
	})

	t.Run("download dota2 replay file should be failed due to invalid replay url", func(t *testing.T) {
		replayURL := "http://replay153.valve.net/570"

		r := NewDefaultReplay()

		ctx := context.Background()

		destination := "/tmp"

		err := r.Download(ctx, replayURL, destination)

		assert.Error(t, err)
		assert.Equal(t, err, ErrDownloadRelpay)
	})
}
