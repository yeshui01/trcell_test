/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:14
 * @LastEditTime: 2022-09-20 10:20:07
 * @FilePath: \trcell\pkg\pb\pbcmd\user_cmd_data_def.go
 */
package pbcmd

import (
	"net"
	"trcell/pkg/evhub"

	"github.com/gorilla/websocket"
)

type CmdTypeWebsocketMessageData struct {
	WsConn     *websocket.Conn
	WsMsgType  int32
	MsgData    []byte
	HubMsg     *evhub.NetMessage
	RecvTimeMs int64
}
type CmdTypeTcpsocketMessageData struct {
	TcpConn    net.Conn
	HubMsg     *evhub.NetMessage
	RecvTimeMs int64
}
