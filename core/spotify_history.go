package spotifyhistory

import (
	"encoding/json"
	"time"
	"fmt"
)


func GetTracksMap(byteData []byte) *map[time.Time]map[Key]int {
	var listenHistory []ListenInstance

	json.Unmarshal(byteData, &listenHistory)

	m := make(map[time.Time]map[Key]int)

	for _, track := range listenHistory {
		date, err := time.Parse(DateOnly, track.EndTime[:len(DateOnly)])
		if err != nil {
			fmt.Println(err)
		}

		var key Key = Key{
			ArtistName: track.ArtistName,
			TrackName:  track.TrackName,
			EndTime:    date,
		}

		if _, exists := m[date]; !exists {
			m[date] = make(map[Key]int)
		}

		m[date][key] += track.MsPlayed
	}

	return &m
}