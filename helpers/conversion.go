package helpers

import "time"

// StringToTime converts a timestamp string in the given layout to time.Time.
// An empty string yields the zero time and a nil error; a non-empty string
// that fails to parse returns the parse error.
func StringToTime(s, layout string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse(layout, s)
}
