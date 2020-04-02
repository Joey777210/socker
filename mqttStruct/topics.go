package mqttStruct

import (
	"strings"
)

type devTopic struct {
	Type int
	Topic string
}

const (
	SysOnLinePub	= iota
	SysDataPub
	SysDataSub
)

var topics = []*devTopic{
	{SysOnLinePub, "sys/{CN}/status/online"},
	{SysDataPub, "sys/{CN}/test/pub"},
	{SysDataSub, "sys/{CN}/test/sub"},
}

func Replace(cn string) {
	for i := range topics {
		topics[i].Topic = strings.Replace(topics[i].Topic, "{CN}", cn, -1)
	}
}

func GetTopic(Type int) string {
	var s string
	for i := range topics {
		if topics[i].Type == Type {
			s = topics[i].Topic
			break
		}
	}
	return s
}