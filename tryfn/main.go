package main

import (
	"chonlatee/fn/fn"
	"fmt"
)

func main() {

	s := fn.SL[int]{
		1, 2, 3, 4, 5, 6, 7,
	}

	double := s.TransForm(func(v int) int {
		return v * 2
	})

	for v := range double {
		fmt.Println(v)
	}

	fmt.Println()

	filter := s.Filter(func(v int) bool {
		return v%2 == 0
	})

	for v := range filter {
		fmt.Println(v)
	}

}
