package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/danilovalente/golangspell/cmd"

	"github.com/danilovalente/golangspell/domain"
	"github.com/spf13/cobra"
)

func loadSpellCommand(spell *domain.Spell, command *domain.Command) {
	spellCMD := &cobra.Command{
		Use:   command.Name,
		Short: command.ShortDescription,
		Long:  command.LongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			library := domain.GolangLibrary{Name: spell.Name, URL: spell.URL}
			spellCommand := exec.Command(library.BinPath(), append([]string{command.Name}, args...)...)
			outputBytes, err := spellCommand.Output()
			if err != nil {
				log.Fatalf("%s failed with %s\n", command.Name, err)
			} else {
				fmt.Println(string(outputBytes))
			}
		},
	}
	cmd.RootCmd.AddCommand(spellCMD)
}

func loadSpellDescription(golangLibrary *domain.GolangLibrary, config *domain.Config) {
	fmt.Printf("Loading Spell %s description ...\n", golangLibrary.Name)
	execCmd := exec.Command(golangLibrary.BinPath(), "build-config")
	outputBytes, err := execCmd.Output()
	if err != nil {
		log.Fatalf("%s build-config failed with %s\n", golangLibrary.BinPath(), err)
	}
	var spell domain.Spell
	err = json.Unmarshal(outputBytes, &spell)
	if err != nil {
		panic(err)
	}
	if nil == config.Spellbook {
		config.Spellbook = make(map[string]domain.Spell, 0)
	}
	config.Spellbook[spell.Name] = spell
	repo := domain.GetConfigRepository()
	repo.Save(config)
	fmt.Printf("Spell %s description loaded\n", golangLibrary.Name)
}

//LoadSpells and configure cmds to call them
func LoadSpells() {
	fmt.Println("Loading Spells ...")
	config := domain.GetConfig()
	if nil == config.Spellbook || len(config.Spellbook) == 0 {
		for _, golangLibrary := range config.DefaultSpells {
			loadSpellDescription(&golangLibrary, &config)
		}
	}
	for _, spell := range config.Spellbook {
		if !spell.Installed {
			InstallSpell(&spell, &config)
		}
		for _, command := range spell.Commands {
			loadSpellCommand(&spell, &command)
		}
	}
	fmt.Println("Spells loaded!")
}
