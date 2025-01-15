package spotifyhistory

import (
	"time"
)

const DateOnly = "2006-01-02" // Format for the date needed in json

const OUTPUTFILE = "../output/Ayat Spotify.txt"

type ListenInstance struct {
	EndTime    string `json:"endTime"`
	ArtistName string `json:"artistName"`
	TrackName  string `json:"trackName"`
	MsPlayed   int    `json:"msPlayed"`
}

type ListenEntry struct {
	ArtistName string
	TrackName  string
	MsPlayed   int
	EndTime    time.Time
}

type Key struct {
	ArtistName string
	TrackName  string
	EndTime    time.Time
}
