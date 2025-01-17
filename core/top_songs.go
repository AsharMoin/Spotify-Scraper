package spotifyhistory

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func GetSpotifyHistory() {
	var year int = 2018
	// walk the resources dir and keep appending to byteData
	err := filepath.WalkDir("../resources", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// if the file is not a dir and has the correct prefix, open it and read it, appending to byteData
		if !d.IsDir() && strings.HasPrefix(d.Name(), "Streaming_History_Audio_") {
			// open the file
			data, err := os.Open(path)
			if err != nil {
				fmt.Printf("Something went wrong while opening the file %v: %v\n", path, err)
				return err
			}
			// defer closing the file to ensure it's closed when we're done
			defer data.Close()
			// read the file contents
			var byteData []ListenInstance
			if err := json.NewDecoder(data).Decode(&byteData); err != nil {
				fmt.Printf("Error decoding JSON in file %v: %v\n", path, err)
				return err
			}
			outputFile := fmt.Sprintf("../output/Streaming_History_%d.txt", year)
			GetTopSongs(outputFile, byteData)
			year++
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func GetTopSongs(outputFile string, spotifyHistory []ListenInstance) {
	// get writer so we can write to output file
	w := MakeWriter(outputFile)
	defer w.Flush()

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
}

func GetDailyFavourite(listenHistory *map[time.Time]map[Entry]int) *map[time.Time]ListenEntry {
	favourites := make(map[time.Time]ListenEntry)
	for dateListened, listenInstance := range *listenHistory {
		var listenTime float64 = 0.0
		var favListen Entry
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
			favourites[dateListened] = ListenEntry{
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
