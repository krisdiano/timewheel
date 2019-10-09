package timewheel

import (
	"math"
	"time"
)

type entry struct {
	id   int
	ts   time.Time
	fn   func()
	next *entry
}

type level struct {
	index  int
	buffer []*entry
	next   *level
}

type timewheel struct {
	root   level
	tick   time.Duration
	length int // the lenth of every layer is the same
	
	cnt    int // number of layers of time wheel
	id     int

	sadd  chan *entry
	sdel  chan int
	exist map[int]struct{}
}

// At the same time, it runs. 
func New(length int, tick time.Duration) *timewheel {
	tw := &timewheel{
		tick:   tick,
		length: length,
		cnt:    1,
		sadd:   make(chan *entry, 1),
		sdel:   make(chan int, 1),
	}
	tw.root.buffer = make([]*entry, length)
	go tw.run()
	
	return tw
}

// which [0,cnt)
func (tw *timewheel) pos(ts time.Time) (which int, index int) {
	now := time.Now()
	ptr := &tw.root
	bound := []time.Time{now.Add(time.Duration(tw.length) * tw.tick)}
	for i := 0; bound[i].Before(ts); i++ {
		which++
		if ptr.next == nil {
			ptr.next = &level{
				buffer: make([]*entry, tw.length),
			}
			tw.cnt++
		}
		ptr = ptr.next
		tem := bound[i].Add(time.Duration(int64(math.Pow(float64(tw.length), float64(i+2)))) * tw.tick)
		bound = append(bound, tem)
	}
	
	if which == 0 {
		index = int((ts.Sub(now))/tw.tick) % tw.length
	} else {
		last := len(bound) - 1
		start, _ := bound[last-1], bound[last]
		for {
			start = start.Add(time.Duration(int64(math.Pow(float64(tw.length), float64(which)))) * tw.tick)
			if !start.Before(ts) {
				break
			}
			index++
		}
	}
	
	return
}

func (tw *timewheel) AddFunc(ts time.Time, fn func()) (id int, err error) {
	if !ts.After(time.Now()) {
		go fn()
	}

	e := &entry{
		id: tw.id,
		ts: ts,
		fn: fn,
	}

	if tw.exist == nil {
		tw.exist = make(map[int]struct{})
	}
	tw.exist[tw.id] = struct{}{}
	tw.id++
	tw.sadd <- e
	
	return
}

func (tw *timewheel) addFunc(e *entry) (id int, err error) {
	ts := e.ts
	which, idx := tw.pos(ts)
	wheel := &tw.root
	base := tw.root.index
	
	if which != 0 {
		base = 0
		for which != 0 {
			which--
			wheel = wheel.next
		}
	}
	
	idx = (base + idx) % tw.length
	if wheel.buffer[idx] == nil {
		wheel.buffer[idx] = e
		return e.id, nil
	}
	e.next, wheel.buffer[idx] = wheel.buffer[idx], e
	
	return
}

func (tw *timewheel) DelFunc(id int) {
	_, ok := tw.exist[id]
	if !ok {
		return
	}
	tw.sdel <- id
}

func (tw *timewheel) delFunc(id int) {
	delete(tw.exist, id)
}

func (tw *timewheel) run() {
	ticker := time.NewTicker(tw.tick)

	for {
		select {
		case <-ticker.C:
			elem := tw.root.buffer[tw.root.index]
			for elem != nil {
				_, ok := tw.exist[elem.id]
				if ok {
					go elem.fn()
				}
				elem = elem.next
			}
			tw.root.buffer[tw.root.index] = nil
			tw.root.index = (tw.root.index + 1) % tw.length
			tw.deliver(&tw.root)
		case a := <-tw.sadd:
			tw.addFunc(a)
		case d := <-tw.sdel:
			tw.delFunc(d)
		}
	}
}

func (tw *timewheel) deliver(wheel *level) {
	if wheel.index != 0 || wheel.next == nil {
		return
	}

	wheel = wheel.next
	elem := wheel.buffer[wheel.index]
	for elem != nil {
		toAdd := elem
		elem = elem.next
		toAdd.next = nil
		tw.addFunc(toAdd)
	}
	wheel.buffer[wheel.index] = nil
	wheel.index = (wheel.index + 1) % tw.length
	tw.deliver(wheel)
}
