package spotifyhistory

import (
	"fmt"
	"time"
)

func GetTracksMap(listenHistory []ListenInstance) *map[time.Time]map[Entry]int {
	m := make(map[time.Time]map[Entry]int)

	for _, track := range listenHistory {
		date, err := time.Parse(DATEONLY, track.TS[:len(DATEONLY)])
		if err != nil {
			fmt.Println(err)
		}

		var key Entry = Entry{
			ArtistName: track.Master_Metadata_Album_Artist_Name,
			AlbumName:  track.Master_Metadata_Album_Album_Name,
			TrackName:  track.Master_Metadata_Track_Name,
			TimeStamp:  date,
			URI:        track.Spotify_Track_Uri,
		}

		if _, exists := m[date]; !exists {
			m[date] = make(map[Entry]int)
		}

		m[date][key] += track.Ms_Played
	}

	return &m
}
