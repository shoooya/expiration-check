/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/hokupod/expiration-check/expchk"
	"github.com/hokupod/expiration-check/expchk/domain"
	"github.com/hokupod/expiration-check/expchk/ssl"
	"github.com/spf13/cobra"
)

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Extracts expiration dates for all supported source",
	Long: `Extracts expiration dates for all supported source.(JSON output)

Example for:
  expiration-check all example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			sh ssl.Holder
			dh domain.Holder
		)

		ec := expchk.New(args[0])
		ec.AddHolder(sh)
		ec.AddHolder(dh)
		res := ec.Run()
		for _, ex := range res.Expirations {
			if ex.Errors != nil {
				for _, err := range ex.Errors {
					fmt.Printf("Error: %v: %v\n", ex.Name, err)
				}
				os.Exit(1)
			}
		}

		jsonStr, err := json.Marshal(res)
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}

		var buf bytes.Buffer
		err = json.Indent(&buf, []byte(jsonStr), "", "  ")
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
		fmt.Println(buf.String())
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires domain")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
}
