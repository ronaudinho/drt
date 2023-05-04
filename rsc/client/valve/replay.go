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

type ReplayDownloader interface {
	// Download
	// downloading the replay file is downloaded directly to the
	// Valve replay proxy using the replayURL.
	//
	// replayURL contains the full raw URI to download the replay to the Valve replay proxy,
	// for example http://replay153.valve.net/570/7132230434_1635105612.dem.bz2
	// Based on that value, we can download the replay DotA2 file directly using that URL.
	//
	// destination parameter specifies the file destination to save the downloaded replay file.
	Download(ctx context.Context, replayURL string, destination string) error
}

type replay struct {
	httpClient *http.Client
}

func NewReplay(httpClient *http.Client) ReplayDownloader {
	return &replay{
		httpClient: httpClient,
	}
}

func (r *replay) Download(ctx context.Context, replayURL string, destination string) error {
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

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	fileName := getFileName(replayURL)

	// create location for destination file
	// example:
	// /tmp/${filename}
	name := fmt.Sprintf("%s/%s", destination, fileName)

	file, err := os.Create(name)

	if err != nil {
		return err
	}

	_, err = file.Write(body)

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
