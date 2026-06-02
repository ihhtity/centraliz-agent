package utils

// StringPtr 返回字符串指针
func StringPtr(s string) *string {
	return &s
}

// Int32Ptr 返回int32指针
func Int32Ptr(i int) *int32 {
	val := int32(i)
	return &val
}

// Int64Ptr 返回int64指针
func Int64Ptr(i int) *int64 {
	val := int64(i)
	return &val
}

// Float64Ptr 返回float64指针
func Float64Ptr(f float64) *float64 {
	return &f
}
