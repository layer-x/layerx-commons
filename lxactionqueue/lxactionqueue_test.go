package lxactionqueue_test

import (
	. "github.com/layer-x/layerx-commons/lxactionqueue"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lxactionqueue", func() {
	Describe("Size", func(){
		It("returns empty when action queue is empty", func(){
			aq := NewActionQueue()
			Expect(aq.Size()).To(Equal(0))
		})
	})
	Describe("Push(func)", func(){
		It("adds an action to the queue", func(){
			fun1 := func(){}
			aq := NewActionQueue()
			aq.Push(fun1)
			Expect(aq.Size()).To(Equal(1))
		})
	})
	Describe("Pop()", func(){
		It("pops the first function off the queue", func(){
			fun1 := func(){}
			aq := NewActionQueue()
			aq.Push(fun1)
			aq.Pop()
			Expect(aq.Size()).To(Equal(0))
		})
	})
	Describe("action.Call", func(){
		It("calls the function", func(){
			strc := make(chan string)
			fun1 := func(){
				go func(){
					strc <- "test_1"
				}()
			}
			fun2 := func(){
				go func(){
					strc <- "test_2"
				}()
			}
			aq := NewActionQueue()
			aq.Push(fun1)
			aq.Push(fun2)
			action := aq.Pop()
			Expect(aq.Size()).To(Equal(1))
			action.Call()
			str := <- strc
			Expect(str).To(Equal("test_1"))
		})
	})
})
