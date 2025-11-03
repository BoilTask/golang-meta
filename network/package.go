package network

import (
	"bytes"
	"encoding/binary"
	"io"
	metaerror "meta/meta-error"
	"net"
)

const PackageReceiveSizeMax = 4 * 1024 * 1024

type Package struct {
	Size int32
	Data []byte
}

func isSizeValid(size int32) bool {
	if size <= 0 || size > PackageReceiveSizeMax {
		return false
	}
	return true
}

func LoadPackage(conn net.Conn) (*Package, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		return nil, err
	}
	var dataSize int32
	if err := binary.Read(bytes.NewReader(header), binary.BigEndian, &dataSize); err != nil {
		return nil, err
	}
	if !isSizeValid(dataSize) {
		return nil, metaerror.New("invalid data size: %d", dataSize)
	}
	data := make([]byte, dataSize)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}
	networkPackage := &Package{Size: dataSize, Data: data}
	return networkPackage, nil
}

func MakeBytes(packets ...*Packet) ([]byte, error) {
	var packetBuffer bytes.Buffer
	for _, packet := range packets {
		err := packet.MakeBytes(&packetBuffer)
		if err != nil {
			return nil, err
		}
	}
	var packageBuffer bytes.Buffer
	err := binary.Write(&packageBuffer, binary.BigEndian, int32(packetBuffer.Len()))
	if err != nil {
		return nil, err
	}
	packageBuffer.Write(packetBuffer.Bytes())
	return packageBuffer.Bytes(), nil
}
