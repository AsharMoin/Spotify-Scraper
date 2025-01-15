package spotifyhistory

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func GetSpotifyHistory(path string) ([]ListenInstance, error) {
	// declare the byte slice we will be using
	var spotifyHistory []ListenInstance

	// walk the resources dir and keep appending to byteData
	err := filepath.WalkDir("../resources", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// if the file is not a dir and has the correct prefix, open it and read it, appending to byteData
		if !d.IsDir() && strings.HasPrefix(d.Name(), "StreamingHistory_music") {
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

			spotifyHistory = append(spotifyHistory, byteData...)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return spotifyHistory, err
}
