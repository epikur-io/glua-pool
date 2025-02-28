package pool

import (
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func TestAqcuireAndRelease(t *testing.T) {
	lpool := NewPool(2, nil)
	lvms := []*lua.LState{}
	for range 2 {
		lvm := lpool.Acquire()
		lvms = append(lvms, lvm)
	}

	plen := lpool.Len()
	if plen != 0 {
		t.Errorf("pool expected to be empty but %d instances remained", plen)
	}

	// should timeout since pool is now empty
	_, err := lpool.AcquireWithTimeout(1 * time.Second)
	if err == nil {
		t.Errorf("expected timout error but got %v", err)
	}

	for _, lvm := range lvms {
		lpool.Release(lvm)
	}

	plen = lpool.Len()
	if plen != 2 {
		t.Errorf("pool expected to be full but got %d instances", plen)
	}
}

func TestUpdate(t *testing.T) {
	lpool := NewPool(2, nil)
	lvms := []*lua.LState{}
	for range 2 {
		lvm := lpool.Acquire()
		lvms = append(lvms, lvm)
	}
	delay := 1 * time.Second
	go func() {
		time.Sleep(delay)
		for _, lvm := range lvms {
			lpool.Release(lvm)
		}
	}()
	start := time.Now()
	lpool.Update()
	duration := time.Since(start)

	if duration < delay {
		// check if channel was blocking
		t.Errorf("expected pool to block for 1 second but got %v", duration)
	}

}

func TestUpdateTimeout(t *testing.T) {
	lpool := NewPool(3, nil)
	for range 3 {
		lpool.Acquire()
	}
	removedInstances, updatedInstances := lpool.UpdateWithTimeout(1 * time.Second)
	if removedInstances != 0 {
		t.Errorf("expected %d removed instances but got %d", lpool.Len(), removedInstances)
	}
	if updatedInstances != 0 {
		t.Errorf("expected %d updated instances but got %d", lpool.Len(), updatedInstances)
	}
}
