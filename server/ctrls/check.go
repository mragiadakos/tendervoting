package ctrls

import (
	"encoding/json"

	"github.com/tendermint/abci/types"
)

func (app *TVApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	tvd := TVDelivery{}
	json.Unmarshal(tx, &tvd)

	code, err := app.verifyDelivery(tvd)
	if err != nil {
		return types.ResponseCheckTx{Code: code, Log: err.Error()}
	}

	return types.ResponseCheckTx{Code: CodeTypeOK}
}
