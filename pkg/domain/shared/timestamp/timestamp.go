package timestamp

import "time"

// Timestamp stores the creation and last update times of an entity.
type Timestamp struct {
	Created time.Time
	Updated time.Time
}

func New() Timestamp {
	now := time.Now().UTC()
	return Timestamp{
		Created: now,
		Updated: now,
	}
}

func (t Timestamp) Update() Timestamp {
	return Timestamp{
		Created: t.Created,
		Updated: time.Now().UTC(),
	}
}
