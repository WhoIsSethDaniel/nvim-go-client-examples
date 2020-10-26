package main

import (
	"log"
	"os"

	"github.com/neovim/go-client/nvim/plugin"
)

func main() {
	// create a log to log to right away. It will help with debugging
	l, _ := os.Create("nvim-go-client-example.log")
	log.SetOutput(l)

	plugin.Main(func(p *plugin.Plugin) error {
		// Functions
		p.HandleFunction(&plugin.FunctionOptions{Name: "Upper"},
			func(args []string) (string, error) {
				log.Print("calling Upper")
				return upper(p, args[0]), nil
			})
		p.HandleFunction(&plugin.FunctionOptions{Name: "UpperCwd", Eval: "getcwd()"},
			func(args []string, dir string) (string, error) {
				log.Print("calling UpperCwd")
				return upper(p, dir), nil
			})
		p.HandleFunction(&plugin.FunctionOptions{Name: "ShowThings", Eval: "[getcwd(),argc()]"},
			func(args []string, eval *someArgs) ([]string, error) {
				log.Print("calling ShowThings")
				return returnArgs(p, eval)
			})
		p.HandleFunction(&plugin.FunctionOptions{Name: "GetVV"},
			func(args []string) ([]string, error) {
				log.Print("calling GetVV")
				return getvv(p, args[0])
			})
		p.HandleFunction(&plugin.FunctionOptions{Name: "ShowFirst"},
			func(args []string) (string, error) {
				log.Print("calling ShowFirst")
				return showfirst(p), nil
			})

		// AutoCommands
		p.HandleAutocmd(&plugin.AutocmdOptions{Event: "BufEnter", Group: "ExmplNvGoClientGrp", Pattern: "*"},
			func() {
				log.Print("Just entered a buffer")
			})
		p.HandleAutocmd(&plugin.AutocmdOptions{Event: "BufAdd", Group: "ExmplNvGoClientGrp", Pattern: "*", Eval: "*"},
			func(eval *autocmdEvalExample) {
				log.Printf("buffer has cwd: %s", eval.Cwd)
			})

		// Commands
		p.HandleCommand(&plugin.CommandOptions{Name: "ExCmd", NArgs: "?", Bang: true, Eval: "[getcwd(),bufname()]"},
			func(args []string, bang bool, eval *cmdEvalExample) {
				log.Print("called command ExCmd")
				exCmd(p, args, bang, eval)
			})
		return nil
	})
}
