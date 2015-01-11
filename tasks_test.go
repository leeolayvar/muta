package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewTasker(t *testing.T) {
	Convey("Should initialize the Tasks map", t, func() {
		tr := NewTasker()
		So(tr.Tasks, ShouldNotBeNil)
	})
}

func TestTaskerTask(t *testing.T) {
	Convey("Should add a task", t, func() {
		ta := NewTasker()
		err := ta.Task("a", []string{}, func() {})
		So(err, ShouldBeNil)
		So(len(ta.Tasks), ShouldEqual, 1)
	})

	Convey("Should not allow replacing tasks", t, func() {
		ta := NewTasker()
		err := ta.Task("a", []string{}, func() {})
		So(err, ShouldBeNil)
		err = ta.Task("a", []string{}, func() {})
		So(err, ShouldNotBeNil)
	})
}

func TestTaskerRunTask(t *testing.T) {
	Convey("Should run a task", t, func() {
		ran := false
		ta := NewTasker()
		ta.Task("a", nil, func() {
			ran = true
		})
		err := ta.RunTask("a")
		So(err, ShouldBeNil)
		So(ran, ShouldBeTrue)
	})

	Convey("Should run task dependencies", t, func() {
		deps := []string{"b", "c"}
		called := []string{}
		ta := NewTasker()
		ta.Task("a", deps, func() {
			called = append(called, "a")
		})
		ta.Task("b", nil, func() {
			called = append(called, "b")
		})
		ta.Task("c", nil, func() {
			called = append(called, "c")
		})
		err := ta.RunTask("a")
		So(err, ShouldBeNil)
		So(called, ShouldContain, "b")
		So(called, ShouldContain, "c")
	})

	Convey("Should error on circular dependencies like", t, func() {
		Convey("a[a]", nil)
		Convey("a[b], b[a]", nil)
		Convey("a[b], b[c], c[a]", nil)
	})
}
