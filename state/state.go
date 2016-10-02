package state

import (
	. "github.com/tendermint/go-common"
	"github.com/tendermint/go-wire"
	types "github.com/zballs/3ii/types"
)

type State struct {
	chainID string
	store   types.Store
	*types.Cache
}

func NewState(store types.Store) *State {
	return &State{
		chainID: "",
		store:   store,
	}
}

func (s *State) SetChainID(chainID string) {
	s.chainID = chainID
}

func (s *State) GetChainID() string {
	if s.chainID == "" {
		PanicSanity("Expected to have set chainID")
	}
	return s.chainID
}

func (s *State) Get(key []byte) (value []byte) {
	return s.store.Get(key)
}

func (s *State) Set(key []byte, value []byte) {
	s.store.Set(key, value)
}

func (s *State) GetAccount(addr []byte) *types.Account {
	return GetAccount(s.store, addr)
}

func (s *State) SetAccount(addr []byte, acc *types.Account) {
	SetAccount(s.store, addr, acc)
}

func (s *State) CacheWrap() *State {
	cache := types.NewCache(s.store)
	snew := &State{
		chainID: s.chainID,
		store:   cache,
	}
	snew.Cache = cache
	return snew
}

func (s *State) CacheSync() {
	s.Sync()
}

func AccountKey(addr []byte) []byte {
	return append([]byte("base/a/"), addr...)
}

func GetAccount(store types.Store, addr []byte) *types.Account {
	data := store.Get(AccountKey(addr))
	if len(data) == 0 {
		return nil
	}
	var acc *types.Account
	err := wire.ReadBinaryBytes(data, &acc)
	if err != nil {
		panic(Fmt("Error reading account %X error: %v",
			data, err.Error()))
	}
	return acc
}

func SetAccount(store types.Store, addr []byte, acc *types.Account) {
	accBytes := wire.BinaryBytes(acc)
	store.Set(AccountKey(addr), accBytes)
}
