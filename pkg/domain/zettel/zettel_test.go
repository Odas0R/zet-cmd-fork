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
		kind        zettel.Kind
		expectedErr error
	}

	testCases := []testCase{
		{
			test:        "should return an error when the title is empty",
			title:       "",
			content:     "content",
			kind:        zettel.Permanent,
			expectedErr: zettel.ErrMissingValues,
		},
		{
			test:        "should return an error when the content is empty",
			title:       "title",
			content:     "",
			kind:        zettel.Permanent,
			expectedErr: zettel.ErrMissingValues,
		},
		{
			test:        "should return an error when the kind is wrong",
			title:       "title",
			content:     "content",
			kind:        "random_kind",
			expectedErr: zettel.ErrInvalidZettelKind,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			_, err := zettel.New(tc.title, tc.content, tc.kind)
			if err != tc.expectedErr {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
