// reader_test
package geoip

import (
	"logger"
	"testing"
	"utils"
)

func TestFindMetadataStart(t *testing.T) {
	logger.StartDevModel()
	database, err := NewIpDataBase("D:/work/goOther/GeoLite2-City.mmdb/GeoLite2-City.mmdb")
	if err != nil {
		t.Fatalf("错误:%s", err)
	}
	//city, err := database.get("172.22.15.85")
	//city, err := database.GetCityByIp("218.89.49.93", "")
	city, err := database.GetCityByIp("122.224.178.61", "")
	if err != nil {
		t.Fatalf("搜索出现错误:%s", err)
	} else {
		logger.Debugf("获取内容:%v", city)
	}
	utils.WaitTimeOut(1)
}
