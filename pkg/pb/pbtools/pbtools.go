/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:14
 * @LastEditTime: 2022-09-23 11:14:08
 * @FilePath: \trcell\pkg\pb\pbtools\pbtools.go
 */
package pbtools

import (
	"trcell/pkg/loghlp"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// 获取消息的全名称
func GetFullNameByMessage(msg proto.Message) string {
	// reflectPB := proto.MessageReflect(msg)
	// descripterPB := reflectPB.Descripter()
	if msg == nil {
		return "nil_protomsg"
	}
	return string(proto.MessageName(msg))
}

// 根据消息的名称获取对象
func GetNewMessageObjByName(messageName string) proto.Message {
	pbMsgType, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(messageName))
	if err != nil {
		loghlp.Errorf("not find pbmessage name:%s", messageName)
		return nil
	}
	newMsg := pbMsgType.New().Interface()
	return proto.Message(newMsg)
}
