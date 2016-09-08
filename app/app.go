package app

import (
	. "github.com/tendermint/go-common"
	"github.com/tendermint/go-merkle"
	types "github.com/tendermint/tmsp/types"
)

type Application struct {
	state merkle.Tree
	cache *Cache
}

func NewApplication() *Application {
	state := merkle.NewIAVLTree(
		0,
		nil,
	)
	return &Application{
		state: state,
		cache: CreateCache(),
	}
}

func (app *Application) Info() string {
	return Fmt("size:%v", app.state.Size())
}

func (app *Application) SetOption(key string, value string) (log string) {
	return ""
}

func (app *Application) AppendTx(tx []byte) types.Result {
	form, _ := MakeForm(tx)
	id := FormID(form)
	go app.cache.NewForm(id, form)
	app.state.Set([]byte(id), tx)
	return types.NewResultOK(nil, id)
}

func (app *Application) CheckTx(tx []byte) types.Result {
	_, err := MakeForm(tx)
	if err != nil {
		return types.NewResult(types.CodeType_InternalError, nil, err.Error())
	}
	return types.NewResultOK(nil, "")
}

func (app *Application) Commit() types.Result {
	hash := app.state.Hash()
	return types.NewResultOK(hash, "")
}

func (app *Application) Query(query []byte) types.Result {
	index, value, exists := app.state.Get(query)
	resStr := Fmt("Index=%v value=%v exists=%v", index, string(value), exists)
	return types.NewResultOK([]byte(resStr), "")
}