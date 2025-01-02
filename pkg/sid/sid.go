package sid

import (
	"github.com/sony/sonyflake"
	"time"
)

type Sid struct {
	sf *sonyflake.Sonyflake
}

// NewSid 创建一个新的 Sid 实例
func NewSid() *Sid {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	if sf == nil {
		panic("sonyflake not created")
	}
	return &Sid{sf}
}

// GenIncrementalUID 基于时间戳生成唯一的用户 ID
func (s *Sid) GenIncrementalUID() int64 {
	// 获取当前时间的毫秒级时间戳
	return time.Now().UnixMilli()
}

// GenSonyflakeID 基于 Sonyflake 生成唯一的用户 ID
func (s *Sid) GenSonyflakeID() (int64, error) {
	id, err := s.sf.NextID()
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}
