package spotifyhistory

import (
	"time"
)

const DATEONLY = "2006-01-02" // Format for the date needed in json

const OUTPUTFILE = "../output/Ayat Spotify.txt"

const DATAFILES = "../resources"

type ListenInstance struct {
	TS                                string `json:"ts"`
	Ms_Played                         int    `json:"ms_played"`
	Master_Metadata_Track_Name        string `json:"master_metadata_track_name"`
	Master_Metadata_Album_Artist_Name string `json:"master_metadata_album_artist_name"`
	Master_Metadata_Album_Album_Name  string `json:"master_metadata_album_album_name"`
	Spotify_Track_Uri                 string `json:"spotify_track_uri"`
}

type ListenEntry struct {
	ArtistName string
	AlbumName  string
	TrackName  string
	MsPlayed   int
	TimeStamp  time.Time
	URI        string
}

type Entry struct {
	ArtistName string
	AlbumName  string
	TrackName  string
	TimeStamp  time.Time
	URI        string
}

type SongEntry struct {
	TrackName string
	MsPlayed  int
}
