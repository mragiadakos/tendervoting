package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/mragiadakos/tendervoting/server/confs"
	"github.com/mragiadakos/tendervoting/server/ctrls"
	absrv "github.com/tendermint/abci/server"
	cmn "github.com/tendermint/tmlibs/common"
	tmlog "github.com/tendermint/tmlibs/log"
)

func main() {
	logger := tmlog.NewTMLogger(kitlog.NewSyncWriter(os.Stdout))
	flagAbci := "socket"
	ipfsDaemon := flag.String("ipfs", "127.0.0.1:5001", "the URL for the IPFS's daemon")
	node := flag.String("node", "tcp://0.0.0.0:46658", "the TCP URL for the ABCI daemon")
	gonvermentPublicKey := flag.String("gonverment", "", "the gonverment's public key")
	flag.Parse()

	if len(*gonvermentPublicKey) == 0 {
		fmt.Println("Error ", errors.New("The gonverment's public key is missing"))
		return
	}

	confs.Conf.GonvermentPublicKeyHex = *gonvermentPublicKey

	confs.Conf.AbciDaemon = *node
	confs.Conf.IpfsConnection = *ipfsDaemon

	app := ctrls.NewTVApplication()
	srv, err := absrv.NewServer(confs.Conf.AbciDaemon, flagAbci, app)
	if err != nil {
		fmt.Println("Error ", err)
		return
	}
	srv.SetLogger(logger.With("module", "abci-server"))
	if err := srv.Start(); err != nil {
		fmt.Println("Error ", err)
		return
	}

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})

}
