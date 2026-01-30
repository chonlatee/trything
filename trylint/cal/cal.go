package cal

import (
	"log"
	"strconv"
)

func Add(a, b int) int {
	return a + b
}

func ParseInt(v string) float64 {
	_, err := strconv.Atoi(v)
	if err != nil {
		log.Fatal(err)
	}

	r1, err := strconv.ParseFloat(v, 64)
	return r1
}
