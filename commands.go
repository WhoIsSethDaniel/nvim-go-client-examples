package main

import (
	"log"

	"github.com/neovim/go-client/nvim/plugin"
)

type cmdEvalExample struct {
	Cwd     string `msgpack:",array"`
	Bufname string
}

func exCmd(p *plugin.Plugin, args []string, bang bool, eval *cmdEvalExample) {
	log.Println("  Args to exCmd:")
	log.Printf("    arg1: %s\n", args[0])
	log.Printf("    bang: %t\n", bang)
	log.Printf("    cwd: %s\n", eval.Cwd)
	log.Printf("    buffer: %s\n", eval.Bufname)
}
