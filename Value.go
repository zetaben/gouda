package gouda

import (
	"fmt"
)

const (
		EOOKind	= iota;
		NumberKind;
		StringKind;
		BooleanKind;
		NullKind;
		IntKind;
		LongKind;
)



type Value interface {
	Kind() int;
	Number() float64;
	String() string;
	Bool() bool;
	Int() int;
}

type WriteValue interface {
	SetNumber( float64);
	SetString( string);
	SetBool( bool);
	SetInt(int);
}

type _Null struct{}

func (*_Null) Kind() int		{ return NullKind }
func (*_Null) Number() float64		{ return 0 }
func (*_Null) String() string		{ return "null" }
func (*_Null) Bool() bool		{ return false }
func (*_Null) Int() int			{ return 0 }
func (*_Null)	SetNumber( float64) {};
func (*_Null)	SetString( string){};
func (*_Null)	SetBool( bool) {};
func (*_Null)	SetInt(int) {};


type SysInt int;
type _Int struct
{
	value SysInt;
	_Null
}

func (e * _Int) Kind() SysInt { return IntKind}
func (e * _Int) Int() SysInt { return e.value}
func (e * _Int) SetInt(i SysInt) { e.value=i}

func (i  SysInt ) Value()(ii * _Int){ ii=new(_Int);ii.SetInt(i); return }


type SysString string;
type _String struct {
	value SysString;
	_Null
}

func (e * _String) Kind() int { return StringKind}
func (e * _String) SetString(s SysString) { e.value=s}
func (e * _String) String() SysString { return e.value}

func (i  SysString ) Value()(ii * _String){ ii=new(_String);ii.SetString(i); return }


func (v * _Int) GoString() string {

	return "gouda.SysInt("+fmt.Sprint(v.Int())+").Value()"
}
func (v * _String) GoString() string {
	return fmt.Sprintf("gouda.SysString(%#v).Value()",v.String())
}

func (v * _Null) GoString() string {
	return "*Hem*"
}
