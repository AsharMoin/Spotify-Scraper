package main

import (
	"bufio"
	"encoding/json"
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

const DateOnly = "2006-01-02" // Format for the date needed in json

const OUTPUTFILE = "./Ayat Spotify.txt"

var FavouriteOfTheYear Listen

type Track struct {
	EndTime    string `json:"endTime"`
	ArtistName string `json:"artistName"`
	TrackName  string `json:"trackName"`
	MsPlayed   int    `json:"msPlayed"`
}

type Listen struct {
	ArtistName, TrackName string
	MsPlayed              int
	EndTime               time.Time
}

type Key struct {
	ArtistName string
	TrackName  string
	EndTime    time.Time
}

func main() {
	err := filepath.WalkDir("./", func(path string, d fs.DirEntry, err error) error {
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

			doStuff(byteData)
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Impossible to walk directories: %s", err)
	}

	w := makeWriter(OUTPUTFILE)
	output := fmt.Sprintf("\nFavourite song of the year was: \nDate: %s | Most Popular: %s, %s | Minutes Listened: %v\n", FavouriteOfTheYear.EndTime.Format(DateOnly), FavouriteOfTheYear.ArtistName, FavouriteOfTheYear.TrackName, (FavouriteOfTheYear.MsPlayed / 60000))
	writeStuff(output, w)
	if err := w.Flush(); err != nil {
		fmt.Println("Error flushing buffer:", err)
		return
	}

}

func doOtherStuff (byteData []byte) {
	
}

func get

func doStuff(byteData []byte) {
	var tracks []Track

	json.Unmarshal(byteData, &tracks)

	m := make(map[string]map[Key]int)

	for _, track := range tracks {
		date := track.EndTime[:len(DateOnly)]

		datetime, err := time.Parse(DateOnly, date)
		if err != nil {
			fmt.Println(err)
		}

		var key Key = Key{
			ArtistName: track.ArtistName,
			TrackName:  track.TrackName,
			EndTime:    datetime,
		}

		if _, exists := m[date]; !exists {
			m[date] = make(map[Key]int)
		}

		m[date][key] += track.MsPlayed
	}

	favourites := make(map[string]Listen)
	for dateListened, items := range m {
		favourite := 0
		var favListen Key
		for item, listen := range items {
			if listen > favourite {
				favListen = item
				favourite = listen
			}
			if listen > FavouriteOfTheYear.MsPlayed {
				FavouriteOfTheYear = Listen{favListen.ArtistName, favListen.TrackName, favourite, favListen.EndTime}
			}
		}
		favourites[dateListened] = Listen{favListen.ArtistName, favListen.TrackName, favourite, favListen.EndTime}
	}

	var keys []time.Time
	for date := range favourites {
		date, err := time.Parse(DateOnly, date)
		if err != nil {
			fmt.Println(err)
		}
		keys = append(keys, date)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	w := makeWriter(OUTPUTFILE)

	for _, date := range keys {
		favSong := favourites[date.Format(DateOnly)]

		output := fmt.Sprintf("Date: %s | Most Popular: %s, %s | Minutes Listened: %v\n", date.Format(DateOnly), favSong.ArtistName, favSong.TrackName, (favSong.MsPlayed / 60000))

		writeStuff(output, w)
	}

	// Flush the buffered writer to ensure all data is written to the file
	if err := w.Flush(); err != nil {
		fmt.Println("Error flushing buffer:", err)
		return
	}
}

func makeWriter (fName string) *bufio.Writer{
	file, err := os.OpenFile(fName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	return bufio.NewWriter(file)	
}

func writeStuff (output string, w *bufio.Writer) {
	// Write the output to the file
	_, err := w.WriteString(output)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
