package gormx

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm/schema"
)

type WindowSelectClause struct {
	Expr      string
	Name      string
	Partition string
	Order     string
	FrameSpec string
	As        string
}

func (c *WindowSelectClause) String() string {
	overClause := make([]string, 0, 4)
	if v := c.Name; 0 != len(v) {
		overClause = append(overClause, v)
	}
	if v := c.Partition; 0 != len(v) {
		overClause = append(overClause, "PARTITION BY "+v)
	}
	if v := c.Order; 0 != len(v) {
		overClause = append(overClause, "ORDER BY "+v)
	}
	if v := c.FrameSpec; 0 != len(v) {
		overClause = append(overClause, v)
	}
	return fmt.Sprintf("%s OVER (%s) AS `%s`", c.Expr, strings.Join(overClause, " "), c.As)
}

var gSelectsCache = map[reflect.Type][]string{}

func Select(v interface{}) []string {
	t := structType(v)
	if nil == t {
		return nil
	}
	if sels, ok := gSelectsCache[t]; ok {
		return sels
	}

	ns := schema.NamingStrategy{}
	sels := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag, ok := f.Tag.Lookup(tagKeySelect)
		if !ok {
			continue
		}
		m := schema.ParseTagSetting(tag, ";")

		var sel string
		switch m["TYPE"] {
		case "column":
			if column := m["COLUMN"]; 0 == len(column) {
				sel = ns.ColumnName("", f.Name)
			} else {
				sel = column
			}
		case "window":
			c := &WindowSelectClause{
				Expr:      m["EXPR"],
				Name:      m["NAME"],
				Partition: m["PARTITION"],
				Order:     m["ORDER"],
				FrameSpec: m["FRAME_SPEC"],
				As:        m["AS"],
			}
			if 0 == len(c.As) {
				c.As = ns.ColumnName("", f.Name)
			}
			sel = c.String()
		case "expr":
			as := m["AS"]
			if 0 == len(as) {
				as = ns.ColumnName("", f.Name)
			}
			sel = fmt.Sprintf("%s AS %s", m["EXPR"], as)
		default:
			continue
		}
		sels = append(sels, sel)
	}
	gSelectsCache[t] = sels
	return sels
}
