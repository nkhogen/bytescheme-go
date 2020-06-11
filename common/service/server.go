package service

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// SendRetryLimit is the max number of send retries
	SendRetryLimit = 5
	// RetryDelay is the retry back-off time
	RetryDelay = time.Millisecond * 120
)

// EventServer is the TCP event server
type EventServer struct {
	host        string
	port        int
	listener    net.Listener
	rwLock      *sync.RWMutex
	connInfoMap map[int]*ConnectionInfo
}

// ConnectionInfo holds the client connection
type ConnectionInfo struct {
	clientID int
	version  int64
	conn     net.Conn
}

// NewEventServer creates an instance of EventServer
func NewEventServer(host string, port int) (*EventServer, error) {
	eventServer := &EventServer{
		host:        host,
		port:        port,
		rwLock:      &sync.RWMutex{},
		connInfoMap: map[int]*ConnectionInfo{},
	}
	err := eventServer.start()
	return eventServer, err
}

// Start starts the TCP server
func (server *EventServer) start() error {
	fmt.Printf("Starting TCP server...")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.host, server.port))
	if err != nil {
		return err
	}
	server.listener = listener
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if server.IsClosed() {
					break
				}
				fmt.Printf("Error in connection. Error: %s\n", err.Error())
				continue
			}
			go server.processClient(conn)
		}
	}()
	return nil
}

// Close closes the event server
func (server *EventServer) Close() error {
	server.rwLock.Lock()
	defer server.rwLock.Unlock()
	if server.listener != nil {
		for _, connInfo := range server.connInfoMap {
			if connInfo.conn != nil {
				connInfo.conn.Close()
			}
		}
		server.connInfoMap = map[int]*ConnectionInfo{}
		err := server.listener.Close()
		server.listener = nil
		return err
	}
	return nil
}

// IsClosed checks if the server is closed
func (server *EventServer) IsClosed() bool {
	server.rwLock.RLock()
	defer server.rwLock.RUnlock()
	return server.listener == nil
}

func (server *EventServer) processClient(conn net.Conn) {
	// will listen for message to process ending in newline (\n)
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Printf("Error in reading data from remote address %s. Error: %s\n", conn.RemoteAddr(), err.Error())
		return
	}
	fmt.Printf("Message Received: %s\n", string(message))
	clientID, err := strconv.Atoi(strings.TrimSpace(message))
	if err != nil {
		fmt.Printf("Malformed message %s from remote address %s. Error: %s\n", message, conn.RemoteAddr(), err.Error())
		return
	}

	server.rwLock.Lock()
	defer server.rwLock.Unlock()
	server.connInfoMap[clientID] = &ConnectionInfo{
		clientID: clientID,
		conn:     conn,
		version:  time.Now().UnixNano(),
	}
}

// Send sends a message to the client
func (server *EventServer) Send(clientID int, message string) (string, error) {
	staleDeleter := func(connInfo *ConnectionInfo) {
		connInfo.conn.Close()
		server.rwLock.Lock()
		defer server.rwLock.Unlock()
		currConnInfo, ok := server.connInfoMap[clientID]
		if ok {
			if connInfo.version == currConnInfo.version {
				delete(server.connInfoMap, clientID)
			}
		}
	}

	for i := 0; i < SendRetryLimit; i++ {
		server.rwLock.RLock()
		connInfo, ok := server.connInfoMap[clientID]
		server.rwLock.RUnlock()
		if !ok {
			return "", fmt.Errorf("Unrecognized client %d", clientID)
		}
		message = strings.TrimSpace(message)
		_, err := connInfo.conn.Write([]byte(message + "\n"))
		if err == nil {
			fmt.Printf("Message %s sent successfully to client %d\n", message, clientID)
			message, err := bufio.NewReader(connInfo.conn).ReadString('\n')
			if err == nil {
				message = strings.TrimSpace(message)
				fmt.Printf("Message %s received from client %d\n", message, clientID)
				return message, nil
			}

		}
		if err != nil {
			fmt.Println("Error occurred in event server for client %d. Error: %s\n", clientID, err.Error())
		}
		staleDeleter(connInfo)
		time.Sleep(RetryDelay)
	}
	return "", fmt.Errorf("Failed to send message %s to client %d", message, clientID)
}
