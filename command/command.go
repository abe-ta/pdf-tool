package command

import (
	"fmt"

	"github.com/urfave/cli"
)

type CLI interface {
	NewCommand() cli.Command
	Action(c *cli.Context) error
}

func NewCommands() []cli.Command {
	clis := []CLI{
		&SplitCommand{},
		&MergeCommand{},
	}

	var cmds []cli.Command
	for _, c := range clis {
		cmds = append(cmds, c.NewCommand())
	}
	return cmds
}

func validate(c *cli.Context, nargMin, nargMax int) error {
	narg := c.NArg()
	if narg < nargMin || (nargMax > 0 && narg > nargMax) {
		return fmt.Errorf("NArg is invalid: %d, NArgMax: %d, NArgMin: %d", narg, nargMax, nargMin)
	}
	return nil
}
