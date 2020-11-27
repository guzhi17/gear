// ------------------
// User: pei
// DateTime: 2020/10/29 14:10
// Description: 
// ------------------

package ims

import (
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
)

const ImPackageTag = 0x1615
const HeaderLength = 18
var(
	ErrUnknownPackage = errors.New("unknown package")
	ErrUnknown = errors.New("error unknown")
	ErrTimeout = errors.New("error timeout")
)


type ErrorPackage struct {
	Package
}

func (p ErrorPackage)Error() string{
	var buf strings.Builder
	buf.WriteString("error unknown package: ")
	p.string(&buf)
	return buf.String()
}


type Header struct {
	Tag uint16
	Ver uint16
	Tp uint8
	Res uint8
	Qid uint32
	Fid uint32
	Code uint32
}

//================================================================
type Package []byte

func (p Package)String() string  {
	var buf strings.Builder
	p.string(&buf)
	return buf.String()
}

func (p Package)string(buf *strings.Builder)  {
	buf.WriteString("[tag: ")
	buf.WriteString(strconv.FormatInt(int64(p.Tag()), 10))
	buf.WriteString(", Ver: ")
	buf.WriteString(strconv.FormatInt(int64(p.Ver()), 10))
	buf.WriteString(", Tp: ")
	buf.WriteString(strconv.FormatInt(int64(p.Tp()), 10))
	buf.WriteString(", Res: ")
	buf.WriteString(strconv.FormatInt(int64(p.Res()), 10))
	buf.WriteString(", Qid: ")
	buf.WriteString(strconv.FormatInt(int64(p.Qid()), 10))
	buf.WriteString(", Fid: ")
	buf.WriteString(strconv.FormatInt(int64(p.Fid()), 10))
	buf.WriteString(", Code: ")
	buf.WriteString(strconv.FormatInt(int64(p.Code()), 10))
	buf.WriteString(", Len: ")
	b := p.Body()
	buf.WriteString(strconv.FormatInt(int64(len(b)), 10))
	buf.WriteString("] ")
	buf.Write(b)
}


func PackageFromBytes(v []byte) (Package, error)  {
	if len(v) < HeaderLength{
		return nil, ErrUnknownPackage
	}
	return Package(v), nil
}
func (p Package)Body() []byte {
	return p[HeaderLength:]
}
func (p Package)Tag() uint16 {
	return binary.BigEndian.Uint16(p[:])
}
func (p Package)SetTag(v uint16) Package {
	binary.BigEndian.PutUint16(p[:], v)
	return p
}
func (p Package)Ver() uint16 {
	return binary.BigEndian.Uint16(p[2:])
}
func (p Package)SetVer(v uint16) Package {
	binary.BigEndian.PutUint16(p[2:], v)
	return p
}
func (p Package)Tp() uint8 {
	return p[4]
}
func (p Package)SetTp(v uint8) Package {
	p[4] = v
	return p
}
func (p Package)Res() uint8 {
	return p[5]
}
func (p Package)SetRes(v uint8)Package {
	p[5] = v
	return p
}
func (p Package)Qid() uint32 {
	return binary.BigEndian.Uint32(p[6:])
}
func (p Package)SetQid(v uint32) Package {
	binary.BigEndian.PutUint32(p[6:], v)
	return p
}
func (p Package)Fid() uint32 {
	return binary.BigEndian.Uint32(p[10:])
}
func (p Package)SetFid(v uint32) Package {
	binary.BigEndian.PutUint32(p[10:], v)
	return p
}
func (p Package)Code() uint32 {
	return binary.BigEndian.Uint32(p[14:])
}
func (p Package)SetCode(v uint32) Package {
	binary.BigEndian.PutUint32(p[14:], v)
	return p
}
func (p Package)Header() Package {
	return p[:HeaderLength]
}
func (p Package)CloneHeader() Package {
	buf := make([]byte, HeaderLength)
	copy(buf, p)
	return buf
}
func MakeHeader(h Header) (p Package){
	p = make([]byte, HeaderLength)
	p.SetTag(h.Tag)
	p.SetVer(h.Ver)
	p.SetTp(h.Tp)
	p.SetRes(h.Res)
	p.SetQid(h.Qid)
	p.SetFid(h.Fid)
	p.SetCode(h.Code)
	return
}