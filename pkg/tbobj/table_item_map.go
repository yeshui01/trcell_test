/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-09 17:10:59
 * @LastEditTime: 2022-10-09 17:29:33
 * @FilePath: \trcell\pkg\tbobj\table_item_map.go
 */
package tbobj

// Map
type ItemMap[K int8 | uint8 | int32 | uint32 | int64 | uint64 | string, V ITableItem] struct {
	elemList map[K]V
}

/**
 * @description: 设置map元素
 * @param {K} key: 元素的键值
 * @param {V} value: 元素值
 * @return {*}
 */
func (obj *ItemMap[K, V]) Set(key K, value V) {
	if obj.elemList == nil {
		obj.elemList = make(map[K]V)
	}
	if _, ok := obj.elemList[key]; ok {
		obj.elemList[key] = value
		value.SetDbStatus(DbStatusUpdate)
	} else {
		obj.elemList[key] = value
		value.SetDbStatus(DbStatusInsert)
	}
}

/**
 * @description: 根据key获取对应的value
 * @param {K} key
 * @param {V} defaultValue: 当key不存在的时候,返回这个默认的值
 * @return {*}
 */
func (obj *ItemMap[K, V]) Get(key K, defaultValue V) V {
	if obj.elemList == nil {
		return defaultValue
	}
	if v, ok := obj.elemList[key]; ok {
		return v
	}
	return defaultValue
}

/**
 * @description: 删除一个元素
 * @param {K} key: 要删除的元素的键
 * @return {*}
 */
func (obj *ItemMap[K, V]) Erase(key K) {
	tbItem, ok := obj.elemList[key]
	if !ok {
		return
	}
	tbItem.SetDbStatus(DbStatusDelete)
	//delete(obj.elemList, key)
}

// /**
//  * @description: 清理Map
//  * @return {*}
//  */
// func (obj *ItemMap[K, V]) Clear() {
// 	for k := range obj.elemList {
// 		delete(obj.elemList, k)
// 	}
// }

/**
 * @description: 获取当前容器元素个数
 * @return {*}
 */
func (obj *ItemMap[K, V]) Size() int32 {
	var count int32 = 0
	for _, v := range obj.elemList {
		if v.GetDbStatus() != DbStatusDelete {
			count++
		}
	}
	return count
	// return int32(len(obj.elemList))
}

/**
 * @description: 遍历Map
 * @param visitor: 访问器函数
 * @return {*}
 */
func (obj *ItemMap[K, V]) ForEach(visitor func(key K, val V)) {
	for k, v := range obj.elemList {
		if v.GetDbStatus() != DbStatusDelete {
			visitor(k, v)
		}
	}
}

/**
 * @description: 遍历Map,条件不成立时可以break终止
 * @param visitor: 访问器函数,返回true时表示继续遍历,否则break
 * @return {*}
 */
func (obj *ItemMap[K, V]) ForEachWithBreak(visitor func(key K, val V) bool) {
	for k, v := range obj.elemList {
		if v.GetDbStatus() != DbStatusDelete {
			if !visitor(k, v) {
				break
			}
		}
	}
}

/**
 * @description: 创建一个Map对象
 * @return {*}
 */
func NewItemMap[K int8 | uint8 | int32 | uint32 | int64 | uint64 | string, V ITableItem]() *ItemMap[K, V] {
	return &ItemMap[K, V]{
		elemList: make(map[K]V),
	}
}
