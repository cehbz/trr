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

	"github.com/charles-haynes/transmission"
	"github.com/spf13/cobra"
)

// peersCmd represents the peers command
var peersCmd = &cobra.Command{
	Use:   "peers",
	Short: "Info about the peers for a torrent",
	Long:  `For each torrent display detailed peer information`,
	Run:   doInfoPeers,
}

func doInfoPeers(cmd *cobra.Command, args []string) {
	x := getServer()
	c := transmission.NewGetTorrentsCmd()
	c.Arguments.Ids = getTorrents()
	c.Arguments.Fields = []string{"peers", "id", "name"}
	res, err := x.ExecuteCommand(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, t := range res.Arguments.Torrents {
		fmt.Printf("Torrent %d: %s\n", t.ID, t.Name)
		fmt.Println("Address         Flags   Done   Down     Up Client")
		for _, p := range t.Peers {
			fmt.Printf("%-15s %5s %5.1f%% %6.1f %6.1f %s\n",
				p.Address,
				p.Flags,
				p.Progress*100.0,
				float64(p.RateToClient)/1000.0,
				float64(p.RateToPeer)/1000.0,
				p.ClientName)
		}
	}
}

func init() {
	infoCmd.AddCommand(peersCmd)
}
