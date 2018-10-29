package tspool

import (
	"testing"
)

func TestDefaultWorkerPool(t *testing.T) {
	wp, err := DefaultWorkerPool()
	if err != nil {
		t.Error(err)
	}
	wp1 := wp.(*defaultWorkerPool)
	if wp1.initCap != defaultWorkerPoolCap {
		t.Errorf("DefaultWorkerPool falied. Got %d, expected %d\n", wp1.initCap, defaultWorkerPoolCap)
	}
	if wp1.maxCap != defaultWorkerPoolCap {
		t.Errorf("DefaultWorkerPool falied. Got %d, expected %d\n", wp1.maxCap, defaultWorkerPoolCap)
	}

	wp, err = DefaultWorkerPool(10)
	if err != nil {
		t.Error(err)
	}
	wp2 := wp.(*defaultWorkerPool)
	if wp2.initCap != 10 {
		t.Errorf("DefaultWorkerPool falied. Got %d, expected 10\n", wp2.initCap)
	}
	if wp2.maxCap != 10 {
		t.Errorf("DefaultWorkerPool falied. Got %d, expected 10\n", wp2.maxCap)
	}

	wp, err = DefaultWorkerPool(10, 20)
	if err != nil {
		t.Error(err)
	}
	wp3 := wp.(*defaultWorkerPool)
	if wp3.initCap != 10 {
		t.Errorf("DefaultWorkerPool falied. Got %d, expected 10\n", wp3.initCap)
	}
	if wp3.maxCap != 20 {
		t.Errorf("DefaultWorkerPool falied. Got %d, expected 20\n", wp3.maxCap)
	}
	_, err = DefaultWorkerPool(10, 9)
	if err == nil {
		t.Error(errCap)
	}
	_, err = DefaultWorkerPool(10, 20, 30)
	if err == nil {
		t.Error(errArgs)
	}
}

func TestGet(t *testing.T) {
	wp, err := DefaultWorkerPool(10)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 10; i++ {
		w, err := wp.Get()
		if err != nil {
			t.Error(err)
		}
		dw := w.(*defaultWorker)
		if dw.pos != uint(i) {
			t.Errorf("Get falied. Got %d, expected %d\n", i, dw.pos)
		}
	}
	_, err = wp.Get()
	if err == nil {
		t.Errorf("DefaultWorkerPool falied. Got nil, expected %s\n", errWorker)
	}
}

func TestPut(t *testing.T) {
	wp, err := DefaultWorkerPool(10)
	var ws []Worker
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 10; i++ {
		w, err := wp.Get()
		if err != nil {
			t.Error(err)
		}
		ws = append(ws, w)
	}
	for _, w := range ws {
		err := wp.Put(w)
		if err != nil {
			t.Error(err)
		}
	}

}
