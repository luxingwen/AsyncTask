package task

import "fmt"

type TaskState int

const (
	TaskStateActive TaskState = iota + 1
	TaskStateDoing
	TaskStateFailed
	TaskStateDone
)

type Task interface {
	Create() error
	State() TaskState
	CheckDone() (bool, error)
	Done() error
}

func TaskRun(t Task) {
	state := t.State()
	if state == TaskStateActive {
		err := t.Create()
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	if state == TaskStateDoing {
		ok, err := t.CheckDone()
		if err != nil || ok {
			if err != nil {
				fmt.Println(err)
			}
			err = t.Done()
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}
		return
	}

	if state == TaskStateFailed || state == TaskStateDone {
		err := t.Done()
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}

}
