package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/mragiadakos/tendervoting/server/ctrls"
	uuid "github.com/satori/go.uuid"
	"github.com/tendermint/abci/types"
	"github.com/urfave/cli"
)

var GenerateKeyCommand = cli.Command{
	Name:    "generate",
	Aliases: []string{"g"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "filename",
			Usage: "the filename that the key will be saved",
		},
	},
	Usage: "generate the key in a file",
	Action: func(c *cli.Context) error {
		filename := c.String("filename")
		if len(filename) == 0 {
			return errors.New("Error: filename is missing")
		}
		privk, _, _ := crypto.GenerateKeyPair(crypto.Ed25519, 0)
		kj := KeyJson{}
		b, _ := privk.GetPublic().Bytes()
		kj.PublicKey = hex.EncodeToString(b)
		kj.PrivateKey, _ = crypto.MarshalPrivateKey(privk)
		b, _ = json.Marshal(kj)
		err := ioutil.WriteFile(filename, b, 0644)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}
		fmt.Println("The generate was successful")
		return nil
	},
}

var CreateElectionCommand = cli.Command{
	Name:    "create-election",
	Aliases: []string{"ce"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "key",
			Usage: "the filename of the key",
		},
		cli.StringFlag{
			Name:  "voters",
			Usage: "the voters' public keys seperated by comma",
		},
	},
	Usage: "create the election and adding the voters",
	Action: func(c *cli.Context) error {
		filename := c.String("key")
		if len(filename) == 0 {
			return errors.New("Error: filename is missing")
		}
		priv, err := fileKey(filename)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}
		strVoters := c.String("voters")
		voters := strings.Split(strVoters, ",")
		edd := ctrls.ElectionDeliveryData{}
		pubB, _ := priv.GetPublic().Bytes()
		edd.From = hex.EncodeToString(pubB)
		edd.ID = uuid.NewV4().String()
		edd.Voters = voters
		b, _ := json.Marshal(edd)
		sigB, err := priv.Sign(b)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}
		tvd := ctrls.TVDelivery{}
		tvd.Data = edd
		tvd.Type = ctrls.ELECTION
		tvd.Signature = sigB
		b, _ = json.Marshal(tvd)
		_, err = deliver(b)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}
		fmt.Println("The election submitted with ID", edd.ID)
		return nil
	},
}

var AddPollCommand = cli.Command{
	Name:    "add-poll",
	Aliases: []string{"cp"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "key",
			Usage: "the filename of the key",
		},
		cli.StringFlag{
			Name:  "hash",
			Usage: "the poll's directory as an IPFS hash",
		},
		cli.StringFlag{
			Name:  "election",
			Usage: "the election's ID",
		},
	},
	Usage: "add the poll to the election",
	Action: func(c *cli.Context) error {
		filename := c.String("key")
		if len(filename) == 0 {
			return errors.New("Error: filename is missing")
		}

		priv, err := fileKey(filename)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}

		hash := c.String("hash")
		if len(hash) == 0 {
			return errors.New("Error: hash is missing")
		}

		election := c.String("election")
		if len(election) == 0 {
			return errors.New("Error: election is missing")
		}

		pdd := ctrls.PollDeliveryData{}
		pubB, _ := priv.GetPublic().Bytes()
		pdd.From = hex.EncodeToString(pubB)
		pdd.PollHash = hash
		pdd.ElectionID = election
		b, _ := json.Marshal(pdd)
		sigB, err := priv.Sign(b)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}
		tvd := ctrls.TVDelivery{}
		tvd.Data = pdd
		tvd.Type = ctrls.POLL
		tvd.Signature = sigB
		b, _ = json.Marshal(tvd)
		_, err = deliver(b)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}
		fmt.Println("The poll submitted")
		return nil
	},
}

var VoteCommand = cli.Command{
	Name:    "vote",
	Aliases: []string{"v"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "key",
			Usage: "the filename of the key",
		},
		cli.StringFlag{
			Name:  "hash",
			Usage: "the poll's directory as an IPFS hash",
		},
		cli.StringFlag{
			Name:  "choice",
			Usage: "the choice's ID from the poll",
		},
	},
	Usage: "vote for a specific poll",
	Action: func(c *cli.Context) error {
		filename := c.String("key")
		if len(filename) == 0 {
			return errors.New("Error: filename is missing")
		}

		priv, err := fileKey(filename)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}

		hash := c.String("hash")
		if len(hash) == 0 {
			return errors.New("Error: hash is missing")
		}

		choice := c.String("choice")
		if len(hash) == 0 {
			return errors.New("Error: choice is missing")
		}

		vdd := ctrls.VoteDeliveryData{}
		pubB, _ := priv.GetPublic().Bytes()
		vdd.From = hex.EncodeToString(pubB)
		vdd.PollHash = hash
		vdd.Choice = choice
		b, _ := json.Marshal(vdd)
		sigB, err := priv.Sign(b)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}
		tvd := ctrls.TVDelivery{}
		tvd.Data = vdd
		tvd.Type = ctrls.VOTE
		tvd.Signature = sigB
		b, _ = json.Marshal(tvd)
		_, err = deliver(b)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}
		fmt.Println("The vote submitted")
		return nil

	},
}

