// Package inspect implements 'kuma token inspect'.
package inspect

import (
	"fmt"
	"log"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/yurinandayona-com/kuma/subcmd/token/config"
)

// Cmd is 'kuma token inspect` command.
var Cmd *cobra.Command

func init() {
	var token string

	Cmd = &cobra.Command{
		Use:   "inspect",
		Short: "Inspect the user token",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if token == "" {
				return flag.ErrHelp
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			jm, err := config.LoadJWTManager(cmd)
			if err != nil {
				log.Fatalf("alert: %s", err)
			}

			claims, valid, err := jm.Parse(token)
			if err != nil {
				log.Printf("error: %s", err)
				valid = false
			}

			if valid {
				_, err := jm.UserDB.Verify(claims.ID, claims.Name)
				if err != nil {
					log.Printf("error: %s", err)
					valid = false
				}
			}

			expire := time.Unix(claims.ExpiresAt, 0)

			fmt.Printf("ID     : %s\n", claims.ID)
			fmt.Printf("Name   : %s\n", claims.Name)
			fmt.Printf("Expire : %s (%s)\n", expire, humanize.Time(expire))
			fmt.Printf("Valid  : %t\n", valid)
		},
	}

	Cmd.Flags().StringVarP(&token, "token", "t", "", "user token (required)")
}