package spotifyhistory

import (
	"time"
	"sort"
)

func GetDailyFavourite(listenHistory *map[time.Time]map[Key]int) *map[time.Time]ListenEntry {
	favourites := make(map[time.Time]ListenEntry)
	for dateListened, listenInstance := range *listenHistory {
		favourite := 0
		var favListen Key
		for item, listen := range listenInstance {
			if listen > favourite {
				favListen = item
				favourite = listen
			}
		}
		favourites[dateListened] = ListenEntry{favListen.ArtistName, favListen.TrackName, favourite, favListen.EndTime}
	}

	return &favourites
}

func GetSortedDates (m *map[time.Time]ListenEntry) *[]time.Time{
	var keys []time.Time
	for date := range *m {
		keys = append(keys, date)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	return &keys
}