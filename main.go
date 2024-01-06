package main

import (
	cli "github.com/Moon1706/node-shell/cmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // required for GKE
)

func main() {
	cli.InitAndExecute()
}
