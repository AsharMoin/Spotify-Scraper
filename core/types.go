package spotifyhistory

import (
	"time"
)

const DATEONLY = "2006-01-02"

const MONTHONLY = "2006-01"

const OUTPUTFILE = "../output/Ayat Spotify.txt"

const DATAFILES = "../resources"

type Music interface {
	Song | Album
}

type ListenInstance struct {
	TS                                string `json:"ts"`
	Ms_Played                         int    `json:"ms_played"`
	Master_Metadata_Track_Name        string `json:"master_metadata_track_name,omitempty"`
	Master_Metadata_Album_Artist_Name string `json:"master_metadata_album_artist_name,omitempty"`
	Master_Metadata_Album_Album_Name  string `json:"master_metadata_album_album_name,omitempty"`
	Spotify_Track_Uri                 string `json:"spotify_track_uri,omitempty"`
}
type Song struct {
	ArtistName string
	AlbumName  string
	TrackName  string
	MsPlayed   float64
	TimeStamp  string
	URI        string
}

type SongEntry struct {
	ArtistName string
	AlbumName  string
	TrackName  string
	TimeStamp  time.Time
	URI        string
}

type Album struct {
	ArtistName  string
	AlbumName   string
	TimesPlayed int
}

type AlbumEntry struct {
	ArtistName string
	AlbumName  string
}
