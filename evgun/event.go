package evgun

// Op describes a set of file operations.
type Op uint32

// These are the generalized file operations that can trigger a notification.
const (
	Create Op = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

// Event represents a single file system notification.
type Event struct {
	Name string
	Op   Op
}
