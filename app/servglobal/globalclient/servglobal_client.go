/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-28 14:53:13
 * @LastEditTime: 2022-10-17 10:28:17
 * @FilePath: \trcell\app\servglobal\globalclient\servglobal_client.go
 */
package globalclient

import (
	"context"
	"trcell/pkg/loghlp"
	"trcell/pkg/pb/pbrpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServGlobalClient struct {
	rpcClient pbrpc.GlobalDataClient
	conn      *grpc.ClientConn
	isOpen    bool
	rpcAddr   string
}

func NewServGlobalClient(servAddr string) *ServGlobalClient {
	return &ServGlobalClient{
		isOpen:  false,
		rpcAddr: servAddr,
	}
}

func (servClient *ServGlobalClient) Connect(servAddr string) error {
	conn, err := grpc.Dial(servAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		loghlp.Errorf("rpcConnect err:%s", err.Error())
		return err
	}
	servClient.conn = conn
	servClient.rpcClient = pbrpc.NewGlobalDataClient(conn)
	servClient.isOpen = true
	return nil
}
func (servClient *ServGlobalClient) Close() {
	if servClient.isOpen {
		servClient.conn.Close()
		servClient.isOpen = false
	}
}

// 接口调用: EchoTest
func (servClient *ServGlobalClient) CallEchoTest(req *pbrpc.EchoReq) (*pbrpc.EchoRep, error) {
	echoRep, repErr := servClient.rpcClient.EchoTest(context.Background(), req)
	if repErr != nil {
		loghlp.Errorf("repError:%s", repErr.Error())
		return nil, repErr
	}
	return echoRep, nil
}
func (servClient *ServGlobalClient) CallGenUid(req *pbrpc.GenUIDReq) (*pbrpc.GenUIDRep, error) {
	rep, repErr := servClient.rpcClient.GenUid(context.Background(), req)
	if repErr != nil {
		loghlp.Errorf("repError:%s", repErr.Error())
		return nil, repErr
	}
	return rep, nil
}
