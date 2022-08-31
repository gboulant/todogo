package todo

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"time"
)

const (
	layoutISO     = "2006-01-02"
	layoutUS      = "January 2, 2006"
	layoutUSshort = "2006-Jan-02"
	layoutUSlong  = "2006-Jan-02 15:04:05 Mon"
	layoutInt     = "20060102"
	layoutDefault = layoutUSshort
)

// timestamp returns the current date (now) in unix format
func timestamp() int64 {
	return time.Now().Unix()
}

// datelabel return a string representation of the given timestamp
func datelabel(timestamp int64) string {
	date := time.Unix(timestamp, 0)
	label := date.Format(layoutDefault)
	return label
}

// hashInt returns an integer hash representation (hash.crc32) of the givent string
func hashInt(s string) uint32 {
	b := []byte(s)
	h := crc32.ChecksumIEEE(b)
	return h
}

// dateInt returns an integer representation (YYYYMMDD) of the given timestamp
func dateInt(timestamp int64) uint64 {
	date := time.Unix(timestamp, 0)
	label := fmt.Sprintf("%.4d%.2d%.2d", date.Year(), date.Month(), date.Day())
	i, _ := strconv.ParseUint(label, 10, 64)
	return i
}

// hashdate returns an integer representation of the composition of the given
// text and timestamp. This can be used as a general invariant and unique index
// for a task, based on the string representation of the task and its timestamp
func hashdate(text string, timestamp int64) uint64 {
	d := dateInt(timestamp)
	h := hashInt(text)
	i := d*10000000000 + uint64(h)
	return i
}
