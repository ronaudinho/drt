//go:build integration
// +build integration

package opendota

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchMatchDetail(t *testing.T) {
	ctx := context.Background()
	openDotaAPI := NewDefaultAPI()

	t.Run("fetch match detail dota2 should be success", func(t *testing.T) {
		// my personal match id
		matchID := 271145478

		matchDetail, err := openDotaAPI.FetchMatchDetail(ctx, int64(matchID))

		assert.NoError(t, err)
		assert.NotEmpty(t, matchDetail)
	})

	t.Run("fetch match detail dota2 should failed because match id is not found", func(t *testing.T) {
		matchID := 0

		matchDetail, err := openDotaAPI.FetchMatchDetail(ctx, int64(matchID))

		assert.Error(t, err)
		assert.Empty(t, matchDetail)
	})
}
