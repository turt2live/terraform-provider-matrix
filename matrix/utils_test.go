package matrix

import (
	"testing"
)

func TestUnitUtilsNilIfEmpty_nilWhenEmpty(t *testing.T) {
	r := nilIfEmptyString("")
	if r != nil {
		t.Errorf("result was not nil, got %#v, expected nil", r)
	}
}

func TestUnitUtilsNilIfEmpty_nilWhenNil(t *testing.T) {
	r := nilIfEmptyString(nil)
	if r != nil {
		t.Errorf("result was not nil, got %#v, expected nil", r)
	}
}

func TestUnitUtilsNilIfEmpty_stringWhenNotEmpty(t *testing.T) {
	r := nilIfEmptyString("hello")
	if r != "hello" {
		t.Errorf("result was not nil, got %#v, expected 'hello'", r)
	}
}

func TestUnitUtilsNilIfEmpty_objectWhenObject(t *testing.T) {
	v := struct{}{}
	r := nilIfEmptyString(v)
	if r != v {
		t.Errorf("result was not nil, got %#v, expected %#v", r, v)
	}
}
