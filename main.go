package main

import (
	"github.com/d4yvid/runet-proxy/logger"
	"github.com/d4yvid/runet-proxy/proxy"
)

func main() {
	terminalLogger := &logger.TerminalLogger{SaveToFile: false}

	proxyServer, err := proxy.CreateProxyServer(proxy.ProxyConfiguration{
		BindAddress:    "0.0.0.0:19132",
		TargetAddress:  "fire.blazebr.com:26575",
		MaxConnections: 8,
		Motd:           "Forwarding to fire.blazebr.com",
	}, *terminalLogger)

	if err != nil {
		terminalLogger.Log(logger.Error, "Main", "couldn't create proxy server: %s", err.Error())
		return
	}

	proxyServer.Start(true)

	// listen, err := raknet.Listen("0.0.0.0:19132")

	// if err != nil {
	// 	terminalLogger.Log(logger.Error, "Main", "couldn't create listener for port 19132: %s", err)
	// 	return
	// }

	// listen.PongData([]byte("MCPE;Vaom se foderein Proxy;84;0.15.10;800;32000;13253860892328930865;world;Survival;0;19132;19133;"))

	// defer listen.Close()

	// terminalLogger.Log(logger.Info, "Main", "waiting for connections")

	// for {
	// 	conn, err := listen.Accept()

	// 	if err != nil {
	// 		terminalLogger.Log(logger.Error, "Main", "couldn't accept connection: %s", err)
	// 		continue
	// 	}

	// 	terminalLogger.Log(logger.Info, "Main", "new connection from %s", conn.RemoteAddr().String())
	// }
}
