package link

import "time"

type Link struct {
	from      string
	to        string
	createdAt time.Time
	updatedAt time.Time
}
