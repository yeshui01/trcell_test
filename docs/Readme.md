<!--
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:13
 * @LastEditTime: 2022-09-27 10:12:13
 * @FilePath: \trcell\docs\Readme.md
-->
客户端消息结构
type Header struct {
	Len      uint32         (消息总长度:消息头+Data的总长度)
	MsgClass uint8          
	MsgType  uint8         
}
type ClientMsg struct {
    Head Header
    Data []byte
}

Len(uint32)+MsgClass(uint8)+MsgType(uint8)+Data
Len表示总长度