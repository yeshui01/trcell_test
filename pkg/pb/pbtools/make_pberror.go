/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-22 13:59:15
 * @LastEditTime: 2022-09-23 11:34:22
 * @FilePath: \trcell\pkg\pb\pbtools\make_pberror.go
 */
package pbtools

import (
	"trcell/pkg/pb/pbclient"
	"trcell/pkg/pb/pbframe"
)

func MakeErrorParams(errInfo string, strParams ...string) *pbframe.SErrorParams {
	paramData := &pbframe.SErrorParams{
		ParamList: make([]string, 0),
	}
	paramData.ErrDesInfo = errInfo
	paramData.ParamList = append(paramData.ParamList, strParams...)
	return paramData
}

func MakeErrorData(msgClass int32, msgType int32, errCode int32, strParams ...string) *pbclient.ECMsgBasePushErrorOptNotify {
	errData := &pbclient.ECMsgBasePushErrorOptNotify{
		ErrCode:   errCode,
		ErrParams: make([]string, 0),
		MsgClass:  msgClass,
		MsgType:   msgType,
	}
	errData.ErrParams = append(errData.ErrParams, strParams...)
	return errData
}
