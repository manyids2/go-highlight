package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/manyids2/go-highlight/highlights"
	"github.com/manyids2/go-highlight/syntax"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-highlight",
	Short: "neovim highlights for tcell",
	Long:  `Print neovim highlights in terminal using tcell.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Parse path, cterm/gui
		config, _ := cmd.Flags().GetString("config")
		useCterm, _ := cmd.Flags().GetBool("use-cterm")

		// Parse file and get highlights
		h, err := highlights.LoadHighlights(config, useCterm)
		if err != nil {
			log.Fatalln("Could not parse highlights file: ", err)
		}
		h.Print()

		// Read and parse input file with tree-sitter
		path, _ := cmd.Flags().GetString("path")
		t, err := syntax.LoadTree(path)
		if err != nil {
			log.Fatalln("Could not parse source file: ", err)
		}

		// Print
		fmt.Println(t.RootNode())
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("path", "p", "./main.go", "Path to file.")
	rootCmd.Flags().StringP("config", "c", "./corpus/md2.hi", "Path to highlights.")
	rootCmd.Flags().Bool("use-cterm", false, "Use cterm colors instead of gui.")
}
