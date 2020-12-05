package spatial

import (
	"bytes"
	"context"
	"encoding/binary"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MultiLineStringData []LineStringData

func (mld *MultiLineStringData) isAGenericGeometryData() {}

func (mld *MultiLineStringData) decode(data *bytes.Reader) error {
	var length uint32
	if err := binary.Read(data, byteOrder, &length); nil != err {
		return err
	}

	*mld = make([]LineStringData, length)
	var l LineString
	for i := uint32(0); i < length; i++ {
		if err := l.decode(data, false); nil != err {
			return err
		}
		(*mld)[i] = l.Data
	}
	return nil
}

func (mld *MultiLineStringData) encode(data *bytes.Buffer) {
	length := uint32(len(*mld))
	binary.Write(data, byteOrder, length)
	for i := uint32(0); i < length; i++ {
		encodeHeader(data, GEOMETRY_TYPE_LINE_STRING)
		(*mld)[i].encode(data)
	}
}

type MultiLineString struct {
	baseGeometry
	Data MultiLineStringData
}

func NewMultiLineString(srid Srid) MultiLineString {
	return MultiLineString{baseGeometry: baseGeometry{srid: srid}}
}

func (ml *MultiLineString) decode(data *bytes.Reader, decodeSrid bool) error {
	if _, err := ml.decodeHeader(data, decodeSrid, GEOMETRY_TYPE_MULTI_LINE_STRING); nil != err {
		return err
	}
	return ml.Data.decode(data)
}

func (ml MultiLineString) encode(data *bytes.Buffer) {
	encodeHeader(data, GEOMETRY_TYPE_MULTI_LINE_STRING)
	ml.Data.encode(data)
}

func (ml *MultiLineString) Scan(v interface{}) error {
	return geometryScan(ml, v)
}

func (ml MultiLineString) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	return geometryGormValue(ml)
}
