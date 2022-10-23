/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 15:31:07
 * @LastEditTime: 2022-09-19 15:49:55
 * @FilePath: \trcell\pkg\tcpserver\tcp_server.go
 */
package tcpserver

import (
	"net"
	"sync"
	"time"
	"trcell/pkg/loghlp"

	"github.com/sirupsen/logrus"
)

type TcpConnectedCallback func(tcpConn net.Conn, err error)
type TcpServer struct {
	connCallback TcpConnectedCallback
	tcpListener  net.Listener
	goWg         sync.WaitGroup // 记录开启的线程,防止优雅退出时发生线程逃逸
}

func NewTcpServer() *TcpServer {
	return &TcpServer{}
}

func (serv *TcpServer) runTcpLiten() {
	serv.goWg.Add(1)
	go func(serv *TcpServer) {
		var tempDelay time.Duration
		for {
			conn, err := serv.tcpListener.Accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					if tempDelay == 0 {
						tempDelay = 5 * time.Millisecond
					} else {
						tempDelay *= 2
					}
					if max := 1 * time.Second; tempDelay > max {
						tempDelay = max
					}
					// log.Release("accept error: %v; retrying in %v", err, tempDelay)
					time.Sleep(tempDelay)
					continue
				}
				break
			} else {
				if serv.connCallback != nil {
					serv.connCallback(conn, err)
				}
			}
		}
		serv.goWg.Done()
	}(serv)
}
func (serv *TcpServer) Run(listenAddr string, stopCh chan bool) error {
	loghlp.Infof("tcpserver run:%s", listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	serv.tcpListener = ln
	serv.runTcpLiten()

	<-stopCh
	serv.tcpListener.Close()
	logrus.Info("Server exiting")
	return nil
}

func (serv *TcpServer) SetupConnCallback(cb TcpConnectedCallback) {
	serv.connCallback = cb
}
