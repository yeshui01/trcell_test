/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 14:53:28
 * @LastEditTime: 2022-09-23 10:56:29
 * @FilePath: \trcell\app\servgate\cellserv_gate.go
 */
package servgate

import (
	"trcell/app/servgate/servgatehandler"
	"trcell/app/servgate/servgatemain"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbtools"
	"trcell/pkg/protocol"
	"trcell/pkg/trframe"
	"trcell/pkg/trframe/iframe"

	"google.golang.org/protobuf/proto"
)

type CellServGate struct {
	servGateGlobal *servgatemain.ServGateGlobal
	UserMgr        *servgatemain.HGateUserManager
	ConnMgr        *servgatemain.HGateClientManager
	CmdHandlerMap  map[int64]func(frameCmd *trframe.TRFrameCommand)
}

func NewCellServGate() *CellServGate {
	s := &CellServGate{
		servGateGlobal: servgatemain.NewServGateGlobal(),
		CmdHandlerMap:  make(map[int64]func(frameCmd *trframe.TRFrameCommand)),
	}
	s.UserMgr = servgatemain.NewHGateUserManager()
	s.ConnMgr = servgatemain.NewHGateClientManager()
	servgatehandler.InitServGateObj(s)
	s.RegisterMsgHandler()
	s.InitCmdHandle()
	return s
}

func (serv *CellServGate) GetGateGlobal() *servgatemain.ServGateGlobal {
	return serv.servGateGlobal
}

func (serv *CellServGate) FrameRun(curTimeMs int64) {
	serv.servGateGlobal.FrameUpdate(curTimeMs)
}
func GetCmdHandlerKey(cmdClass int32, cmdType int32) int64 {
	k := int64(cmdClass)<<32 + int64(cmdType)
	return k
}
func (serv *CellServGate) InitCmdHandle() {
	serv.RegisterCmdHandler(protocol.CellCmdClassTcpsocket,
		protocol.CmdTypeWebsocketConnect,
		serv.HandleCmdTcpConnect)
	serv.RegisterCmdHandler(protocol.CellCmdClassTcpsocket,
		protocol.CmdTypeWebsocketClosed,
		serv.HandleCmdTcpClose)
	serv.RegisterCmdHandler(protocol.CellCmdClassTcpsocket,
		protocol.CmdTypeWebsocketMessage,
		serv.HandleCmdTcpMessage)

}

func (serv *CellServGate) RegisterCmdHandler(cmdClass int32, cmdType int32, handle func(frameCmd *trframe.TRFrameCommand)) {
	hKey := GetCmdHandlerKey(cmdClass, cmdType)
	if _, ok := serv.CmdHandlerMap[hKey]; ok {
		loghlp.Errorf("repeated register cmd(%d_%d) handle", cmdClass, cmdType)
	} else {
		serv.CmdHandlerMap[hKey] = handle
	}
}

func (serv *CellServGate) SendTcpClientReplyMessage(okCode int32, cltRep proto.Message, env *iframe.TRRemoteMsgEnv) {
	if okCode != protocol.ECodeSuccess {
		if cltRep != nil {
			var pbMsgName = pbtools.GetFullNameByMessage(cltRep)
			if pbMsgName != "SErrorParams" {
				loghlp.Warnf("SendTcpClientReplyMessage but okCode is not succ, pbdata not match,watch msg(%d_%d)[%s]",
					env.SrcMessage.Head.MsgClass,
					env.SrcMessage.Head.MsgType,
					pbMsgName)
			}
		}
	}
	if cltRep != nil {
		msgData, _ := proto.Marshal(cltRep)
		serv.SendTcpClientReplyMessage2(okCode, msgData, env)
	} else if env.SrcMessage != nil {
		serv.SendTcpClientReplyMessage2(okCode, env.SrcMessage.Data, env)
	} else {
		serv.SendTcpClientReplyMessage2(okCode, make([]byte, 0), env)
	}
}
func (serv *CellServGate) SendTcpClientReplyMessage2(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
	hgsession, ok := env.UserData.(*servgatemain.HGateConnction)
	if !ok {
		hgsession = nil
		loghlp.Errorf("SendTcpClientReplyMessage2,hgsession convert fail")
		// 根据玩家ID查找
		if env.SrcMessage != nil {
			if env.SrcMessage.SecondHead.ID > 0 {
				gateUser := serv.GetUserManager().GetGateUser(env.SrcMessage.SecondHead.ID)
				if gateUser != nil {
					hgsession = gateUser.GetGateConnect()
				}
			}
		}
		return
	}
	if hgsession == nil {
		return
	}
	env.SrcMessage.Data = msgData
	env.SrcMessage.Head.Result = uint16(okCode)
	clientReqID := uint64(0)
	if env.SrcMessage.SecondHead != nil {
		if env.SrcMessage.SecondHead.ReqID > 0 {
			env.SrcMessage.SecondHead.RepID = env.SrcMessage.SecondHead.ReqID
			clientReqID = env.SrcMessage.SecondHead.ReqID
		}
	}
	hgsession.SendMsg(env.SrcMessage)
	loghlp.Debugf("SendTcpClientReplyMessage2(%d_%d)(clientReqNo:%d)",
		env.SrcMessage.Head.MsgClass,
		env.SrcMessage.Head.MsgType,
		clientReqID,
	)
}
func (hg *CellServGate) GetGateConnMgr() *servgatemain.HGateClientManager {
	return hg.ConnMgr
}
