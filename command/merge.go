package command

import (
	"errors"
	"os"

	pdf "github.com/unidoc/unidoc/pdf/model"
	"github.com/urfave/cli"
)

type MergeCommand struct {
	output string
}

func (cmd *MergeCommand) NewCommand() cli.Command {
	return cli.Command{
		Name:      "merge",
		Aliases:   []string{"m"},
		Usage:     "Merge PDF files.",
		ArgsUsage: "input1.pdf input2.pdf [input3.pdf ...]",
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

func (cmd *MergeCommand) Action(c *cli.Context) error {
	if err := validate(c, 2, -1); err != nil {
		return err
	}

	return cmd.mergePdf(c.Args(), cmd.output)
}

func (cmd *MergeCommand) mergePdf(inputPaths []string, outputPath string) error {
	pdfWriter := pdf.NewPdfWriter()

	for _, inputPath := range inputPaths {
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
			auth, err := pdfReader.Decrypt([]byte(""))
			if err != nil {
				return err
			}
			if !auth {
				return errors.New("Cannot merge encrypted, password protected document")
			}
		}

		numPages, err := pdfReader.GetNumPages()
		if err != nil {
			return err
		}

		for i := 0; i < numPages; i++ {
			pageNum := i + 1

			page, err := pdfReader.GetPage(pageNum)
			if err != nil {
				return err
			}

			err = pdfWriter.AddPage(page)
			if err != nil {
				return err
			}
		}
	}

	fWrite, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer fWrite.Close()

	err = pdfWriter.Write(fWrite)
	if err != nil {
		return err
	}

	return nil
}
