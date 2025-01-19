package spotifyhistory

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var year string
var month Month

type Month int

const (
	January Month = iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

func (m Month) String() string {
	months := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	if m >= January && m <= December {
		return months[m]
	}
	return "Unknown"
}

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

func FormatSongOutput(entry Song) string {
	var output strings.Builder
	boxWidth := 100
	labelWidth := 10 // Width for labels (Title, Artist, etc.)

	// if its a new year, append to string
	entryYear := entry.TimeStamp[:4]
	if entryYear != year {
		yearDivider := fmt.Sprintf("\n---%s---\n\n", entryYear)
		output.WriteString(yearDivider)
		year = entryYear
	}

	// if its a new month, append to string
	entryMonthStr := entry.TimeStamp[5:7]
	entryMonthInt, err := strconv.Atoi(entryMonthStr)
	if err != nil {
		fmt.Println("Error converting month:", err)
		return "Error Converting"
	}
	entryMonth := Month(entryMonthInt - 1)
	if entryMonth != month {
		monthDivide := fmt.Sprintf("\n---%s---\n\n", entryMonth)
		output.WriteString(monthDivide)
		month = entryMonth
	}

	// make the box design the data gets printed out in
	output.WriteString(strings.Repeat("-", boxWidth) + "\n")

	v := reflect.ValueOf(entry)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := fmt.Sprint(v.Field(i))

		line := fmt.Sprintf("| %-*s : %-*s |\n",
			labelWidth,
			fieldName,
			boxWidth-labelWidth-7,
			fieldValue,
		)
		output.WriteString(line)
	}

	output.WriteString(strings.Repeat("-", boxWidth) + "\n")

	return output.String()
}

func FormatMonthOutput(entry Album, curMonth string) string {
	var output strings.Builder
	boxWidth := 50
	labelWidth := 12 // Width for labels (Title, Artist, etc.)

	entryMonthStr := curMonth[5:7]
	entryMonthInt, err := strconv.Atoi(entryMonthStr)
	if err != nil {
		fmt.Println("Error converting month:", err)
		return "Error Converting"
	}
	entryMonth := Month(entryMonthInt - 1)
	if entryMonth != month {
		monthDivide := fmt.Sprintf("\n---%s---\n\n", entryMonth)
		output.WriteString(monthDivide)
		month = entryMonth
	}

	output.WriteString(strings.Repeat("-", boxWidth) + "\n")

	v := reflect.ValueOf(entry)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := fmt.Sprint(v.Field(i))

		line := fmt.Sprintf("| %-*s : %-*s |\n",
			labelWidth,
			fieldName,
			boxWidth-labelWidth-7,
			fieldValue,
		)
		output.WriteString(line)
	}

	output.WriteString(strings.Repeat("-", boxWidth) + "\n")

	return output.String()
}
