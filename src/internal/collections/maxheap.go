package maxheap

import (
	"sort"
)

type Interface interface {
	sort.Interface
	Push(x interface{}) // add x as element Len()
	Pop() interface{}   // remove and return element Len() - 1.
	Peek() interface{}
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func Push(heap Interface, x interface{}) {
	heap.Push(x)

	i := heap.Len() - 1
	// Fix the max heap property if it is violated
	for i != 0 && heap.Less(parent(i), i){
		heap.Swap(i, parent(i))
		i = parent(i)
	}
}

func Pop(heap Interface) interface{} {
	heap.Swap(0, heap.Len()-1)
	i := 0
	end := heap.Len() - 1
	for i != end {
		chosenI := leftChild(i)
		if leftChild(i) < end && rightChild(i) < end {
			if heap.Less(leftChild(i), rightChild(i)) {
				chosenI = rightChild(i)
			}
		} else if rightChild(i) < end {
			chosenI = rightChild(i)
		} else if leftChild(i) >= end {
			break
		}

		if heap.Less(i, chosenI) {
			heap.Swap(i, chosenI)
		}
		i = chosenI
	}
	return heap.Pop()
}

func Peek(heap Interface) interface{} {
	return heap.Peek()
}

func leftChild(i int) int {
	return 2 * i  + 1
}

func rightChild(i int) int {
	return 2 * i  + 2
}

func parent(i int) int {
	return int(float64((i-1) /2))
}

