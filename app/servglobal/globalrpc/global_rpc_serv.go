/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-28 13:43:54
 * @LastEditTime: 2022-10-17 10:32:14
 * @FilePath: \trcell\app\servglobal\globalrpc\global_rpc_serv.go
 */
package globalrpc

import (
	"context"
	"trcell/app/servglobal/iservglobal"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbrpc"
)

type GlobalDataServer struct {
	pbrpc.UnimplementedGlobalDataServer
	globalServ iservglobal.IServGlobal
}

func NewGlobalDataServer(iservGlobal iservglobal.IServGlobal) *GlobalDataServer {
	return &GlobalDataServer{
		globalServ: iservGlobal,
	}
}

func (dataServ *GlobalDataServer) EchoTest(ctx context.Context, req *pbrpc.EchoReq) (*pbrpc.EchoRep, error) {
	loghlp.Debugf("serv recv EchoTest:%s", req.SendText)
	return &pbrpc.EchoRep{
		SendText: req.SendText,
	}, nil
}

func (dataServ *GlobalDataServer) GenUid(ctx context.Context, req *pbrpc.GenUIDReq) (*pbrpc.GenUIDRep, error) {
	loghlp.Debugf("serv recv GenUid:%d", req.Num)
	rep := &pbrpc.GenUIDRep{}
	if req.Num < 1 {
		req.Num = 1
	}
	rep.UIDs = make([]int64, req.Num)
	for i := int32(0); i < req.Num; i++ {
		rep.UIDs[i] = dataServ.globalServ.GetIDGenerator().GetID()
	}

	return rep, nil
}
