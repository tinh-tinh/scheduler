package scheduler

import (
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type Task struct {
	core.DynamicProvider
	module core.Module
}

func NewTask(module core.Module) *Task {
	return &Task{
		module: module,
	}
}

func (t *Task) Cron(spec string, fnc func()) {
	registry := Inject(t.module)

	registry.AddFunc(spec, fnc)
	registry.Start()
}
