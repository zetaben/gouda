package gouda

import (
	"container/vector"
	"fmt"
	"reflect"
	"strings"
)

type RequestKind int
type OperandKind int

type Condition struct {
	field   string
	operand OperandKind
	value   Value
	or      vector.Vector
}

const (
	NULL OperandKind = iota
	EQUAL
	NOTEQUAL
	ISNOTNULL
	ISNULL
	GREATER
	LOWER
	GREATEROREQUAL
	LOWEROREQUAL
	OR
)
const (
	SELECT RequestKind = iota
	UPDATE
	INSERT
	COUNT
	DELETE
)

type Relation struct {
	conditions       vector.Vector
	table            string
	limit_offset     int
	limit_count      int
	order_field      vector.StringVector
	order_direction  vector.StringVector
	kind             RequestKind
	attributes       vector.StringVector
	values           map[string]Value
	id               int
	identifier_field string
}

func (r *Relation) Count(fields []string) *Relation {
	for _, s := range fields {
		r.attributes.Push(s)
	}
	r.kind = COUNT
	return r
}

func (r *Relation) Delete() *Relation {
	r.kind = DELETE
	return r
}

func (r *Relation) Insert(mp map[string]Value) *Relation {
	r.kind = INSERT
	r.values = mp
	return r
}

func (r *Relation) Update(mp map[string]Value, identifier string, id int) *Relation {
	r.kind = UPDATE
	r.values = mp
	r.id = id
	r.identifier_field = identifier
	return r
}

func (r *Relation) Where(s *Condition) *Relation {

	r.conditions.Push(s)

	return r

}

func (r *Relation) Table(t string) *Relation {
	r.table = t
	return r
}

func (r *Relation) Limit(offset, count int) *Relation {
	r.limit_offset = offset
	r.limit_count = count
	return r
}

func (r *Relation) First() *Relation {
	r.Limit(0, 1)
	return r
}

func (r *Relation) Order(field, direction string) *Relation {
	r.order_field.Push(field)
	r.order_direction.Push(direction)
	return r
}

func (r *Relation) String() string {

	s := " Conditions : \n\t ("

	for _, ss := range r.conditions {
		s += ss.(*Condition).String()
		if ss != r.conditions.Last() {
			s += ", "
		}
	}
	s += ") \n"
	if r.order_field.Len() > 0 {
		s += "Order :\n\t ("
		for i, ss := range r.order_field {
			s += ss
			s += " "
			s += r.order_direction[i]
			if ss != r.order_field.Last() {
				s += ", "
			}
		}
		s += ") \n"
	}
	if r.limit_count > 0 {
		s += "Offset : " + fmt.Sprint(r.limit_offset) + "\n"
		s += "Count : " + fmt.Sprint(r.limit_count) + "\n"
	}
	return "Une relation" + s
}

func From(t interface{}) (r *Relation) { return NewRelation(t) }

func NewRelation(t interface{}) (r *Relation) {
	r = new(Relation)
	r.kind = SELECT
	switch typ := t.(type) {
	case string:
		r.Table(t.(string))
	default:
		tab := strings.Split(reflect.Typeof(t).String(), ".", 0)
		tablename := strings.ToLower(tab[len(tab)-1])
		r.Table(tablename)
	}
	return
}


/** Condition **/

func F(field string) (c *Condition) {
	c = new(Condition)
	c.operand = NULL
	c.field = strings.ToLower(field)
	return
}

func (c *Condition) Value(v interface{}) *Condition {
	switch v.(type) {
	case int:
		c.value = SysInt(v.(int)).Value()
	case string:
		c.value = SysString(v.(string)).Value()
	}
	return c
}

func (c *Condition) Eq(v interface{}) *Condition {
	c.operand = EQUAL
	return c.Value(v)
}

func (c *Condition) NEq(v interface{}) *Condition {
	c.operand = NOTEQUAL
	return c.Value(v)
}

func (c *Condition) Gt(v interface{}) *Condition {
	c.operand = GREATER
	return c.Value(v)
}

func (c *Condition) Lt(v interface{}) *Condition {
	c.operand = LOWER
	return c.Value(v)
}

func (c *Condition) GtEq(v interface{}) *Condition {
	c.operand = GREATEROREQUAL
	return c.Value(v)
}

func (c *Condition) LtEq(v interface{}) *Condition {
	c.operand = LOWEROREQUAL
	return c.Value(v)
}


func (c *Condition) IsNotNull() *Condition {
	c.operand = ISNOTNULL
	return c
}

func (c *Condition) IsNull() *Condition {
	c.operand = ISNULL
	return c
}

func (c *Condition) Or(c2 *Condition) *Condition {
	if c.operand != OR {
		orc := new(Condition)
		orc.operand = OR
		orc.or.Push(c)
		c = orc
	}
	c.or.Push(c2)
	return c
}

func (c *Condition) String() string {
	if c.operand == OR {
		ret := " ( "
		for i := range c.or {
			ret += c.or.At(i).(*Condition).String()
			if i != c.or.Len()-1 {
				ret += " ) OR ("
			}
		}
		ret += ")"
		return ret
	}
	ret := c.field
	switch c.operand {
	case EQUAL:
		ret += " = "
	case NOTEQUAL:
		ret += " != "
	case ISNOTNULL:
		ret += " IS NOT NULL "
	case ISNULL:
		ret += " IS NULL "
	case LOWER:
		ret += " < "
	case GREATER:
		ret += " > "
	case LOWEROREQUAL:
		ret += " <= "
	case GREATEROREQUAL:
		ret += " >= "
	default:
		ret += " NO SE OP"
	}
	if c.operand != ISNOTNULL && c.operand != ISNULL {
		switch c.value.Kind() {
		case IntKind:
			ret += fmt.Sprint(int(c.value.Int()))
		case StringKind:
			ret += fmt.Sprint(string(c.value.String()))
		}
	}
	return ret
}
