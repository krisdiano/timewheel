package timewheel

import "time"

type entry struct {
	id int
	ts time.Time
	fn func()

	middlewares []func()
	index       int

	next *entry
}

func (e *entry) nextHandler() {
	e.index++

	if e.index-1 == len(e.middlewares) {
		e.fn()
		return
	}

	/*
		for i := e.index - 1; i < len(e.middlewares); i++ {
			e.middlewares[i]()
		}
	*/
	e.middlewares[e.index-1]()

}

func (e *entry) run() {
	e.nextHandler()
}

type EntryInfo struct {
	Id   int
	Next func()
}
