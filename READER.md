GEOIP2 的解析库 go语言实现
======
***


	```
	func TestFindMetadataStart(t *testing.T) {
	logger.StartDevModel()
	database, err := NewIpDataBase("D:/work/goOther/GeoLite2-City.mmdb/GeoLite2-City.mmdb")
	if err != nil {
		t.Fatalf("错误:%s", err)
	}
	city, err := database.get("172.22.15.85")
	//city, err := database.GetCityByIp("172.22.15.85", "")
	if err != nil {
		t.Fatalf("搜索出现错误:%s", err)
	} else {
		logger.Debugf("获取内容:%v", city)
	}
	utils.WaitTimeOut(1)
	}
	```

使用很简单,会写go的就会用.

