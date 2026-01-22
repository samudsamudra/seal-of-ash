package utils

import (
	"time"
	"fmt"
)

func FormatIndoTime(t time.Time) string {
	months := []string{
		"Januari",
		"Februari",
		"Maret",
		"April",
		"Mei",
		"Juni",
		"Juli",
		"Agustus",
		"September",
		"Oktober",
		"November",
		"Desember",
	}

	day := t.Day()
	month := months[int(t.Month())-1]
	year := t.Year()
	hour := t.Hour()
	min := t.Minute()
	sec := t.Second()

	return fmt.Sprintf(
		"%d %s %d pada jam %02d.%02d.%02d",
		day, month, year, hour, min, sec,
	)
}
