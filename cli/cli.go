package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/aaqaishtyaq/gotopdf/topdf"
	"github.com/urfave/cli/v2"
)

// Gotopdf cli
func Gotopdf() {
  app := cli.NewApp()

  app.Flags = []cli.Flag {
    &cli.StringFlag{
      Name: "html",
      Value: "example.com",
      Usage: "html for conversion",
      Required: true,
    },
  }

  app.Action = func(c *cli.Context) error {
    var output string
    output = "PDF Created"

    fmt.Println(output)
    return nil
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}
