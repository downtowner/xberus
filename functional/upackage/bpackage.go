package upackage

import (
	"encoding/binary"
	"log"
	"regexp"
	"strconv"
)

/**
authoer: icarus
data: 2020-10-26

note:
This package provides encoding and decoding services for byte protocol, suport id and cmd mode. so its header length is not fixed.
in id mode: header length is 7 bytes.
in cmd mode: header length is 13 bytes.
important note: u must use it in given order.

steps to use:

create pkg step:
p :=NewBPackage{}
p.AddPackageID(12)
p.AddUint8(3)
...
p.Done()
p.GetData()

read pkg step：
p :=NewBPackage{}
p.AddBytes(...)
p.ReadPackageID()/p.ReadPackageCmd()
p.ReadDataLength()
p.Read...


header note:
first field: mark bit,0: uninit.use error, 1: id mode, 2: cmd mode
second field: id/cmd
third field: data size 4

**/

//HeaderType header mark
type HeaderType int8

const (
	//HeaderMarkUninit 0 Illegal package
	HeaderMarkUninit HeaderType = iota
	//HeaderMarkID 1 ID package
	HeaderMarkID
	//HeaderMarkCmd 2 command package
	HeaderMarkCmd
)

//BPackage Ceate a bit package
type BPackage struct {
	buf []byte
}

//NewBPackage create a bpackage obj,eg: it happened receive net data
func NewBPackage(params ...[]byte) *BPackage {

	p := BPackage{nil}

	for _, v := range params {

		p.AddBytes(v)
	}

	return &p
}

//AddPackageID ID mode
/**
note:
Structure of the communication protocol->1+2+4
**/
func (b *BPackage) AddPackageID(id int16) {

	if nil != b.buf {

		log.Panic("buf is not nil. AddPackageID must be call at first step")
	}

	//flag
	b.buf = append(b.buf, byte(HeaderMarkID))
	//id
	b.AddInt16(id)
	//length
	b.AddUint32(0)
}

//AddPackageCmd CMD mode
/**
note:
Structure of the communication protocol->1+8+4
**/
func (b *BPackage) AddPackageCmd(cmd string) {

	matched, err := regexp.MatchString("^[0-9a-zA-Z_]{1,8}$", cmd)

	if !matched || err != nil {

		log.Panicf("cmd must be string of 1-8 bytes. Support numbers. underscores. uppercase and lowercase letters.cmd = [%s]", cmd)
	}

	if nil != b.buf {

		log.Panic("buf is not nil. addpackagecmd must be call at first step")
	}

	//flag
	b.buf = append(b.buf, byte(HeaderMarkCmd))

	//cmd
	cmd = cmd + "       "
	hex := []byte(cmd)
	var n = binary.LittleEndian.Uint64(hex)
	b.AddUint64(n)

	//length
	b.AddUint32(0)
}

//Done The package has assembled data
func (b *BPackage) Done() {

	b.check()

	var datalen int
	switch HeaderType(b.buf[0]) {
	case HeaderMarkID: //1+2+4=7

		datalen = len(b.buf) - 7
		b.buf[3] = uint8(datalen)
		b.buf[4] = uint8(datalen >> 8)
		b.buf[5] = uint8(datalen >> 16)
		b.buf[6] = uint8(datalen >> 24)
	case HeaderMarkCmd: //1+8+4=13

		datalen = len(b.buf) - 13
		b.buf[9] = uint8(datalen)
		b.buf[10] = uint8(datalen >> 8)
		b.buf[11] = uint8(datalen >> 16)
		b.buf[12] = uint8(datalen >> 24)
	}
}

//GetData get data of bytes chunk
func (b *BPackage) GetData() []byte {

	return b.buf
}

//AddUint8 add data in uint8 format
func (b *BPackage) AddUint8(d uint8) {

	b.check()

	b.buf = append(b.buf, d)
}

//AddUint16 add data in uint16 format
func (b *BPackage) AddUint16(d uint16) {

	b.check()

	buf := make([]byte, 2)
	buf[0] = uint8(d)
	buf[1] = uint8(d >> 8)
	b.buf = append(b.buf, buf...)
}

