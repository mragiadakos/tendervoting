package ctrls

import (
	"encoding/json"
	"errors"

	dbm "github.com/tendermint/tmlibs/db"
)

var (
	stateKey            = []byte("stateKey")
	electionKey         = []byte("election:")
	pollKey             = []byte("poll:")
	voteKey             = []byte("vote:")
	currentElectionsKey = []byte("currentElections")
	currentPollsKey     = []byte("currentPolls")
	latestElectionKey   = []byte("latestElection")
	latestPollKey       = []byte("latestPoll")
)

func prefixElection(uuid string) []byte {
	b := []byte(uuid)
	return append(electionKey, b...)
}

func prefixPoll(hash string) []byte {
	b := []byte(hash)
	return append(pollKey, b...)
}

func prefixVote(vd VoteDeliveryData) []byte {
	b := []byte(vd.From + "-" + vd.PollHash)
	return append(voteKey, b...)
}

type State struct {
	db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

type ElectionState struct {
	ID     string
	Voters []string
}

func (s *State) GetElection(uuid string) (*ElectionState, error) {
	has := s.db.Has(prefixElection(uuid))
	if !has {
		return nil, errors.New("Could not find the election " + uuid + ".")
	}
	b := s.db.Get(prefixElection(uuid))
	es := ElectionState{}
	err := json.Unmarshal(b, &es)
	if err != nil {
		return nil, errors.New("The election " + uuid + " didnt have a correct json format: " + err.Error())
	}
	return &es, nil
}

func (s *State) CreateElection(ed ElectionDeliveryData) {
	es := ElectionState{}
	es.ID = ed.ID
	es.Voters = ed.Voters
	b, _ := json.Marshal(es)
	s.db.Set(prefixElection(es.ID), b)
	s.db.Set(latestElectionKey, []byte(es.ID))

	curElsB := s.db.Get(currentElectionsKey)
	curEls := []ElectionQuery{}
	json.Unmarshal(curElsB, &curEls)
	curEls = append(curEls, ElectionQuery{ID: ed.ID, NumberOfVoters: len(ed.Voters)})
	curElsBRes, _ := json.Marshal(curEls)
	s.db.Set(currentElectionsKey, curElsBRes)
}

func (s *State) GetElections() []ElectionQuery {
	curElsB := s.db.Get(currentElectionsKey)
	curEls := []ElectionQuery{}
	json.Unmarshal(curElsB, &curEls)
	return curEls
}

func (s *State) GetLatestElection() (string, error) {
	has := s.db.Has(latestElectionKey)
	if !has {
		return "", errors.New("There is not any latest election.")
	}
	uuid := s.db.Get(latestElectionKey)
	return string(uuid), nil
}
func (s *State) IsLatestElection(id string) bool {
	uuid := s.db.Get(latestElectionKey)
	return string(uuid) == id
}

type PollState struct {
	ElectionID   string
	PollHash     string
	VotedAlready []string
	Choices      map[string]int
}

func (s *State) GetPoll(hash string) (*PollState, error) {
	has := s.db.Has(prefixPoll(hash))
	if !has {
		return nil, errors.New("Could not find the poll " + hash + ".")
	}
	b := s.db.Get(prefixPoll(hash))
	ps := PollState{}
	err := json.Unmarshal(b, &ps)
	if err != nil {
		return nil, errors.New("The poll " + hash + " didnt have a correct json format: " + err.Error())
	}
	return &ps, nil
}

func (s *State) HasPoll(hash string) bool {
	return s.db.Has(prefixPoll(hash))
}

func (s *State) AddVoteToThePoll(vd VoteDeliveryData) error {
	ps, err := s.GetPoll(vd.PollHash)
	if err != nil {
		return err
	}
	ps.VotedAlready = append(ps.VotedAlready, vd.From)
	_, ok := ps.Choices[vd.Choice]
	if !ok {
		ps.Choices[vd.Choice] = 1
	} else {
		ps.Choices[vd.Choice] += 1
	}

	b, _ := json.Marshal(ps)
	s.db.Set(prefixPoll(vd.PollHash), b)
	return nil
}

func (s *State) CreatePoll(pd PollDeliveryData) {
	ps := PollState{}
	ps.PollHash = pd.PollHash
	ps.ElectionID = pd.ElectionID
	ps.VotedAlready = []string{}
	ps.Choices = map[string]int{}
	pj, _ := pd.GetPollJsonFromPollHash()
	for k, _ := range pj.Choices {
		ps.Choices[k] = 0
	}
	b, _ := json.Marshal(ps)
	s.db.Set(prefixPoll(ps.PollHash), b)
	s.db.Set(latestPollKey, []byte(ps.PollHash))

	curPlsB := s.db.Get(currentPollsKey)
	curPls := []PollQuery{}
	json.Unmarshal(curPlsB, &curPls)
	curPls = append(curPls, PollQuery{PollHash: ps.PollHash})
	curPlsBRes, _ := json.Marshal(curPls)
	s.db.Set(currentPollsKey, curPlsBRes)
}

func (s *State) GetPolls() []PollQuery {
	curPlsB := s.db.Get(currentPollsKey)
	curPls := []PollQuery{}
	json.Unmarshal(curPlsB, &curPls)
	return curPls
}

func (s *State) GetLatestPoll() (string, error) {
	has := s.db.Has(latestPollKey)
	if !has {
		return "", errors.New("There is not any latest poll.")
	}
	hash := s.db.Get(latestPollKey)
	return string(hash), nil
}

func (s *State) IsLatestPoll(id string) bool {
	hash := s.db.Get(latestPollKey)
	return string(hash) == id
}

func (s *State) CreateVote(vd VoteDeliveryData) {
	s.db.Set(prefixVote(vd), nil)
}

func (s *State) HasVote(vd VoteDeliveryData) bool {
	return s.db.Has(prefixVote(vd))
}

func loadState(db dbm.DB) State {
	stateBytes := db.Get(stateKey)
	var state State
	if len(stateBytes) != 0 {
		err := json.Unmarshal(stateBytes, &state)
		if err != nil {
			panic(err)
		}
	}
	state.db = db
	return state
}

func saveState(state State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.db.Set(stateKey, stateBytes)
}
