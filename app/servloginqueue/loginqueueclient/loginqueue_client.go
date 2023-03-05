package loginqueueclient

import (
	"context"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServLoginQueueClient struct {
	rpcClient pbrpc.LoginQueueBackendClient
	conn      *grpc.ClientConn
	isOpen    bool
	rpcAddr   string
}

func NewServLoginQueueClient(servAddr string) *ServLoginQueueClient {
	return &ServLoginQueueClient{
		isOpen:  false,
		rpcAddr: servAddr,
	}
}

func (servClient *ServLoginQueueClient) Connect() error {
	conn, err := grpc.Dial(servClient.rpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		loghlp.Errorf("rpcConnect err:%s", err.Error())
		return err
	}
	servClient.conn = conn
	servClient.rpcClient = pbrpc.NewLoginQueueBackendClient(conn)
	servClient.isOpen = true
	return nil
}
func (servClient *ServLoginQueueClient) Close() {
	if servClient.isOpen {
		servClient.conn.Close()
		servClient.isOpen = false
	}
}

// 接口调用
func (servClient *ServLoginQueueClient) CallGetLoginSeqNo(req *pbrpc.LoginSeqNoReq) (*pbrpc.LoginSeqNoRep, error) {
	echoRep, repErr := servClient.rpcClient.GetLoginSeqNo(context.Background(), req)
	if repErr != nil {
		loghlp.Errorf("repError:%s", repErr.Error())
		return nil, repErr
	}
	return echoRep, nil
}

func (servClient *ServLoginQueueClient) CallGetFinishSeqNo(req *pbrpc.QueryCurLoginFinishNoReq) (*pbrpc.QueryCurLoginFinishNoRep, error) {
	finishNumRep, repErr := servClient.rpcClient.GetFinishSeqNo(context.Background(), req)
	if repErr != nil {
		loghlp.Errorf("repError:%s", repErr.Error())
		return nil, repErr
	}
	return finishNumRep, nil
}

func (servClient *ServLoginQueueClient) CallIncrementLoginFinish(req *pbrpc.IncrementLoginFinishReq) (*pbrpc.IncrementLoginFinishRep, error) {
	finishNumRep, repErr := servClient.rpcClient.IncrementLoginFinish(context.Background(), req)
	if repErr != nil {
		loghlp.Errorf("repError:%s", repErr.Error())
		return nil, repErr
	}
	return finishNumRep, nil
}