//AddUint32 add data in uint32 format
func (b *BPackage) AddUint32(d uint32) {

	b.check()

	buf := make([]byte, 4)
	buf[0] = uint8(d)
	buf[1] = uint8(d >> 8)
	buf[2] = uint8(d >> 16)
	buf[3] = uint8(d >> 24)
	b.buf = append(b.buf, buf...)
}

//AddUint64 add data in uint64 format
func (b *BPackage) AddUint64(d uint64) {

	b.check()

	buf := make([]byte, 8)
	buf[0] = uint8(d)
	buf[1] = uint8(d >> 8)
	buf[2] = uint8(d >> 16)
	buf[3] = uint8(d >> 24)
	buf[4] = uint8(d >> 32)
	buf[5] = uint8(d >> 40)
	buf[6] = uint8(d >> 48)
	buf[7] = uint8(d >> 56)
	b.buf = append(b.buf, buf...)
}

//AddInt8 add data in int8 format
func (b *BPackage) AddInt8(d int8) {

	b.check()

	b.buf = append(b.buf, uint8(d))
}

//AddInt16 add data in int8 format
func (b *BPackage) AddInt16(d int16) {

	b.check()

	buf := make([]byte, 2)
	buf[0] = uint8(d)
	buf[1] = uint8(d >> 8)
	b.buf = append(b.buf, buf...)
}

//AddInt32 add data in int8 format
func (b *BPackage) AddInt32(d int32) {

	b.check()

	buf := make([]byte, 4)
	buf[0] = uint8(d)
	buf[1] = uint8(d >> 8)
	buf[2] = uint8(d >> 16)
	buf[3] = uint8(d >> 24)
	b.buf = append(b.buf, buf...)
}

//AddInt64 add data in int16 format
func (b *BPackage) AddInt64(d int64) {

	b.check()

	buf := make([]byte, 8)
	buf[0] = uint8(d)
	buf[1] = uint8(d >> 8)
	buf[2] = uint8(d >> 16)
	buf[3] = uint8(d >> 24)
	buf[4] = uint8(d >> 32)
	buf[5] = uint8(d >> 40)
	buf[6] = uint8(d >> 48)
	buf[7] = uint8(d >> 56)
	b.buf = append(b.buf, buf...)
}

//AddStringL add data of string type with length of string
func (b *BPackage) AddStringL(s string, l int) {

	b.check()

	buf := []byte(s)
	if len(buf) > l {

		buf = buf[:l]
	} else if len(buf) < l {

		zeroCount := l - len(buf)
		for i := 0; i < zeroCount; i++ {

			buf = append(buf, 0)
		}
	}

	b.buf = append(b.buf, buf...)
}

//AddString add data of string
func (b *BPackage) AddString(s string) {

	buf := []byte(s)
	b.buf = append(b.buf, buf...)
	b.buf = append(b.buf, 0)
}

//AddBool add data in bool format
func (b *BPackage) AddBool(istrue bool) {

	b.check()

	if istrue {

		b.AddUint8(1)
	} else {

		b.AddUint8(0)
	}
}

//AddFloat32 add data of float32
func (b *BPackage) AddFloat32(f float32, prec int) {

	b.check()

	//func FormatFloat(f float64, fmt byte, prec, bitSize int) string
	s := strconv.FormatFloat(float64(f), 'f', prec, 32)
	b.AddString(s)
}

//AddFloat64 add data of float64
func (b *BPackage) AddFloat64(f float64, prec int) {

	b.check()

	s := strconv.FormatFloat(f, 'f', prec, 64)
	b.AddString(s)
}

//AddBytes add data of bytes
func (b *BPackage) AddBytes(bs []byte) {

	b.buf = append(b.buf, bs...)
}

func (b *BPackage) check() {

	if 0 == int8(b.buf[0]) {

		log.Panic("BPackage uninit. please use it by given steps.")
	}
}

//ReadHeaderMark read header mark from bytes chunk
func (b *BPackage) ReadHeaderMark() HeaderType {

	if 0 == b.buf[0] {

		log.Panic("read header error, mark: 0")
	}

	return HeaderType(b.ReadUint8())
}

//ReadPackageID read package id from bytes chunk
func (b *BPackage) ReadPackageID() int16 {

	return b.ReadInt16()
}

