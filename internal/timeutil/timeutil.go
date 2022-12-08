package timeutil

import "time"

func CurrentDate(layout string) string {
	return time.Now().Format(layout)
}
