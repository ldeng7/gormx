package spatial

import (
	"bytes"
	"context"
	"encoding/binary"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GeometryCollectionData []GenericGeometry

func (cd *GeometryCollectionData) isAGenericGeometryData() {}

func (cd *GeometryCollectionData) decode(data *bytes.Reader, srid Srid) error {
	var length uint32
	if err := binary.Read(data, byteOrder, &length); nil != err {
		return err
	}

	*cd = make([]GenericGeometry, length)
	for i := uint32(0); i < length; i++ {
		g := &((*cd)[i])
		g.srid = srid
		if err := g.decode(data, false); nil != err {
			return err
		}
	}
	return nil
}

func (cd *GeometryCollectionData) encode(data *bytes.Buffer) {
	length := uint32(len(*cd))
	binary.Write(data, byteOrder, length)
	for i := uint32(0); i < length; i++ {
		(*cd)[i].encode(data)
	}
}

type GeometryCollection struct {
	baseGeometry
	Data GeometryCollectionData
}

func NewGeometryCollection(srid Srid) GeometryCollection {
	return GeometryCollection{baseGeometry: baseGeometry{srid: srid}}
}

func (c *GeometryCollection) decode(data *bytes.Reader, decodeSrid bool) error {
	if _, err := c.decodeHeader(data, decodeSrid, GEOMETRY_TYPE_COLLECTION); nil != err {
		return err
	}
	return c.Data.decode(data, c.srid)
}

func (c GeometryCollection) encode(data *bytes.Buffer) {
	encodeHeader(data, GEOMETRY_TYPE_COLLECTION)
	c.Data.encode(data)
}

func (c *GeometryCollection) Scan(v interface{}) error {
	return geometryScan(c, v)
}

func (c GeometryCollection) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	return geometryGormValue(c)
}