//ReadPackageCmd read package cmd from bytes chunk
func (b *BPackage) ReadPackageCmd() string {

	if len(b.buf) < 8 {

		log.Panicf("Expect to read 8 bytes, but the buffer has only %d bytes", len(b.buf))
	}

	var cmd = string(b.buf[0:8])

	for i := 0; i < 7; i++ {

		if cmd[len(cmd)-1] == 0 || cmd[len(cmd)-1] == ' ' {

			cmd = cmd[0 : len(cmd)-1]
		}
	}

	b.buf = b.buf[8:]

	return cmd
}

//ReadDataLength data length without header
func (b *BPackage) ReadDataLength() int32 {

	return b.ReadInt32()
}

//ReadUint8 read data in uint8 format
func (b *BPackage) ReadUint8() uint8 {

	re := b.buf[0]

	b.buf = b.buf[1:]

	return re
}

//ReadUint16 read data in uint16 format
func (b *BPackage) ReadUint16() uint16 {

	var re uint16
	re = uint16(b.buf[0])
	re |= uint16(b.buf[1]) << 8

	b.buf = b.buf[2:]

	return re
}

//ReadUint32 read data in uint32 format
func (b *BPackage) ReadUint32() uint32 {

	var re uint32
	re = uint32(b.buf[0])
	re |= uint32(b.buf[1]) << 8
	re |= uint32(b.buf[2]) << 16
	re |= uint32(b.buf[3]) << 24

	b.buf = b.buf[4:]

	return re
}

//ReadUint64 read data in uint64 format
func (b *BPackage) ReadUint64() uint64 {

	var re uint64
	re = uint64(b.buf[0])
	re |= uint64(b.buf[1]) << 8
	re |= uint64(b.buf[2]) << 16
	re |= uint64(b.buf[3]) << 24
	re |= uint64(b.buf[4]) << 32
	re |= uint64(b.buf[5]) << 40
	re |= uint64(b.buf[6]) << 48
	re |= uint64(b.buf[7]) << 56

	b.buf = b.buf[8:]

	return re
}

//ReadInt8 read data in int8 format
func (b *BPackage) ReadInt8() int8 {

	re := b.buf[0]

	b.buf = b.buf[1:]

	return int8(re)
}

//ReadInt16 read data in int8 format
func (b *BPackage) ReadInt16() int16 {

	var re int16
	re = int16(b.buf[0])
	re |= int16(b.buf[1]) << 8

	b.buf = b.buf[2:]

	return re
}

//ReadInt32 read data in int8 format
func (b *BPackage) ReadInt32() int32 {

	var re int32
	re = int32(b.buf[0])
	re |= int32(b.buf[1]) << 8
	re |= int32(b.buf[2]) << 16
	re |= int32(b.buf[3]) << 24

	b.buf = b.buf[4:]

	return re
}

//ReadInt64 read data in int16 format
func (b *BPackage) ReadInt64() int64 {

	var re int64
	re = int64(b.buf[0])
	re |= int64(b.buf[1]) << 8
	re |= int64(b.buf[2]) << 16
	re |= int64(b.buf[3]) << 24
	re |= int64(b.buf[4]) << 32
	re |= int64(b.buf[5]) << 40
	re |= int64(b.buf[6]) << 48
	re |= int64(b.buf[7]) << 56

	b.buf = b.buf[8:]

	return re
}

//ReadString read data in string format
func (b *BPackage) ReadString() string {

	for i := 0; i < len(b.buf); i++ {

		if 0 == b.buf[i] {

			str := string(b.buf[:i])
			b.buf = b.buf[i+1:]

			return str
		}
	}

	return ""
}

//ReadStringL read string with fix length
func (b *BPackage) ReadStringL(l int) string {

	buf := b.buf[:l]
	for i := 0; i < len(buf); i++ {

		if buf[i] == 0 {

			buf = buf[:i]

			break
		}
	}

	b.buf = b.buf[l:]

	return string(buf)
}

//ReadBool read data in bool format
func (b *BPackage) ReadBool() bool {

	if 0 == b.ReadUint8() {

		return false
	}

	return true
}

//ReadFloat32 add data of float32
func (b *BPackage) ReadFloat32() (float32, error) {

	f, err := strconv.ParseFloat(b.ReadString(), 32)

	return float32(f), err
}

//ReadFloat64 add data of float64
func (b *BPackage) ReadFloat64() (float64, error) {

	return strconv.ParseFloat(b.ReadString(), 64)
}
