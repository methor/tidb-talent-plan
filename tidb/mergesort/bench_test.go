package main

import (
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"testing"
)

func BenchmarkMergeSort(b *testing.B) {
	f, err := os.Create("mergesort.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	numElements := 16 << 20
	src := make([]int64, numElements)
	original := make([]int64, numElements)
	prepare(original)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		copy(src, original)
		b.StartTimer()
		MergeSort(src)
	}
}

func BenchmarkNormalSort(b *testing.B) {
		f, err := os.Create("normalsort.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	numElements := 16 << 20
	src := make([]int64, numElements)
	original := make([]int64, numElements)
	prepare(original)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		copy(src, original)
		b.StartTimer()
		sort.Slice(src, func(i, j int) bool { return src[i] < src[j] })
	}
}
