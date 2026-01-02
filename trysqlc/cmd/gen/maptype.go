package main

var mapType = map[string]string{
	"pgtype.UUID":        "string",
	"pgtype.Text":        "string",
	"pgtype.Timestamptz": "time.Time",
	"pgtype.Numeric":     "float64",
}
