package cmd

import (
	"log"
	"os"

	"github.com/manyids2/go-highlight/highlights"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-highlight",
	Short: "neovim highlights for tcell",
	Long:  `Print neovim highlights in terminal using tcell.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse path, cterm/gui
		path, _ := cmd.Flags().GetString("path")
		useCterm, _ := cmd.Flags().GetBool("use-cterm")

		// Parse file and get highlights
		h, err := highlights.LoadHighlights(path, useCterm)
		if err != nil {
			log.Fatalln("Could not parse highlights file: ", err)
		}
		h.Print()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("path", "p", "./corpus/md2.hi", "Path to highlights file.")
	rootCmd.Flags().BoolP("use-cterm", "c", false, "Use cterm colors instead of gui.")
}
