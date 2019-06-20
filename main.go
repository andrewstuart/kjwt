package main

import (
	"os/user"

	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	u, _ := user.Current()
	config, err := clientcmd.BuildConfigFromFlags("", ".kube/config")
}
