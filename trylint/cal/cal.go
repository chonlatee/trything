package cal

import (
	"log"
	"strconv"
)

func Add(a, b int) int {
	return a + b
}

func ParseInt(v string) int64 {
	r, err := convertInt(v)
	if err != nil {
		log.Fatal(err)
	}

	return r
}

func convertInt(v string) (int64, error) {
	return strconv.ParseInt(v, 10, 64)
}
