package evgun

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
)

// NotifyClient is a client to connect to the EventGun server.
type NotifyClient struct {
	c net.Conn

	Events chan Event
}

// NewNotifyClient creates a new notify client.
func NewNotifyClient(address string) (*NotifyClient, error) {
	c, err := net.Dial("tcp", address)

	if err != nil {
		log.Println("Could not connect to EventGun server at", address+":", err)
		return nil, err
	}

	events := make(chan Event, 8)

	go func() {
		r := bufio.NewReader(c)

		for {
			bop := [8]byte{}
			r.Read(bop[:])

			l, err := r.ReadBytes(0)
			if err == io.EOF {
				return
			} else if err != nil {
				log.Println("Lost connection to EventGun server:", err)
			}

			// Remove delimiter
			l = l[:len(l)-1]

			op := binary.BigEndian.Uint64(bop[:])
			fn := string(l[:len(l)-1])

			events <- Event{
				Op:   Op(op),
				Name: fn,
			}
		}
	}()

	return &NotifyClient{c, events}, nil
}

// AddWatch adds a path to be watched.
func (n *NotifyClient) AddWatch(path string) {
	log.Println("Adding", path)
	n.c.Write(append(append([]byte{1}, path...), 0))
}

// RemoveWatch removes a path from the watch.
func (n *NotifyClient) RemoveWatch(path string) {
	log.Println("Removing", path)
	n.c.Write(append(append([]byte{2}, path...), 0))
}

// Close disconnects from the watch server.
func (n *NotifyClient) Close() error {
	n.c.Write([]byte{3, 0})
	return n.c.Close()
}
