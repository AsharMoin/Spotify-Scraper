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
		output := fmt.Sprintf("Date %s | Most Popular: %s, %s | Album: %s | Minutes Listened: %v minutes | Spotify URI: %s\n", date.Format(DATEONLY), favSong.ArtistName, favSong.TrackName, favSong.AlbumName, (favSong.MsPlayed / 60000), favSong.URI)
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
		for item, listen := range listenInstance {
			if listen > listenTime {
				favListen = item
				listenTime = listen
			}
		}
		favourites[dateListened] = ListenEntry{favListen.ArtistName, favListen.AlbumName, favListen.TrackName, listenTime, favListen.TimeStamp, favListen.URI}
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
