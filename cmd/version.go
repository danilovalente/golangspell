package cmd

import (
	"fmt"

	"github.com/golangspell/golangspell/config"
	"github.com/golangspell/golangspell/domain"
	"github.com/golangspell/golangspell/usecase"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Golang Spell version number",
	Long: `Shows the Golang Spell current installed version, or the specified Spell version
Usage: 
golangspell version

or 

golangspell version [Spell name]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			spellConfig := domain.GetConfig()
			err := usecase.ShowSpellVersion(&domain.GolangLibrary{Name: args[0]}, spellConfig)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Printf("Golang Spell v%s -- HEAD\n", config.Version)
		}
	},
}
