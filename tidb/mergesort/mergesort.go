package main

import (
	"sort"
	"sync"
)

type I64MultiLevelPool struct {
	Pools       []*sync.Pool
	InitialSize int
}

func NewI64MultiLevelPool(initialSize int, steps int) *I64MultiLevelPool {
	mPool := &I64MultiLevelPool{
		InitialSize: initialSize,
	}
	for ; steps >= 0; initialSize, steps = initialSize*2, steps-1 {
		curSize := initialSize
		mPool.Pools = append(mPool.Pools, &sync.Pool{
			New: func() interface{} {
				return make([]int64, curSize, curSize)
			},
		})
	}
	return mPool
}

func (this *I64MultiLevelPool) Get(size int) []int64 {
	if size <= this.InitialSize {
		return this.Pools[0].Get().([]int64)[:size]
	}
	i := 0
	for curSize := this.InitialSize; curSize < size && i < len(this.Pools); curSize, i = curSize<<1, i+1 {
	}
	if i < len(this.Pools) {
		return this.Pools[i].Get().([]int64)[:size]
	}
	return make([]int64, size, size)
}

func (this *I64MultiLevelPool) Put(slice []int64) {
	sl := cap(slice)
	i := 0
	for curSize := this.InitialSize; curSize < sl && i < len(this.Pools); curSize, i = curSize<<1, i+1 {
	}
	if i < len(this.Pools) {
		this.Pools[i].Put(slice)
	}
}

// MergeSort performs the merge sort algorithm.
// Please supplement this function to accomplish the home work.
func MergeSort(src []int64) {
	sl := len(src)
	for ; (sl >> 1) > 128; sl = sl >> 1 {}
	i64MultiLevelPool := NewI64MultiLevelPool(sl, 12)
	mergeSortImpl(src, i64MultiLevelPool)
}

func mergeSortImpl(src []int64, pools *I64MultiLevelPool) {
	if len(src) < 256 {
		sort.Slice(src, func(i, j int) bool { return src[i] < src[j] })
		return
	}
	pivot := len(src) >> 1
	auxLen := len(src) - pivot
	childs := [][]int64{src[:pivot], src[pivot:]}
	wg := &sync.WaitGroup{}
	for _, child := range childs {
		wg.Add(1)
		go childMergeSort(child, wg, pools)
	}
	wg.Wait()
	i, j, k := pivot-1, auxLen-1, len(src)-1
	//aux := pools.Get(auxLen)
	aux := make([]int64, auxLen, auxLen)
	copy(aux, childs[1])
	for ; i >= 0 && j >= 0; {
		if src[i] > aux[j] {
			src[k] = src[i]
			k--
			i--
		} else {
			src[k] = aux[j]
			k--
			j--
		}
	}
	for ; j >= 0; j-- {
		src[k] = aux[j]
		k--
	}
	//pools.Put(aux)
}

func childMergeSort(childSrc []int64, g *sync.WaitGroup, pools *I64MultiLevelPool) {
	mergeSortImpl(childSrc, pools)
	g.Done()
}
