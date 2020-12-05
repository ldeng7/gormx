package spatial

import (
	"bytes"
	"context"
	"encoding/binary"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PointData struct {
	X, Y float64
}

func (pd *PointData) isAGenericGeometryData() {}

func (pd *PointData) decode(data *bytes.Reader) error {
	if err := binary.Read(data, byteOrder, &pd.X); nil != err {
		return err
	}
	return binary.Read(data, byteOrder, &pd.Y)
}

func (pd *PointData) encode(data *bytes.Buffer) {
	binary.Write(data, byteOrder, pd.X)
	binary.Write(data, byteOrder, pd.Y)
}

type Point struct {
	baseGeometry
	Data PointData
}

func NewPoint(srid Srid) Point {
	return Point{baseGeometry: baseGeometry{srid: srid}}
}

func (p *Point) decode(data *bytes.Reader, decodeSrid bool) error {
	if _, err := p.decodeHeader(data, decodeSrid, GEOMETRY_TYPE_POINT); nil != err {
		return err
	}
	return p.Data.decode(data)
}

func (p Point) encode(data *bytes.Buffer) {
	encodeHeader(data, GEOMETRY_TYPE_POINT)
	p.Data.encode(data)
}

func (p *Point) Scan(v interface{}) error {
	return geometryScan(p, v)
}

func (p Point) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	return geometryGormValue(p)
}
