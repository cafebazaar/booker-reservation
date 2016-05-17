package proto

import "github.com/cafebazaar/booker-reservation/common"

func ReplyPropertiesTemplate() *ReplyProperties {
	return &ReplyProperties{
		ServerVersion: common.Version,
	}
}
