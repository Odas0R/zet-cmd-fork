package zettel_test

import (
	"testing"

	"github.com/odas0r/zet/pkg/domain/zettel"
)

func TestZettel_NewZettel(t *testing.T) {
	type testCase struct {
		test        string
		title       string
		content     string
		kind        zettel.ZettelKind
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "should return an error when the title is empty",
			title:       "",
			content:     "content",
			kind:        Permanent,
			expectedErr: ErrMissingValues,
		},
	}
}
