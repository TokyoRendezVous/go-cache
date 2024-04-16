package DistributedCache

// ByteView 保存不可变的字节视图
type ByteView struct {
	b []byte
}

// Len 返回视图的长度
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 以byte slice的形式返回视图的副本
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String 以字符串的形式返回数据
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
