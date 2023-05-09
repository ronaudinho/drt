package opendota

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrFetchMatchDetail = fmt.Errorf("failed to fetch match detail")
)

type OpenDotaAPI struct {
	httpClient *http.Client
	BaseURL    string
}

func NewMatchAPI(httpClient *http.Client, BaseURL string) *OpenDotaAPI {
	return &OpenDotaAPI{
		httpClient: httpClient,
		BaseURL:    BaseURL,
	}
}

func (m *OpenDotaAPI) FetchMatchDetail(ctx context.Context, matchID int64) (MatchDetail, error) {
	matchDetailURL := fmt.Sprintf("%s/matches/%d", m.BaseURL, matchID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, matchDetailURL, nil)
	if err != nil {
		return MatchDetail{}, err
	}

	res, err := m.httpClient.Do(req)
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
