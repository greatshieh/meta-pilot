package v1

import (
	"crypto/sha1"
	"fmt"
	"server/internal/api"
	"server/pkg/global"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
)

type WechatService struct{}

func (*WechatService) Get(query api.ServerValidateReqModel) string {
	sortList := []string{global.MPA_CONFIG.WeChat.Token, query.TimeStamp, query.Nonce}

	slice.Sort(sortList)

	sha := sha1.New()
	sha.Write([]byte(strings.Join(sortList, "")))

	hashcode := fmt.Sprintf("%x", sha.Sum(nil))

	if query.Signature == hashcode {
		return query.Echostr
	}

	return ""
}
