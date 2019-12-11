/*
Copyright © 2019 Shi Han NG <shihanng@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/shihanng/gig/internal/repo"
	"github.com/spf13/cobra"
)

func Execute(w io.Writer, version string) {
	command := &command{
		output:       w,
		templatePath: filepath.Join(xdg.CacheHome(), `gig`),
		version:      version,
	}

	rootCmd := newRootCmd(command)

	rootCmd.PersistentFlags().StringVarP(&command.commitHash, "commit-hash", "c", "",
		"use templates from a specific commit hash of github.com/toptal/gitignore")

	rootCmd.AddCommand(
		newListCmd(command),
		newGenCmd(command),
		newVersionCmd(command),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func newRootCmd(c *command) *cobra.Command {
	return &cobra.Command{
		Use:   "gig",
		Short: "A tool that generates .gitignore",
		Long: `gig is a command line tool to help you create useful .gitignore files
for your project. It is inspired by gitignore.io and make use of
the large collection of useful .gitignore templates of the web service.`,
		PersistentPreRunE: c.RootRunE,
	}
}

type command struct {
	output       io.Writer
	commitHash   string
	templatePath string
	version      string
}

func (c *command) RootRunE(cmd *cobra.Command, args []string) error {
	r, err := repo.New(c.templatePath, repo.SourceRepo)
	if err != nil {
		return err
	}

	ch, err := repo.Checkout(r, c.commitHash)
	if err != nil {
		return err
	}

	c.commitHash = ch

	return nil
}
