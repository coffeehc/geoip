// geoip project geoip.go
package geoip

import (
	"errors"
	"fmt"
)

func (this *IpDataBase) GetNodeByIp(ip string, language string) (*Node, error) {
	entry, err := this.get(ip)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf("未知地区,Ip:%s", ip)
	}
	if language != "" {
		var support bool = false
		for _, _language := range this.metadata.languages {
			if _language == language {
				support = true
				break
			}
		}
		if !support {
			return nil, fmt.Errorf("不支持%s的语言编码", language)
		}
	} else {
		language = LANG_EN
	}
	if data, ok := entry.node.(map[interface{}]interface{}); ok {
		node := parseNode(data, language)
		node.Ip = ip
		return node, nil
	}
	return nil, errors.New("数据类型错误")
}
