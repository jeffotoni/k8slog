package fmts

import (
	"bytes"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"time"

	gcat "github.com/jeffotoni/gconcat"
)

// Stdout func
func Stdout(strs ...interface{}) {
	str := gcat.Build(strs...)
	_, err := io.Copy(os.Stdout, strings.NewReader(str))
	if err != nil {
		log.Println(err)
	}
}

// Concat
func Concat(strs ...interface{}) string {
	return gcat.Concat(strs...)
}

// ConcatStr
func ConcatStr(strs ...string) string {
	return gcat.ConcatStr(strs...)
}

// ConcatStr
func ConcatStrInt(strs ...interface{}) string {
	return gcat.ConcatStringInt(strs...)
}

// substr (text, 2, 30)
func Substr(value string, leni, lenf int) string {

	if len(value) < leni {
		return ""
	}

	lenx := len(value) // max of text
	leny := lenf       // amount of character

	if lenx < lenf { // amount of character
		leny = lenx
	}

	// return..
	return value[leni:leny]
}

// replaces "\t" and "\n" to ""
func ReplaceSpaces(b *[]byte) {
	*b = bytes.ReplaceAll(*b, []byte("\t"), []byte(""))
	*b = bytes.ReplaceAll(*b, []byte("\n"), []byte(""))
}

/*
CalculateDaysDate - This function returns number of days of difference between two dates;

Case returns: -(some number) // it's before start date;

Case returns: (some number)  // it's after start date;
*/
func CalculateDaysDate(initialDate, fDate string, layout string) (int, error) {

	startDate, err := time.Parse(layout, initialDate)

	if err != nil {
		return 0, err
	}

	finalDate, err := time.Parse(layout, fDate)

	if err != nil {
		return 0, err
	}

	duration := startDate.Sub(finalDate)
	durationDays := math.Abs(duration.Hours() / 24)

	if finalDate.After(startDate) {
		return int(durationDays), nil
	}

	return -int(durationDays), nil
}
