package main

import (
	"fmt"
	"time"
)

func UnixtimeToHTMLDateString(unixtime int64) string {
	t := time.Unix(unixtime, 0)
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func HTMLDateStringToUnixtime(date string) (int64, error) {
	t, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return -1, err
	}
	return t.Unix(), nil
}
