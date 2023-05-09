package opendota

type MatchDetail struct {
	ReplayURL string `json:"replay_url"`
	MatchID   int64  `json:"match_id"`
	Region    int32  `json:"region"`
}
