//go:build !tools

package main

import (
	"fmt"
	"io"
	"messanger/internal/model"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
    stmts, err := gormschema.New("postgres").
        Load(
            &model.Chat{},
            &model.ChatMember{},
            &model.Message{},
            &model.Attachment{},
            &model.Reaction{},
            &model.PinnedMessages{},
        )
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
        os.Exit(1)
    }
    io.WriteString(os.Stdout, stmts)
}