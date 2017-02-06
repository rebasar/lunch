package lunchy

import (
	"time"
)

// Stolen shamelessly from: http://stackoverflow.com/questions/25087960/json-unmarshal-time-that-isnt-in-rfc-3339-format

type Date struct {
	time.Time
}

const dateLayout = "2006-01-02"

var nilTime = (time.Time{}).UnixNano()

func (d *Date) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	timestring := string(b)
	if timestring == "null" {
		return nil
	}
	d.Time, err = time.Parse(dateLayout, timestring)
	return err
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return []byte(d.Time.Format(dateLayout)), nil
}

func (d *Date) IsSet() bool {
	return d.UnixNano() != nilTime
}

func (d Date) String() string {
	return d.Time.Format(dateLayout)
}
