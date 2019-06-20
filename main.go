package main

import (
	"fmt"
	"log"
	"os/user"
	"path"
	"time"

	"github.com/dgrijalva/jwt-go"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	u, _ := user.Current()
	config, err := clientcmd.BuildConfigFromFlags("", path.Join(u.HomeDir, ".kube", "config"))
	if err != nil {
		log.Fatal(err)
	}

	id, ref := config.AuthProvider.Config["id-token"], config.AuthProvider.Config["refresh-token"]

	p := &jwt.Parser{}

	idT, _, err := p.ParseUnverified(id, &jwt.StandardClaims{})
	if err != nil {
		log.Fatal(err)
	}

	refT, _, err := p.ParseUnverified(ref, &jwt.StandardClaims{})
	if err != nil {
		log.Fatal(err)
	}

	n := time.Now()
	exp := time.Unix(idT.Claims.(*jwt.StandardClaims).ExpiresAt, 0)
	refExp := time.Unix(refT.Claims.(*jwt.StandardClaims).ExpiresAt, 0)
	switch {
	case n.Before(exp):
		fmt.Printf("id token valid until %s (%s)\n", exp, exp.Sub(n))
	case n.Before(refExp):
		fmt.Printf("id token expired, but refresh valid until %s\n", exp)
	}
}
