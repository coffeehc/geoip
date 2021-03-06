// reader_test
package geoip

import (
	"testing"
	"time"

	"github.com/coffeehc/logger"
)

func TestFindMetadataStart(t *testing.T) {
	database, err := NewIpDataBase("/Users/coffee/coder/gowork/coffee/src/logagent/server/config/GeoLite2-City.mmdb")
	if err != nil {
		t.Fatalf("错误:%s", err)
	}
	//	city, err := database.get("172.22.15.85")

	node, err := database.GetNodeByIp("14.17.37.144", LANG_CN)
	logger.Info("获取内容:%#v", node)
	node, err = database.GetNodeByIp("112.105.54.153", LANG_CN)
	logger.Info("获取内容:%#v", node)
	node, err = database.GetNodeByIp("175.45.20.138", LANG_CN)
	logger.Info("获取内容:%#v", node)
	node, err = database.GetNodeByIp("122.100.160.253", LANG_CN)
	if err != nil {
		t.Fatalf("搜索出现错误:%s", err)
	} else {
		logger.Info("获取内容:%#v", node)
	}
	time.Sleep(time.Second * 5)
}

