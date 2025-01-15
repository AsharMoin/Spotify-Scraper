package spotifyhistory

import (
	"fmt"
	"time"
)

func GetTracksMap(listenHistory []ListenInstance) *map[time.Time]map[Entry]int {
	m := make(map[time.Time]map[Entry]int)

	for _, track := range listenHistory {
		date, err := time.Parse(DATEONLY, track.EndTime[:len(DATEONLY)])
		if err != nil {
			fmt.Println(err)
		}

		var key Entry = Entry{
			ArtistName: track.ArtistName,
			TrackName:  track.TrackName,
			EndTime:    date,
		}

		if _, exists := m[date]; !exists {
			m[date] = make(map[Entry]int)
		}

		m[date][key] += track.MsPlayed
	}

	return &m
}
