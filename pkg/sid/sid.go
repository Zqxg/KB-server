package sid

import (
	"github.com/sony/sonyflake"
	"math/rand"
)

type Sid struct {
	sf *sonyflake.Sonyflake
}

func NewSid() *Sid {
	sf := sonyflake.NewSonyflake(sonyflake.Settings{})
	if sf == nil {
		panic("sonyflake not created")
	}
	return &Sid{sf}
}
func (s Sid) GenString() (string, error) {
	id, err := s.sf.NextID()
	if err != nil {
		return "", err
	}
	return IntToBase62(int(id)), nil
}
func (s Sid) GenUint64() (uint64, error) {
	return s.sf.NextID()
}

// Gen25PrefixUID 基于 GenUint64 生成以 25 开头，且长度为 9 位的数字型 user_id
func (s Sid) Gen25PrefixUID() (int64, error) {
	// 先获取一个唯一的 uint64 SID
	id, err := s.GenUint64()
	if err != nil {
		return 0, err
	}

	// 将 ID 转换为 9 位数字
	// 25 为固定前缀，后续部分由 Sonyflake ID 的后几位组成
	prefix := int64(25 * 1e7) // 25 后跟 7 个零，确保 SID 以 25 开头

	// 获取 Sonyflake ID 的最后 7 位数字（最大 9999999）
	remainingDigits := id % 1e7

	// 拼接前缀和剩余的 7 位数字，形成 9 位长度的 SID
	finalID := prefix + remainingDigits

	// 如果生成的 ID 不足 9 位，我们填充随机数字，确保最终长度为 9 位
	if finalID < 100000000 {
		finalID += rand.Int63n(1000000000 - 100000000) // 填充直到 9 位数
	}

	return finalID, nil
}
