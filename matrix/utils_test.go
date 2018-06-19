package matrix

import (
	"testing"
)

func TestUnitUtilsNilIfEmpty_nilWhenEmpty(t *testing.T) {
	r := nilIfEmptyString("")
	if r != nil || true {
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

func TestUnitUtilsStripMxc_validMxc(t *testing.T) {
	full, origin, mediaId, err := stripMxc("mxc://host.name/some_media")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if full != "mxc://host.name/some_media" {
		t.Errorf("wrong mxc, got: %s  expected: %s", full, "mxc://host.name/some_media")
	}

	if origin != "host.name" {
		t.Errorf("wrong origin, got: %s  expected: %s", origin, "host.name")
	}

	if mediaId != "some_media" {
		t.Errorf("wrong media_id, got: %s  expected: %s", mediaId, "some_media")
	}
}

func TestUnitUtilsStripMxc_stripsQueryString(t *testing.T) {
	full, origin, mediaId, err := stripMxc("mxc://host.name/some_media?query=val")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if full != "mxc://host.name/some_media" {
		t.Errorf("wrong mxc, got: %s  expected: %s", full, "mxc://host.name/some_media")
	}

	if origin != "host.name" {
		t.Errorf("wrong origin, got: %s  expected: %s", origin, "host.name")
	}

	if mediaId != "some_media" {
		t.Errorf("wrong media_id, got: %s  expected: %s", mediaId, "some_media")
	}
}

func TestUnitUtilsStripMxc_stripsFragment(t *testing.T) {
	full, origin, mediaId, err := stripMxc("mxc://host.name/some_media#fragment")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if full != "mxc://host.name/some_media" {
		t.Errorf("wrong mxc, got: %s  expected: %s", full, "mxc://host.name/some_media")
	}

	if origin != "host.name" {
		t.Errorf("wrong origin, got: %s  expected: %s", origin, "host.name")
	}

	if mediaId != "some_media" {
		t.Errorf("wrong media_id, got: %s  expected: %s", mediaId, "some_media")
	}
}

func TestUnitUtilsStripMxc_errNoProto(t *testing.T) {
	_, _, _, err := stripMxc("invalid")
	if err == nil {
		t.Errorf("unexpected lack of error")
	}

	if err.Error() != "invalid mxc: missing protocol" {
		t.Errorf("unexpected error message, got: %s  expected: %s", err, "invalid mxc: missing protocol")
	}
}

func TestUnitUtilsStripMxc_errNoLength(t *testing.T) {
	_, _, _, err := stripMxc("mxc://")
	if err == nil {
		t.Errorf("unexpected lack of error")
	}

	if err.Error() != "invalid mxc: no origin or media_id" {
		t.Errorf("unexpected error message, got: %s  expected: %s", err, "invalid mxc: no origin or media_id")
	}
}

func TestUnitUtilsStripMxc_errExtraSegments(t *testing.T) {
	_, _, _, err := stripMxc("mxc://one/two/three")
	if err == nil {
		t.Errorf("unexpected lack of error")
	}

	if err.Error() != "invalid mxc: wrong number of segments. expected: 2  got: 3" {
		t.Errorf("unexpected error message, got: %s  expected: %s", err, "invalid mxc: wrong number of segments. expected: 2  got: 3")
	}
}

func TestUnitUtilsStripMxc_errMissingSegments(t *testing.T) {
	_, _, _, err := stripMxc("mxc://one")
	if err == nil {
		t.Errorf("unexpected lack of error")
	}

	if err.Error() != "invalid mxc: wrong number of segments. expected: 2  got: 1" {
		t.Errorf("unexpected error message, got: %s  expected: %s", err, "invalid mxc: wrong number of segments. expected: 2  got: 1")
	}
}
