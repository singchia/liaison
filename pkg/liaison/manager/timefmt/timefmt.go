package timefmt

import "time"

var shanghaiLocation = func() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.FixedZone("CST", 8*60*60)
	}
	return loc
}()

func FormatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(shanghaiLocation).Format(time.DateTime)
}
