package ctrls

import (
	"encoding/json"
	"errors"

	"github.com/tendermint/abci/types"
)

func (app *TVApplication) verifyDelivery(tvd TVDelivery) (uint32, error) {
	ver, err := tvd.VerifySignature()
	if err != nil {
		return CodeTypeEncodingError, err
	}
	if !ver {
		return CodeTypeUnauthorized, errors.New("The signature does not verify the data.")
	}

	switch tvd.Type {
	case ELECTION:
		d := tvd.GetElectionDeliveryData()
		err := d.ValidateGonverment()
		if err != nil {
			return CodeTypeUnauthorized, err
		}
		err = d.ValidateTime()
		if err != nil {
			return CodeTypeUnauthorized, err
		}
		err = d.ValidateVoters()
		if err != nil {
			return CodeTypeUnauthorized, err
		}
	}
	return CodeTypeOK, nil
}

func (app *TVApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {
	tvd := TVDelivery{}
	json.Unmarshal(tx, &tvd)

	code, err := app.verifyDelivery(tvd)
	if err != nil {
		return types.ResponseDeliverTx{Code: code, Log: err.Error()}
	}

	switch tvd.Type {
	case ELECTION:
		d := tvd.GetElectionDeliveryData()
		_, err := app.state.GetElection(d.ID)
		if err == nil {
			return types.ResponseDeliverTx{Code: CodeTypeUnauthorized, Log: "The election's ID exists"}
		}
		app.state.CreateElection(d)

	}
	return types.ResponseDeliverTx{Code: CodeTypeOK}
}
