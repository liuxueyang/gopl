package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Task 结构体表示一个任务
type Task struct {
	ID          int       `json:"id"`
	Name        string    `json:"name,omitempty"`
	Duration    int       `json:"duration"`    // 总需要时间(秒)
	Remaining   int       `json:"remaining"`   // 剩余时间(秒)
	ArrivalTime time.Time `json:"arrivalTime"` // 到达时间
	StartTime   time.Time `json:"startTime"`   // 开始时间
	EndTime     time.Time `json:"endTime"`     // 结束时间
	Status      string    `json:"status"`      // "queued", "running", "completed"
}

// TaskManager 管理所有任务
type TaskManager struct {
	mu             sync.Mutex
	tasks          map[int]*Task
	nextID         int
	fifoQueue      []*Task
	srtfHeap       SRTFHeap
	currentTask    *Task
	schedulingMode string // "FIFO" 或 "SRTF"
	clockTicker    *time.Ticker
	stopChan       chan struct{}
}

// SRTFHeap 实现最小堆接口用于SRTF调度
type SRTFHeap []*Task

func (h SRTFHeap) Len() int           { return len(h) }
func (h SRTFHeap) Less(i, j int) bool { return h[i].Remaining < h[j].Remaining }
func (h SRTFHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *SRTFHeap) Push(x interface{}) {
	*h = append(*h, x.(*Task))
}

func (h *SRTFHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// NewTaskManager 创建新的任务管理器
func NewTaskManager() *TaskManager {
	tm := &TaskManager{
		tasks:          make(map[int]*Task),
		nextID:         1,
		schedulingMode: "FIFO",
		stopChan:       make(chan struct{}),
	}

	// 启动调度时钟
	tm.clockTicker = time.NewTicker(1 * time.Second)
	go tm.schedulerLoop()

	return tm
}

// Stop 停止任务管理器
func (tm *TaskManager) Stop() {
	close(tm.stopChan)
	tm.clockTicker.Stop()
}

// schedulerLoop 调度器循环，每秒触发一次
func (tm *TaskManager) schedulerLoop() {
	for {
		select {
		case <-tm.clockTicker.C:
			tm.tick()
			tm.printStatus()
		case <-tm.stopChan:
			return
		}
	}
}

// tick 处理一个时间单位的调度
func (tm *TaskManager) tick() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// 更新当前任务的剩余时间
	if tm.currentTask != nil && tm.currentTask.Status == "running" {
		tm.currentTask.Remaining--
		if tm.currentTask.Remaining <= 0 {
			tm.currentTask.Status = "completed"
			tm.currentTask.EndTime = time.Now()
			tm.currentTask = nil
		}
	}

	// 如果当前没有运行的任务，尝试获取下一个
	if tm.currentTask == nil || tm.currentTask.Status != "running" {
		tm.getNextTask()
	}
}

// getNextTask 获取下一个要执行的任务
func (tm *TaskManager) getNextTask() {
	var nextTask *Task

	switch tm.schedulingMode {
	case "FIFO":
		if len(tm.fifoQueue) > 0 {
			nextTask = tm.fifoQueue[0]
			tm.fifoQueue = tm.fifoQueue[1:]
		}
	case "SRTF":
		if tm.srtfHeap.Len() > 0 {
			nextTask = heap.Pop(&tm.srtfHeap).(*Task)
		}
	}

	if nextTask != nil {
		nextTask.StartTime = time.Now()
		nextTask.Status = "running"
		tm.currentTask = nextTask
	}
}

