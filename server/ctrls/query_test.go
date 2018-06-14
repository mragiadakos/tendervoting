package ctrls

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"testing"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/abci/types"
)

func TestQueryListOfElections(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()
	pubHex := hex.EncodeToString(pubB)
	forTestCreateElection(t, app, privk, []string{pubHex})
	forTestCreateElection(t, app, privk, []string{pubHex})
	forTestCreateElection(t, app, privk, []string{pubHex})

	qreq := types.RequestQuery{}
	qreq.Path = "/elections"
	qresp := app.Query(qreq)
	list := ListElectionQuery{}
	json.Unmarshal(qresp.Value, &list)
	assert.Equal(t, 3, len(list))
}

func TestQueryLatestElection(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()
	pubHex := hex.EncodeToString(pubB)
	forTestCreateElection(t, app, privk, []string{pubHex})
	forTestCreateElection(t, app, privk, []string{pubHex})
	latest := forTestCreateElection(t, app, privk, []string{pubHex})

	qreq := types.RequestQuery{}
	qreq.Path = "/elections/latest"
	qresp := app.Query(qreq)
	eq := ElectionQuery{}
	json.Unmarshal(qresp.Value, &eq)
	assert.Equal(t, latest, eq.ID)
}

func TestQueryListOfPolls(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()
	pubHex := hex.EncodeToString(pubB)
	electionID := forTestCreateElection(t, app, privk, []string{pubHex})

	forTestCreatePoll(t, app, privk, electionID, map[string]string{"a": "a"})
	forTestCreatePoll(t, app, privk, electionID, map[string]string{"b": "b"})
	forTestCreatePoll(t, app, privk, electionID, map[string]string{"c": "c"})

	qreq := types.RequestQuery{}
	qreq.Path = "/polls"
	qresp := app.Query(qreq)
	list := ListPollQuery{}
	json.Unmarshal(qresp.Value, &list)
	assert.Equal(t, 3, len(list))
}

func TestQueryLatestPoll(t *testing.T) {
	app := NewTVApplication()
	privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)

	pubB, _ := privk.GetPublic().Bytes()
	pubHex := hex.EncodeToString(pubB)
	electionID := forTestCreateElection(t, app, privk, []string{pubHex})

	forTestCreatePoll(t, app, privk, electionID, map[string]string{"a": "a"})
	forTestCreatePoll(t, app, privk, electionID, map[string]string{"b": "b"})
	latest := forTestCreatePoll(t, app, privk, electionID, map[string]string{"c": "c"})

	qreq := types.RequestQuery{}
	qreq.Path = "/polls/latest"
	qresp := app.Query(qreq)
	pq := PollQuery{}
	json.Unmarshal(qresp.Value, &pq)
	assert.Equal(t, latest, pq.PollHash)
}

func TestQueryVotes(t *testing.T) {
	app := NewTVApplication()
	gonPrivk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	assert.Nil(t, err)
	voters := []crypto.PrivKey{}
	voterHexs := []string{}
	for i := 0; i < 100; i++ {
		privk, _, err := crypto.GenerateEd25519Key(rand.Reader)
		assert.Nil(t, err)
		pubB, _ := privk.GetPublic().Bytes()
		pubHex := hex.EncodeToString(pubB)
		voters = append(voters, privk)
		voterHexs = append(voterHexs, pubHex)
	}

	electionID := forTestCreateElection(t, app, gonPrivk, voterHexs)
	pollHash := forTestCreatePoll(t, app, gonPrivk, electionID, map[string]string{"a": "a", "b": "b", "c": "c"})

	for _, v := range voters {
		forTestCreateVote(t, app, v, electionID, pollHash, "a")
	}

	qreq := types.RequestQuery{}
	qreq.Path = "/votes"
	qreq.Data, _ = json.Marshal(PollQuery{PollHash: pollHash})
	qresp := app.Query(qreq)
	pvq := PollVotesQuery{}
	json.Unmarshal(qresp.Value, &pvq)
	assert.Equal(t, 100, pvq.Choices["a"])
	assert.Equal(t, 100, pvq.NumberOfVotes)
}
