package ctrls

import (
	"encoding/json"
	"errors"
	"time"

	dbm "github.com/tendermint/tmlibs/db"
)

var (
	stateKey       = []byte("stateKey")
	electionKey    = []byte("election:")
	pollKey        = []byte("poll:")
	voteKey        = []byte("vote:")
	latestElection = []byte("latestElection")
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
	ID        string
	Voters    []string
	StartTime time.Time
	EndTime   time.Time
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
	es.StartTime = ed.StartTime
	es.EndTime = ed.EndTime
	b, _ := json.Marshal(es)
	s.db.Set(prefixElection(es.ID), b)
	s.db.Set(latestElection, []byte(es.ID))
}

func (s *State) GetLatestElection() (string, error) {
	has := s.db.Has(latestElection)
	if !has {
		return "", errors.New("There is not any latest election.")
	}
	uuid := s.db.Get(latestElection)
	return string(uuid), nil
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
