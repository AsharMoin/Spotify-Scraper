package main

import (
	spotify "github.com/AsharMoin/Spotifystats/core"
)

func main() {
	spotify.GetListenHistory()
	spotify.GetAlbumHistory()
}