var QueryElectionsCommand = cli.Command{
	Name:    "elections",
	Aliases: []string{"e"},
	Usage:   "list the election's IDs",
	Action: func(c *cli.Context) error {
		req := types.RequestQuery{}
		req.Path = "/elections"
		resp, err := RpcQuery(req)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}
		if resp.Code > CodeTypeOK {
			return errors.New("Error :" + resp.Log)
		}
		les := ctrls.ListElectionQuery{}
		json.Unmarshal(resp.Value, &les)
		for _, v := range les {
			fmt.Println("Election ID:", v.ID)
			fmt.Println("Latest:", v.Latest)
			fmt.Println("Number of voters:", v.NumberOfVoters)
			fmt.Println()
		}
		return nil
	},
}

var QueryLatestElectionCommand = cli.Command{
	Name:    "latest-election",
	Aliases: []string{"le"},
	Usage:   "get the latest election's IDs",
	Action: func(c *cli.Context) error {
		req := types.RequestQuery{}
		req.Path = "/elections/latest"
		resp, err := RpcQuery(req)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}

		if resp.Code > CodeTypeOK {
			return errors.New("Error :" + resp.Log)
		}
		v := ctrls.ItemElectionQuery{}
		json.Unmarshal(resp.Value, &v)
		fmt.Println("Election ID:", v.ID)
		fmt.Println("Latest:", v.Latest)
		fmt.Println("Number of voters:", v.NumberOfVoters)
		fmt.Println()

		return nil
	},
}

var QueryPollsCommand = cli.Command{
	Name:    "polls",
	Aliases: []string{"p"},
	Usage:   "list the poll's hash keys",
	Action: func(c *cli.Context) error {
		req := types.RequestQuery{}
		req.Path = "/polls"
		resp, err := RpcQuery(req)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}

		if resp.Code > CodeTypeOK {
			return errors.New("Error :" + resp.Log)
		}
		pes := ctrls.ListPollQuery{}
		json.Unmarshal(resp.Value, &pes)
		for _, v := range pes {
			fmt.Println("Poll's Hash:", v.PollHash)
			fmt.Println("Latest:", v.Latest)
			fmt.Println()
		}

		return nil
	},
}

var QueryLatestPollCommand = cli.Command{
	Name:    "latest-poll",
	Aliases: []string{"lp"},
	Usage:   "get the latest poll's hash key",
	Action: func(c *cli.Context) error {
		req := types.RequestQuery{}
		req.Path = "/polls/latest"
		resp, err := RpcQuery(req)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}

		if resp.Code > CodeTypeOK {
			return errors.New("Error :" + resp.Log)
		}
		v := ctrls.ItemPollQuery{}
		json.Unmarshal(resp.Value, &v)
		fmt.Println("Poll's Hash:", v.PollHash)
		fmt.Println("Latest:", v.Latest)
		fmt.Println()

		return nil
	},
}

var QueryResultsCommand = cli.Command{
	Name:    "results",
	Aliases: []string{"r"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "hash",
			Usage: "the poll's directory as an IPFS hash",
		},
	},
	Usage: "get the latest results based on poll's hash key",
	Action: func(c *cli.Context) error {
		hash := c.String("hash")
		if len(hash) == 0 {
			return errors.New("Error: hash is missing")
		}
		req := types.RequestQuery{}
		req.Path = "/votes"
		b, _ := json.Marshal(ctrls.PollQuery{PollHash: hash})
		req.Data = b
		resp, err := RpcQuery(req)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}

		if resp.Code > CodeTypeOK {
			return errors.New("Error :" + resp.Log)
		}
		v := ctrls.PollVotesQuery{}
		json.Unmarshal(resp.Value, &v)
		for k, n := range v.Choices {
			fmt.Println("Votes for choice '"+k+"':", n)
		}
		fmt.Println("Number of voters:", v.NumberOfVotes)
		fmt.Println()

		return nil
	},
}
