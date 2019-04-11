package stats

// OrderlyKey 排序工具
type OrderlyKey []int64

// Len OrderlyKey 长度
func (o OrderlyKey) Len() int {
	return len(o)
}

// Swap 交换
func (o OrderlyKey) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

// Less 对比
func (o OrderlyKey) Less(i, j int) bool {
	return o[i] < o[j]
}
