package spatial

import (
	"bytes"
	"context"
	"encoding/binary"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MultiPolygonData []PolygonData

func (mpd *MultiPolygonData) isAGenericGeometryData() {}

func (mpd *MultiPolygonData) decode(data *bytes.Reader) error {
	var length uint32
	if err := binary.Read(data, byteOrder, &length); nil != err {
		return err
	}

	*mpd = make([]PolygonData, length)
	var p Polygon
	for i := uint32(0); i < length; i++ {
		if err := p.decode(data, false); nil != err {
			return err
		}
		(*mpd)[i] = p.Data
	}
	return nil
}

func (mpd *MultiPolygonData) encode(data *bytes.Buffer) {
	length := uint32(len(*mpd))
	binary.Write(data, byteOrder, length)
	for i := uint32(0); i < length; i++ {
		encodeHeader(data, GEOMETRY_TYPE_POLYGON)
		(*mpd)[i].encode(data)
	}
}

type MultiPolygon struct {
	baseGeometry
	Data MultiPolygonData
}

func NewMultiPolygon(srid Srid) MultiPolygon {
	return MultiPolygon{baseGeometry: baseGeometry{srid: srid}}
}

func (mp *MultiPolygon) decode(data *bytes.Reader, decodeSrid bool) error {
	if _, err := mp.decodeHeader(data, decodeSrid, GEOMETRY_TYPE_MULTI_POLYGON); nil != err {
		return err
	}
	return mp.Data.decode(data)
}

func (mp MultiPolygon) encode(data *bytes.Buffer) {
	encodeHeader(data, GEOMETRY_TYPE_MULTI_POLYGON)
	mp.Data.encode(data)
}

func (mp *MultiPolygon) Scan(v interface{}) error {
	return geometryScan(mp, v)
}

func (mp MultiPolygon) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	return geometryGormValue(mp)
}
