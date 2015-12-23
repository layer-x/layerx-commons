package lxactionqueue

type ActionQueue interface {
	Size() int
	Push(func())
	Pop () *action
	ExecuteNext()
}

//action
type action struct {
	callback func()
}

func (a *action) Call() {
	a.callback()
}

//queue
type actionQueue struct {
	actions []action
}

func NewActionQueue() *actionQueue {
	return &actionQueue{}
}

func (a *actionQueue) Size() int {
	return len(a.actions)
}

func (a *actionQueue) Push(callback func()) {
	a.actions = append(a.actions, action{callback: callback})
}

func (a *actionQueue) Pop () *action {
	if len(a.actions) < 1 {
		return nil
	}
	if len(a.actions) < 2 {
		actn := a.actions[0]
		a.actions = []action{}
		return &actn
	}
	actn := a.actions[0]
	a.actions = a.actions[1:]
	return &actn
}

func (a *actionQueue) ExecuteNext() {
	actn := a.Pop()
	if actn != nil {
		actn.Call()
	}
}