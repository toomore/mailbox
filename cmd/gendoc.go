// Copyright © 2017 Toomore Chiang
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var genbash *bool

// gendocCmd represents the gendoc command
var gendocCmd = &cobra.Command{
	Use:   "doc",
	Short: "Gen docs md",
	Run: func(cmd *cobra.Command, args []string) {
		if *genbash {
			RootCmd.GenBashCompletionFile("./mailbox")
			cmd.Println("Gen Bash Completion File ...")
		} else {
			doc.GenMarkdownTree(RootCmd, "./")
			cmd.Println("Gen Markdown Tree ...")
		}
	},
	Hidden: true,
}

func init() {
	genbash = gendocCmd.Flags().BoolP("genbash", "b", false, "Gen bash completion")
	RootCmd.AddCommand(gendocCmd)
}
