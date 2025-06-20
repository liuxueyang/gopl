package main

import (
	"container/heap"
	"fmt"
	"testing"
)

// test for SRTF scheduling policy
func TestSRTFScheduling(t *testing.T) {
	tasks := []Task{
		{ID: 0, Remaining: 2},
		{ID: 1, Remaining: 4},
		{ID: 2, Remaining: 100},
		{ID: 3, Remaining: 6},
		{ID: 4, Remaining: 10},
		{ID: 5, Remaining: 90},
		{ID: 6, Remaining: 1},
		{ID: 7, Remaining: 1},
		{ID: 8, Remaining: 10},
		{ID: 9, Remaining: 2},
		{ID: 10, Remaining: 15},
		{ID: 11, Remaining: 30},
		{ID: 12, Remaining: 1},
		{ID: 13, Remaining: 5},
		{ID: 14, Remaining: 9},
		{ID: 15, Remaining: 10},
	}

	taskHeap := &TaskHeap{}
	for _, task := range tasks {
		heap.Push(taskHeap, task)
	}

	fmt.Println("SRTF Scheduling")

	for len(*taskHeap) > 0 {
		executedTasks := SRTFScheduling(taskHeap)
		if len(executedTasks) == 0 {
			break
		}

		printExecutedTasks(executedTasks)
	}
}
