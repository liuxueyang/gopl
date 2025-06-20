package main

func front(tasks []Task) *Task {
	if len(tasks) == 0 {
		return nil
	}
	return &tasks[0]
}

func pop(tasks []Task) []Task {
	if len(tasks) == 0 {
		return tasks
	}
	return tasks[1:]
}

func push_front(tasks []Task, task Task) []Task {
	if len(tasks) == 0 {
		return []Task{task}
	}
	return append([]Task{task}, tasks...)
}

func FIFOScheduling(tasks []Task) []Task {
	var currentTimeSlice int = timeSlice
	executedTasks := make([]Task, 0)

	for len(tasks) > 0 && currentTimeSlice > 0 {
		task := front(tasks)
		tasks = pop(tasks)

		if task == nil || task.IsComplete() {
			continue
		}

		currentTimeSlice = task.Execute(currentTimeSlice)
		executedTasks = append(executedTasks, *task)

		if !task.IsComplete() {
			tasks = push_front(tasks, *task)
		}
	}

	return executedTasks
}
