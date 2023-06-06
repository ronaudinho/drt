package opendota

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	DefaultBaseURL = "https://api.opendota.com/api"
)

var (
	ErrFetchMatchDetail = fmt.Errorf("failed to fetch match detail")
)

type API struct {
	httpClient *http.Client
	BaseURL    string
}

func NewAPI(httpClient *http.Client, BaseURL string) *API {
	return &API{
		httpClient: httpClient,
		BaseURL:    BaseURL,
	}
}

// NewDefaultAPI creates a OpenDotaAPI with default params.
func NewDefaultAPI() *API {
	return NewAPI(http.DefaultClient, DefaultBaseURL)
}

func (api *API) FetchMatchDetail(ctx context.Context, matchID int64) (MatchDetail, error) {
	matchDetailURL := fmt.Sprintf("%s/matches/%d", api.BaseURL, matchID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, matchDetailURL, nil)
	if err != nil {
		return MatchDetail{}, err
	}

	res, err := api.httpClient.Do(req)
	if err != nil {
		return MatchDetail{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return MatchDetail{}, ErrFetchMatchDetail
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return MatchDetail{}, err
	}

	var matchDetail MatchDetail

	err = json.Unmarshal(body, &matchDetail)
	if err != nil {
		return MatchDetail{}, err
	}

	return matchDetail, nil
}
