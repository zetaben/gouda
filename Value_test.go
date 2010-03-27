package gouda_test

import (
	"gouda"
	"testing"
//	"fmt"
)

func TestIntValue(t *testing.T) {
	f:=gouda.SysInt(12).Value()
	if f.Int()!=12 {
		t.Error("Failed to recover 12")
	}
	f.SetInt(13)
	if   f.Int()!=13 {
		t.Error("Failed to recover 13")
	}
	if f.Kind()!=gouda.IntKind {
		t.Error("Not of IntKind")
	}
//	fmt.Printf("%#v\n",f)
}

func TestStringValue(t *testing.T) {
	f:=gouda.SysString("toto").Value()
	if f.String()!="toto" {
		t.Error("Failed to recover toto")
	}
	f.SetString("titi")
	if   f.String()!="titi" {
		t.Error("Failed to recover titi")
	}
	if f.Kind()!=gouda.StringKind {
		t.Error("Not of StringKind")
	}
//	fmt.Printf("%#v\n",f)
}
