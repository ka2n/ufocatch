package main

import (
	"github.com/ka2n/ufocatch/client"
	"github.com/ka2n/ufocatch/command"
	"github.com/mitchellh/cli"
)

func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"list": func() (cli.Command, error) {
			return &command.ListCommand{
				Meta:   *meta,
				Client: client.Client{},
			}, nil
		},
		"get": func() (cli.Command, error) {
			return &command.GetCommand{
				Meta:   *meta,
				Client: client.Client{},
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Meta:     *meta,
				Version:  Version,
				Revision: GitCommit,
				Name:     Name,
			}, nil
		},
	}
}
