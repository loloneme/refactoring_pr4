package time

import "time"

func LastDays(n int) (string, string) {
	to := time.Now().UTC()
	from := to.AddDate(0, 0, -n)
	return from.Format("2006-01-02"), to.Format("2006-01-02")
}
