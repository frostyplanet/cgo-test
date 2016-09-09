package main
/*
#cgo CFLAGS: -I.
#include <foo.h>
#include <cfuncs.h>
#include <stdint.h>
*/
import "C"
import (
	"os"
	"sync"
	"unsafe"
	"encoding/binary"
	//"runtime"
)

type Job struct {
	wg sync.WaitGroup
}

var jobPool *sync.Pool

func init () {
	jobPool = &sync.Pool{
		New: func() interface{} {
			return new(Job)
		},
	}
}

//export JobDoneCallback
func JobDoneCallback(job unsafe.Pointer) {
	j := (*Job)(job)
	//runtime.LockOSThread()
	j.wg.Done()
	//runtime.UnlockOSThread()
}

type Bar struct {
	foo C.foo_t
	input *os.File
	lock sync.Mutex
}

func (b *Bar) Init() {
	res := C.foo_init(&(b.foo))
	if res < 0 {
		panic("Init error")
	}
	b.input = os.NewFile(uintptr(b.foo.fds[1]), "")
}

func (b *Bar) CallSync() {
	C.foo_work()
}

func (b* Bar) CallAsync() {
	j := jobPool.Get().(*Job)
	var buf [8]byte
	var paddr uint64
	//runtime.LockOSThread()
	j.wg.Add(1)
	paddr = uint64(uintptr(unsafe.Pointer(j)))
	binary.LittleEndian.PutUint64(buf[0:8], paddr)
	//println("write", paddr)
	b.lock.Lock()
	_, err := b.input.Write(buf[:])
	b.lock.Unlock()
	if err != nil {
		panic(err)
	}
	j.wg.Wait()
	//runtime.UnlockOSThread()
	jobPool.Put(j)
}
