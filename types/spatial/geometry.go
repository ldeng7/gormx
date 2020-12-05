package spatial

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/ldeng7/gormx"
	"gorm.io/gorm/clause"
)

type GeometryType uint32
type Srid uint32

const (
	GEOMETRY_TYPE_GENERIC GeometryType = iota
	GEOMETRY_TYPE_POINT
	GEOMETRY_TYPE_LINE_STRING
	GEOMETRY_TYPE_POLYGON
	GEOMETRY_TYPE_MULTI_POINT
	GEOMETRY_TYPE_MULTI_LINE_STRING
	GEOMETRY_TYPE_MULTI_POLYGON
	GEOMETRY_TYPE_COLLECTION
)

const (
	SRID_WGS_84 Srid = 4326
)

var (
	byteOrder                   = binary.LittleEndian
	geometryInstantiableMinType = GEOMETRY_TYPE_POINT
	geometryInstantiableMaxType = GEOMETRY_TYPE_COLLECTION
)

type baseGeometry struct {
	srid Srid
}

func (g baseGeometry) Srid() Srid {
	return g.srid
}

func (g baseGeometry) GormDataType() string {
	return "geometry"
}

func (g *baseGeometry) decodeHeader(data *bytes.Reader,
	decodeSrid bool, expectedType GeometryType) (GeometryType, error) {
	if decodeSrid {
		if err := binary.Read(data, byteOrder, &g.srid); nil != err {
			return 0, err
		}
	}

	if _, err := data.ReadByte(); nil != err {
		return 0, err
	}

	var typ GeometryType
	if err := binary.Read(data, byteOrder, &typ); nil != err {
		return 0, err
	} else if (typ < geometryInstantiableMinType && typ > geometryInstantiableMaxType) ||
		(GEOMETRY_TYPE_GENERIC != expectedType && typ != expectedType) {
		return 0, errors.New("unexpected geometry type")
	}
	return typ, nil
}

func encodeHeader(data *bytes.Buffer, typ GeometryType) {
	data.WriteByte(0x01)
	binary.Write(data, byteOrder, typ)
}

type geometryDecoder interface {
	decode(data *bytes.Reader, decodeSrid bool) error
}

type geometryEncoder interface {
	encode(data *bytes.Buffer)
	Srid() Srid
}

func geometryScan(g geometryDecoder, v interface{}) error {
	bs, ok := v.([]byte)
	if !ok {
		return errors.New("Failed to scan a geometry field: invalid data type")
	}

	data := bytes.NewReader(bs)
	if err := g.decode(data, true); nil != err {
		return errors.New("Failed to scan a geometry field: " + err.Error())
	}
	return nil
}

func geometryGormValue(g geometryEncoder) clause.Expr {
	data := bytes.NewBuffer(make([]byte, 0, 25))
	binary.Write(data, byteOrder, g.Srid())
	g.encode(data)
	return clause.Expr{SQL: gormx.BytesToSql(data.Bytes())}
}
