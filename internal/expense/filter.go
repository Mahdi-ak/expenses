package expense

import "time"

type Filter struct {
	Category *string
	Date     *time.Time
}
