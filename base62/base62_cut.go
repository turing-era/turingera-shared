package base62

// AlphanumericSet 字母数字集
var AlphanumericSet = []rune{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

// GetInvCodeByUIDUnique 获取指定长度的邀请码
func GetInvCodeByUIDUnique(uid uint64, l int) string {
	var code []rune
	for i := 0; i < l; i++ {
		idx := uid % uint64(len(AlphanumericSet))
		code = append(code, AlphanumericSet[idx])
		uid = uid / uint64(len(AlphanumericSet)) // 相当于右移一位（62进制）
	}
	return string(code)
}
