package spotifyhistory

import (
	"fmt"
	"strings"
	"time"
)

func GetAlbumHistory() {
	err := ProcessMusicHistory(
		"../resources",             // Base directory
		"Streaming_History_Audio_", // File prefix
		"../output",                // Output directory
		2018,                       // Starting year
		func(outputFile string, data []ListenInstance) error {
			getTopAlbums(outputFile, data)
			return nil
		},
	)
	if err != nil {
		fmt.Println(err)
	}
}

func getTopAlbums(outputFile string, spotifyHistory []ListenInstance) {
	// make writer
	w := MakeWriter(outputFile)
	defer w.Flush()

	boxWidth := 50
	var output strings.Builder

	output.WriteString("\n" + strings.Repeat("-", boxWidth/2-9))
	output.WriteString(" Favourite Albums ")
	output.WriteString(strings.Repeat("-", boxWidth/2-9) + "\n")

	WriteStuff(output.String(), w)

	// get the data into a map where months are the key
	history := GetAlbumMap(spotifyHistory)

	favouriteAlbums := getMonthlyFavouriteAlbum(history)
	months := GetSortedKeys(favouriteAlbums)

	for _, month := range *months {
		favAlbum := (*favouriteAlbums)[month]
		output := FormatMonthOutput(favAlbum, month.Format(MONTHONLY))
		WriteStuff(output, w)
	}
}

func getMonthlyFavouriteAlbum(listenHistory *map[time.Time]map[AlbumEntry]int) *map[time.Time]Album {
	favouriteAlbums := make(map[time.Time]Album)

	// iterate over every album entry in each month
	for month, albumMap := range *listenHistory {
		var favAlbum AlbumEntry
		maxListens := 0
		foundValid := false

		for album, timesListened := range albumMap {
			if timesListened > maxListens {
				favAlbum = album
				maxListens = timesListened
				foundValid = true
			}
		}
		if foundValid {
			favouriteAlbums[month] = Album{favAlbum.ArtistName, favAlbum.AlbumName, maxListens}
		}
	}

	return &favouriteAlbums
}
