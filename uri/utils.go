package uri

import "time"

func newInt(i int) *int {
	return &i
}
func newInt32(i int32) *int32 {
	return &i
}
func newInt64(i int64) *int64 {
	return &i
}
func newFloat32(f float32) *float32 {
	return &f
}
func newFloat64(f float64) *float64 {
	return &f
}

func mTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
