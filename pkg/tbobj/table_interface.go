/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-26 14:18:49
 * @LastEditTime: 2022-10-09 17:13:00
 * @FilePath: \trcell\pkg\tbobj\table_interface.go
 */
package tbobj

type ITableBlobField interface {
	ToBytes() []byte
	FromBytes(binaryData []byte)
}

type ITableItem interface {
	SetDbStatus(s int8)
	GetDbStatus() int8
	ClearDbStatus()
}
