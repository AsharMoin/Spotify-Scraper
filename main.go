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

type ListenInstance struct {
	EndTime    string `json:"endTime"`
	ArtistName string `json:"artistName"`
	TrackName  string `json:"trackName"`
	MsPlayed   int    `json:"msPlayed"`
}

type ListenEntry struct {
	ArtistName string
	TrackName  string
	MsPlayed   int
	EndTime    time.Time
}

type Key struct {
	ArtistName string
	TrackName  string
	EndTime    time.Time
}

func main() {
	w := makeWriter(OUTPUTFILE)
	err := filepath.WalkDir("./resources", func(path string, d fs.DirEntry, err error) error {
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

			listenHistory := getTracksMap(byteData)

			favourites, dates := doStuff(listenHistory)

			for _, date := range *dates {
				favSong := (*favourites)[date]
		
				output := fmt.Sprintf("Date: %s | Most Popular: %s, %s | Minutes Listened: %v\n", date.Format(DateOnly), favSong.ArtistName, favSong.TrackName, (favSong.MsPlayed / 60000))
		
				writeStuff(output, w)
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

func doOtherStuff(byteData []byte) {

}

func getTracksMap(byteData []byte) *map[time.Time]map[Key]int {
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

func getDailyFavourite(listenHistory *map[time.Time]map[Key]int) *map[time.Time]ListenEntry {
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

func getSortedDates (m *map[time.Time]ListenEntry) *[]time.Time{
	var keys []time.Time
	for date := range *m {
		keys = append(keys, date)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	return &keys
}

func doStuff(listenHistory *map[time.Time]map[Key]int) (*map[time.Time]ListenEntry, *[]time.Time){
	favourites := getDailyFavourite(listenHistory)
	
	dates := getSortedDates(favourites)

	return favourites, dates
}

func makeWriter(fName string) *bufio.Writer {
	file, err := os.OpenFile(fName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	return bufio.NewWriter(file)
}

func writeStuff(output string, w *bufio.Writer) {
	// Write the output to the file
	_, err := w.WriteString(output)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
