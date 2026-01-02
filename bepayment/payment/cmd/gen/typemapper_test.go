package main

import (
	"testing"
)

func Test_pgxTypeMapper(t *testing.T) {

	testCases := []struct {
		name string
		src  string
		dst  string
	}{
		{
			name: "string to string",
			src:  "string",
			dst:  "string",
		},
		{
			name: "int to int",
			src:  "int",
			dst:  "int",
		},
		{
			name: "pgtype.UUID to string",
			src:  "pgtype.UUID",
			dst:  "string",
		},
		{
			name: "[]pgtype.UUID to []string",
			src:  "[]pgtype.UUID",
			dst:  "[]string",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := pgxTypeMapper(tc.src)
			if got != tc.dst {
				t.Errorf("got %s want %s", got, tc.dst)
			}
		})
	}

}
