package main

import (
	"errors"

	client "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
)

func deliver(b []byte) (uint32, error) {
	cli := client.NewHTTP(Conf.NodeDaemon, "/websocket")
	btc, err := cli.BroadcastTxCommit(types.Tx(b))
	if err != nil {
		return CodeTypeClientError, errors.New("Error: " + err.Error())
	}
	if btc.CheckTx.Code > CodeTypeOK {
		return btc.CheckTx.Code, errors.New("Error: " + btc.CheckTx.Log)
	}
	return CodeTypeOK, nil
}

func query(path string, data []byte) ([]byte, error) {
	cli := client.NewHTTP(Conf.NodeDaemon, "/websocket")
	q, err := cli.ABCIQuery(path, data)
	if err != nil {
		return nil, errors.New("Error:" + err.Error())
	}
	if q.Response.Code > CodeTypeOK {
		return nil, errors.New("Error: " + q.Response.Log)
	}
	return q.Response.Value, nil
}
