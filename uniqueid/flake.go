package uniqueid

import (
	"fmt"

	"github.com/sony/sonyflake"
)

var flake *sonyflake.Sonyflake

// InitUniqueID 初始化唯一ID生成器
func InitUniqueID() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

// NewID 创建唯一ID
func NewID() (uint64, error) {
	id, err := flake.NextID()
	if err != nil {
		return 0, fmt.Errorf("flake.NextID err: %v", err)
	}
	return id, nil
}
