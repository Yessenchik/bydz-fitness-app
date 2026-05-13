package usecase

import "time"

func nowUTC() time.Time {
	return time.Now().UTC()
}
