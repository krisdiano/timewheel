package timewheel

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var opt []string

func f(s string) {
	opt = append(opt, fmt.Sprintf("%v %s\n", time.Now(), s))
}

func TestRun(t *testing.T) {
	ids := []int{}
	now := time.Now()
	w := New(3, time.Second)
	e1, _ := w.AddFunc(now.Add(2*time.Second), func() { f("1st") })
	e2, _ := w.AddFunc(now.Add(10*time.Second), func() { f("2nd") })
	e3, _ := w.AddFunc(now.Add(11*time.Second), func() { f("3rd") })
	e4, _ := w.AddFunc(now.Add(13*time.Second), func() { f("4th") })

	ids = append(ids, e1.Id, e2.Id, e3.Id, e4.Id)
	idx := rand.Intn(len(ids))
	w.DelFunc(idx)

	time.Sleep(time.Second * 15)
	if len(opt) != len(ids)-1 {
		t.Fatal("add or del failed")
	}
}

func TestWrapper(t *testing.T) {
	var rst string
	var excpted string = "m1im2itaskm2om1o"

	w := New(3, time.Second)
	einfo, err := w.AddFunc(time.Now().Add(time.Second*2), func() { rst += "task" })
	if err != nil {
		t.Fatal("add func failed")
	}

	m1 := func() {
		rst += ("m1i")
		einfo.Next()
		rst += ("m1o")
	}

	m2 := func() {
		rst += ("m2i")
		einfo.Next()
		rst += ("m2o")
	}
	w.AddWrappers(einfo.Id, m1, m2)

	time.Sleep(4 * time.Second)
	if rst == excpted {
		t.Log("TestWrapper Ok")
	} else {
		t.Fatalf("TestWrapper Failed, real is %s, excpted is %s\n", rst, excpted)
	}
}
