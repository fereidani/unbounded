package unbounded

import (
	"sync"
	"testing"
	"time"
)

// TestUnboundedChannel tests the basic functionality of the unbounded channel.
func TestUnboundedChannel(t *testing.T) {
	ch := New[int]()

	// Test sending and receiving.
	ch.Send(1)
	ch.Send(2)

	if val, ok := ch.Receive(); !ok || val != 1 {
		t.Errorf("expected 1, got %v", val)
	}

	if val, ok := ch.Receive(); !ok || val != 2 {
		t.Errorf("expected 2, got %v", val)
	}

	// Test receiving from an empty channel (should block).
	go func() {
		time.Sleep(1 * time.Second)
		ch.Send(3)
	}()

	if val, ok := ch.Receive(); !ok || val != 3 {
		t.Errorf("expected 3, got %v", val)
	}

	// Test closing the channel.
	ch.Send(4)
	ch.Close()

	if val, ok := ch.Receive(); !ok || val != 4 {
		t.Errorf("expected 4, got %v", val)
	}

	// Test panic on sending to a closed channel.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on send to closed channel")
		}
	}()
	ch.Send(5)
}

// TestManySend test sending 1million message to channel
func TestManySend(t *testing.T) {
	ch := New[int]()
	for i := 0; i < 10000000; i++ {
		ch.Send(i)
	}
	for i := 0; i < 10000000; i++ {
		val, _ := ch.Receive()
		if val != i {
			t.Fail()
		}
	}
}

// TestConcurrentSendReceive tests concurrent send and receive operations.
func TestConcurrentSendReceive(t *testing.T) {
	ch := New[int]()
	var wg sync.WaitGroup

	// Start multiple senders.
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ch.Send(i)
		}(i)
	}

	// Start multiple receivers.
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			val, ok := ch.Receive()
			if !ok {
				t.Errorf("expected value, got nothing")
			}
			t.Logf("received %v", val)
		}()
	}

	wg.Wait()
}

// TestReceiveFromEmptyChannel tests receiving from an empty channel.
func TestReceiveFromEmptyChannel(t *testing.T) {
	ch := New[int]()
	done := make(chan struct{})

	go func() {
		defer close(done)
		val, ok := ch.Receive()
		if !ok {
			t.Errorf("expected value, got nothing")
		}
		if val != 42 {
			t.Errorf("expected 42, got %v", val)
		}
	}()

	time.Sleep(1 * time.Second) // Ensure the goroutine is waiting.
	ch.Send(42)
	<-done
}

// TestCloseWithPendingReceives tests closing the channel with pending receives.
func TestCloseWithPendingReceives(t *testing.T) {
	ch := New[int]()
	done := make(chan struct{})

	go func() {
		defer close(done)
		_, ok := ch.Receive()
		if ok {
			t.Errorf("expected closed channel, got value")
		}
	}()

	time.Sleep(1 * time.Second) // Ensure the goroutine is waiting.
	ch.Close()
	<-done
}
