package uniqueid

import (
	"github.com/sony/sonyflake"

	"github.com/turing-era/turingera-shared/log"
)

var flake *sonyflake.Sonyflake

// InitUniqueID 初始化唯一ID生成器
func InitUniqueID() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

// NewID 创建唯一ID
func NewID() int64 {
	id, err := flake.NextID()
	if err != nil {
		log.Errorf("flake.NextID err: %v", err)
		return 0
	}
	return int64(id)
}
