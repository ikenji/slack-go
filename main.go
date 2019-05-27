package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/ikenji/slack-go/slacklog"
)

var SLACK_TOKEN string = os.Getenv("SLACK_TOKEN")
var CHANNEL_ID string = os.Getenv("SAAS_NOTIFICATION_ID")

const (
	ENDPOINT string = "https://slack.com/api/conversations.history"
	LIMIT    string = "1000"
)

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		panic("Arguments count invalid")
	}

	beginTime, err := time.Parse("20060102", flag.Arg(0))
	endTime, err := time.Parse("20060102", flag.Arg(1))
	// add 23:59:59
	endTime = endTime.Add(time.Duration(86399) * time.Second)

	if err != nil {
		fmt.Printf("Some arguments are invalid. \nError info:\n%v\n", err)
		os.Exit(1)
	}

	params := url.Values{}
	params.Set("token", SLACK_TOKEN)
	params.Add("channel", CHANNEL_ID)
	params.Add("oldest", strconv.FormatInt(beginTime.Unix(), 10))
	params.Add("latest", strconv.FormatInt(endTime.Unix(), 10))
	params.Add("limit", LIMIT)

	res, _ := http.Get(ENDPOINT + "?" + params.Encode())
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var sa slacklog.SlackApi
	json.Unmarshal(body, &sa)

	raws := slacklog.Format(sa.Messages)

	for _, raw := range raws {
		fmt.Print(raw.Time.Format("2006-01-02 15:04:05") + ",")
		fmt.Print(raw.Store + ",")
		fmt.Print(raw.Mail + ",")
		fmt.Print(raw.Ip + ",")
		fmt.Print(raw.Kind + "\n")
	}
	os.Exit(0)
}
