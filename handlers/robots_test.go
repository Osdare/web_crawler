package handlers

import (
	"testing"
)

func TestRobotsToDomain(t *testing.T) {
	youtubeRobots, err := DownloadRobots("https://youtube.com")
	if err != nil {
		t.Error(err)
	}

	domain, err := RobotsToDomain("https://youtube.com", youtubeRobots)
	if err != nil {
		t.Error(err)
	}

	if domain.Disallowed[0] != "/api/" {
		t.Error("expected /api/ got: ", domain.Disallowed[0])
	}
}
