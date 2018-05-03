package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	pdf "github.com/unidoc/unidoc/pdf/model"
	"github.com/urfave/cli"
)

type SplitCommand struct {
	output string
}

func (cmd *SplitCommand) NewCommand() cli.Command {
	return cli.Command{
		Name:      "split",
		Aliases:   []string{"s"},
		Usage:     "Split PDF pages.",
		ArgsUsage: "input.pdf <FROM> <TO>",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "output, o",
				Usage:       "output file `path`.",
				Destination: &cmd.output,
			},
		},
		Action: cmd.Action,
	}
}

func (cmd *SplitCommand) Action(c *cli.Context) error {
	if err := validate(c, 3, 3); err != nil {
		return err
	}
	from, err := strconv.Atoi(c.Args().Get(1))
	if err != nil {
		return fmt.Errorf("`from` is not a number: %v", err)
	}
	to, err := strconv.Atoi(c.Args().Get(2))
	if err != nil {
		return fmt.Errorf("`to` is not a number: %v", err)
	}
	if from > to {
		return fmt.Errorf("`from` is bigger than `to`: from: %d, to: %d", from, to)
	}
	input := c.Args().Get(0)
	output := cmd.output
	if output == "" {
		ext := filepath.Ext(input)
		output = strings.Split(input, ext)[0] + "_splitted" + ext
	}

	return cmd.splitPDF(input, output, from, to)
}

func (cmd *SplitCommand) splitPDF(inputPath string, outputPath string, from, to int) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}
	if isEncrypted {
		_, err = pdfReader.Decrypt([]byte(""))
		if err != nil {
			return err
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil || numPages < to {
		return err
	}

	pdfWriter := pdf.NewPdfWriter()
	for i := from; i <= to; i++ {
		pageNum := i

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}
		if err = pdfWriter.AddPage(page); err != nil {
			return err
		}
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer fWrite.Close()

	return pdfWriter.Write(fWrite)
}
