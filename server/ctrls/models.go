package ctrls

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/satori/go.uuid"

	"github.com/ipfs/go-ipfs-api"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/mragiadakos/tendervoting/server/confs"
)

const (
	CodeTypeOK            uint32 = 0
	CodeTypeEncodingError uint32 = 1
	CodeTypeBadNonce      uint32 = 2
	CodeTypeUnauthorized  uint32 = 3
	CodeTypeServerError   uint32 = 4
)

type DeliveryType string

type DeliveryDataInterface interface {
	GetFrom() string
}

const (
	ELECTION = DeliveryType("election")
	POLL     = DeliveryType("poll")
	VOTE     = DeliveryType("vote")
)

type TVDelivery struct {
	Signature []byte
	Type      DeliveryType
	Data      interface{}
}

func (v *TVDelivery) GetFrom() (string, error) {
	pubHex := ""
	switch v.Type {
	case ELECTION:
		d := v.GetElectionDeliveryData()
		pubHex = d.From
	case POLL:
		d := v.GetPollDeliveryData()
		pubHex = d.From
	case VOTE:
		d := v.GetVoteDeliveryData()
		pubHex = d.From
	default:
		return "", errors.New("The type for the delivery can only be 'election', 'poll' or 'vote'.")
	}
	return pubHex, nil
}

func (v *TVDelivery) GetElectionDeliveryData() ElectionDeliveryData {
	b, _ := json.Marshal(v.Data)
	d := ElectionDeliveryData{}
	json.Unmarshal(b, &d)
	return d
}

func (v *TVDelivery) GetPollDeliveryData() PollDeliveryData {
	b, _ := json.Marshal(v.Data)
	d := PollDeliveryData{}
	json.Unmarshal(b, &d)
	return d
}

func (v *TVDelivery) GetVoteDeliveryData() VoteDeliveryData {
	b, _ := json.Marshal(v.Data)
	d := VoteDeliveryData{}
	json.Unmarshal(b, &d)
	return d
}

func (v *TVDelivery) GetDataInStructureOrder() ([]byte, error) {
	b, _ := json.Marshal(v.Data)
	out := []byte("")
	switch v.Type {
	case ELECTION:
		d := ElectionDeliveryData{}
		json.Unmarshal(b, &d)
		out, _ = json.Marshal(d)
	case POLL:
		d := PollDeliveryData{}
		json.Unmarshal(b, &d)
		out, _ = json.Marshal(d)
	case VOTE:
		d := VoteDeliveryData{}
		json.Unmarshal(b, &d)
		out, _ = json.Marshal(d)
	default:
		return out, errors.New("The type for the delivery can only be 'election', 'poll' or 'vote'.")
	}
	return out, nil
}

func (v *TVDelivery) VerifySignature() (bool, error) {
	pubHex, err := v.GetFrom()
	if err != nil {
		return false, err
	}
	pubB, err := hex.DecodeString(pubHex)
	if err != nil {
		return false, errors.New("The public key is not correct hex: " + err.Error())
	}
	pub, err := crypto.UnmarshalPublicKey(pubB)
	if err != nil {
		return false, errors.New("The public key is not correct")
	}
	b, _ := v.GetDataInStructureOrder()
	ver, err := pub.Verify(b, v.Signature)
	if err != nil {
		return false, errors.New("The signature's format is not correct.")
	}
	return ver, nil
}

type VoteDeliveryData struct {
	From     string
	PollHash string
	Choice   string
}

func (self *VoteDeliveryData) GetFrom() string {
	return self.From
}

type PollDeliveryData struct {
	From       string
	PollHash   string
	ElectionID string
}

func (self *PollDeliveryData) GetFrom() string {
	return self.From
}

func (p *PollDeliveryData) ValidateGonverment() error {
	if p.From != confs.Conf.GonvermentPublicKeyHex {
		return errors.New("You are not a gonverment.")
	}
	return nil
}

type PollJson struct {
	Description string
	Choices     map[string]string
}

func (p *PollDeliveryData) GetPollJsonFromPollHash() (*PollJson, error) {
	sh := shell.NewShell(confs.Conf.IpfsConnection)
	objLs, err := sh.FileList(p.PollHash)
	if err != nil {
		return nil, errors.New("Folder from IPFS hash " + p.PollHash + " could not list files: " + err.Error())
	}
	pollJsonFound := false
	hash := ""
	for _, v := range objLs.Links {
		if v.Name == "poll.json" {
			hash = v.Hash
			pollJsonFound = true
		}
	}
	if !pollJsonFound {
		return nil, errors.New("Folder from IPFS hash " + p.PollHash + " could not find the poll.json in the folder.")
	}

	pollFile := uuid.NewV4().String()
	err = sh.Get(hash, pollFile)
	if err != nil {
		return nil, errors.New("Failed to get file poll.json with hash " + p.PollHash + ": " + err.Error())
	}
	b, err := ioutil.ReadFile(pollFile)
	if err != nil {
		return nil, errors.New("Failed to read file poll.json with hash " + p.PollHash + ": " + err.Error())
	}
	os.RemoveAll(pollFile)
	pj := PollJson{}
	err = json.Unmarshal(b, &pj)
	if err != nil {
		return nil, errors.New("The poll.json has not the correct JSON format: " + err.Error())
	}
	if len(pj.Description) == 0 {
		return nil, errors.New("The poll.json has empty description.")
	}
	if len(pj.Choices) == 0 {
		return nil, errors.New("The poll.json has empty choices.")
	}
	return &pj, nil
}

type ElectionDeliveryData struct {
	ID     string
	From   string
	Voters []string
}

func (self *ElectionDeliveryData) GetFrom() string {
	return self.From
}

func (e *ElectionDeliveryData) ValidateGonverment() error {
	if e.From != confs.Conf.GonvermentPublicKeyHex {
		return errors.New("You are not a gonverment.")
	}
	return nil
}

func (e *ElectionDeliveryData) ValidateVoters() error {
	voters := map[string]int{}
	for _, v := range e.Voters {
		pubB, err := hex.DecodeString(v)
		if err != nil {
			return errors.New("The voter " + v + " is not a correct hex: " + err.Error())
		}
		_, err = crypto.UnmarshalPublicKey(pubB)
		if err != nil {
			return errors.New("The voter " + v + " has not a correct public key: " + err.Error())
		}
		_, ok := voters[v]
		if ok {
			return errors.New("The voter " + v + " exists already in the list.")
		}
		voters[v] = 0
	}
	return nil
}
