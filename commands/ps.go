package commands

import (
	"strings"

	"github.com/codegangsta/cli"
	"github.com/vinceprignano/orchestra/services"
	"github.com/wsxiaoys/terminal"
)

var PsCommand = &cli.Command{
	Name:   "ps",
	Usage:  "Outputs the status of all services",
	Action: BeforeAfterWrapper(PsAction),
}

// PsAction checks the status for every service and output
func PsAction(c *cli.Context) {
	for name, service := range FilterServices(c) {
		spacing := strings.Repeat(" ", services.MaxServiceNameLength+2-len(service.Name))
		if service.Process != nil {
			terminal.Stdout.Colorf("@{g}%s", name).Reset().Colorf("%s|", spacing).Print(" running ").Colorf("  %d\n", service.Process.Pid)
		} else {
			terminal.Stdout.Colorf("@{r}%s", name).Reset().Colorf("%s|", spacing).Reset().Print(" aborted\n")
		}
	}
}
