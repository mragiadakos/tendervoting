package ctrls

import (
	"encoding/json"

	"github.com/tendermint/abci/types"
)

func (tva *TVApplication) queryListElections() ListElectionQuery {
	list := ListElectionQuery{}
	latest, err := tva.state.GetLatestElection()
	if err != nil {
		return list
	}
	for _, v := range tva.state.GetElections() {
		item := ItemElectionQuery{}
		item.ElectionQuery = v
		if v.ID == latest {
			item.Latest = true
		}
		list = append(list, item)
	}
	return list
}

func (tva *TVApplication) queryListPolls() ListPollQuery {
	list := ListPollQuery{}
	latest, err := tva.state.GetLatestPoll()
	if err != nil {
		return list
	}
	for _, v := range tva.state.GetPolls() {
		item := ItemPollQuery{}
		item.PollQuery = v
		if v.PollHash == latest {
			item.Latest = true
		}
		list = append(list, item)
	}
	return list
}

func (tva *TVApplication) queryVotes(pollHash string) (*PollVotesQuery, error) {
	ps, err := tva.state.GetPoll(pollHash)
	if err != nil {
		return nil, err
	}
	pvq := new(PollVotesQuery)
	pvq.Choices = ps.Choices
	pvq.NumberOfVotes = len(ps.VotedAlready)
	return pvq, nil
}

func (tva *TVApplication) Query(qreq types.RequestQuery) types.ResponseQuery {
	switch qreq.Path {
	case "/elections":
		list := tva.queryListElections()
		b, _ := json.Marshal(list)
		resp := types.ResponseQuery{Code: CodeTypeOK, Value: b}
		return resp
	case "/elections/latest":
		list := tva.queryListElections()
		for _, v := range list {
			if v.Latest {
				b, _ := json.Marshal(v)
				resp := types.ResponseQuery{Code: CodeTypeOK, Value: b}
				return resp
			}
		}
	case "/polls":
		list := tva.queryListPolls()
		b, _ := json.Marshal(list)
		resp := types.ResponseQuery{Code: CodeTypeOK, Value: b}
		return resp
	case "/polls/latest":
		list := tva.queryListPolls()
		for _, v := range list {
			if v.Latest {
				b, _ := json.Marshal(v)
				resp := types.ResponseQuery{Code: CodeTypeOK, Value: b}
				return resp
			}
		}
	case "/votes":
		pq := PollQuery{}
		err := json.Unmarshal(qreq.Data, &pq)
		if err != nil {
			resp := types.ResponseQuery{Code: CodeTypeEncodingError, Log: "The JSON for the poll hash is incorrect."}
			return resp
		}
		pvq, err := tva.queryVotes(pq.PollHash)
		if err != nil {
			resp := types.ResponseQuery{Code: CodeTypeUnauthorized, Log: err.Error()}
			return resp
		}
		b, _ := json.Marshal(pvq)
		resp := types.ResponseQuery{Code: CodeTypeOK, Value: b}
		return resp
	}

	resp := types.ResponseQuery{Code: CodeTypeOK}
	return resp
}
