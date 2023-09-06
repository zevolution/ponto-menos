package client

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ponto-menos",
	Short: "Ponto Menos CLI in GoLang",
	Long:  `Ponto Menos CLI application written in Go.`,
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
