package etcdsync

import (
	"flag"
	"log"
	"testing"
	"time"

	"github.com/coreos/go-etcd/etcd"
)

var key string = "test/mutex"

func init() {
	flag.Parse()
}

func TestTwoNoKey(t *testing.T) {

	//etcd.SetLogger(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))

	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	client.Delete(key, true)

	quit1 := make(chan bool)
	quit2 := make(chan bool)
	errchan := make(chan bool)

	progress := make(chan bool)

	// first thread
	go func() {

		mutex := NewMutexFromClient(client, key, 0)
		err := mutex.Lock()
		if err != nil {

			errchan <- true
		}

		progress <- true

		// sleep for 5 seconds, ttl should be refreshed after 3 seconds
		time.Sleep(5 * time.Second)
		mutex.Unlock()

		quit1 <- true
	}()

	select {
	case <-progress:
	case <-errchan:
		t.Fatal("could not acquire lock, is etcd running?")
	}

	// second thread
	go func() {

		mutex := NewMutexFromClient(client, key, 0)
		//mutex := NewMutexFromServers([]string{"http://127.0.0.1:4001/"}, key, 0)
		// should take us 5 seconds to acquire the lock
		now := time.Now()

		err := mutex.Lock()
		if err != nil {

			t.Fatal("could not acquire lock, is etcd running?", err)
			errchan <- true
		}

		timeToLock := time.Since(now)
		if timeToLock < 5*time.Second {

			t.Fatalf("mutex TTL was not refreshed, lock acquired after %v seconds", timeToLock)
		}

		mutex.Unlock()
		quit2 <- true
	}()

	var (
		q1 bool
		q2 bool
	)

	for !q1 || !q2 {
		select {
		case <-quit1:
			q1 = true
		case <-quit2:
			q2 = true
		case <-errchan:
			t.Fatal("could not acquire lock, is etcd running?")
		}
	}
}

func TestTwoExistingKey(t *testing.T) {

	//etcd.SetLogger(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))

	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	client.Set(key, "released", 0)

	quit1 := make(chan bool)
	quit2 := make(chan bool)
	errchan := make(chan bool)

	progress := make(chan bool)

	// first thread
	go func() {

		mutex := NewMutexFromServers([]string{"http://127.0.0.1:4001"}, key, 0)
		err := mutex.Lock()
		if err != nil {

			errchan <- true
		}

		progress <- true

		// sleep for 5 seconds, ttl should be refreshed after 3 seconds
		time.Sleep(5 * time.Second)
		mutex.Unlock()

		quit1 <- true
	}()

	select {
	case <-progress:
	case <-errchan:
		t.Fatal("could not acquire lock, is etcd running?")
	}

	// second thread
	go func() {

		mutex := NewMutexFromClient(client, key, 0)
		//mutex := NewMutexFromServers([]string{"http://127.0.0.1:4001/"}, key, 0)
		// should take us 5 seconds to acquire the lock
		now := time.Now()

		err := mutex.Lock()
		if err != nil {

			errchan <- true
		}

		timeToLock := time.Since(now)
		if timeToLock < 5*time.Second {

			t.Fatalf("mutex TTL was not refreshed, lock acquired after %v seconds", timeToLock)
		}

		mutex.Unlock()
		quit2 <- true
	}()

	var (
		q1 bool
		q2 bool
	)

	for !q1 || !q2 {
		select {
		case <-quit1:
			q1 = true
		case <-quit2:
			q2 = true
		case <-errchan:
			t.Fatal("could not acquire lock, is etcd running?")
		}
	}
}

func TestUnlockReleased(t *testing.T) {

	//etcd.SetLogger(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))

	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	client.Delete(key, true)

	mutex := NewMutexFromClient(client, key, 0)

	defer func() {
		if msg := recover(); msg == nil {

			t.Fatalf("panic not initiated")
		}
	}()
	mutex.Unlock()
}

func TestUnlockNoKey(t *testing.T) {

	//etcd.SetLogger(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))

	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	client.Delete(key, true)

	mutex := NewMutexFromClient(client, key, 0)

	err := mutex.Lock()
	if err != nil {

		t.Fatal("could not acquire lock, is etcd running?", err)
	}

	client.Delete(key, false)
	time.Sleep(2 * time.Second)
	mutex.Unlock()
}

