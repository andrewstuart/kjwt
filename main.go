package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/dgrijalva/jwt-go"
	"k8s.io/client-go/tools/clientcmd"
)

type SillyTime struct {
	time.Time
}

func (s *SillyTime) UnmarshalJSON(bs []byte) error {
	var st int64
	err := json.Unmarshal(bs, &st)
	if err != nil {
		return err
	}
	*s = SillyTime{time.Unix(st, 0)}
	return nil
}

type claims struct {
	jwt.StandardClaims
	Groups    []string
	Email     string
	ExpiresAt *SillyTime `json:"exp,omitempty"`
	IssuedAt  *SillyTime `json:"iat,omitempty"`
	NotBefore *SillyTime `json:"nbf,omitempty"`
}

func first(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

func main() {
	u, _ := user.Current()
	config, err := clientcmd.BuildConfigFromFlags("", first(os.Getenv("KUBECONFIG"), path.Join(u.HomeDir, ".kube", "cache")))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("config = %+v\n", config)
	fmt.Printf("config.AuthProvider = %+v\n", config.AuthProvider)
	id, ref := config.AuthProvider.Config["id-token"], config.AuthProvider.Config["refresh-token"]

	p := &jwt.Parser{}

	idT, _, err := p.ParseUnverified(id, &claims{})
	if err != nil {
		log.Fatal(err)
	}

	refT, _, err := p.ParseUnverified(ref, &claims{})
	if err != nil {
		log.Fatal(err)
	}

	idC, refC := idT.Claims.(*claims), refT.Claims.(*claims)
	n := time.Now()

	tw := tabwriter.NewWriter(os.Stdout, 1, 2, 1, ' ', 0)

	// 	red := color.New(color.FgRed).SprintFunc()
	red := func(s string) string {
		return s
	}

	fmt.Fprintf(tw, red("ID\t\n"))
	fmt.Fprintf(tw, "Expiry\t%s\n", idC.ExpiresAt.Format(time.Stamp))
	fmt.Fprintf(tw, "Time Left\t%s\n", idC.ExpiresAt.Sub(n).Round(time.Minute))
	fmt.Fprintf(tw, "Groups\t%s\n", strings.Join(idC.Groups, ", "))
	fmt.Fprintf(tw, red("\t\nRefresh\t\n"))
	fmt.Fprintf(tw, "Expiry\t%s\n", refC.ExpiresAt.Format(time.Stamp))
	fmt.Fprintf(tw, "Time Left\t%s\n", refC.ExpiresAt.Sub(n).Round(time.Minute))

	tw.Flush()
}
