package link

import (
	"github.com/odas0r/zet/pkg/domain/shared"
)

type Link struct {
	from      string
	to        string
	createdAt shared.Time
	updatedAt shared.Time
}
