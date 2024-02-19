package proxy

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"

	"github.com/d4yvid/runet-proxy/logger"
	"github.com/d4yvid/runet-proxy/raknet"
)

type ProxyConfiguration struct {
	Motd           string
	BindAddress    string
	TargetAddress  string
	MaxConnections uint
}

type ProxyServer struct {
	configuration      ProxyConfiguration
	listener           *raknet.Listener
	connections        sync.Map
	currentConnections uint
	logger             logger.Logger
	closed             bool
	running            atomic.Bool
}

func CreateProxyServer(configuration ProxyConfiguration, logger logger.Logger) (*ProxyServer, error) {
	listener, err := raknet.Listen(configuration.BindAddress)

	if err != nil {
		return nil, err
	}

	return &ProxyServer{
		listener:           listener,
		configuration:      configuration,
		running:            atomic.Bool{},
		connections:        sync.Map{},
		currentConnections: 0,
		logger:             logger,
	}, nil
}

func (server *ProxyServer) createConnectionLink(conn *raknet.Conn) (*ProxyConnection, error) {
	dial, err := raknet.Dial(server.configuration.TargetAddress)

	if err != nil {
		server.logger.Log(logger.Error, "ProxyServer::createConnectionLink", "couldn't create connection link: %s", err.Error())
		return nil, err
	}

	return &ProxyConnection{
		id:          rand.Int(),
		server:      dial,
		client:      conn,
		proxyServer: server,
	}, nil
}

func (server *ProxyServer) updatePongData() {
	format := "MCPE;%s;1;1.0;%d;%d;13253860892328930865;RunetProxy;Survival;0;19132;19133;"
	str := fmt.Sprintf(format, server.configuration.Motd, server.currentConnections, server.configuration.MaxConnections)

	server.listener.PongData([]byte(str))
}

func (server *ProxyServer) start() {
	server.updatePongData()

	server.running.Store(true)

	for server.running.Load() {
		conn, err := server.listener.Accept()

		if err != nil {
			server.logger.Log(logger.Error, "ProxyServer loop", "couldn't accept connection: %s", err.Error())
			continue
		}

		raknetConn := conn.(*raknet.Conn)

		if server.currentConnections >= server.configuration.MaxConnections {
			server.logger.Log(logger.Info, "ProxyServer loop", "the proxy is full, disconnecting %s", conn.RemoteAddr().String())

			// TODO: send a DisconnectPacket (from the game) to indicate that the server is full
			//       or redirect to another proxy instance
			raknetConn.Close()
			continue
		}

		connection, err := server.createConnectionLink(raknetConn)

		if err != nil {
			server.logger.Log(logger.Error, "ProxyServer loop", "couldn't create connection link: %s", err.Error())
			conn.Close()

			continue
		}

		connection.StartLink()

		server.AddConnection(connection)
	}
}

func (server *ProxyServer) AddConnection(conn *ProxyConnection) error {
	server.connections.Store(conn.id, conn)
	server.currentConnections++

	server.updatePongData()
	return nil
}

func (server *ProxyServer) RemoveConnection(conn *ProxyConnection) error {
	_, loaded := server.connections.LoadAndDelete(conn.id)

	if !loaded {
		return errors.New("couldn't find connection in connections")
	}

	server.currentConnections--

	server.updatePongData()
	return nil
}

func (server *ProxyServer) Start(sync bool) error {
	if server.closed {
		return errors.New("the proxy server is closed")
	}

	if server.running.Load() {
		return errors.New("the server is already running")
	}

	if !sync {
		go server.start()
		return nil
	}

	server.start()

	return nil
}

func (server *ProxyServer) Close() {
	server.connections.Range(func(key any, value any) bool {
		connection := value.(*ProxyConnection)

		connection.server.Close()
		connection.client.Close()

		server.connections.Delete(key)

		return true
	})

	server.listener.Close()

	server.listener = nil
	server.closed = true
}

func (server *ProxyServer) Closed() bool {
	return server.closed
}
