package entity

import "time"

type Time struct {
	time.Time
}

func (t *Time) MarshalJSON() ([]byte, error) {
	t.UnmarshalJSON()

	stamp := string(t.Unix())
	return []byte(stamp), nil
}
