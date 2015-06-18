package commands

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/mondough/orchestra/services"
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
			terminal.Stdout.Colorf("@{g}%s", name).Reset().Colorf("%s|", spacing).Print(" running ").Colorf("  %d  %s\n", service.Process.Pid, getPorts(service))
		} else {
			terminal.Stdout.Colorf("@{r}%s", name).Reset().Colorf("%s|", spacing).Reset().Print(" aborted\n")
		}
	}
}

func getPorts(service *services.Service) string {
	re := regexp.MustCompile("LISTEN")
	cmd := exec.Command("lsof", "-p", fmt.Sprintf("%d", service.Process.Pid))
	output := bytes.NewBuffer([]byte{})
	cmd.Stdout = output
	cmd.Stderr = output
	err := cmd.Run()
	if err != nil {
		return ""
	}
	lsofOutput := ""
	for {
		s, err := output.ReadString('\n')
		if err == io.EOF {
			break
		}
		matched := re.MatchString(s)
		if matched {
			fields := strings.Fields(s)
			lsofOutput += fmt.Sprintf("%s/%s ", fields[8], strings.ToLower(fields[7]))
		}
	}
	return lsofOutput
}
