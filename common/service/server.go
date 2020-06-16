package service

import (
	"bufio"
	"bytescheme/common/log"
	"context"
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
	// ReadDeadline is the read deadline for messages
	ReadDeadline = time.Millisecond * 500
	// ConnectionRefreshInterval is the interval for force refresh client connection
	ConnectionRefreshInterval = time.Hour * 1
)

// EventServer is the TCP event server
type EventServer struct {
	host              string
	port              int
	onConnectCallback func(clientID int) error
	listener          net.Listener
	rwLock            *sync.RWMutex
	connInfoMap       map[int]*ConnectionInfo
	ctx               context.Context
	cancel            context.CancelFunc
}

// ConnectionInfo holds the client connection
type ConnectionInfo struct {
	clientID int
	version  int64
	conn     net.Conn
	lock     *sync.Mutex
}

// NewEventServer creates an instance of EventServer
func NewEventServer(host string, port int, onConnectCallback func(clientID int) error) (*EventServer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	eventServer := &EventServer{
		host:              host,
		port:              port,
		onConnectCallback: onConnectCallback,
		rwLock:            &sync.RWMutex{},
		connInfoMap:       map[int]*ConnectionInfo{},
		ctx:               ctx,
		cancel:            cancel,
	}
	err := eventServer.start()
	return eventServer, err
}

// Start starts the TCP server
func (server *EventServer) start() error {
	log.Infof("Starting TCP server...")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.host, server.port))
	if err != nil {
		return err
	}
	server.listener = listener
	go server.clientCleaner()
	go func() {
		for {
			select {
			case <-server.ctx.Done():
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					if server.IsClosed() {
						break
					}
					log.Errorf("Error in connection. Error: %s", err.Error())
					continue
				}
				go server.processClient(conn)
			}
		}
	}()
	return nil
}

// Close closes the event server
func (server *EventServer) Close() error {
	server.cancel()
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
	select {
	case <-server.ctx.Done():
		return true
	default:
		return false
	}
}

// clientCleaner removes stale connections
func (server *EventServer) clientCleaner() {
	fn := func() {
		server.rwLock.Lock()
		defer server.rwLock.Unlock()
		now := time.Now()
		for clientID, connInfo := range server.connInfoMap {
			connTime := time.Unix(0, connInfo.version)
			if now.Sub(connTime) > ConnectionRefreshInterval {
				log.Warnf("Declaring stale connection for client %d", clientID)
				connInfo.conn.Close()
				// Safe to delete
				delete(server.connInfoMap, clientID)
			}
		}
	}

	ticker := time.NewTicker(ConnectionRefreshInterval)
	for {
		select {
		case <-server.ctx.Done():
			return
		case <-ticker.C:
			fn()
		}
	}
}

func (server *EventServer) processClient(conn net.Conn) {
	// No timeout
	conn.SetReadDeadline(time.Time{})
	// will listen for message to process ending in newline (\n)
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		conn.Close()
		log.Errorf("Error in reading data from remote address %s. Error: %s", conn.RemoteAddr(), err.Error())
		return
	}
	log.Infof("Message Received: %s", string(message))
	clientID, err := strconv.Atoi(strings.TrimSpace(message))
	if err != nil {
		conn.Close()
		log.Errorf("Malformed message %s from remote address %s. Error: %s", message, conn.RemoteAddr(), err.Error())
		return
	}
	if server.onConnectCallback != nil {
		err = server.onConnectCallback(clientID)
		if err != nil {
			conn.Close()
			log.Errorf("Error in connection callback for client %d. Error: %s", clientID, err.Error())
			return
		}
	}
	server.rwLock.Lock()
	defer server.rwLock.Unlock()
	oldClientInfo, ok := server.connInfoMap[clientID]
	if ok && oldClientInfo != nil {
		oldClientInfo.conn.Close()
	}
	server.connInfoMap[clientID] = &ConnectionInfo{
		clientID: clientID,
		conn:     conn,
		version:  time.Now().UnixNano(),
		lock:     &sync.Mutex{},
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

	sender := func(connInfo *ConnectionInfo, message string) (string, error) {
		connInfo.lock.Lock()
		defer connInfo.lock.Unlock()
		now := time.Now()
		deadline := now.Add(ReadDeadline)
		connInfo.conn.SetReadDeadline(deadline)
		message = strings.TrimSpace(message)
		_, err := connInfo.conn.Write([]byte(message + "\n"))
		if err != nil {
			return "", err
		}
		log.Infof("Message %s sent successfully to client %d", message, clientID)
		message, err = bufio.NewReader(connInfo.conn).ReadString('\n')
		if err != nil {
			return "", err
		}
		message = strings.TrimSpace(message)
		log.Infof("Message %s received from client %d", message, clientID)
		return message, nil
	}

	for i := 0; i < SendRetryLimit; i++ {
		server.rwLock.RLock()
		connInfo, ok := server.connInfoMap[clientID]
		server.rwLock.RUnlock()
		if !ok {
			return "", fmt.Errorf("Unrecognized client %d", clientID)
		}
		message, err := sender(connInfo, message)
		if err != nil {
			log.Errorf("Failed to send message %s to client %d. Error: %s", message, clientID, err.Error())
			staleDeleter(connInfo)
			time.Sleep(RetryDelay)
			continue
		}
		return message, nil
	}
	return "", fmt.Errorf("Failed to send message %s to client %d", message, clientID)
}
