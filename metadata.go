// metadata
package geoip

import "fmt"

type metadata struct {
	binaryFormatMajorVersion int //uint16
	binaryFormatMinorVersion int //uint16
	buildEpoch               uint64
	databaseType             string
	description              map[interface{}]interface{}
	ipVersion                int //uint16
	nodeCount                int //uint32
	recordSize               int //uint32
	nodeByteSize             int //uint32
	searchTreeSize           int //uint32
	languages                []interface{}
}

func newMetadata(entry *entry_data) (*metadata, error) {
	if v, ok := entry.node.(map[interface{}]interface{}); ok {
		m := new(metadata)
		var data interface{}
		data = v["binary_format_major_version"]
		if v1, ok := data.(uint16); ok {
			m.binaryFormatMajorVersion = int(v1)
		} else {
			return nil, fmt.Errorf("binaryFormatMajorVersion不是uint16类型:%v", data)
		}
		data = v["binary_format_minor_version"]
		if v2, ok := data.(uint16); ok {
			m.binaryFormatMinorVersion = int(v2)
		} else {
			return nil, fmt.Errorf("binaryFormatMinorVersion不是uint16类型:%v", data)
		}
		data = v["build_epoch"]
		if v3, ok := data.(uint64); ok {
			m.buildEpoch = v3
		} else {
			return nil, fmt.Errorf("buildEpoch不是uint64类型:%v", data)
		}
		data = v["database_type"]
		if v4, ok := data.(string); ok {
			m.databaseType = v4
		} else {
			return nil, fmt.Errorf("databaseType不是string类型:%v", data)
		}
		data = v["description"]
		if v5, ok := data.(map[interface{}]interface{}); ok {
			m.description = v5
		} else {
			return nil, fmt.Errorf("description不是map[string]string类型:%v", data)
		}
		data = v["ip_version"]
		if v6, ok := data.(uint16); ok {
			m.ipVersion = int(v6)
		} else {
			return nil, fmt.Errorf("ipVersion不是uint16类型:%v", data)
		}
		data = v["languages"]
		if v7, ok := data.([]interface{}); ok {
			m.languages = v7
		} else {
			return nil, fmt.Errorf("languages不是[]string类型:%v", data)
		}
		data = v["node_count"]
		if v8, ok := data.(uint32); ok {
			m.nodeCount = int(v8)
		} else {
			return nil, fmt.Errorf("nodeCount不是uint32类型:%v", data)
		}
		data = v["record_size"]
		if v9, ok := data.(uint16); ok {
			m.recordSize = int(v9)
		} else {
			return nil, fmt.Errorf("recordSize不是uint32类型:%v", data)
		}
		m.nodeByteSize = m.recordSize / 4
		m.searchTreeSize = m.nodeCount * m.nodeByteSize
		return m, nil
	} else {
		return nil, fmt.Errorf("entry不是Map类型,内容:%v", entry.node)
	}
}
