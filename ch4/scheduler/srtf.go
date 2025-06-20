package main

import (
	"container/heap"
)

type TaskHeap []Task

func (h TaskHeap) Len() int {
	return len(h)
}

func (h TaskHeap) Less(i, j int) bool {
	if h[i].Remaining != h[j].Remaining {
		return h[i].Remaining < h[j].Remaining
	}
	return h[i].ID < h[j].ID
}

func (h TaskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *TaskHeap) Push(x any) {
	*h = append(*h, x.(Task))
}

func (h *TaskHeap) Pop() any {
	res := (*h)[len(*h)-1]
	*h = (*h)[:len(*h)-1]
	return res
}

func SRTFScheduling(taskHeap *TaskHeap) []Task {
	var currentTimeSlice int = timeSlice
	executedTasks := make([]Task, 0)

	for taskHeap.Len() > 0 && currentTimeSlice > 0 {
		task := heap.Pop(taskHeap).(Task)

		if task.IsComplete() {
			continue
		}

		currentTimeSlice = task.Execute(currentTimeSlice)
		executedTasks = append(executedTasks, task)

		if !task.IsComplete() {
			heap.Push(taskHeap, task)
		}
	}

	return executedTasks
}
