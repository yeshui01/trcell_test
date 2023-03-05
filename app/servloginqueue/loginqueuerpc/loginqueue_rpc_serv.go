package loginqueuerpc

import (
	"context"
	"errors"
	"trcell/app/servloginqueue/iservloginqueue"
	"trcell/pkg/pb/pbrpc"
)

type LoginQueueServer struct {
	pbrpc.UnimplementedLoginQueueBackendServer
	loginQueueServ iservloginqueue.IServLoginQueue
}

func NewLoginQueueServer(iservLoginQueue iservloginqueue.IServLoginQueue) *LoginQueueServer {
	return &LoginQueueServer{
		loginQueueServ: iservLoginQueue,
	}
}

func (serv *LoginQueueServer) GetLoginSeqNo(ctx context.Context, req *pbrpc.LoginSeqNoReq) (*pbrpc.LoginSeqNoRep, error) {
	if len(req.AccountName) < 1 {
		return nil, errors.New("account param error!!!!")
	}

	// 玩家当前的seq
	loginQueue := serv.loginQueueServ.GetLoginQueue()
	curUserSeq := loginQueue.GetUserSeq(req.AccountName)
	curQueuingNum := int32(0)
	if curUserSeq == 0 {
		curUserSeq, curQueuingNum = loginQueue.GenLoginSeqNo()
		loginQueue.CacheUserSeq(req.AccountName, curUserSeq)
	}
	rep := &pbrpc.LoginSeqNoRep{
		SeqNo:      curUserSeq,
		CurQueueNo: loginQueue.GetCurLoginSeqNo(),
		QueueNum:   curQueuingNum,
	}

	return rep, nil
}

func (serv *LoginQueueServer) GetFinishSeqNo(ctx context.Context, req *pbrpc.QueryCurLoginFinishNoReq) (*pbrpc.QueryCurLoginFinishNoRep, error) {
	// return nil, status.Errorf(codes.Unimplemented, "method GetFinishSeqNo not implemented")
	loginQueue := serv.loginQueueServ.GetLoginQueue()
	rep := &pbrpc.QueryCurLoginFinishNoRep{
		CurFinishNo: loginQueue.GetFinishNum(),
	}

	return rep, nil
}

func (serv *LoginQueueServer) IncrementLoginFinish(ctx context.Context, req *pbrpc.IncrementLoginFinishReq) (*pbrpc.IncrementLoginFinishRep, error) {
	// return nil, status.Errorf(codes.Unimplemented, "method IncrementLoginFinish not implemented")
	rep := &pbrpc.IncrementLoginFinishRep{}
	loginQueue := serv.loginQueueServ.GetLoginQueue()
	loginQueue.AddFinishNum()
	loginQueue.DelUserSeq(req.AccountName)
	return rep, nil
}

func (serv *LoginQueueServer) GetQueueNum(ctx context.Context, req *pbrpc.GetQueueNumReq) (*pbrpc.GetQueueNumRep, error) {
	loginQueue := serv.loginQueueServ.GetLoginQueue()
	queuingNum := loginQueue.GetFinishNum()
	rep := &pbrpc.GetQueueNumRep{
		QueueNum: queuingNum,
	}
	return rep, nil
}
