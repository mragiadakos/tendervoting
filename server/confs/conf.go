package confs

type configuration struct {
	IpfsConnection         string
	AbciDaemon             string
	GonvermentPublicKeyHex string
}

var Conf = configuration{}

func init() {
	Conf.IpfsConnection = "127.0.0.1:5001"
	Conf.AbciDaemon = "tcp://0.0.0.0:46658"
	Conf.GonvermentPublicKeyHex = ""
}
