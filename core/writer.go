package spotifyhistory

import (
	"bufio"
	"fmt"
	"os"
)

func MakeWriter(fName string) *bufio.Writer {
	file, err := os.OpenFile(fName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	return bufio.NewWriter(file)
}

func WriteStuff(output string, w *bufio.Writer) {
	// Write the output to the file
	_, err := w.WriteString(output)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func FormatOutput(entry ListenEntry) {

}
