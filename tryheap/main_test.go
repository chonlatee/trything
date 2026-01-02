package main

import "testing"

func BenchmarkFilterBefore(b *testing.B) {
	l := targetList{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for b.Loop() {
		l.filterBefore()
	}
}

func BenchmarkFilterAfter(b *testing.B) {
	l := targetList{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	for b.Loop() {
		l.filterAfter()
	}
}
