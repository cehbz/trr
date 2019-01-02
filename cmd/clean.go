// Copyright © 2017 Charles Haynes <ceh@ceh.bz>
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
	"fmt"
	"strconv"
	"strings"

	"github.com/charles-haynes/transmission"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up unregistered torrents",
	Long: `Clean up unregistered torrents

Remove any unregistered torrents that have the same name as a registered 
torrent. Print a list of all other unregistered torrents. Takes list of torrent 
specifiers, Defaults to all.`,
	Run: doClean,
}

func doClean(cmd *cobra.Command, args []string) {
	x := getServer()
	c := transmission.NewGetTorrentsCmd()
	c.Arguments.Ids = getTorrents()
	c.Arguments.Fields = []string{"errorString", "hashString", "id", "name", "status"}
	res, err := x.ExecuteCommand(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	ts := res.Arguments.Torrents
	r := map[string]interface{}{}
	for _, t := range ts {
		if t.ErrorString != "Unregistered torrent" {
			r[t.Name] = nil
		}
	}
	fmt.Printf("# %d registered torrents\n", len(r))
	d := []string{}
	for _, t := range ts {
		if t.ErrorString == "Unregistered torrent" {
			d = append(d, strconv.Itoa(t.ID))
			if _, ok := r[t.Name]; !ok {
				l := fmt.Sprintf(
					`cp -v "%s.%s.torrent" ~/uploadable`,
					t.Name,
					t.Hash[0:16])
				fmt.Println(l)
			}
		}
	}
	fmt.Printf("transmission-remote %s -t %s -r", server, strings.Join(d, ","))
}

func init() {
	RootCmd.AddCommand(cleanCmd)
}