func _TestUnlockBadIndex(t *testing.T) {

	//etcd.SetLogger(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))

	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	client.Delete(key, true)

	mutex := NewMutexFromClient(client, key, 0)

	err := mutex.Lock()
	if err != nil {

		t.Fatal("could not acquire lock, is etcd running?", err)
	}

	client.Update(key, "locked", 0)
	mutex.Unlock()

	trigger := make(chan bool)
	errchan := make(chan bool)
	go func() {

		err := mutex.Lock()
		if err != nil {

			errchan <- true
		}

		trigger <- true
		mutex.Unlock()
	}()

	tick := time.Tick(time.Second)

	select {
	case <-errchan:
		t.Fatal("could not acquire lock, is etcd running?", err)
	case <-trigger:
		t.Fatalf("managed to get a lock on an out of sync mutex")
		break
	case <-tick:
		// release the blocked goroutine
		client.Delete(key, true)
	}
}

func HammerMutex(m *EtcdMutex, loops int, cdone chan bool, errchan chan error, t *testing.T) {
	log.Printf("starting %d iterations", loops)
	for i := 0; i < loops; i++ {
		err := m.Lock()
		if err != nil {

			errchan <- err
			return
		}

		m.Unlock()
	}
	log.Printf("completed all iterations")
	cdone <- true
}

func TestConcurrentSingleMutex(t *testing.T) {
	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	client.Delete(key, true)

	m := NewMutexFromClient(client, key, 0)
	c := make(chan bool)
	e := make(chan error)
	for i := 0; i < 10; i++ {
		go HammerMutex(m, 100, c, e, t)
	}
	for i := 0; i < 10; i++ {
		select {
		case <-c:
		case err := <-e:
			t.Fatal("could not acquire lock, is etcd running?", err)
		}
	}
}

func TestConcurrentMultipleMutex(t *testing.T) {
	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	client.Delete(key, true)

	c := make(chan bool)
	e := make(chan error)
	for i := 0; i < 10; i++ {
		m := NewMutexFromClient(client, key, 0)
		go HammerMutex(m, 100, c, e, t)
	}
	for i := 0; i < 10; i++ {
		select {
		case <-c:
		case err := <-e:
			t.Fatal("could not acquire lock, is etcd running?", err)
		}
	}
}

func TestMutexPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("unlock of unlocked mutex did not panic")
		}
	}()

	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	client.Delete(key, true)

	mu := NewMutexFromClient(client, key, 0)
	err := mu.Lock()
	if err != nil {

		t.Fatal("could not acquire lock, is etcd running?", err)
	}

	mu.Unlock()
	mu.Unlock()
}

func BenchmarkMutexUncontended(b *testing.B) {
	type PaddedMutex struct {
		*EtcdMutex
		pad [128]uint8
	}

	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	client.Delete(key, true)

	b.RunParallel(func(pb *testing.PB) {
		mu := PaddedMutex{EtcdMutex: NewMutexFromClient(client, key, 0)}

		for pb.Next() {
			err := mu.Lock()
			if err != nil {

				b.Fatal("could not acquire lock, is etcd running?", err)
			}

			mu.Unlock()
		}
	})
}

func benchmarkMutex(b *testing.B, slack, work bool) {
	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	client.Delete(key, true)

	mu := NewMutexFromClient(client, key, 0)
	if slack {
		b.SetParallelism(10)
	}
	b.RunParallel(func(pb *testing.PB) {
		foo := 0
		for pb.Next() {
			err := mu.Lock()
			if err != nil {

				b.Fatal("could not acquire lock, is etcd running?", err)
			}

			mu.Unlock()
			if work {
				for i := 0; i < 100; i++ {
					foo *= 2
					foo /= 2
				}
			}
		}
		_ = foo
	})
}

func BenchmarkMutex(b *testing.B) {
	benchmarkMutex(b, false, false)
}

func BenchmarkMutexSlack(b *testing.B) {
	benchmarkMutex(b, true, false)
}

func BenchmarkMutexWork(b *testing.B) {
	benchmarkMutex(b, false, true)
}

func BenchmarkMutexWorkSlack(b *testing.B) {
	benchmarkMutex(b, true, true)
}
