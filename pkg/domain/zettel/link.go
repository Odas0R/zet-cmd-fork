package zettel

import (
	"github.com/google/uuid"
	"github.com/odas0r/zet/pkg/domain/shared/timestamp"
)

// Link represents a connection between two Zettels
type Link struct {
	From      uuid.UUID 
	To        uuid.UUID
	Timestamp timestamp.Timestamp
}

// NewLink creates a new link from one Zettel to another
func NewLink(from, to uuid.UUID) Link {
	return Link{
		From:      from,
		To:        to,
		Timestamp: timestamp.New(),
	}
}
