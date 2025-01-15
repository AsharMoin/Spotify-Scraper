package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	spotify "github.com/AsharMoin/Spotify-Scraper/core"
)

func main() {
	w := spotify.MakeWriter(spotify.OUTPUTFILE)
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

			listenHistory := spotify.GetTracksMap(byteData)

			favourites := spotify.GetDailyFavourite(listenHistory)

			dates := spotify.GetSortedDates(favourites)

			for _, date := range *dates {
				favSong := (*favourites)[date]

				output := fmt.Sprintf("Date: %s | Most Popular: %s, %s | Minutes Listened: %v\n", date.Format(spotify.DateOnly), favSong.ArtistName, favSong.TrackName, (favSong.MsPlayed / 60000))

				spotify.WriteStuff(output, w)
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
