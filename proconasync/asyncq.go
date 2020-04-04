package proconasync

import "log"

// Task interface
type Task interface {
	Perform()
}

var (
	//TaskQueue channel
	TaskQueue chan Task;
	//TaskWorkerQueue channel
	TaskWorkerQueue chan chan Task 
)

// TaskWorker is a struct 
type TaskWorker struct {
	ID int
	TaskChannel chan Task
	TaskWorkerQueue chan chan Task

}

// NewTaskWorker takes and id int and taskWorkerQue  and returns a TaskWorker
func NewTaskWorker(id int, taskWorkerQueue chan chan Task) TaskWorker{
	taskWorker := TaskWorker{
		ID: id,
		TaskChannel: make(chan Task),
		TaskWorkerQueue: taskWorkerQueue,
	}
	return taskWorker 
}
//Start method has a TaskWorker Receiver
func (t *TaskWorker) Start(){
	go func ()  {
		for  {
			t.TaskWorkerQueue <- t.TaskChannel
			select {
			case task := <-t.TaskChannel:
				log.Printf("Asyncv task workier #%d performing a task.\n", t.ID)
				task.Perform()
			}
		}
	}()
	
}
// StartTaskDispatcher function takes a taskWorkerSize and fills a TaskWorkerQue to the
func StartTaskDispatcher(taskWorkerSize int)  {

	TaskWorkerQueue = make(chan chan Task, taskWorkerSize)
	for i := 0; i < taskWorkerSize; i++ {
		log.Print("Starting async task worker #", i+1)
		taskWorker := NewTaskWorker(i+1,TaskWorkerQueue)
		taskWorker.Start() 
	}

	go func() {
		for {
			select {
			case task := <-TaskQueue:
				go func() {
					taskChannel := <-TaskWorkerQueue
					taskChannel <- task
				}()
			}
		}
		
	}()
}

func init() {
	TaskQueue = make(chan Task, 100)
}