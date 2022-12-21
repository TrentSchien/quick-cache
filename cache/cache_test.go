package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestGetCheckIfKeyAutoDelete(t *testing.T) {

	t.Run("Validate key will get removed during get call if expired", func(t *testing.T) {
		timeLimit := TimeLimit{Seconds: 1}

		InitCache(nil, timeLimit)
		Add("test", "test")
		time.Sleep(3 * time.Second)
		_, hasSomethingPreGet := resp.Load("test")

		if hasSomethingPreGet == false {
			t.Errorf("Key Should Exist Prior to Get call")
		}

		_, hasSomethingPostGet := Get("test")

		if hasSomethingPostGet != false {
			t.Errorf("Key Should Not Exist Post to Get call")
		}
	})

	t.Run("Testing if the cache does not delete items that have not expired", func(t *testing.T) {

		autoClean := TimeLimit{Seconds: 1}
		timeLimit := TimeLimit{Seconds: 60}

		InitCache(&autoClean, timeLimit)
		Add("test", "test")
		time.Sleep(3 * time.Second)

		_, ok := resp.Load("test")
		if ok == false {
			t.Errorf("Key should had been detected")
		}
	})

	t.Run("Testing if the cache will auto clear", func(t *testing.T) {

		autoClean := TimeLimit{Seconds: 1}
		timeLimit := TimeLimit{Seconds: 2}

		InitCache(&autoClean, timeLimit)
		Add("test", "test")
		time.Sleep(3 * time.Second)

		_, ok := resp.Load("test")
		if ok != false {
			t.Errorf("Key should had been deleted")
		}
	})

}

func TestForceCollisions(t *testing.T) {
	//This test will automatically fail if there is an issue. It tries to cause
	//an issue. No issues are a good thing. A normal map could have issues when
	//a map reads and writes at the same time.
	t.Run("This will try to force an issue", func(t *testing.T) {
		InitCache(nil, TimeLimit{Hours: 6})
		readTesting := make(chan bool)
		writeTesting := make(chan bool)
		readWorking := make(chan bool)
		writeWorking := make(chan bool)

		go func() {
			go doManyWrites("testing", "testing", writeTesting)
			go doManyWrites("working", "working", writeWorking)
			go doManyReads("testing", readTesting)
			go doManyReads("working", readWorking)
		}()
		<-readTesting
		<-writeTesting
		<-readWorking
		<-writeWorking
	})
}

func TestAddAndGet(t *testing.T) {
	//This test adds 10,000 items into cache and validates that all are able be gotten
	t.Run("This will validate things get added to cache", func(t *testing.T) {
		InitCache(nil, TimeLimit{Hours: 6})

		for i := 0; i < 10_000; i++ {
			Add(fmt.Sprint("key", i), fmt.Sprint("value", i))
		}
		for i := 0; i < 10_000; i++ {
			value, ok := Get(fmt.Sprint("key", i))
			if !ok {
				t.Errorf("Should Contain Value of %s", fmt.Sprint("value", i))
			}
			if value.(string) != fmt.Sprint("value", i) {
				t.Errorf("Value Should be %s", fmt.Sprint("value", i))
			}

		}
	})
}

func doManyReads(str string, working chan bool) {
	for i := 0; i < 1_000_000; i++ {
		Get(str)
	}
	working <- true
}
func doManyWrites(input, output string, working chan bool) {
	for i := 0; i < 1_000_000; i++ {
		Add(input, output)
	}
	working <- true
}
