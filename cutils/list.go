package cutils

// InStringList 是否在列表中
func InStringList(items []string, id string) bool {
	for _, v := range items {
		if v == id {
			return true
		}
	}
	return false
}

// InUint64List 是否在列表中
func InUint64List(items []uint64, id uint64) bool {
	for _, v := range items {
		if v == id {
			return true
		}
	}
	return false
}

// InInt64List 是否在列表中
func InInt64List(items []int64, id int64) bool {
	for _, v := range items {
		if v == id {
			return true
		}
	}
	return false
}

// InUint32List 是否在列表中
func InUint32List(items []uint32, id uint32) bool {
	for _, v := range items {
		if v == id {
			return true
		}
	}
	return false
}

// InInt32List 是否在列表中
func InInt32List(items []int32, id int32) bool {
	for _, v := range items {
		if v == id {
			return true
		}
	}
	return false
}
