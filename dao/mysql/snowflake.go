package mysql

import (
	"github.com/bwmarrin/snowflake"
	"time"
)

var node *snowflake.Node //包含生成雪花ID的基础信息的结构体
// 妈的，不知道这些雪花算法到底在干什么，等我基础算法学完了再来深入研究
//
// 对于我们来说machineID填1就可以了
func InitSnowFlack(startTime string, machineID int64) error {
	var st time.Time
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		return err
	}
	return nil
}

// 雪花算法实现创建UID
func GenUID() int64 {
	return node.Generate().Int64()
}
