package spatial

import (
	"bytes"
	"context"
	"encoding/binary"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PolygonData []LineStringData

func (pd *PolygonData) isAGenericGeometryData() {}

func (pd *PolygonData) decode(data *bytes.Reader) error {
	var length uint32
	if err := binary.Read(data, byteOrder, &length); nil != err {
		return err
	}

	*pd = make([]LineStringData, length)
	for i := uint32(0); i < length; i++ {
		if err := (*pd)[i].decode(data); nil != err {
			return err
		}
	}
	return nil
}

func (pd *PolygonData) encode(data *bytes.Buffer) {
	length := uint32(len(*pd))
	binary.Write(data, byteOrder, length)
	for i := uint32(0); i < length; i++ {
		(*pd)[i].encode(data)
	}
}

type Polygon struct {
	baseGeometry
	Data PolygonData
}

func NewPolygon(srid Srid) Polygon {
	return Polygon{baseGeometry: baseGeometry{srid: srid}}
}

func (p *Polygon) decode(data *bytes.Reader, decodeSrid bool) error {
	if _, err := p.decodeHeader(data, decodeSrid, GEOMETRY_TYPE_POLYGON); nil != err {
		return err
	}
	return p.Data.decode(data)
}
func (p Polygon) encode(data *bytes.Buffer) {
	encodeHeader(data, GEOMETRY_TYPE_POLYGON)
	p.Data.encode(data)
}

func (p *Polygon) Scan(v interface{}) error {
	return geometryScan(p, v)
}

func (p Polygon) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	return geometryGormValue(p)
}
