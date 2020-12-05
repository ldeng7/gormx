package spatial

import (
	"bytes"
	"context"
	"encoding/binary"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LineStringData []PointData

func (ld *LineStringData) isAGenericGeometryData() {}

func (ld *LineStringData) decode(data *bytes.Reader) error {
	var length uint32
	if err := binary.Read(data, byteOrder, &length); nil != err {
		return err
	}

	*ld = make([]PointData, length)
	for i := uint32(0); i < length; i++ {
		if err := (*ld)[i].decode(data); nil != err {
			return err
		}
	}
	return nil
}

func (ld *LineStringData) encode(data *bytes.Buffer) {
	length := uint32(len(*ld))
	binary.Write(data, byteOrder, length)
	for i := uint32(0); i < length; i++ {
		(*ld)[i].encode(data)
	}
}

type LineString struct {
	baseGeometry
	Data LineStringData
}

func NewLineString(srid Srid) LineString {
	return LineString{baseGeometry: baseGeometry{srid: srid}}
}

func (l *LineString) decode(data *bytes.Reader, decodeSrid bool) error {
	if _, err := l.decodeHeader(data, decodeSrid, GEOMETRY_TYPE_LINE_STRING); nil != err {
		return err
	}
	return l.Data.decode(data)
}

func (l LineString) encode(data *bytes.Buffer) {
	encodeHeader(data, GEOMETRY_TYPE_LINE_STRING)
	l.Data.encode(data)
}

func (l *LineString) Scan(v interface{}) error {
	return geometryScan(l, v)
}

func (l LineString) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	return geometryGormValue(l)
}
