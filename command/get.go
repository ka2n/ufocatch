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

	"github.com/ka2n/ufocatch/ufocatcher"
)

// GetCommand impliments `ufocatch get <id>` command
type GetCommand struct {
	Meta
}

// Run get command
func (c *GetCommand) Run(args []string) int {
	id, format, err := parseArgs(args)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*30))
	defer cancel()
	name, err := ufocatcher.Download(ctx, ufocatcher.DefaultEndpoint, format, id)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output(fmt.Sprintf("saved: %v", name))
	return 0
}

// Synopsis for get command
func (c *GetCommand) Synopsis() string {
	return "Get resources by ID"
}

// Help for get command
func (c *GetCommand) Help() string {
	helpText := `Get resources by ID.
	
This command searches ID string from args, then retrieve a resource on your filesystem.

To get XBRL zip archive(default)
	ufocatch get ED2014121600183 --format=xbrl

To get PDF
	ufocatch get ED2014121600183 --format=pdf
	
Also, you can use standard input like this.
	ufocatch list 'Search string' | head -n1 | ufocatch get --format=pdf		
	`
	return strings.TrimSpace(helpText)
}

func parseArgs(args []string) (string, ufocatcher.Format, error) {
	var rawID string
	if isaStdin() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			rawID = scanner.Text()
			break
		}
	}

	var format string
	opt := flag.NewFlagSet("get", flag.ContinueOnError)
	opt.StringVar(&format, "format", "xbrl", "Format to retrieve")
	if err := opt.Parse(args); err != nil {
		return "", "", err
	}
	if rawID == "" {
		rawID = opt.Arg(0)
	}

	if rawID == "" {
		return "", "", errors.New("query is mandatory")
	}

	id := parseRawID(rawID)
	if id == "" {
		return "", "", errors.New("invalid id: " + rawID)
	}

	var dataFormat ufocatcher.Format
	switch format {
	case "xbrl":
		dataFormat = ufocatcher.FormatData
	case "pdf":
		dataFormat = ufocatcher.FormatPDF
	}
	if dataFormat == "" {
		return "", "", errors.New("format is invalid: " + format)
	}

	return id, dataFormat, nil
}

func isaStdin() bool {
	stat, _ := os.Stdin.Stat()
	return stat.Mode()&os.ModeCharDevice == 0
}

func parseRawID(rawID string) string {
	exp := regexp.MustCompile(`(ED|TD|CG)\d+`)
	return exp.FindString(rawID)
}
