//
//
//
package muta

import (
	"errors"
	"fmt"
	"reflect"
)

type Handler func()
type ErrorHandler func() error
type ContextHandler func(Ctx *interface{}) error
type StreamHandler func() (*Stream, error)

var DefaultTasker *Tasker = NewTasker()

func Task(name string, args ...interface{}) error {
	return DefaultTasker.Task(name, args...)
}

func Run() {
	DefaultTasker.Run()
}

func NewTasker() *Tasker {
	return &Tasker{
		Tasks: make(map[string]*TaskerTask),
	}
}

type Tasker struct {
	Tasks map[string]*TaskerTask
}

type TaskerTask struct {
	Name           string
	Dependencies   []string
	Handler        Handler
	ErrorHandler   ErrorHandler
	StreamHandler  StreamHandler
	ContextHandler ContextHandler
}

func (tr *Tasker) Task(n string, args ...interface{}) error {
	if tr.Tasks[n] != nil {
		return errors.New("Task already exists")
	}

	ds := []string{}

	var (
		h  Handler
		er ErrorHandler
		sh StreamHandler
		ch ContextHandler
	)

	for _, arg := range args {
		v := reflect.ValueOf(arg)
		switch v.Type().String() {
		case "string":
			ds = append(ds, v.String())
		case "[]string":
			ds = append(ds, v.Interface().([]string)...)
		case "func()":
			h = v.Interface().(func())
			break
		default:
			return errors.New(fmt.Sprintf(
				"unsupported task argument type '%s'", v.Type().String(),
			))
		}
	}

	tr.Tasks[n] = &TaskerTask{
		Name:           n,
		Dependencies:   ds,
		Handler:        h,
		ErrorHandler:   er,
		StreamHandler:  sh,
		ContextHandler: ch,
	}

	return nil
}

func (tr *Tasker) Run() error {
	return tr.RunTask("default")
}

func (tr *Tasker) RunTask(tn string) error {
	t := tr.Tasks[tn]
	if t == nil {
		return errors.New(fmt.Sprintf("Task \"%s\" does not exist.", tn))
	}

	if t.Dependencies != nil {
		for _, d := range t.Dependencies {
			tr.RunTask(d)
		}
	}

	if t.Handler != nil {
		t.Handler()
	}
	return nil
}
