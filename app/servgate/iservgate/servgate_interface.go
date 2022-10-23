/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 14:58:02
 * @LastEditTime: 2022-09-23 16:50:15
 * @FilePath: \trcell\app\servgate\iservgate\servgate_interface.go
 */
package iservgate

import (
	"trcell/app/servgate/servgatemain"
	"trcell/pkg/trframe/iframe"

	"google.golang.org/protobuf/proto"
)

type IServGate interface {
	GetGateGlobal() *servgatemain.ServGateGlobal
	GetUserManager() *servgatemain.HGateUserManager
	GetGateConnMgr() *servgatemain.HGateClientManager
	SendTcpClientReplyMessage(okCode int32, cltRep proto.Message, env *iframe.TRRemoteMsgEnv)
	SendTcpClientReplyMessage2(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv)
}
