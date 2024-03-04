package cache

import (
	"testing"
	"time"
	"fmt"
)

func TestAddCacheKey(t *testing.T) {
	const interval = 2 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "www.example.com",
			val: []byte("somedata"),
		},
		{
			key: "www.google.com",
			val: []byte("lots of other data"),
		},

	}

	for index, testCase := range cases { 
		t.Run(fmt.Sprintf("Running test case %v: ", index), func(*testing.T) {
			cache := NewCache(interval)
			
			cache.Add(testCase.key, testCase.val)
			data, ok := cache.Get(testCase.key)

			if !ok {
				t.Errorf("Unable to find %v key in cache", testCase.key)
				return
			}

			if string(data) != string(testCase.val) {
				t.Errorf("Key data did not match. Got %v wanted %v", data, testCase.val)
				return
			}
		})
	}
}

func TestCacheReap(t *testing.T) {
	const reaper = 100 * time.Millisecond
	const wait = 150 * time.Millisecond 

	cache := NewCache(reaper)

	cache.Add("www.google.com", []byte("some data"))

	_, ok := cache.Get("www.google.com")
	
	if !ok {
		t.Errorf("Key not found in cache")
		return
	}

	time.Sleep(wait)
	
	_, ok = cache.Get("www.gooogle.com")

	if ok {
		t.Errorf("Expected key to be removed")
		return
	}

}
