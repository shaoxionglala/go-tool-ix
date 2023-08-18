package snowflake

import (
	"sync"
	"time"
)

var snowflakeInstance Snowflake

// InitSnowflake 初始化雪花算法
func InitSnowflake(workerID, datacenterID int64) {
	snowflakeInstance = Snowflake{
		timestamp:    0,
		workerID:     workerID,
		datacenterID: datacenterID,
		sequence:     0,
	}
}

type Snowflake struct {
	sync.Mutex         // 锁
	timestamp    int64 // 时间戳 ，毫秒
	workerID     int64 // 工作节点
	datacenterID int64 // 数据中心机房ID
	sequence     int64 // 序列号
}

const (
	epoch             = int64(1577808000000)                           // 设置起始时间(时间戳/毫秒)：2020-01-01 00:00:00，有效期69年
	timestampBits     = uint(41)                                       // 时间戳占用位数
	datacenterIDBits  = uint(2)                                        // 数据中心ID所占位数
	workerIDBits      = uint(7)                                        // 机器ID所占位数
	sequenceBits      = uint(12)                                       // 序列所占的位数
	timestampMax      = int64(-1 ^ (-1 << timestampBits))              // 时间戳最大值
	datacenterIDMax   = int64(-1 ^ (-1 << datacenterIDBits))           // 支持的最大数据中心ID数量
	workerIDMax       = int64(-1 ^ (-1 << workerIDBits))               // 支持的最大机器ID数量
	sequenceMask      = int64(-1 ^ (-1 << sequenceBits))               // 支持的最大序列ID数量
	workerIDShift     = sequenceBits                                   // 机器ID左移位数
	datacenterIDShift = sequenceBits + workerIDBits                    // 数据中心ID左移位数
	timestampShift    = sequenceBits + workerIDBits + datacenterIDBits // 时间戳左移位数
)

// NextVal 包级别的函数，不需要初始化结构体
func NextVal() int64 {
	return snowflakeInstance.nextVal()
}

func (s *Snowflake) nextVal() int64 {
	s.Lock()
	now := time.Now().UnixNano() / 1000000 // 转毫秒
	if s.timestamp == now {
		// 当同一时间戳（精度：毫秒）下多次生成ID会增加序列号
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 如果当前序列超出12bit长度，则需要等待下一毫秒
			// 下一毫秒将使用sequence:0
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// 不同时间戳（精度：毫秒）下直接使用序列号：0
		s.sequence = 0
	}
	t := now - epoch
	if t > timestampMax {
		s.Unlock()
		return 0
	}
	s.timestamp = now
	r := int64((t)<<timestampShift | (s.datacenterID << datacenterIDShift) | (s.workerID << workerIDShift) | (s.sequence))
	s.Unlock()
	return r
}
