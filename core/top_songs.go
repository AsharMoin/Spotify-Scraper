package spotifyhistory

import (
	"fmt"
	"sort"
	"time"
)

func GetTopSongs() {
	// get writer so we can write to output file
	w := MakeWriter(OUTPUTFILE)

	// get spotify history
	spotifyHistory, err := GetSpotifyHistory(DATAFILES)
	if err != nil {
		fmt.Println(err)
	}

	// get structured data
	listenHistory := GetTracksMap(spotifyHistory)

	// get favourite songs
	// also get months so printing it is organized
	favourites := GetDailyFavourite(listenHistory)
	dates := GetSortedDates(favourites)

	// write to output file
	for _, date := range *dates {
		favSong := (*favourites)[date]
		output := FormatOutput(favSong)
		WriteStuff(output, w)
	}

	if err := w.Flush(); err != nil {
		fmt.Println("Error flushing buffer:", err)
	}

}

func GetDailyFavourite(listenHistory *map[time.Time]map[Entry]int) *map[time.Time]ListenEntry {
	favourites := make(map[time.Time]ListenEntry)
	for dateListened, listenInstance := range *listenHistory {
		listenTime := 0
		var favListen Entry
		foundValid := false

		for item, listen := range listenInstance {
			if listen > listenTime {
				favListen = item
				listenTime = listen
				foundValid = true
			}
		}

		// Only add to favorites if we found a valid entry
		if foundValid {
			favourites[dateListened] = ListenEntry{
				ArtistName: favListen.ArtistName,
				AlbumName:  favListen.AlbumName,
				TrackName:  favListen.TrackName,
				MsPlayed:   listenTime / 60000,
				TimeStamp:  favListen.TimeStamp.Format(DATEONLY),
				URI:        favListen.URI,
			}
		}
	}

	return &favourites
}

func GetSortedDates(m *map[time.Time]ListenEntry) *[]time.Time {
	var keys []time.Time
	for date := range *m {
		keys = append(keys, date)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	return &keys
}
