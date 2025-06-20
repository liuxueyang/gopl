package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type SchedulePolicy int

const (
	FIFO SchedulePolicy = iota
	SRTF
)

type TaskManager struct {
	lock      sync.Mutex
	tasks     []Task
	tasksHeap *TaskHeap
	nextID    int
	policy    SchedulePolicy
	ticker    *time.Ticker
	stop      chan bool
	ts        int
}

func NewTaskManager() *TaskManager {
	res := &TaskManager{
		tasks:  make([]Task, 0),
		nextID: 0,
		policy: FIFO,
		ticker: time.NewTicker(1 * time.Second),
		stop:   make(chan bool),
	}

	res.ticker = time.NewTicker(1 * time.Second)
	go res.TaskLoop()

	return res
}

func (tm *TaskManager) Stop() {
	close(tm.stop)
	tm.ticker.Stop()
}

func (tm *TaskManager) TaskLoop() {
	for {
		select {
		case <-tm.ticker.C:
			tm.schedule()
		case <-tm.stop:
			return
		}
	}
}

func (tm *TaskManager) schedule() {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	if len(tm.tasks) == 0 && (tm.tasksHeap == nil || len(*tm.tasksHeap) == 0) {
		return
	}

	var executedTasks []Task

	switch tm.policy {
	case FIFO:
		executedTasks = FIFOScheduling(tm.tasks)
	case SRTF:
		executedTasks = SRTFScheduling(tm.tasksHeap)
	}

	if len(executedTasks) > 0 {
		fmt.Printf("%d\t", tm.ts)
		printExecutedTasks(executedTasks)

		tm.ts++
	}
}

func (tm *TaskManager) SetPolicy(policy SchedulePolicy) {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	if policy == tm.policy {
		return
	}

	if policy == FIFO {
		for len(*tm.tasksHeap) > 0 {
			task := heap.Pop(tm.tasksHeap).(Task)
			tm.tasks = append(tm.tasks, task)
		}

		sort.Slice(tm.tasks, func(i, j int) bool {
			return tm.tasks[i].ID < tm.tasks[j].ID
		})
		tm.tasksHeap = nil
	} else if policy == SRTF {
		tm.tasksHeap = &TaskHeap{}
		for _, task := range tm.tasks {
			heap.Push(tm.tasksHeap, task)
		}
		heap.Init(tm.tasksHeap)
		tm.tasks = nil
	}
	tm.policy = policy
}

func (tm *TaskManager) AddTasks(task []Task) {
	tm.lock.Lock()
	defer tm.lock.Unlock()

	for _, task := range task {
		task.ID = tm.nextID
		tm.nextID++

		if tm.policy == FIFO {
			tm.tasks = append(tm.tasks, task)
		} else if tm.policy == SRTF {
			if tm.tasksHeap == nil {
				tm.tasksHeap = &TaskHeap{}
			}
			heap.Push(tm.tasksHeap, task)
		}
	}
}

func (tm *TaskManager) AddTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type AddTasksRequest struct {
		Tasks []Task `json:"tasks"`
	}
	var addTasksReq AddTasksRequest

	if err := json.NewDecoder(r.Body).Decode(&addTasksReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		http.Error(w, fmt.Sprintf("body: %v, err: %v", r.Body, err),
			http.StatusBadRequest)
		return
	}

	tm.AddTasks(addTasksReq.Tasks)
	fmt.Fprintf(w, "Added %d tasks\n", len(addTasksReq.Tasks))
}

func (tm *TaskManager) SetPolicyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type PolicyRequest struct {
		Policy string `json:"policy"`
	}

	var policy PolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&policy); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	switch policy.Policy {
	case "FIFO", "fifo":
		tm.SetPolicy(FIFO)
	case "SRTF", "srtf":
		tm.SetPolicy(SRTF)
	default:
		http.Error(w, "Unknown policy", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Set scheduling policy to %s\n", policy.Policy)
}
