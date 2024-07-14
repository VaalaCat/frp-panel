package utils

import (
	"time"
)

func IsSameDay(first time.Time, second time.Time) bool {
	return first.YearDay() == second.YearDay() && first.Year() == second.Year()
}
