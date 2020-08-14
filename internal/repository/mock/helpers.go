package mock

import "time"

func TimeToTimePtr(t time.Time) *time.Time {
	return &t
}

func TimePtrToTime(t *time.Time) (emptyTime time.Time) {
	if t != nil {
		return *t
	}
	return emptyTime
}
