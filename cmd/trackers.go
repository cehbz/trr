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
	"math"
	"time"

	"github.com/charles-haynes/transmission"
	"github.com/spf13/cobra"
)

// trackersCmd represents the trackers command
var trackersCmd = &cobra.Command{
	Use:   "trackers",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: doInfoTrackers,
}

const (
	secsPerMin  = 60
	secsPerHr   = secsPerMin * 60
	secsPerDay  = secsPerHr * 24
	secsPerYear = secsPerDay*365 + secsPerHr*6
)

func myDurationTill(t int64) string {
	switch true {
	case t < 0, t == math.MaxInt64:
		return "∞"
	case t == 0:
		return "never"
	default:
		return myDuration(t - time.Now().Unix())
	}
}

func myDurationSince(t int64) string {
	switch true {
	case t < 0, t == math.MaxInt64:
		return "∞"
	case t == 0:
		return "never"
	default:
		return myDuration(time.Now().Unix() - t)
	}
}

func myDuration(d int64) string {
	switch true {
	case d < 0, d == math.MaxInt64:
		return "∞"
	case d == 0:
		return ""
	case d < 100:
		return fmt.Sprintf("%d secs", d)
	case d < 100*secsPerMin:
		return fmt.Sprintf("%4.1f mins", float64(d)/float64(secsPerMin))
	case d < 100*secsPerHr:
		return fmt.Sprintf("%4.1f hrs", float64(d)/float64(secsPerHr))
	case d < secsPerYear:
		return fmt.Sprintf("%5.1f days", float64(d)/float64(secsPerDay))
	case d < 100*secsPerYear:
		return fmt.Sprintf("%5.1f yrs", float64(d)/float64(secsPerYear))
	default:
		return fmt.Sprintf("%d secs", d)
	}
}

func doInfoTrackers(cmd *cobra.Command, args []string) {
	x := getServer()
	c := transmission.NewGetTorrentsCmd()
	c.Arguments.Ids = getTorrents()
	c.Arguments.Fields = []string{"trackerStats", "id", "name"}
	res, err := x.ExecuteCommand(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, t := range res.Arguments.Torrents {
		fmt.Printf("Torrent %d: %s\n", t.ID, t.Name)
		fmt.Println("Tier Peers Se Le    Last Sc    Next Sc   Last Ann   Next Ann Name")
		for _, t := range t.TrackerStats {
			fmt.Printf("%4d %5d %2d %2d %10s %10s %10s %10s %s\n",
				t.Tier,
				t.LastAnnouncePeerCount,
				t.SeederCount,
				t.LeecherCount,
				myDurationSince(t.LastScrapeTime),
				myDurationTill(t.NextScrapeTime),
				myDurationSince(t.LastAnnounceTime),
				myDurationTill(t.NextAnnounceTime),
				t.Host)
		}
	}
}

func init() {
	infoCmd.AddCommand(trackersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// trackersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// trackersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
