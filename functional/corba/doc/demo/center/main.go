package main

import (
	"git.vnnox.net/ncp/xframe/functional/corba"
)

func main() {

	corba.RunCenterNode(":8080")

	select {}
}
