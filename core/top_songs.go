package spotifyhistory

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func GetTopSongs() {
	w := MakeWriter(OUTPUTFILE)
	err := filepath.WalkDir("../resources", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasPrefix(d.Name(), "StreamingHistory_music") {
			// Open the file
			data, err := os.Open(path)
			if err != nil {
				fmt.Printf("Something went wrong while opening the file %v: %v\n", path, err)
				return err
			}

			// Defer closing the file to ensure it's closed when we're done
			defer data.Close()

			// Read the file contents
			byteData, err := io.ReadAll(data)
			if err != nil {
				fmt.Printf("Error reading file %v: %v\n", path, err)
				return err
			}

			listenHistory := GetTracksMap(byteData)

			favourites := GetDailyFavourite(listenHistory)

			dates := GetSortedDates(favourites)

			for _, date := range *dates {
				favSong := (*favourites)[date]

				output := fmt.Sprintf("Date: %s | Most Popular: %s, %s | Minutes Listened: %v\n", date.Format(DateOnly), favSong.ArtistName, favSong.TrackName, (favSong.MsPlayed / 60000))

				WriteStuff(output, w)
			}

			if err := w.Flush(); err != nil {
				fmt.Println("Error flushing buffer:", err)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Impossible to walk directories: %s", err)
	}
}

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
