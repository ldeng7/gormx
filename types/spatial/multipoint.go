package spatial

import (
	"bytes"
	"context"
	"encoding/binary"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MultiPointData []PointData

func (mpd *MultiPointData) isAGenericGeometryData() {}

func (mpd *MultiPointData) decode(data *bytes.Reader) error {
	var length uint32
	if err := binary.Read(data, byteOrder, &length); nil != err {
		return err
	}

	*mpd = make([]PointData, length)
	var p Point
	for i := uint32(0); i < length; i++ {
		if err := p.decode(data, false); nil != err {
			return err
		}
		(*mpd)[i] = p.Data
	}
	return nil
}

func (mpd *MultiPointData) encode(data *bytes.Buffer) {
	length := uint32(len(*mpd))
	binary.Write(data, byteOrder, length)
	for i := uint32(0); i < length; i++ {
		encodeHeader(data, GEOMETRY_TYPE_POINT)
		(*mpd)[i].encode(data)
	}
}

type MultiPoint struct {
	baseGeometry
	Data MultiPointData
}

func NewMultiPoint(srid Srid) MultiPoint {
	return MultiPoint{baseGeometry: baseGeometry{srid: srid}}
}

func (mp *MultiPoint) decode(data *bytes.Reader, decodeSrid bool) error {
	if _, err := mp.decodeHeader(data, decodeSrid, GEOMETRY_TYPE_MULTI_POINT); nil != err {
		return err
	}
	return mp.Data.decode(data)
}

func (mp MultiPoint) encode(data *bytes.Buffer) {
	encodeHeader(data, GEOMETRY_TYPE_MULTI_POINT)
	mp.Data.encode(data)
}

func (mp *MultiPoint) Scan(v interface{}) error {
	return geometryScan(mp, v)
}

func (mp MultiPoint) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	return geometryGormValue(mp)
}
