package command

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"regexp"

	"flag"

	"net/http"

	"github.com/ka2n/ufocatch/ufocatch"
	"github.com/ka2n/ufocatch/util"
	"golang.org/x/sync/errgroup"
)

// GetCommand impliments `ufocatch get <id>` command
type GetCommand struct {
	Meta
}

// Run get command
func (c *GetCommand) Run(args []string) int {
	ids, format, err := parseArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	client, err := ufocatch.New(ufocatch.DefaultEndpoint, http.DefaultClient)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*10))
	defer cancel()
	done := make(chan error)

	wg := errgroup.Group{}
	for _, id := range ids {
		id := id
		wg.Go(func() error {
			name, err := client.Download(ctx, format, id)
			if err != nil {
				return err
			}
			c.Ui.Output(fmt.Sprintf("saved: %v", name))
			return nil
		})
	}

	go func() {
		err := wg.Wait()
		done <- err
	}()

	if err := waitSignal(ctx, cancel, done); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}
	return 0
}

// Synopsis for get command
func (c *GetCommand) Synopsis() string {
	return "Get resources by ID"
}

// Help for get command
func (c *GetCommand) Help() string {
	helpText := `
Usage: ufocatch get <ID> [OPTIONS]

Get resources by ID.
This command searches ID string from args, then retrieve a resource on your filesystem.

Options:
    --format=xbrl        File format to download. 'pdf' or 'xbrl'

To get XBRL zip archive(default)
	ufocatch get ED2014121600183 --format=xbrl

To get PDF
	ufocatch get ED2014121600183 --format=pdf
	
Also, you can use standard input like this.
	ufocatch list 'Search string' | head -n1 | ufocatch get --format=pdf		
	`
	return strings.TrimSpace(helpText)
}

func parseArgs(args []string) ([]string, string, error) {
	var rawID []string
	if isaStdin() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			txt := scanner.Text()
			txt = util.StripANSISequence(txt)
			rawID = append(rawID, txt)
		}
	}

	var format string
	opt := flag.NewFlagSet("get", flag.ContinueOnError)
	opt.StringVar(&format, "format", "xbrl", "Format to retrieve")
	if err := opt.Parse(args); err != nil {
		return nil, "", err
	}

	if len(rawID) == 0 {
		for _, arg := range opt.Args() {
			rawID = append(rawID, arg)
		}
	}

	if len(rawID) == 0 {
		return nil, "", errors.New("query is mandatory")
	}

	ids := make([]string, len(rawID))
	for i, raw := range rawID {
		id := parseRawID(raw)
		if id == "" {
			return nil, "", errors.New("invalid id: " + raw)
		}
		ids[i] = id
	}

	var dataFormat string
	switch format {
	case "xbrl":
		dataFormat = ufocatch.FormatData
	case "pdf":
		dataFormat = ufocatch.FormatPDF
	}
	if dataFormat == "" {
		return nil, "", errors.New("format is invalid: " + format)
	}

	return ids, dataFormat, nil
}

func isaStdin() bool {
	stat, _ := os.Stdin.Stat()
	return stat.Mode()&os.ModeCharDevice == 0
}

func parseRawID(rawID string) string {
	exp := regexp.MustCompile(`(ED|TD|CG)\d+`)
	return exp.FindString(rawID)
}
