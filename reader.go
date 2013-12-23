// reader
package geoip

import (
	"fmt"
	"github.com/coffeehc/logger"
	"net"
	"os"
)

var DATA_SECTION_SEPARATOR_SIZE int = 16
var METADATA_START_MARKER = []byte{0xAB, 0xCD, 0xEF, 'M', 'a', 'x', 'M', 'i', 'n', 'd', '.', 'c', 'o', 'm'}

type IpDataBase struct {
	path          string
	data          []byte
	metadataStart int
	decoder       *decoder
	metadata      *metadata
	ipV4Start     int
}

func NewIpDataBase(path string) (*IpDataBase, error) {
	database := new(IpDataBase)
	database.ipV4Start = 0
	database.path = path
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	database.data = make([]byte, fileSize)
	_, err = file.Read(database.data)
	if err != nil {
		return nil, err
	}
	err = database.findMetadataStart()
	if err != nil {
		return nil, err
	}
	database.decoder = newDecoder(database, database.metadataStart)
	entry, err := database.decoder.decode(database.metadataStart)
	if err != nil {
		logger.Errorf("解析出现异常:%s", err)
		return nil, err
	}
	database.metadata, err = newMetadata(entry)
	if err != nil {
		logger.Errorf("构建MetaData数据失败,%v", entry)
		return nil, err
	}
	database.decoder = newDecoder(database, database.metadata.searchTreeSize+DATA_SECTION_SEPARATOR_SIZE)
	logger.Infof("解析GeoIp数据库[%s]成功,metadata:%v", path, database.metadata)
	return database, nil
}

func (this *IpDataBase) findMetadataStart() error {
	markerSize := len(METADATA_START_MARKER)
	fileSize := len(this.data)
	goto FILE
FILE:
	for i := 0; i < fileSize-markerSize+1; i++ {
		for j := 0; j < markerSize; j++ {
			if this.data[fileSize-i-j-1] != METADATA_START_MARKER[markerSize-j-1] {
				continue FILE
			}
		}
		this.metadataStart = fileSize - i
		return nil
	}
	return fmt.Errorf("不能找到数据库标记,路径:%s,这可能不是一个标准的数据库", this.path)
}

func (this *IpDataBase) get(ip string) (*entry_data, error) {
	pointer, err := this.findAddressInTree(ip)
	if err != nil {
		logger.Errorf("获取pointer失败,%s", err)
		return nil, err
	}
	if pointer == 0 {
		return nil, nil
	}
	return this.resolveDataPointer(pointer)
}
func (this *IpDataBase) resolveDataPointer(pointer int) (*entry_data, error) {
	resolved := pointer - this.metadata.nodeCount + this.metadata.searchTreeSize
	if resolved >= len(this.data) {
		return nil, fmt.Errorf("指定的索引已经越界,%d", resolved)
	}
	return this.decoder.decode(resolved)
}

func (this *IpDataBase) findAddressInTree(ipAddr string) (int, error) {
	ip := net.ParseIP(ipAddr)
	if ip4 := ip.To4(); ip4 != nil {
		ip = ip4
	}
	bitLength := len(ip) * 8
	record, err := this.startNode(bitLength)
	if err != nil {
		logger.Errorf("StartNode失败:原因:%s", err)
		return 0, err
	}
	for i := 0; i < bitLength; i++ {
		if record >= this.metadata.nodeCount {
			break
		}
		b := int(0xff & ip[i/8])
		bit := 1 & (b >> uint32(7-(i%8)))
		record, err = this.readNode(record, bit)
		if err != nil {
			return 0, err
		}
	}
	if record == this.metadata.nodeCount {
		return 0, nil
	} else if record > this.metadata.nodeCount {
		return record, nil
	}

	return 0, fmt.Errorf("出现了未知错误")
}

func (this *IpDataBase) startNode(bitLength int) (int, error) {
	if this.metadata.ipVersion == 6 && bitLength == 32 {
		return this.ipV4StartNode()
	}
	return 0, nil
}

func (this *IpDataBase) ipV4StartNode() (int, error) {
	if this.metadata.ipVersion == 4 {
		return 0, nil
	}
	if this.ipV4Start != 0 {
		return this.ipV4Start, nil
	}
	node := 0
	var err error
	for i := 0; i < 96 && node < int(this.metadata.nodeCount); i++ {
		node, err = this.readNode(node, 0)
		if err != nil {
			return 0, err
		}
	}
	this.ipV4Start = node
	return node, nil
}

func (this *IpDataBase) readNode(nodeNumber int, index int) (int, error) {
	baseOffset := nodeNumber * this.metadata.nodeByteSize
	switch this.metadata.recordSize {
	case 24:
		offset := baseOffset + index*3
		return int(get_uintX(this.data[offset:offset+3], 0, 3)), nil
	case 28:
		middle := this.data[baseOffset+3]
		if index == 0 {
			middle = byte(uint(0xF0&middle) >> 4)
		} else {
			middle = 0x0F & middle
		}
		offset := baseOffset + index*4
		return int(get_uintX(this.data[offset:offset+3], uint64(middle), 3)), nil
	case 32:
		offset := baseOffset + index*4
		return int(get_uintX(this.data[offset:offset+4], 0, 4)), nil
	default:
		return 0, fmt.Errorf("recordSize不能找到recordSize匹配的类型")
	}
}
