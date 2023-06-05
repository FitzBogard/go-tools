package task_scheduler

type PriorityQueue struct {
	arr []*Task
}

func NewTaskQueue() *PriorityQueue {
	return &PriorityQueue{arr: []*Task{}}
}

func CreateTaskQueue(arr []*Task) {

}
