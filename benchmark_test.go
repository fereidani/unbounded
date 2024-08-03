package unbounded

import (
	"sync"
	"testing"
)

const COROUTINES_PER_SIDE = 8

func BenchmarkUnboundedChannelMPMC(b *testing.B) {
	ch := New[int]()
	for i := 0; i < COROUTINES_PER_SIDE; i++ {
		go func() {
			for i := 0; i < b.N/COROUTINES_PER_SIDE; i++ {
				ch.Send(i)
			}
		}()
	}
	var wg sync.WaitGroup
	wg.Add(COROUTINES_PER_SIDE)
	for i := 0; i < COROUTINES_PER_SIDE; i++ {
		go func() {
			for i := 0; i < b.N/COROUTINES_PER_SIDE; i++ {
				ch.Receive()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkGoChannelMPMC(b *testing.B) {
	ch := make(chan int, 128)
	for i := 0; i < COROUTINES_PER_SIDE; i++ {
		go func() {
			for i := 0; i < b.N/COROUTINES_PER_SIDE; i++ {
				ch <- i
			}
		}()
	}
	var wg sync.WaitGroup
	wg.Add(COROUTINES_PER_SIDE)
	for i := 0; i < COROUTINES_PER_SIDE; i++ {
		go func() {
			for i := 0; i < b.N/COROUTINES_PER_SIDE; i++ {
				<-ch
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkUnboundedChannelSPSC(b *testing.B) {
	ch := New[int]()
	go func() {
		for i := 0; i < b.N; i++ {
			ch.Send(i)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < b.N; i++ {
			ch.Receive()
		}
		wg.Done()
	}()
	wg.Wait()
}

func BenchmarkGoChannelSPSC(b *testing.B) {
	ch := make(chan int, 128)
	go func() {
		for i := 0; i < b.N; i++ {
			ch <- i
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < b.N; i++ {
			<-ch
		}
		wg.Done()
	}()
	wg.Wait()
}
