package mocks

import (
	"context"

	"cosmossdk.io/core/store"
	"cosmossdk.io/errors"
	"cosmossdk.io/store/types"
	db "github.com/cosmos/cosmos-db"
)

var ErrorStoreAccess = errors.New("store", 1, "error accessing store")

type FailingMethod string

type StoreService struct {
	failingMethod FailingMethod
	original      types.KVStore
}

type testStore struct {
	db            db.DB
	failingMethod FailingMethod
	original      types.KVStore
}

type contextStoreKey struct{}

const (
	Get             FailingMethod = "get"
	Has             FailingMethod = "has"
	Set             FailingMethod = "set"
	Delete          FailingMethod = "delete"
	Iterator        FailingMethod = "iterator"
	ReverseIterator FailingMethod = "reverseIterator"
)

// FailingStore returns a store.KVStoreService that can be used to test specific errors within collections.
func FailingStore(failingMethod FailingMethod, original types.KVStore) *StoreService {
	return &StoreService{failingMethod, original}
}

func (s StoreService) OpenKVStore(_ context.Context) store.KVStore {
	return testStore{
		failingMethod: s.failingMethod,
		original:      s.original,
	}
}

func (s StoreService) NewStoreContext() context.Context {
	kv := db.NewMemDB()
	return context.WithValue(context.Background(), contextStoreKey{}, &testStore{kv, s.failingMethod, s.original})
}

func (t testStore) Get(key []byte) ([]byte, error) {
	if t.failingMethod == Get {
		return nil, ErrorStoreAccess
	}
	return t.original.Get(key), nil
}

func (t testStore) Has(key []byte) (bool, error) {
	if t.failingMethod == Has {
		return false, ErrorStoreAccess
	}
	return t.original.Has(key), nil
}

func (t testStore) Set(key, value []byte) error {
	if t.failingMethod == Set {
		return ErrorStoreAccess
	}
	t.original.Set(key, value)
	return nil
}

func (t testStore) Delete(key []byte) error {
	if t.failingMethod == Delete {
		return ErrorStoreAccess
	}
	t.original.Delete(key)
	return nil
}

func (t testStore) Iterator(start, end []byte) (store.Iterator, error) {
	if t.failingMethod == Iterator {
		return nil, ErrorStoreAccess
	}
	return t.original.Iterator(start, end), nil
}

func (t testStore) ReverseIterator(start, end []byte) (store.Iterator, error) {
	if t.failingMethod == ReverseIterator {
		return nil, ErrorStoreAccess
	}
	return t.original.ReverseIterator(start, end), nil
}

var _ store.KVStore = testStore{}
