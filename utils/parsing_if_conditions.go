package utils

import (
	"errors"
	"strconv"
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
