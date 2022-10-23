/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-20 16:14:38
 * @LastEditTime: 2022-09-26 15:09:26
 * @FilePath: \trcell\pkg\ustd\ustd_container.go
 */
package ustd

// 切片容器
//type innerVec[T any] [
// map容器
//type innerMap[T int8 | uint8 | int32 | uint32 | int64 | uint64 | string, V any] map[T
type Vector[T any] struct {
	dataSlice []T
}

/**
 * @description: 创建一个vector对象
 * @return {*}
 */
func NewVector[T any]() *Vector[T] {
	return &Vector[T]{
		dataSlice: make([]T, 0),
	}
}

/**
 * @description: 容器大小
 * @return {*}
 */
func (obj *Vector[T]) Size() int32 {
	return int32(len(obj.dataSlice))
}

/**
 * @description: 加入元素
 * @param {T} elem
 * @return {*}
 */
func (obj *Vector[T]) PushBack(elem ...T) {
	obj.dataSlice = append(obj.dataSlice, elem...)
}

/**
 * @description: 设置元素
 * @param {int32} idx:数组下标
 * @param {T} elem:新的元素
 * @return {*}
 */
func (obj *Vector[T]) Set(idx int32, elem T) {
	if idx < 0 || idx >= obj.Size() {
		return
	}
	obj.dataSlice[idx] = elem
}

/**
 * @description: 获取元素值
 * @param {int32} idx:数组下表
 * @param {T} defaultVal:找不到元素时的默认值
 * @return {*}
 */
func (obj *Vector[T]) Get(idx int32, defaultVal T) T {
	if idx < 0 || idx >= obj.Size() {
		return defaultVal
	}
	return obj.dataSlice[idx]
}

/**
 * @description: 清理数据
 * @return {*}
 */
func (obj *Vector[T]) Clear() {
	obj.dataSlice = nil
}

/**
 * @description: 遍历数组
 * @param visitor: 访问器函数
 * @return {*}
 */
func (obj *Vector[T]) ForEach(visitor func(idx int32, e T)) {
	if obj.dataSlice != nil {
		for i, v := range obj.dataSlice {
			visitor(int32(i), v)
		}
	}
}

/**
 * @description: 遍历数组
 * @param visitor: 访问器函数,返回true时表示继续,否则break终止遍历
 * @return {*}
 */
func (obj *Vector[T]) ForEachWithBreak(visitor func(idx int32, e T) bool) {
	if obj.dataSlice != nil {
		for i, v := range obj.dataSlice {
			if !visitor(int32(i), v) {
				break
			}
		}
	}
}

// Map
type Map[K int8 | uint8 | int32 | uint32 | int64 | uint64 | string, V any] struct {
	elemList map[K]V
}

/**
 * @description: 设置map元素
 * @param {K} key: 元素的键值
 * @param {V} value: 元素值
 * @return {*}
 */
func (obj *Map[K, V]) Set(key K, value V) {
	if obj.elemList == nil {
		obj.elemList = make(map[K]V)
	}
	obj.elemList[key] = value
}

/**
 * @description: 创建一个Map对象
 * @return {*}
 */
func NewMap[K int8 | uint8 | int32 | uint32 | int64 | uint64 | string, V any]() *Map[K, V] {
	return &Map[K, V]{
		elemList: make(map[K]V),
	}
}

/**
 * @description: 根据key获取对应的value
 * @param {K} key
 * @param {V} defaultValue: 当key不存在的时候,返回这个默认的值
 * @return {*}
 */
func (obj *Map[K, V]) Get(key K, defaultValue V) V {
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
func (obj *Map[K, V]) Erase(key K) {
	delete(obj.elemList, key)
}

/**
 * @description: 清理Map
 * @return {*}
 */
func (obj *Map[K, V]) Clear() {
	for k := range obj.elemList {
		delete(obj.elemList, k)
	}
}

/**
 * @description: 获取当前容器元素个数
 * @return {*}
 */
func (obj *Map[K, V]) Size() int32 {
	return int32(len(obj.elemList))
}

/**
 * @description: 遍历Map
 * @param visitor: 访问器函数
 * @return {*}
 */
func (obj *Map[K, V]) ForEach(visitor func(key K, val V)) {
	for k, v := range obj.elemList {
		visitor(k, v)
	}
}

/**
 * @description: 遍历Map,条件不成立时可以break终止
 * @param visitor: 访问器函数,返回true时表示继续遍历,否则break
 * @return {*}
 */
func (obj *Map[K, V]) ForEachWithBreak(visitor func(key K, val V) bool) {
	for k, v := range obj.elemList {
		if !visitor(k, v) {
			break
		}
	}
}
