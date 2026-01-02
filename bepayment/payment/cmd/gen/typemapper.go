package main

import "strings"

var pgxtype = map[string]string{
	"pgtype.UUID":        "string",
	"pgtype.Numeric":     "float64",
	"pgtype.Timestamptz": "time.Timestamptz",
	"pgtype.Bool":        "bool",
}

func pgxTypeMapper(src string) string {
	if v, ok := pgxtype[src]; ok {
		return v
	}

	// []pgtype.xxx
	if strings.HasPrefix(src, "[]") {
		return "[]" + pgxTypeMapper(src[2:])
	}

	return src
}
