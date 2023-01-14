package stuk

import (
	"fmt"
	"testing"
	"time"
)

func TestCase1(t *testing.T) {
	c1 := New(time.Second * 1)
	c3 := New(time.Second * 3)
	c1.Set(1, 2)
	c1.Set(2, 5)
	c3.Set(11, 12)
	c3.Set(12, 15)
	fmt.Println(c1.Get(1), c1.Get(2), c3.Get(11), c3.Get(12))

	select {
	case <-time.After(time.Second * 2):
	}

	fmt.Println(c1.Get(1), c1.Get(2), c3.Get(11), c3.Get(12))

	select {
	case <-time.After(time.Second * 6):
	}

	fmt.Println(c1.Get(1), c1.Get(2), c3.Get(11), c3.Get(12))
}

func BenchmarkNewCase2(b *testing.B) {
	c := New(time.Second * 1)
	for i := uint64(0); i < 10000; i++ {
		c.Set(i, i)
	}
}
