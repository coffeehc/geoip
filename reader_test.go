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
	//city, err := database.get("172.22.15.85")
	//city, err := database.GetCityByIp("218.89.49.93", "")
	city, err := database.GetCityByIp("210.076.200.1", "")
	if err != nil {
		t.Fatalf("搜索出现错误:%s", err)
	} else {
		logger.Info("获取内容:%v", city)
	}
	time.Sleep(time.Second * 5)
}
