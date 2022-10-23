/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-19 15:59:04
 * @LastEditTime: 2022-09-27 16:28:17
 * @FilePath: \trcell\app\servgate\servgatemain\servgate_connection.go
 */
package servgatemain

import (
	"io"
	"net"
	"strings"
	"time"
	"trcell/pkg/evhub"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbclient"
	"trcell/pkg/pb/pbcmd"
	"trcell/pkg/pb/pbframe"
	"trcell/pkg/protocol"
	"trcell/pkg/sconst"
	"trcell/pkg/timeutil"
	"trcell/pkg/trframe"

	"google.golang.org/protobuf/proto"
)

type HGateConnction struct {
	UserID        int64
	UserInfo      *HGateUser
	TcpPeerConn   net.Conn
	LastHeartTime int64 // 最近心跳时间
	sendMsgChan   chan *evhub.NetMessage
	Closed        bool
}

func (hgc *HGateConnction) SendMsg(msg *evhub.NetMessage) bool {
	if hgc.Closed {
		return false
	}
	hgc.sendMsgChan <- msg
	return true
}

func (hgc *HGateConnction) runRead() {
	headBytes := make([]byte, sconst.ClientMsgHeadSize)
	for {
		hgc.TcpPeerConn.SetReadDeadline(time.Now().Add(time.Second * 40))
		_, err := io.ReadFull(hgc.TcpPeerConn, headBytes)
		if err != nil {
			var errInfo string = err.Error()
			if strings.Contains(errInfo, "closed by the remote host") {
				loghlp.Warnf("remote client(%s)[%d] closed!!!", hgc.TcpPeerConn.RemoteAddr(), hgc.UserID)
			} else {
				loghlp.Errorf("tcpRecv error:%s", err.Error())
				// if err.Error() == "EOF" && nRead == 0 {
				// 	loghlp.Warnf("recv zero bytes")
				// 	time.Sleep(time.Millisecond * 1000)
				// 	continue
				// }
			}
			break
		}
		var clientHeader ClientMsgHeader
		clientHeader.Decode(headBytes)
		// 消息体
		var clientMsg ClientMessage
		clientMsg.Head = &clientHeader
		bodyLen := clientHeader.Len - ClientHeaderSize
		if bodyLen >= 1024*16 {
			loghlp.Warnf("recv client msg,but body is too long,size(%d)", bodyLen)
		}
		if bodyLen >= 1024*1024*2 {
			loghlp.Warnf("recv client msg,but body is too long >= 1024*1024*2 ,size(%d)", bodyLen)
			break
		}
		clientMsg.Data = make([]byte, bodyLen)
		if bodyLen > 0 {
			_, err = io.ReadFull(hgc.TcpPeerConn, clientMsg.Data)
			if err != nil {
				var errInfo string = err.Error()
				if strings.Contains(errInfo, "closed by the remote host") {
					loghlp.Warnf("remote client(%s)[%d] closed!!!", hgc.TcpPeerConn.RemoteAddr(), hgc.UserID)
				} else {
					loghlp.Errorf("tcpRecv body error:%s", err.Error())
				}
				break
			}
		}
		// 特殊hook一下(ECMsgClassBase,ECMsgBaseWrapOpt)消息
		if clientMsg.Head.MainType == protocol.ECMsgClassBase && clientMsg.Head.SubType == protocol.ECMsgBaseWrapOpt {
			// 消息转换
			wantMsg := &pbclient.ECMsgBaseWrapOptReq{}
			errConvert := proto.Unmarshal(clientMsg.Data, wantMsg)
			if errConvert == nil {
				clientMsg.Head.MainType = uint8(wantMsg.MsgClass)
				clientMsg.Head.SubType = uint8(wantMsg.MsgType)
				clientMsg.Head.ID = uint32(wantMsg.ReqID)
				clientMsg.Data = wantMsg.Data
				clientMsg.Head.Len = ClientHeaderSize + uint32(len(clientMsg.Data))
			}
		}
		msgCmdData := &pbcmd.CmdTypeTcpsocketMessageData{
			TcpConn:    hgc.TcpPeerConn,
			RecvTimeMs: timeutil.NowTimeMs(),
		}
		hubMsg := evhub.MakeEmptyMessage()
		hubMsg.Head.Len = uint32(len(clientMsg.Data)) + ClientHeaderSize
		hubMsg.Head.MsgClass = uint16(clientMsg.Head.MainType)
		hubMsg.Head.MsgType = uint16(clientMsg.Head.SubType)
		if clientMsg.Head.ID > 0 {
			hubMsg.Head.HasSecond = 1
			hubMsg.SecondHead = &evhub.NetMsgSecondHead{
				ReqID: uint64(clientMsg.Head.ID),
			}
		}
		hubMsg.Data = clientMsg.Data // 这里直接使用这个data
		msgCmdData.HubMsg = hubMsg
		trframe.PostUserCommand(protocol.CellCmdClassTcpsocket,
			protocol.CmdTypeTcpsocketMessage,
			msgCmdData)
	}
	loghlp.Debugf("HGateConnection[%d](%s) exitRead", hgc.UserID, hgc.TcpPeerConn.RemoteAddr().String())
	trframe.PostUserCommand(protocol.CellCmdClassTcpsocket, protocol.CmdTypeTcpsocketClosed, hgc.TcpPeerConn)
}
func (hgc *HGateConnction) runWrite() {
	for sMsg := range hgc.sendMsgChan {
		if sMsg == nil {
			break
		}
		clientMsg := MakeClientTcpMessage(uint8(sMsg.Head.MsgClass),
			uint8(sMsg.Head.MsgType),
			0, sMsg.Head.Result)
		clientMsg.Data = sMsg.Data
		// // --------- temp test code begin -----------
		// testData := make([]byte, 10240)
		// clientMsg.Data = append(clientMsg.Data, testData...)
		// // --------- temp test code end -------------
		var convertErrorMsg bool = false
		if sMsg.Head.HasSecond > 0 && sMsg.SecondHead != nil {
			clientMsg.Head.ID = uint32(sMsg.SecondHead.ReqID)
			// 进行消息转换
			if clientMsg.Head.ID > 0 {
				loghlp.Debugf("gatesend, trans convert wrap message(%d_%d)->(1_255),ack_id(%d)",
					sMsg.Head.MsgClass,
					sMsg.Head.MsgType,
					clientMsg.Head.ID)

				wrapMsg := &pbclient.ECMsgBaseWrapOptRsp{}
				wrapMsg.MsgClass = int32(sMsg.Head.MsgClass)
				wrapMsg.MsgType = int32(sMsg.Head.MsgType)
				wrapMsg.AckID = int32(sMsg.SecondHead.ReqID)
				wrapMsg.Data = sMsg.Data
				wrapMsg.ErrCode = int32(sMsg.Head.Result)
				wrapData, errWrap := proto.Marshal(wrapMsg)
				if errWrap == nil {
					clientMsg.Head.MainType = protocol.ECMsgClassBase
					clientMsg.Head.SubType = protocol.ECMsgBaseWrapOpt
					clientMsg.Data = wrapData
				}
			} else if sMsg.Head.Result != protocol.ECodeSuccess {
				convertErrorMsg = true
			}
		} else if sMsg.Head.Result != protocol.ECodeSuccess {
			convertErrorMsg = true
		}
		if convertErrorMsg {
			// 转换成通用错误码消息
			clientMsg.Head.MainType = protocol.ECMsgClassBase
			clientMsg.Head.SubType = protocol.ECMsgBasePushErrorOpt
			convertMsg := &pbclient.ECMsgBasePushErrorOptNotify{
				MsgClass:  int32(sMsg.Head.MsgClass),
				MsgType:   int32(sMsg.Head.MsgType),
				ErrCode:   int32(sMsg.Head.Result),
				ErrParams: make([]string, 0),
			}
			paramsData := &pbframe.SErrorParams{}
			if proto.Unmarshal(sMsg.Data, paramsData) == nil {
				convertMsg.ErrParams = paramsData.ParamList // 这里直接赋值使用即可,不用重新拷贝
			}
			clientMsg.Data, _ = proto.Marshal(convertMsg)
		}

		sendData := clientMsg.Encode()
		hgc.TcpPeerConn.SetWriteDeadline(time.Now().Add(time.Second * 40))
		// 特殊处理测试消息
		if sMsg.Head.MsgClass == 0 && sMsg.Head.MsgType == 1 {
			sendN, errSend := hgc.TcpPeerConn.Write(sendData)
			if sendN != len(sendData) {
				loghlp.Errorf("sendN(%d) != len(sendData)(%d)", sendN, len(sendData))
			}
			if errSend != nil {
				loghlp.Errorf("hgateconnection send error:%s", errSend.Error())
			}
		} else if sMsg.Head.MsgClass > 0 && sMsg.Head.MsgType > 0 {
			lenData := len(sendData)
			for lenData > 0 {
				var sendN int
				var errSend error
				if lenData <= sconst.TcpPackUnit {
					sendN, errSend = hgc.TcpPeerConn.Write(sendData)
					//loghlp.Debugf("sendclient packet1,sendN(%d)totalN(%d)", sendN, lenData)
				} else {
					sendN, errSend = hgc.TcpPeerConn.Write(sendData[0:sconst.TcpPackUnit])
					//loghlp.Debugf("sendclient packet2,sendN(%d)totalN(%d)", sendN, lenData)
				}

				if errSend == nil {
					if sendN > 0 {
						lenData = lenData - sendN
						if sendN != len(sendData) {
							sendData = sendData[sendN:]
							loghlp.Warnf("sendToClient sendN(%d) != len(sendData)(%d)", sendN, len(sendData))
						}
						if lenData < 1 {
							//loghlp.Debugf("sendclient packet finish")
							break
						}
					}
				} else {
					loghlp.Errorf("sendToClient hgateconnection send error2:%s", errSend.Error())
					break
				}
				time.Sleep(time.Millisecond * 1)
			}
		} else if sMsg.Head != nil && sMsg.SecondHead != nil && sMsg.Head.MsgClass == 0 && sMsg.Head.MsgType == 0 && sMsg.Head.HasSecond == 1 {
			// 关闭
			loghlp.Warnf("server send close message to client, role(%d)", sMsg.SecondHead.ID)
			trframe.PostUserCommand(protocol.CellCmdClassTcpsocket, protocol.CmdTypeTcpsocketClosed, hgc.TcpPeerConn)
			//break
		}
	}
	loghlp.Debugf("HGateConnection[%d](%s) exitWrite", hgc.UserID, hgc.TcpPeerConn.RemoteAddr().String())
	hgc.TcpPeerConn.Close()
}
func (hgc *HGateConnction) Start() {
	go hgc.runRead()
	go hgc.runWrite()
}
func (hgc *HGateConnction) Stop() {
	hgc.Closed = true
	close(hgc.sendMsgChan)
}

type HGateClientManager struct {
	connMap map[net.Conn]*HGateConnction
}

func NewHGateClientManager() *HGateClientManager {
	return &HGateClientManager{
		connMap: make(map[net.Conn]*HGateConnction),
	}
}
func (mgr *HGateClientManager) AddConnection(tcpConn net.Conn) {
	hgc := &HGateConnction{
		TcpPeerConn: tcpConn,
		UserID:      0,
		sendMsgChan: make(chan *evhub.NetMessage, 1),
		Closed:      false,
	}
	mgr.connMap[tcpConn] = hgc
	hgc.Start()
}
func (mgr *HGateClientManager) RemoveConnection(tcpConn net.Conn) {
	delete(mgr.connMap, tcpConn)
}

func (mgr *HGateClientManager) GetConnection(tcpConn net.Conn) *HGateConnction {
	if hgc, ok := mgr.connMap[tcpConn]; ok {
		return hgc
	}
	return nil
}
