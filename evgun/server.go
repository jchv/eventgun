package evgun

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// NotifyServer is the EventGun server.
type NotifyServer struct {
	from, to string
}

// NewNotifyServer creates a new server instance.
func NewNotifyServer(from, to string) *NotifyServer {
	return &NotifyServer{from, to}
}

// PathToHost converts a client path to a host path.
func (n *NotifyServer) PathToHost(path string) string {
	path = strings.Replace(path, n.to, n.from, -1)
	path = strings.Replace(path, "/", "\\", -1) // TODO: only windows
	return path
}

// PathToClient converts a host path to a client path.
func (n *NotifyServer) PathToClient(path string) string {
	path = strings.Replace(path, n.from, n.to, -1)
	path = strings.Replace(path, "\\", "/", -1)
	return path
}

// Listen starts listening for new connections.
func (n *NotifyServer) Listen(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	log.Printf("Now listening on %s (%s)\n", address, l.Addr())

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}

		log.Println("Connection from", c.RemoteAddr())
		go func() {
			err := n.HandleConn(c)
			if err == io.EOF {
				return
			} else if err != nil {
				log.Println("Error in connection:", err)
			}
		}()
	}
}

// HandleConn handles a single connection.
func (n *NotifyServer) HandleConn(c net.Conn) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	go func() {
		for event := range watcher.Events {
			log.Println("Broadcasting", event)

			line := make([]byte, 8)
			binary.BigEndian.PutUint64(line, uint64(event.Op))
			line = append(append(line, n.PathToClient(event.Name)...), 0)
			c.Write(line)
		}
	}()

	r := bufio.NewReader(c)

	for {
		l, err := r.ReadBytes(0)
		if err != nil {
			return err
		}

		// Remove delimiter
		l = l[:len(l)-1]

		path := n.PathToHost(string(l[1:]))
		switch l[0] {
		case 1:
			log.Println("Adding", path)
			watcher.Add(path)
		case 2:
			log.Println("Removing", path)
			watcher.Remove(path)
		case 3:
			return c.Close()
		}
	}
}
