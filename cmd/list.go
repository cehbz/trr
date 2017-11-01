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
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charles-haynes/transmission"
	units "github.com/docker/go-units"
	"github.com/spf13/cobra"
)

var sortBy string

type myTorrents transmission.Torrents

type sorter func(t myTorrents, i, j int) bool

var less sorter

func (t myTorrents) Len() int           { return len(t) }
func (t myTorrents) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t myTorrents) Less(i, j int) bool { return less(t, i, j) }

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List torrents",
	Long: `List torrents for a tracker, allows filtering and sorting
For Example:

trr list - list all torrents for default tracker
trr list -sort active - list all torrents sorted by the time they were last active
trr list -filter uploading -sort added,name - list uloading torrents sorted by
    when they were added, and then by name`,
	Run: doList,
}

func getServer() *transmission.TransmissionClient {
	a := fmt.Sprintf("http://%s/transmission/rpc", server)
	x, err := transmission.New(a, user, pass)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return x
}

func getTorrents() []int {
	if torrents == "all" {
		return nil
	}
	idStrings := strings.Split(torrents, ",")
	ids := make([]int, len(idStrings))
	for i, id := range idStrings {
		t, err := strconv.Atoi(id)
		if err != nil {
			log.Print(err)
			return nil
		}
		ids[i] = t
	}
	return ids
}

func doList(cmd *cobra.Command, args []string) {
	x := getServer()
	c := transmission.NewGetTorrentsCmd()
	c.Arguments.Ids = getTorrents()
	c.Arguments.Fields = []string{
		"addedDate",
		"error",
		"errorString",
		"eta",
		"haveUnchecked",
		"haveValid",
		"id",
		"isFinished",
		"leftUntilDone",
		"name",
		"peersGettingFromUs",
		"peersSendingToUs",
		"percentDone",
		"rateDownload",
		"rateUpload",
		"sizeWhenDone",
		"status",
		"uploadRatio",
	}
	res, err := x.ExecuteCommand(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	ts := res.Arguments.Torrents
	if sortBy != "" {
		less = getSorter(sortBy)
		sort.Sort(myTorrents(ts))
	}
	// ID     Done       Have  ETA           Up    Down  Ratio  Status       Name
	//   11    16%   618.8 MB  Unknown      0.0     7.0    0.0  Up & Down    Leo Kottke
	fmt.Println("   ID Done      Have       ETA      Up    Down Ratio Status        Name")
	for _, t := range ts {
		fmt.Printf("%5d %3.0f%% %9s %10s %8s %8s %5.1f %-13s %s\n",
			t.ID,
			t.PercentDone*100.0,
			units.HumanSize(float64(t.Have())),
			myDuration(myETA(t)),
			units.HumanSize(float64(t.RateUpload)),
			units.HumanSize(float64(t.RateDownload)),
			t.UploadRatio,
			Status(t),
			t.Name)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getSorter(s string) sorter {
	switch s {
	case "id":
		return func(t myTorrents, i, j int) bool { return t[i].ID < t[j].ID }
	case "-id":
		return func(t myTorrents, i, j int) bool { return t[i].ID > t[j].ID }
	case "name":
		return func(t myTorrents, i, j int) bool { return t[i].Name < t[j].Name }
	case "-name":
		return func(t myTorrents, i, j int) bool { return t[i].Name > t[j].Name }
	case "age":
		return func(t myTorrents, i, j int) bool { return t[i].AddedDate < t[j].AddedDate }
	case "-age":
		return func(t myTorrents, i, j int) bool { return t[i].AddedDate > t[j].AddedDate }
	case "size":
		return func(t myTorrents, i, j int) bool { return t[i].SizeWhenDone < t[j].SizeWhenDone }
	case "-size":
		return func(t myTorrents, i, j int) bool { return t[i].SizeWhenDone > t[j].SizeWhenDone }
	case "progress":
		return func(t myTorrents, i, j int) bool { return t[i].PercentDone < t[j].PercentDone }
	case "-progress":
		return func(t myTorrents, i, j int) bool { return t[i].PercentDone > t[j].PercentDone }
	case "downspeed":
		return func(t myTorrents, i, j int) bool { return t[i].RateDownload < t[j].RateDownload }
	case "-downspeed":
		return func(t myTorrents, i, j int) bool { return t[i].RateDownload > t[j].RateDownload }
	case "upspeed":
		return func(t myTorrents, i, j int) bool { return t[i].RateUpload < t[j].RateUpload }
	case "-upspeed":
		return func(t myTorrents, i, j int) bool { return t[i].RateUpload > t[j].RateUpload }
	case "downloaded":
		return func(t myTorrents, i, j int) bool { return t[i].DownloadedEver < t[j].DownloadedEver }
	case "-downloaded":
		return func(t myTorrents, i, j int) bool { return t[i].DownloadedEver > t[j].DownloadedEver }
	case "uploaded":
		return func(t myTorrents, i, j int) bool { return t[i].UploadedEver < t[j].UploadedEver }
	case "-uploaded":
		return func(t myTorrents, i, j int) bool { return t[i].UploadedEver > t[j].UploadedEver }
	case "ratio":
		return func(t myTorrents, i, j int) bool { return t[i].UploadRatio < t[j].UploadRatio }
	case "-ratio":
		return func(t myTorrents, i, j int) bool { return t[i].UploadRatio > t[j].UploadRatio }
	case "eta":
		return func(t myTorrents, i, j int) bool { return myETA(t[i]) < myETA(t[j]) }
	case "-eta":
		return func(t myTorrents, i, j int) bool { return myETA(t[i]) > myETA(t[j]) }
	default:
		return func(t myTorrents, i, j int) bool { return i < j }
	}
}

func myETA(t *transmission.Torrent) int64 {
	switch true {
	case t.LeftUntilDone == 0 || t.Eta == 0:
		return 0
	case t.Eta > 0:
		return int64(t.Eta)
	case t.RateDownload > 0:
		return int64(t.LeftUntilDone / t.RateDownload)
	case t.PercentDone > 0.0:
		return int64((1.0/t.PercentDone - 1.0) * float64(time.Now().Unix()-t.AddedDate))
	default:
		return math.MaxInt64
	}
}

// Status prints a human readable status for the torrent
func Status(t *transmission.Torrent) string {
	if t.ErrorString != "" {
		return t.ErrorString
	}
	if t.Error != 0 {
		return "Error"
	}
	switch t.Status {
	case transmission.Stopped:
		return "Stopped"
	case transmission.CheckPending:
		return "Wait Verify"
	case transmission.Checking:
		return "Verifying"
	case transmission.DownloadPending:
		return "Wait Download"
	case transmission.SeedPending:
		return "Wait Seed"
	case transmission.Seeding, transmission.Downloading:
		switch true {
		case t.RateDownload == 0 && t.RateUpload == 0:
			return "Idle"
		case t.RateDownload > 0 && t.RateUpload > 0:
			return "Both"
		case t.RateDownload > 0:
			return "Downloading"
		case t.RateUpload > 0:
			return "Uploading"
		default:
			return "Wat?"
		}
	default:
		return "unknown"
	}
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")
	listCmd.PersistentFlags().StringVar(&sortBy, "sort", "", "what field to sort on")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