// printStatus 打印当前状态
func (tm *TaskManager) printStatus() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	now := time.Now().Format("15:04:05")
	fmt.Printf("\n=== 调度状态 [%s] ===\n", now)
	fmt.Printf("当前调度算法: %s\n", tm.schedulingMode)

	if tm.currentTask != nil {
		fmt.Printf("当前运行任务: %s (ID: %d, 剩余: %ds)\n",
			tm.currentTask.Name, tm.currentTask.ID, tm.currentTask.Remaining)
	} else {
		fmt.Println("当前运行任务: 无")
	}

	// 打印队列中的任务
	switch tm.schedulingMode {
	case "FIFO":
		fmt.Println("\nFIFO 队列:")
		for i, task := range tm.fifoQueue {
			fmt.Printf("%d. %s (ID: %d, 剩余: %ds)\n",
				i+1, task.Name, task.ID, task.Remaining)
		}
	case "SRTF":
		fmt.Println("\nSRTF 优先队列:")
		for i, task := range tm.srtfHeap {
			fmt.Printf("%d. %s (ID: %d, 剩余: %ds)\n",
				i+1, task.Name, task.ID, task.Remaining)
		}
	}

	fmt.Println("====================")
}

// AddTask 添加新任务
func (tm *TaskManager) AddTask(name string, duration int) *Task {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	task := &Task{
		ID:          tm.nextID,
		Name:        name,
		Duration:    duration,
		Remaining:   duration,
		ArrivalTime: time.Now(),
		Status:      "queued",
	}
	tm.tasks[task.ID] = task
	tm.nextID++

	// 根据调度模式添加到相应队列
	switch tm.schedulingMode {
	case "FIFO":
		tm.fifoQueue = append(tm.fifoQueue, task)
	case "SRTF":
		heap.Push(&tm.srtfHeap, task)
	}

	return task
}

// AddTasks 批量添加任务
func (tm *TaskManager) AddTasks(taskSpecs []TaskSpec) []*Task {
	var tasks []*Task
	for _, spec := range taskSpecs {
		task := tm.AddTask(spec.Name, spec.Duration)
		tasks = append(tasks, task)
	}
	return tasks
}

// GetAllTasks 获取所有任务
func (tm *TaskManager) GetAllTasks() []*Task {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tasks := make([]*Task, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// SetSchedulingMode 设置调度模式
func (tm *TaskManager) SetSchedulingMode(mode string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if mode == "FIFO" || mode == "SRTF" {
		tm.schedulingMode = mode
	}
}

// TaskSpec 用于接收客户端任务描述
type TaskSpec struct {
	Name     string `json:"name"`
	Duration int    `json:"duration"`
}

// HTTP 处理函数

func addTasksHandler(tm *TaskManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var taskSpecs []TaskSpec
		if err := json.NewDecoder(r.Body).Decode(&taskSpecs); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		for _, spec := range taskSpecs {
			if spec.Duration <= 0 {
				http.Error(w, "Duration must be positive", http.StatusBadRequest)
				return
			}
		}

		tasks := tm.AddTasks(taskSpecs)
		response := map[string]interface{}{
			"message":    "Tasks added successfully",
			"task_count": len(tasks),
			"first_task": tasks[0].ID,
			"last_task":  tasks[len(tasks)-1].ID,
		}
		json.NewEncoder(w).Encode(response)
	}
}

func getTasksHandler(tm *TaskManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tasks := tm.GetAllTasks()
		json.NewEncoder(w).Encode(tasks)
	}
}

func setSchedulerHandler(tm *TaskManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			Algorithm string `json:"algorithm"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if request.Algorithm != "FIFO" && request.Algorithm != "SRTF" {
			http.Error(w, "Invalid algorithm", http.StatusBadRequest)
			return
		}

		tm.SetSchedulingMode(request.Algorithm)
		response := map[string]string{"message": fmt.Sprintf("Scheduler set to %s", request.Algorithm)}
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	tm := NewTaskManager()
	defer tm.Stop()

	// 设置路由
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			addTasksHandler(tm)(w, r)
		case "GET":
			getTasksHandler(tm)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/scheduler", setSchedulerHandler(tm))

	// 启动服务器
	fmt.Println("任务调度服务器启动，每秒打印一次调度状态...")
	fmt.Println("API 端点:")
	fmt.Println("POST /tasks - 添加任务列表")
	fmt.Println("GET /tasks - 获取所有任务状态")
	fmt.Println("POST /scheduler - 设置调度算法 (FIFO 或 SRTF)")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
