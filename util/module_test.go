package util

import (
	"bytes"
	"testing"

	"github.com/kosmosJS/engine"
	"github.com/kosmosJS/engine-node/require"
)

func TestUtil_Format(t *testing.T) {
	vm := engine.New()
	util := New(vm)

	var b bytes.Buffer
	util.Format(&b, "Test: %% %ะด %s %d, %j", vm.ToValue("string"), vm.ToValue(42), vm.NewObject())

	if res := b.String(); res != "Test: % %ะด string 42, {}" {
		t.Fatalf("Unexpected result: '%s'", res)
	}
}

func TestUtil_Format_NoArgs(t *testing.T) {
	vm := engine.New()
	util := New(vm)

	var b bytes.Buffer
	util.Format(&b, "Test: %s %d, %j")

	if res := b.String(); res != "Test: %s %d, %j" {
		t.Fatalf("Unexpected result: '%s'", res)
	}
}

func TestUtil_Format_LessArgs(t *testing.T) {
	vm := engine.New()
	util := New(vm)

	var b bytes.Buffer
	util.Format(&b, "Test: %s %d, %j", vm.ToValue("string"), vm.ToValue(42))

	if res := b.String(); res != "Test: string 42, %j" {
		t.Fatalf("Unexpected result: '%s'", res)
	}
}

func TestUtil_Format_MoreArgs(t *testing.T) {
	vm := engine.New()
	util := New(vm)

	var b bytes.Buffer
	util.Format(&b, "Test: %s %d, %j", vm.ToValue("string"), vm.ToValue(42), vm.NewObject(), vm.ToValue(42.42))

	if res := b.String(); res != "Test: string 42, {} 42.42" {
		t.Fatalf("Unexpected result: '%s'", res)
	}
}

func TestJSNoArgs(t *testing.T) {
	vm := engine.New()
	new(require.Registry).Enable(vm)

	if util, ok := require.Require(vm, "util").(*engine.Object); ok {
		if format, ok := engine.AssertFunction(util.Get("format")); ok {
			res, err := format(util)
			if err != nil {
				t.Fatal(err)
			}
			if v := res.Export(); v != "" {
				t.Fatalf("Unexpected result: %v", v)
			}
		}
	}
}
