package stuk

import (
	"fmt"
	"testing"
	"time"
)

func TestCase1(t *testing.T) {
	c1 := New(time.Second*1, dealTimeoutCallback)
	c3 := New(time.Second*3, nil)
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

func dealTimeoutCallback(k uint64, v any) {
	fmt.Println("timeout", "key:", k, "value:", v)
}

func BenchmarkNewCase2(b *testing.B) {
	c := New(time.Second*1, nil)
	for i := uint64(0); i < 10000; i++ {
		c.Set(i, i)
	}
}

func TestCase3(t *testing.T) {
	c := New(time.Second*3, nil)
	c.Set(1, 2)
	c.Set(2, 5)
	fmt.Println(c.Get(1), c.Get(2))

	select {
	case <-time.After(time.Second):
	}
	fmt.Println(c.Get(1), c.Pull(2))

	select {
	case <-time.After(time.Second):
	}
	fmt.Println(c.Get(1), c.Pull(2))

	select {
	case <-time.After(time.Second):
	}
	fmt.Println(c.Get(1), c.Pull(2))

	select {
	case <-time.After(time.Second):
	}
	fmt.Println(c.Get(1), c.Pull(2))

	select {
	case <-time.After(time.Second):
	}
	fmt.Println(c.Get(1), c.Pull(2))
}
