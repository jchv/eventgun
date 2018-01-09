package main

// #cgo LDFLAGS: -ldl
import "C"
import "log"

//export go_inotify_init
func go_inotify_init() C.int {
	log.Printf("inotify_init()\n")
	return 0
}

//export go_inotify_init1
func go_inotify_init1(flags C.int) C.int {
	log.Printf("inotify_init1(%d)\n", flags)
	return 0
}

//export go_inotify_add_watch
func go_inotify_add_watch(fd C.int, pathname *C.char, mask uint32) C.int {
	log.Printf("inotify_add_watch(%d, %s, %d)\n", fd, C.GoString(pathname), mask)
	return 0
}

//export go_inotify_rm_watch
func go_inotify_rm_watch(fd C.int, wd C.int) C.int {
	log.Printf("inotify_init1(%d, %d)\n", fd, wd)
	return 0
}

func main() {
	log.Println("Main Go")
}