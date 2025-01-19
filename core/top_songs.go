package spotifyhistory

import (
	"fmt"
	"math"
	"time"
)

func GetListenHistory() {
	err := ProcessMusicHistory(
		"../resources",             // Base directory
		"Streaming_History_Audio_", // File prefix
		"../output",                // Output directory
		2018,                       // Starting year
		func(outputFile string, data []ListenInstance) error {
			getTopSongs(outputFile, data)
			return nil
		},
	)
	if err != nil {
		fmt.Println(err)
	}
}

func getTopSongs(outputFile string, spotifyHistory []ListenInstance) {
	// get writer so we can write to output file
	w := MakeWriter(outputFile)
	defer w.Flush()

	// get structured data
	listenHistory := GetSongMap(spotifyHistory)

	// get favourite songs
	// also get months so printing it is organized
	favouriteSongs := getDailyFavourite(listenHistory)
	dates := GetSortedKeys(favouriteSongs)

	// write to output file
	for _, date := range *dates {
		favSong := (*favouriteSongs)[date]
		output := FormatSongOutput(favSong)
		WriteStuff(output, w)
	}
}

func getDailyFavourite(listenHistory *map[time.Time]map[SongEntry]int) *map[time.Time]Song {
	favourites := make(map[time.Time]Song)
	for dateListened, listenInstance := range *listenHistory {
		var favListen SongEntry
		listenTime := 0.0
		foundValid := false

		for item, listen := range listenInstance {
			if float64(listen) > listenTime {
				favListen = item
				listenTime = float64(listen)
				foundValid = true
			}
		}

		// Only add to favorites if we found a valid entry
		if foundValid {
			favourites[dateListened] = Song{
				ArtistName: favListen.ArtistName,
				AlbumName:  favListen.AlbumName,
				TrackName:  favListen.TrackName,
				MsPlayed:   math.Round((listenTime * 100 / 100) / 60000),
				TimeStamp:  favListen.TimeStamp.Format(DATEONLY),
				URI:        favListen.URI,
			}
		}
	}

	return &favourites
}
