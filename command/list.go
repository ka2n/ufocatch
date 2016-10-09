package command

import (
	"context"
	"errors"
	"flag"
	"strings"
	"time"

	"github.com/ka2n/ufocatch/ufocatcher"
)

// ListCommand implements `ufocatch list <query>` command
type ListCommand struct {
	Meta
}

// Run list command
func (c *ListCommand) Run(args []string) int {
	query, cat, err := c.parseListArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Error("Searching...: " + query)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*30))
	defer cancel()
	feed, err := ufocatcher.Get(ctx, ufocatcher.DefaultEndpoint, cat, query)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	for _, entry := range feed.Entries {
		c.Ui.Output(entry.ID + ": " + entry.Title)
	}
	return 0
}

// Synopsis for list command
func (c *ListCommand) Synopsis() string {
	return "List resources"
}

// Help for list command
func (c *ListCommand) Help() string {
	helpText := `List resources by query.

To find EDINET(with XBRL) resources (default)
	ufocatch list <query> --source=edinetx
	`
	return strings.TrimSpace(helpText)
}

func (c *ListCommand) parseListArgs(args []string) (string, ufocatcher.Category, error) {
	var source string
	var ask bool

	opt := flag.NewFlagSet("list", flag.ContinueOnError)
	opt.StringVar(&source, "source", "edinetx", "Source category. [edinet, edinetx, tdnet, tdnetx, cg], default 'edientx'")
	opt.BoolVar(&ask, "ask", false, "input source by prompt")
	if err := opt.Parse(args); err != nil {
		return "", "", err
	}

	var query string
	if ask {
		ans, err := c.Ui.Ask("query:")
		if err != nil {
			return "", "", err
		}
		query = ans
	} else {
		query = opt.Arg(0)
	}

	if query == "" {
		return "", "", errors.New("query is mandatory")
	}

	var cat ufocatcher.Category
	switch source {
	case "":
		cat = ufocatcher.CategoryEdinetx
	case string(ufocatcher.CategoryEdinet):
		cat = ufocatcher.CategoryEdinet
	case string(ufocatcher.CategoryEdinetx):
		cat = ufocatcher.CategoryEdinetx
	case string(ufocatcher.CategoryTdnet):
		cat = ufocatcher.CategoryTdnet
	case string(ufocatcher.CategoryTdnetx):
		cat = ufocatcher.CategoryTdnetx
	case string(ufocatcher.CategoryCg):
		cat = ufocatcher.CategoryCg
	}
	if cat == "" {
		return "", "", errors.New("source is invalid: " + source)
	}
	return query, cat, nil
}
