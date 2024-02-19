package proxy

import (
	"context"
	"net"
	"sync/atomic"
	"time"

	"github.com/d4yvid/runet-proxy/logger"
	"github.com/d4yvid/runet-proxy/raknet"
)

type ProxyConnection struct {
	id          int
	linkRunning atomic.Bool

	client *raknet.Conn
	server *raknet.Conn

	proxyServer *ProxyServer

	clientProtocol uint32
}

func (conn *ProxyConnection) startLink() {
	var err error
	var netErr *net.OpError

	maxReadTime, err := time.ParseDuration("10ms")

	if err != nil {
		conn.proxyServer.logger.Log(logger.Error, "Connection link", "couldn't parse duration: %s", err)
		return
	}

	conn.linkRunning.Store(true)

	for conn.linkRunning.Load() {
		conn.client.SetReadDeadline(time.Now().Add(maxReadTime))
		conn.server.SetReadDeadline(time.Now().Add(maxReadTime))

		clientPacket, err := conn.client.ReadPacket()

		if err != nil && (err.(*net.OpError).Err != context.DeadlineExceeded) {
			netErr = err.(*net.OpError)

			conn.proxyServer.logger.Log(logger.Error, "Client -> Server", "an error occurred while reading: %s", netErr.Error())

			conn.Close()
			break
		}

		serverPacket, err := conn.server.ReadPacket()

		if err != nil && (err.(*net.OpError).Err != context.DeadlineExceeded) {
			netErr = err.(*net.OpError)

			conn.proxyServer.logger.Log(logger.Error, "Server -> Client", "an error occurred while reading: %s", netErr.Error())

			conn.Close()
			break
		}

		if serverPacket != nil {
			_, err := conn.client.Write(serverPacket)

			if err != nil {
				conn.proxyServer.logger.Log(logger.Error, "Server -> Client.WRITE", "an error occurred while writing: %s", err)
				conn.Close()
				break
			}
		}

		if clientPacket != nil {
			_, err := conn.server.Write(clientPacket)

			if err != nil {
				conn.proxyServer.logger.Log(logger.Error, "Client -> Server.WRITE", "an error occurred while writing: %s", err)
				conn.Close()
				break
			}
		}
	}

	conn.linkRunning.Store(false)
}

func (conn *ProxyConnection) StartLink() {
	go conn.startLink()
}

func (conn *ProxyConnection) Close() error {
	var err error

	err = conn.server.Close()

	if err != nil {
		return err
	}

	err = conn.client.Close()

	if err != nil {
		return err
	}

	err = conn.proxyServer.RemoveConnection(conn)

	return err
}
