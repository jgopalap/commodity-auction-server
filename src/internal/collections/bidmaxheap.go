package maxheap

import "internal/models"

type BidMaxHeap []*models.BidDetailed

func (heap BidMaxHeap) Len() int { return len(heap) }

func (heap *BidMaxHeap) Push(x interface{}) {
	*heap = append(*heap, x.(*models.BidDetailed))
}

func (heap *BidMaxHeap) Pop() interface{} {
	old := *heap
	n := len(old)
	x := old[n-1]
	*heap = old[0 : n-1]

	return x
}

func (heap *BidMaxHeap) Peek() interface{} {
	old := *heap
   	return old[0]
}

func (heap BidMaxHeap) Less(i, j int) bool { return heap[i].Value < heap[j].Value }
func (heap BidMaxHeap) Swap(i, j int)      { heap[i], heap[j] = heap[j], heap[i] }

