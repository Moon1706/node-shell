package main

import (
	cli "github.com/Tinkoff/node-shell/cmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // required for GKE
)

func main() {
	cli.InitAndExecute()
}
