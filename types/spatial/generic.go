package spatial

import (
	"bytes"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GenericGeometryData interface {
	isAGenericGeometryData()
}

type GenericGeometry struct {
	baseGeometry
	typ  GeometryType
	Data GenericGeometryData
}

func NewGenericGeometry(srid Srid) GenericGeometry {
	return GenericGeometry{baseGeometry: baseGeometry{srid: srid}}
}

func (g *GenericGeometry) Type() GeometryType {
	return g.typ
}

func (g *GenericGeometry) decode(data *bytes.Reader, decodeSrid bool) error {
	var err error
	if g.typ, err = g.decodeHeader(data, decodeSrid, GEOMETRY_TYPE_GENERIC); nil != err {
		return err
	}

	switch g.typ {
	case GEOMETRY_TYPE_POINT:
		gd := &PointData{}
		g.Data, err = gd, gd.decode(data)
	case GEOMETRY_TYPE_LINE_STRING:
		gd := &LineStringData{}
		g.Data, err = gd, gd.decode(data)
	case GEOMETRY_TYPE_POLYGON:
		gd := &PolygonData{}
		g.Data, err = gd, gd.decode(data)
	case GEOMETRY_TYPE_MULTI_POINT:
		gd := &MultiPointData{}
		g.Data, err = gd, gd.decode(data)
	case GEOMETRY_TYPE_MULTI_LINE_STRING:
		gd := &MultiLineStringData{}
		g.Data, err = gd, gd.decode(data)
	case GEOMETRY_TYPE_MULTI_POLYGON:
		gd := &MultiPolygonData{}
		g.Data, err = gd, gd.decode(data)
	case GEOMETRY_TYPE_COLLECTION:
		gd := &GeometryCollectionData{}
		g.Data, err = gd, gd.decode(data, g.srid)
	}
	return err
}

func (g GenericGeometry) encode(data *bytes.Buffer) {
	switch gd := g.Data.(type) {
	case *PointData:
		encodeHeader(data, GEOMETRY_TYPE_POINT)
		gd.encode(data)
	case *LineStringData:
		encodeHeader(data, GEOMETRY_TYPE_LINE_STRING)
		gd.encode(data)
	case *PolygonData:
		encodeHeader(data, GEOMETRY_TYPE_POLYGON)
		gd.encode(data)
	case *MultiPointData:
		encodeHeader(data, GEOMETRY_TYPE_MULTI_POINT)
		gd.encode(data)
	case *MultiLineStringData:
		encodeHeader(data, GEOMETRY_TYPE_MULTI_LINE_STRING)
		gd.encode(data)
	case *MultiPolygonData:
		encodeHeader(data, GEOMETRY_TYPE_MULTI_POLYGON)
		gd.encode(data)
	case *GeometryCollectionData:
		encodeHeader(data, GEOMETRY_TYPE_COLLECTION)
		gd.encode(data)
	}
}

func (g *GenericGeometry) Scan(v interface{}) error {
	return geometryScan(g, v)
}

func (g GenericGeometry) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	return geometryGormValue(g)
}
