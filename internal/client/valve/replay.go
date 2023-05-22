package valve

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	ErrDownloadRelpay = fmt.Errorf("failed to download the dota2 replay")
)

const (
	DefaultDestination = "/tmp"
)

type Replay struct {
	httpClient *http.Client
}

func NewReplay(httpClient *http.Client) *Replay {
	return &Replay{
		httpClient: httpClient,
	}
}

// Download
// downloading the replay file is downloaded directly to the
// Valve replay proxy using the replayURL.
//
// replayURL contains the full raw URI to download the replay to the Valve replay proxy,
// for example http://replay153.valve.net/570/7132230434_1635105612.dem.bz2
// Based on that value, we can download the replay DotA2 file directly using that URL.
//
// destination parameter specifies the file destination to save the downloaded replay file.
func (r *Replay) Download(ctx context.Context, replayURL string, destination string) error {
	fileName := getFileName(replayURL)

	// create location for destination file
	// example:
	// /tmp/${filename}
	name := fmt.Sprintf("%s/%s", destination, fileName)

	isExist := isReplayAlreadyExist(name)

	if isExist {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, replayURL, nil)
	if err != nil {
		return err
	}

	res, err := r.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ErrDownloadRelpay
	}

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func getFileName(replayURL string) string {
	s := strings.Split(replayURL, "/")

	fileName := s[len(s)-1]

	return fileName
}

// isReplayAlreadyExist
//
// given path and file name, check if file is already exist in the path or not
// if file replay is already exist, then just return immediately.
// otherwise, need to download it first
func isReplayAlreadyExist(file string) bool {
	_, err := os.Stat(file)

	return err == nil
}
