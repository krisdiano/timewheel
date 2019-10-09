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
	i1, _ := w.AddFunc(now.Add(2*time.Second), func() { f("1st") })
	i2, _ := w.AddFunc(now.Add(10*time.Second), func() { f("2nd") })
	i3, _ := w.AddFunc(now.Add(11*time.Second), func() { f("3rd") })
	i4, _ := w.AddFunc(now.Add(13*time.Second), func() { f("4th") })

	ids = append(ids, i1, i2, i3, i4)
	idx := rand.Intn(len(ids))
	w.DelFunc(idx)

	time.Sleep(time.Second * 15)
	if len(opt) != len(ids)-1 {
		t.Fatal("add or del failed")
	}
}
