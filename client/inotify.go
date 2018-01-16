package main

// #cgo LDFLAGS: -ldl
// #include <sys/inotify.h>
import "C"
import (
	"encoding/binary"
	"os"
	"log"
	"sync"
	"bytes"
	"syscall"

	"github.com/jchv/eventgun/evgun"
)

var watchcount = uintptr(1)

type NotifyHandle struct {
	client *evgun.NotifyClient
	files map[uintptr]string
	watches map[string]map[uintptr]struct{}

	w, r *os.File
	mutex sync.Mutex
}

func NewNotifyHandle() *NotifyHandle {
	r, w, err := os.Pipe()
	wd := r.Fd()
	syscall.SetNonblock(int(wd), true)
	if err != nil {
		panic(err)
	}

	addr := os.Getenv("EVENTGUN_ADDR")
	if addr == "" {
		addr = "localhost:4547"
	}

	client, err := evgun.NewNotifyClient(addr)
	if err != nil {
		log.Fatalln("Unable to connect to EventGun server at", addr+":", err)
	}

	handle := &NotifyHandle{
		client: client,
		files: map[uintptr]string{},
		watches: map[string]map[uintptr]struct{}{},
		w: w,
		r: r,
	}

	go func() {
		for event := range client.Events {
			mask := uint32(0)

			if event.Op&evgun.Create == evgun.Create {
				mask |= C.IN_CREATE
			}
			if event.Op&evgun.Rename == evgun.Rename {
				mask |= C.IN_MOVE_SELF
			}
			if event.Op&evgun.Remove == evgun.Remove {
				mask |= C.IN_DELETE
			}
			if event.Op&evgun.Write == evgun.Write {
				mask |= C.IN_MODIFY
			}
			if event.Op&evgun.Chmod == evgun.Chmod {
				mask |= C.IN_ATTRIB
			}

			handle.mutex.Lock()
			buf := bytes.Buffer{}
			log.Println(event.Name)
			log.Println(handle.watches[event.Name])
			for fd := range handle.watches[event.Name] {
				log.Println(fd)
				fn := []byte(event.Name)
				binary.Write(&buf, binary.LittleEndian, int32(fd))
				binary.Write(&buf, binary.LittleEndian, uint32(mask))
				binary.Write(&buf, binary.LittleEndian, uint32(0))
				binary.Write(&buf, binary.LittleEndian, uint32(len(fn)+1))
				buf.Write(fn)
				buf.Write([]byte{0})
			}
			handle.mutex.Unlock()
			w.Write(buf.Bytes())
		}
	}()

	return handle
}

var handles = map[uintptr]*NotifyHandle{}

//export go_inotify_init
func go_inotify_init() C.int {
	log.Printf("inotify_init()\n")

	handle := NewNotifyHandle()
	fd := handle.r.Fd()
	handles[fd] = handle

	return C.int(fd)
}

//export go_inotify_init1
func go_inotify_init1(flags C.int) C.int {
	log.Printf("inotify_init1(%d)\n", flags)

	handle := NewNotifyHandle()
	fd := handle.r.Fd()
	handles[fd] = handle

	return C.int(fd)
}

//export go_inotify_add_watch
func go_inotify_add_watch(fd C.int, pathname *C.char, mask uint32) C.int {
	log.Printf("inotify_add_watch(%d, %s, %d)\n", fd, C.GoString(pathname), mask)

	handle, ok := handles[uintptr(fd)]
	if !ok {
		return -1
	}

	fn := C.GoString(pathname)
	
	handle.mutex.Lock()
	defer handle.mutex.Unlock()

	// Get new watch descriptor
	wd := watchcount
	watchcount++

	// Add to files map
	handle.files[wd] = fn

	// Add to watches map
	wm := handle.watches[fn]
	if (wm == nil) {
		wm = map[uintptr]struct{}{}
		handle.client.AddWatch(fn)
	}
	wm[wd] = struct{}{}
	handle.watches[fn] = wm

	return C.int(wd)
}

//export go_inotify_rm_watch
func go_inotify_rm_watch(fd C.int, cwd C.int) C.int {
	wd := uintptr(cwd)
	log.Printf("inotify_init1(%d, %d)\n", fd, wd)

	handle, ok := handles[uintptr(fd)]
	if !ok {
		return -1
	}

	handle.mutex.Lock()
	defer handle.mutex.Unlock()

	fn := handle.files[wd]
	if fn == "" {
		return -1
	}

	wm := handle.watches[fn]

	delete(wm, wd)

	if len(wm) == 0 {
		delete(handle.files, wd)
		handle.client.RemoveWatch(fn)
	}

	return 0
}

func main() {
}
