package memstore_test

import (
	"testing"

	"selihc.com/glaive/httpz/httpthrottle/memstore"
	"selihc.com/glaive/httpz/httpthrottle/storetest"
)

func TestMemStoreLRU(t *testing.T) {
	st, err := memstore.New(10)
	if err != nil {
		t.Fatal(err)
	}
	storetest.TestGCRAStore(t, st)
}

func TestMemStoreUnlimited(t *testing.T) {
	st, err := memstore.New(10)
	if err != nil {
		t.Fatal(err)
	}
	storetest.TestGCRAStore(t, st)
}

func BenchmarkMemStoreLRU(b *testing.B) {
	st, err := memstore.New(10)
	if err != nil {
		b.Fatal(err)
	}
	storetest.BenchmarkGCRAStore(b, st)
}

func BenchmarkMemStoreUnlimited(b *testing.B) {
	st, err := memstore.New(0)
	if err != nil {
		b.Fatal(err)
	}
	storetest.BenchmarkGCRAStore(b, st)
}
