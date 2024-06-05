package hybrid

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestWarpSlice(t *testing.T) {
	a := []string{"a", "b", "c"}
	w := NewWarp(a)
	b, err := json.Marshal(&w)
	if err != nil {
		t.Fatal(err)
	}
	b2 := "{\"value\":[\"a\",\"b\",\"c\"]}"
	if string(b) != b2 {
		t.Fatalf("got %s, want %s", b, b2)
	}
	w2 := NewWarpEmpty[[]string]()
	err = json.Unmarshal(b, &w2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, w2.Value) {
		t.Errorf("got %v, want %v", w2.Value, a)
	}
}

func TestWarpMap(t *testing.T) {
	a := map[string]string{"a": "b", "c": "d"}
	w := NewWarp(a)
	b, err := json.Marshal(&w)
	if err != nil {
		t.Fatal(err)
	}
	b2 := "{\"value\":{\"a\":\"b\",\"c\":\"d\"}}"
	if string(b) != b2 {
		t.Fatalf("got %s, want %s", b, b2)
	}
	w2 := NewWarpEmpty[map[string]string]()
	err = json.Unmarshal(b, &w2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, w2.Value) {
		t.Errorf("got %v, want %v", w2.Value, a)
	}
}

func TestWarpStruct(t *testing.T) {
	type T struct {
		W string `json:"w"`
	}
	a := T{W: "a"}
	w := NewWarp(a)
	b, err := json.Marshal(&w)
	if err != nil {
		t.Fatal(err)
	}
	b2 := "{\"value\":{\"w\":\"a\"}}"
	if string(b) != b2 {
		t.Fatalf("got %s, want %s", b, b2)
	}
	w2 := NewWarpEmpty[T]()
	err = json.Unmarshal(b, &w2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, w2.Value) {
		t.Errorf("got %v, want %v", w2.Value, a)
	}
}

func TestWarpPrt(t *testing.T) {
	type T struct {
		W string `json:"w"`
	}
	a := &T{W: "a"}
	w := NewWarp(a)
	b, err := json.Marshal(&w)
	if err != nil {
		t.Fatal(err)
	}
	b2 := "{\"value\":{\"w\":\"a\"}}"
	if string(b) != b2 {
		t.Fatalf("got %s, want %s", b, b2)
	}
	w2 := NewWarpEmpty[*T]()
	err = json.Unmarshal(b, &w2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, w2.Value) {
		t.Errorf("got %v, want %v", w2.Value, a)
	}
}

func TestWarpPrtNil(t *testing.T) {
	type T struct {
		W string `json:"w"`
	}
	var a *T = nil
	w := NewWarp(a)
	b, err := json.Marshal(&w)
	if err != nil {
		t.Fatal(err)
	}
	b2 := "{\"value\":null}"
	if string(b) != b2 {
		t.Fatalf("got %s, want %s", b, b2)
	}
	w2 := NewWarpEmpty[*T]()
	err = json.Unmarshal(b, &w2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, w2.Value) {
		t.Errorf("got %v, want %v", w2.Value, a)
	}
}

func TestWarpSimple(t *testing.T) {
	a := 12
	w := NewWarp(a)
	b, err := json.Marshal(&w)
	if err != nil {
		t.Fatal(err)
	}
	b2 := "{\"value\":12}"
	if string(b) != b2 {
		t.Fatalf("got %s, want %s", b, b2)
	}
	w2 := NewWarpEmpty[int]()
	err = json.Unmarshal(b, &w2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, w2.Value) {
		t.Errorf("got %v, want %v", w2.Value, a)
	}
}
