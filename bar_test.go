package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	RUN_CONCURRENT = 5000
	RUN_COUNT = 100
)

func TestSync(t *testing.T) {
	var wg sync.WaitGroup
	b := new(Bar)
	b.Init()
	start_ts := time.Now()
	var count uint64
	for i:=0; i<RUN_CONCURRENT; i++ {
		wg.Add(1)
		go func() {
			for j:=0; j<RUN_COUNT; j++ {
				b.CallSync()
				atomic.AddUint64(&count, uint64(1))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	duration := time.Now().Sub(start_ts)
	println("sync")
	fmt.Printf("%.2f seconds\n", duration.Seconds())
	fmt.Printf("%.2f ops\n", float64(count) / duration.Seconds())
}

func TestAsync(t *testing.T) {
	var wg sync.WaitGroup
	b := new(Bar)
	b.Init()
	start_ts := time.Now()
	var count uint64
	for i:=0; i<RUN_CONCURRENT; i++ {
		wg.Add(1)
		go func() {
			for j:=0; j<RUN_COUNT; j++ {
				b.CallAsync()
				atomic.AddUint64(&count, uint64(1))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	duration := time.Now().Sub(start_ts)
	println("async")
	fmt.Printf("%.2f seconds\n", duration.Seconds())
	fmt.Printf("%.2f ops\n", float64(count) / duration.Seconds())
}
