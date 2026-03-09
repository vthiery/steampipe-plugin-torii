package torii

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// FlexTime can unmarshal both ISO 8601 date strings and Unix epoch millisecond
// integers from JSON, as different Torii API endpoints use both formats.
type FlexTime struct {
	Value time.Time
}

// UnmarshalJSON implements json.Unmarshaler for FlexTime.
func (f *FlexTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	// Try as a JSON string (ISO 8601)
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		for _, layout := range []string{
			time.RFC3339Nano,
			time.RFC3339,
			"2006-01-02T15:04:05.000Z",
		} {
			if t, err := time.Parse(layout, s); err == nil {
				f.Value = t.UTC()
				return nil
			}
		}
		return nil // Unparseable string — leave zero
	}
	// Try as a number (Unix epoch milliseconds)
	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		ms, err := strconv.ParseInt(string(n), 10, 64)
		if err == nil && ms > 0 {
			f.Value = time.Unix(ms/1000, (ms%1000)*int64(time.Millisecond)).UTC()
		}
		return nil
	}
	return nil
}

// flexTimeTransform is a steampipe transform that converts a FlexTime field
// to a *time.Time suitable for a TIMESTAMP column, returning nil for zero values.
func flexTimeTransform(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}
	switch v := d.Value.(type) {
	case FlexTime:
		if v.Value.IsZero() {
			return nil, nil
		}
		return v.Value, nil
	case *FlexTime:
		if v == nil || v.Value.IsZero() {
			return nil, nil
		}
		return v.Value, nil
	}
	return nil, nil
}
