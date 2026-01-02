package main

import "fmt"

// type data struct {
// 	value int
// }

// func newData() *data {
// 	return &data{value: 42}
// }

// escape. alloc 1
// func main() {
// 	data := newData()
// 	_ = data
// }

// do not escape. not alloc
// func main() {
// 	data := &data{value: 42}
// 	_ = data
// }

// not escape on build. but escape on run time.
// func main() {
// 	num := 5
// 	s := make([]int, 0)

// 	for i := range num {
// 		s = append(s, i)
// 	}
// 	_ = s
// }
//

func main() {
	list := targetList{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// list.filterBefore()
	n := list.filterAfter()
	fmt.Printf("%+v\n", n)
}

type targetList []int

func (list targetList) filterBefore() targetList {
	filtered := make(targetList, 0, len(list))

	for _, e := range list {
		if e%2 == 0 {
			filtered = append(filtered, e)
		}
	}

	return filtered
}

func (list targetList) filterAfter() targetList {
	n := 0
	for _, e := range list {
		if e%2 == 0 {
			list[n] = e
			n++
		}
	}
	list = list[:n]
	return list
}
