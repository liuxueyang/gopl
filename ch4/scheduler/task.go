package main

import "fmt"

type Task struct {
	ID        int
	Remaining int `json:"cost"`
	LastRun   int
}

// return the remaining duration after executing the task
func (t *Task) Execute(duration int) int {
	if t.Remaining >= duration {
		t.Remaining -= duration
		t.LastRun = duration
		duration = 0
	} else {
		t.LastRun = t.Remaining
		duration -= t.Remaining
		t.Remaining = 0
	}
	return duration
}

func (t *Task) IsComplete() bool {
	return t.Remaining <= 0
}

func (t Task) String() string {
	return fmt.Sprintf("Task %d: %d-%d=%d", t.ID, t.Remaining+t.LastRun, t.LastRun, t.Remaining)
}

func printExecutedTasks(executedTasks []Task) {
	if len(executedTasks) == 0 {
		return
	}

	for i, executedTask := range executedTasks {
		fmt.Printf("%d", executedTask.ID)
		if i < len(executedTasks)-1 {
			fmt.Printf(", ")
		} else {
			fmt.Printf("\t")
		}
	}

	for i, executedTask := range executedTasks {
		fmt.Printf("%d", executedTask.Remaining)
		if i < len(executedTasks)-1 {
			fmt.Printf(", ")
		} else {
			fmt.Printf(" ")
		}
	}

	fmt.Printf("(")
	for i, executedTask := range executedTasks {
		fmt.Printf("%s", executedTask)
		if i < len(executedTasks)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Printf(")\n")
}
