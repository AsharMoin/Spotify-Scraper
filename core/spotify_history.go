package spotifyhistory

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type ProcessMusicHistoryFunc func(outputFile string, data []ListenInstance) error

func ProcessMusicHistory(baseDir string, filePrefix string, outputDir string, startYear int, processFunc ProcessMusicHistoryFunc) error {
	year := startYear

	// walk directory containing streaming history files
	err := filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasPrefix(d.Name(), filePrefix) {
			data, err := os.Open(path)
			if err != nil {
				fmt.Printf("Something went wrong while opening the file %v: %v\n", path, err)
				return err
			}
			defer data.Close()

			// decode the JSON data
			var byteData []ListenInstance
			if err := json.NewDecoder(data).Decode(&byteData); err != nil {
				fmt.Printf("Error decoding JSON in file %v: %v\n", path, err)
				return err
			}

			// define the output file path
			outputFile := fmt.Sprintf("%s/Streaming_History_%d.txt", outputDir, year)

			// call the processing function
			if err := processFunc(outputFile, byteData); err != nil {
				fmt.Printf("Error processing data for file %v: %v\n", path, err)
				return err
			}

			// increment the year
			year++
		}
		return nil
	})

	return err
}

func GetSongMap(listenHistory []ListenInstance) *map[time.Time]map[SongEntry]int {
	m := make(map[time.Time]map[SongEntry]int)

	// iterate over the array of structs
	for _, track := range listenHistory {
		// if json was null (happens if audio book and not song json entry)
		if track.Master_Metadata_Track_Name == "" {
			continue
		}

		date, err := time.Parse(DATEONLY, track.TS[:len(DATEONLY)])
		if err != nil {
			fmt.Println(err)
		}

		// create song entry
		var key SongEntry = SongEntry{
			ArtistName: track.Master_Metadata_Album_Artist_Name,
			AlbumName:  track.Master_Metadata_Album_Album_Name,
			TrackName:  track.Master_Metadata_Track_Name,
			TimeStamp:  date,
			URI:        track.Spotify_Track_Uri,
		}

		// if does not exist make a new entry
		if _, exists := m[date]; !exists {
			m[date] = make(map[SongEntry]int)
		}

		// otherwise add the time played so we have a cumulative time this song was played for the day
		m[date][key] += track.Ms_Played
	}

	return &m
}

func GetAlbumMap(listenHistory []ListenInstance) *map[time.Time]map[AlbumEntry]int {
	m := make(map[time.Time]map[AlbumEntry]int)

	for _, track := range listenHistory {
		if track.Master_Metadata_Album_Album_Name == "" {
			continue
		}
		month, err := time.Parse(MONTHONLY, track.TS[:len(MONTHONLY)])
		if err != nil {
			fmt.Println(err)
		}

		var key AlbumEntry = AlbumEntry{
			ArtistName: track.Master_Metadata_Album_Artist_Name,
			AlbumName:  track.Master_Metadata_Album_Album_Name,
		}

		if _, exists := m[month]; !exists {
			m[month] = make(map[AlbumEntry]int)
		}

		m[month][key]++
	}

	return &m
}

func GetSortedKeys[M Music](m *map[time.Time]M) *[]time.Time {
	// return a pointer to the time.Time keys sorted in desc order
	var keys []time.Time
	for date := range *m {
		keys = append(keys, date)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	return &keys
}
