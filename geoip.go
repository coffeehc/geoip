// geoip project geoip.go
package geoip

import (
	"fmt"
)

/*
格式:
map[city:map[geoname_id:5037649 names:map[de:Minneapolis en:Minneapolis es:Mineápolis fr:Minneapolis ja:ミネアポリス pt-BR:Minneapolis ru:Миннеаполис zh-CN:明尼阿波利斯]] continent:map[code:NA geoname_id:6255149 names:map[de:Nordamerika en:North America es:Norteamérica fr:Amérique du Nord ja:北アメリカ pt-BR:América do Norte ru:Северная Америка zh-CN:北美洲]] country:map[geoname_id:6252001 iso_code:US names:map[de:USA en:United States es:Estados Unidos fr:États-Unis ja:アメリカ合衆国 pt-BR:Estados Unidos ru:США zh-CN:美国]] location:map[latitude:44.9759 longitude:-93.2166 metro_code:613 time_zone:America/Chicago] postal:map[code:55414] registered_country:map[geoname_id:6252001 iso_code:US names:map[de:USA en:United States es:Estados Unidos fr:États-Unis ja:アメリカ合衆国 pt-BR:Estados Unidos ru:США zh-CN:美国]] subdivisions:[map[geoname_id:5037779 iso_code:MN names:map[en:Minnesota es:Minnesota ja:ミネソタ州 ru:Миннесота]]]]
*/
type GeoIp_City struct {
	Id           int     `json:"id"`
	Ip           string  `json:"ip"`
	City         string  `json:"city"`
	Continent    string  `json:"continent"`
	Country      string  `json:"country"`
	Postal       string  `json:"postal"`
	IsoCode      string  `json:"isoCode"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	TimeZone     string  `json:"timeZone"`
	Subdivisions string  `json:"subdivisions"`
}

func (this *IpDataBase) GetCityByIp(ip string, isoCode string) (*GeoIp_City, error) {
	entry, err := this.get(ip)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf("未知地区,Ip:%s", ip)
	}
	if isoCode != "" {
		var support bool = false
		for _, language := range this.metadata.languages {
			if language == isoCode {
				support = true
				break
			}
		}
		if !support {
			return nil, fmt.Errorf("不支持%s的语言编码", isoCode)
		}
	} else {
		isoCode = "zh-CN"
	}
	city := new(GeoIp_City)
	city.IsoCode = isoCode
	city.Ip = ip
	var value map[interface{}]interface{}
	node := entry.node.(map[interface{}]interface{})
	if node["city"] != nil {
		value = node["city"].(map[interface{}]interface{})
		if value["geoname_id"] != nil {
			city.Id = int(value["geoname_id"].(uint32))
		}
		value = value["names"].(map[interface{}]interface{})
		city.City = getName(value, isoCode)
	}
	if node["continent"] != nil {
		value = node["continent"].(map[interface{}]interface{})
		value = value["names"].(map[interface{}]interface{})
		city.Continent = getName(value, isoCode)
	}
	if node["country"] != nil { //国家是有可能为空的
		value = node["country"].(map[interface{}]interface{})
		if value["names"] != nil {
			value = value["names"].(map[interface{}]interface{})
			city.Country = getName(value, isoCode)
		}
	}
	if node["postal"] != nil {
		value = node["postal"].(map[interface{}]interface{})
		city.Postal = value["code"].(string)
	}
	if node["location"] != nil {
		value = node["location"].(map[interface{}]interface{})
		city.Latitude = value["latitude"].(float64)
		city.Longitude = value["longitude"].(float64)
		if value["time_zone"] != nil {
			city.TimeZone = value["time_zone"].(string)
		}
	}
	if node["subdivisions"] != nil {
		values := node["subdivisions"].([]interface{})
		if len(values) > 0 {
			value = values[0].(map[interface{}]interface{})
			if value["iso_code"] != nil {
				city.Subdivisions = value["iso_code"].(string)
			}
		}
	}
	return city, nil
}

func getName(value map[interface{}]interface{}, isoCode string) string {
	if value[isoCode] != nil {
		return value[isoCode].(string)
	} else {
		if value["en"] != nil {
			return value["en"].(string)
		} else {
			return "未知"
		}
	}
}
