package zettel

import "github.com/odas0r/zet/pkg/domain/shared/timestamp"

// TODO
type Link struct {
	From      Zettel
	To        Zettel
	Timestamp timestamp.Timestamp
}
