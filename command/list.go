package command

import (
	"context"
	"errors"
	"flag"
	"strings"
	"time"

	"fmt"
	"os"

	"bufio"

	"github.com/ka2n/ufocatch/ufocatch"
)

// ListCommand implements `ufocatch list <query>` command
type ListCommand struct {
	Meta
	Client ufocatch.Client
}

// Run list command
func (c *ListCommand) Run(args []string) int {
	query, cat, err := c.parseListArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Error("Searching...: " + query)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*10))
	defer cancel()
	done := make(chan error)

	output := bufio.NewWriter(os.Stdout)
	defer output.Flush()
	go func() {
		feed, err := c.Client.Get(ctx, ufocatch.DefaultEndpoint, cat, query)
		if err != nil {
			done <- err
			return
		}

		for _, entry := range feed.Entries {
			fmt.Fprintln(output, entry.ID+": "+entry.Title)
		}
		close(done)
	}()

	if err := waitSignal(ctx, cancel, done); err != nil {
		c.Ui.Error(err.Error())
		return 1
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

func (c *ListCommand) parseListArgs(args []string) (string, ufocatch.Category, error) {
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

	var cat ufocatch.Category
	switch source {
	case "":
		cat = ufocatch.CategoryEdinetx
	case string(ufocatch.CategoryEdinet):
		cat = ufocatch.CategoryEdinet
	case string(ufocatch.CategoryEdinetx):
		cat = ufocatch.CategoryEdinetx
	case string(ufocatch.CategoryTdnet):
		cat = ufocatch.CategoryTdnet
	case string(ufocatch.CategoryTdnetx):
		cat = ufocatch.CategoryTdnetx
	case string(ufocatch.CategoryCg):
		cat = ufocatch.CategoryCg
	}
	if cat == "" {
		return "", "", errors.New("source is invalid: " + source)
	}
	return query, cat, nil
}
