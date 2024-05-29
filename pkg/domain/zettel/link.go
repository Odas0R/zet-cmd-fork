package zettel

import "github.com/odas0r/zet/pkg/domain/shared/timestamp"

// Link represents a connection between two Zettels
type Link struct {
	From      Zettel
	To        Zettel
	Timestamp timestamp.Timestamp
}

// NewLink creates a new link from one Zettel to another
func NewLink(from, to Zettel) Link {
	return Link{
		From:      from,
		To:        to,
		Timestamp: timestamp.New(),
	}
}
