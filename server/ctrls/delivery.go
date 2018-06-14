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
		err = d.ValidateVoters()
		if err != nil {
			return CodeTypeUnauthorized, err
		}
		_, err = app.state.GetElection(d.ID)
		if err == nil {
			return CodeTypeUnauthorized, errors.New("The election's ID exists.")
		}
	case POLL:
		d := tvd.GetPollDeliveryData()
		err := d.ValidateGonverment()
		if err != nil {
			return CodeTypeUnauthorized, err
		}
		_, err = app.state.GetElection(d.ElectionID)
		if err != nil {
			return CodeTypeUnauthorized, errors.New("The election's ID does not exists.")
		}
		if len(d.PollHash) == 0 {
			return CodeTypeUnauthorized, errors.New("Missing the IPFS hash for the poll.")
		}
		_, err = d.GetPollJsonFromPollHash()
		if err != nil {
			return CodeTypeUnauthorized, err
		}
		_, err = app.state.GetPoll(d.PollHash)
		if err == nil {
			return CodeTypeUnauthorized, errors.New("The poll's hash exists.")
		}
		if !app.state.IsLatestElection(d.ElectionID) {
			return CodeTypeUnauthorized, errors.New("The election's ID is not the latest.")
		}
	case VOTE:
		d := tvd.GetVoteDeliveryData()
		if len(d.PollHash) == 0 {
			return CodeTypeUnauthorized, errors.New("The poll's hash is empty.")
		}
		ps, err := app.state.GetPoll(d.PollHash)
		if err != nil {
			return CodeTypeUnauthorized, errors.New("The poll's hash does not exists.")
		}
		es, err := app.state.GetElection(ps.ElectionID)
		if err != nil {
			return CodeTypeServerError, errors.New("Could not find the election from the poll that exists.")
		}
		foundVoter := false
		for _, v := range es.Voters {
			if v == d.From {
				foundVoter = true
				break
			}
		}
		if !foundVoter {
			return CodeTypeUnauthorized, errors.New("You don't exist in the list of voters.")
		}
		_, ok := ps.Choices[d.Choice]
		if !ok {
			return CodeTypeUnauthorized, errors.New("The choice " + d.Choice + " does not exists for poll " + d.PollHash + ".")
		}
		if app.state.HasVote(d) {
			return CodeTypeUnauthorized, errors.New("You voted already for the specific poll.")
		}
		if !app.state.IsLatestPoll(d.PollHash) {
			return CodeTypeUnauthorized, errors.New("The poll's hash is not the latest.")
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
		app.state.CreateElection(d)
	case POLL:
		d := tvd.GetPollDeliveryData()
		app.state.CreatePoll(d)
	case VOTE:
		d := tvd.GetVoteDeliveryData()
		app.state.CreateVote(d)
		app.state.AddVoteToThePoll(d)
	}
	return types.ResponseDeliverTx{Code: CodeTypeOK}
}
