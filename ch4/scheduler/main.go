package main

import (
	"log"
	"net/http"
)

var timeSlice int = 5

func main() {
	tm := NewTaskManager()
	defer tm.Stop()

	http.HandleFunc("/add_tasks", tm.AddTasksHandler)
	http.HandleFunc("/set_policy", tm.SetPolicyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
