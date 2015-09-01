// define
package geoip

const (
	LANG_EN = "en"
	LANG_CN = "zh-CN"
	LANG_JA = "ja"
	LANG_BR = "pt-BR"
	LANG_RU = "ru"
	LANG_DE = "de"
	LANG_ES = "es"
	LANG_FR = "fr"
)

type Node struct {
	Ip                 string             `json:"ip"`
	City               City               `json:"city"`
	Subdivisions       Subdivisions       `json:"subdivisions"`
	Country            Country            `json:"country"`
	Continent          Continent          `json:"continent"`
	Registered_country Registered_country `json:"registered_country"`
	Location           Location           `json:"location"`
	//postal
}

func parseNode(data map[interface{}]interface{}, language string) *Node {
	node := new(Node)
	if value, ok := data["city"].(map[interface{}]interface{}); ok {
		node.City = parseCity(value, language)
	}
	if value, ok := data["subdivisions"].([]interface{}); ok {
		if len(value) > 0 {
			if v, ok1 := value[0].(map[interface{}]interface{}); ok1 {
				node.Subdivisions = parseSubdivisions(v, language)
			}
		}
	}
	if value, ok := data["country"].(map[interface{}]interface{}); ok {
		node.Country = parseCountry(value, language)
	}
	if value, ok := data["continent"].(map[interface{}]interface{}); ok {
		node.Continent = parseContinent(value, language)
	}
	if value, ok := data["registered_country"].(map[interface{}]interface{}); ok {
		node.Registered_country = parseRegistered_country(value, language)
	}
	if value, ok := data["location"].(map[interface{}]interface{}); ok {
		node.Location = parseLocation(value)
	}
	return node
}

type City struct {
	Geoname_id int    `json:"geoname_id"`
	Name       string `json:"name"`
}

func parseName(data map[interface{}]interface{}, language string) string {
	if names, ok := data["names"]; ok {
		if namesMap, ok := names.(map[interface{}]interface{}); ok {
			if name, ok := namesMap[language]; ok {
				_name, _ := name.(string)
				return _name
			}
			_name, _ := namesMap[LANG_EN].(string)
			return _name
		}
	}
	return ""
}

func parseCity(data map[interface{}]interface{}, language string) City {
	city := new(City)
	if geoname_id, ok := data["geoname_id"]; ok {
		city.Geoname_id = int(geoname_id.(uint32))
	} else {
		return *city
	}
	city.Name = parseName(data, language)
	return *city
}

type Subdivisions struct {
	Geoname_id int    `json:"geoname_id"`
	Name       string `json:"name"`
	Iso_code   string `json:"iso_code"`
}

func parseSubdivisions(data map[interface{}]interface{}, language string) Subdivisions {
	subdivisions := new(Subdivisions)
	if geoname_id, ok := data["geoname_id"]; ok {
		subdivisions.Geoname_id = int(geoname_id.(uint32))
	} else {
		return *subdivisions
	}
	subdivisions.Name = parseName(data, language)
	subdivisions.Iso_code, _ = data["iso_code"].(string)
	return *subdivisions
}

type Country struct {
	Geoname_id int    `json:"geoname_id"`
	Name       string `json:"name"`
	Iso_code   string `json:"iso_code"`
}

func parseCountry(data map[interface{}]interface{}, language string) Country {
	country := new(Country)
	if geoname_id, ok := data["geoname_id"]; ok {
		country.Geoname_id = int(geoname_id.(uint32))
	} else {
		return *country
	}
	country.Name = parseName(data, language)
	country.Iso_code, _ = data["iso_code"].(string)
	return *country
}

type Continent struct {
	Geoname_id int    `json:"geoname_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
}

func parseContinent(data map[interface{}]interface{}, language string) Continent {
	continent := new(Continent)
	if geoname_id, ok := data["geoname_id"]; ok {
		continent.Geoname_id = int(geoname_id.(uint32))
	} else {
		return *continent
	}
	continent.Name = parseName(data, language)
	continent.Code, _ = data["code"].(string)
	return *continent
}

type Registered_country struct {
	Geoname_id int    `json:"geoname_id"`
	Name       string `json:"name"`
	Iso_code   string `json:"iso_code"`
}

func parseRegistered_country(data map[interface{}]interface{}, language string) Registered_country {
	registered_country := new(Registered_country)
	if geoname_id, ok := data["geoname_id"]; ok {
		registered_country.Geoname_id = int(geoname_id.(uint32))
	} else {
		return *registered_country
	}
	registered_country.Name = parseName(data, language)
	registered_country.Iso_code, _ = data["iso_code"].(string)
	return *registered_country
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Time_zone string  `json:"time_zone"`
}

func parseLocation(data map[interface{}]interface{}) Location {
	location := new(Location)
	location.Latitude, _ = data["latitude"].(float64)
	location.Longitude, _ = data["longitude"].(float64)
	location.Time_zone, _ = data["time_zone"].(string)
	return *location
}
