package utils

import (
	"errors"
	"strconv"
	"time"
)

func SetIsNotEmpty(dest **string, val string) {

	if val != "" {
		*dest = &val
	}

}

func SetIsNotEmptyAbsen(dest **int, val string) error {

	if val != "" {
		val_absen, err := strconv.Atoi(val)
		if err != nil {
			return errors.New("Failed to convert into a integer!")
		}
		*dest = &val_absen
	}

	return nil

}

func SetResponseTime(format string) string {

	format_created_at_and_updated_at := time.Now().UTC().Format(format)
	return format_created_at_and_updated_at

}
