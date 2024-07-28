package store

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestCASPathTransportFromFun(t *testing.T) {
	pathKey := CASPathTransportFromFun("lsjhtang")
	pathName := "b2e03/95f49/af5ea/07c95/073a1/fbf57/18612/615d1"
	original := "b2e0395f49af5ea07c95073a1fbf5718612615d1"
	if pathKey.PathName != pathName {
		t.Fatalf("want pathName: %s, got pathName:%s", pathName, pathKey.PathName)
	}
	if pathKey.Original != original {
		t.Fatalf("want original: %s, got original:%s", original, pathKey.Original)
	}
}

func TestStore(t *testing.T) {
	opts := Opts{
		Root:                 "root",
		PathTransportFromFun: nil,
	}
	store := NewFileStore(opts)
	key := "lsjhtang"
	buf := bytes.NewBuffer([]byte(key))
	err := store.Write(key, buf)
	if err != nil {
		t.Fatal(err)
	}
	bl := store.Has(key)
	if !bl {
		t.Fatalf("store does not have key %s", key)
	}

	r, err := store.Read(key)
	if err != nil {
		t.Fatalf("store read error %s", err)
	}
	_, err = io.Copy(os.Stdout, r)
	if err != nil {
		t.Fatalf("store read error %s", err)
	}

	err = store.RemoveAll(key)
	if err != nil {
		t.Fatal(err)
	}

	bl = store.Has(key)
	if bl {
		t.Fatalf("removell fail have key %s", key)
	}
}
