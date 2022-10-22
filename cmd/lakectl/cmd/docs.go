package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// language=markdown
const cliReferenceHeader = `---
layout: default
title: Command (CLI) Reference
description: lakeFS comes with its own native CLI client. Here you can see the complete command reference.
parent: Reference
nav_order: 3
has_children: false
---

# Commands (CLI) Reference
{:.no_toc}

{% include toc.html %}

## Installing the lakectl command locally

` + "`lakectl`" + ` is distributed as a single binary, with no external dependencies - and is available for MacOS, Windows and Linux.

[Download lakectl](../index.md#downloads){: .btn .btn-green target="_blank"}


### Configuring credentials and API endpoint

Once you've installed the lakectl command, run:

` + "```" + `bash
lakectl config
# output:
# Config file /home/janedoe/.lakectl.yaml will be used
# Access key ID: AKIAIOSFODNN7EXAMPLE
# Secret access key: ****************************************
# Server endpoint URL: http://localhost:8000/api/v1
` + "```" + `

This will setup a ` + "`$HOME/.lakectl.yaml`" + ` file with the credentials and API endpoint you've supplied.
When setting up a new installation and creating initial credentials (see [Quick start](../quickstart/index.md)), the UI
will provide a link to download a preconfigured configuration file for you.

` + "`lakectl`" + ` configuration items can each be controlled by an environment variable. The variable name will have a prefix of
*LAKECTL_*, followed by the name of the configuration, replacing every '.' with a '_'. Example: ` + "`LAKECTL_SERVER_ENDPOINT_URL`" + ` 
controls ` + "`server.endpoint_url`" + `.
`

func printOptions(buf *bytes.Buffer, cmd *cobra.Command) error {
	flags := cmd.NonInheritedFlags()
	flags.SetOutput(buf)
	if flags.HasAvailableFlags() {
		buf.WriteString("#### Options\n{:.no_toc}\n\n```\n")
		flags.PrintDefaults()
		buf.WriteString("```\n\n")
	}
	parentFlags := cmd.InheritedFlags()
	parentFlags.SetOutput(buf)

	if cmd == rootCmd {
		buf.WriteString("**note:** The `base-uri` option can be controlled with the `LAKECTL_BASE_URI` environment variable.\n{: .note .note-warning }\n\n")
		buf.WriteString("#### Example usage\n{:.no_toc}\n\n")
		buf.WriteString("```shell\n$ export LAKECTL_BASE_URI=\"lakefs://my-repo/my-branch\"\n# Once set, use relative lakefs uri's:\n$ lakectl fs ls /path\n```")
	}
	return nil
}

// this is a version of github.com/spf13/cobra/doc that fits better
// into lakeFS' docs.
func genMarkdownForCmd(cmd *cobra.Command, w io.Writer) error {
	cmd.InitDefaultHelpCmd()
	cmd.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()

	buf.WriteString("### " + name + "\n\n")

	if cmd.Hidden {
		buf.WriteString("**note:** This command is a lakeFS plumbing command. Don't use it unless you're really sure you know what you're doing.\n{: .note .note-warning }\n\n")
	}

	buf.WriteString(cmd.Short + "\n\n")
	if len(cmd.Long) > 0 {
		buf.WriteString("#### Synopsis\n{:.no_toc}\n\n")
		buf.WriteString(cmd.Long + "\n\n")
	}

	if cmd.Runnable() {
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.UseLine()))
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("#### Examples\n{:.no_toc}\n\n")
		buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", cmd.Example))
	}

	if err := printOptions(buf, cmd); err != nil {
		return err
	}

	buf.WriteString("\n\n")
	_, err := buf.WriteTo(w)
	if err != nil {
		return err
	}

	// recurse to children
	for _, c := range cmd.Commands() {
		err := genMarkdownForCmd(c, w)
		if err != nil {
			return err
		}
	}

	return nil
}

var docsCmd = &cobra.Command{
	Use:    "docs",
	Hidden: true,
	Args:   cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		output := MustString(cmd.Flags().GetString("output"))
		format := MustString(cmd.Flags().GetString("format"))

		switch format {
		case "md":
			// output is stdout or file specified by output flag
			var writer io.Writer
			if output == "" {
				writer = os.Stdout
			} else {
				f, err := os.Create(output)
				if err != nil {
					return err
				}
				defer func() {
					err := f.Close()
					if err != nil {
						DieErr(err)
					}
				}()
				writer = f
			}
			return GenMarkdown(writer)
		case "man":
			// write man files into output folder
			return doc.GenManTree(rootCmd, &doc.GenManHeader{
				Title:   "MINE",
				Section: "3",
				Source:  "Auto generated by lakectl docs command",
			}, output)
		default:
			DieFmt("Unknown output format '%s'", format)
		}
		return nil
	},
}

func GenMarkdown(writer io.Writer) error {
	_, err := io.WriteString(writer, cliReferenceHeader)
	if err != nil {
		return err
	}
	return genMarkdownForCmd(rootCmd, writer)
}

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.Flags().StringP("output", "o", "", "Output location (based on the format file or directory)")
	docsCmd.Flags().String("format", "md", "Output format (md markdown or man files)")
}
