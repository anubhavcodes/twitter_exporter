package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	twitterFollowers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "twitter_followers",
		Help: "Total number of twitter followers",
	})
	twitterFollowing = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "twitter_following",
		Help: "Total number of accounts followed",
	})
	tweetCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "tweet_count",
		Help: "Total number of tweets",
	})
	listedCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "listed_count",
		Help: "Number of list the user belongs to",
	})
)

func recordMetrics() {
	go func() {
		for {
			err, body := getTwitterData()
			if err != nil {
				panic(err)
			}
			metrics := getPublicMetrics(body)
			twitterFollowers.Set(float64(metrics.FollowersCount))
			twitterFollowing.Set(float64(metrics.FollowingCount))
			tweetCount.Set(float64(metrics.TweetCount))
			listedCount.Set(float64(metrics.ListedCount))
			fmt.Println(metrics.FollowingCount)
			time.Sleep(10 * time.Second)
		}
	}()
}

func main() {
	recordMetrics()
	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":8080", nil)

}

func getTwitterData() (error, []byte) {
	var TwitterToken = os.Getenv("TWITTER_TOKEN")
	var TwitterHandle = os.Getenv("TWITTER_HANDLE")
	if TwitterToken == "" || TwitterHandle == "" {
		log.Println("Please define TWITTER_TOKEN and the TWITTER_HANDLE environment variable")
		os.Exit(100)
	}
	client := &http.Client{}
	url := "https://api.twitter.com/2/users/by/username/" + TwitterHandle + "?user.fields=public_metrics"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Bearer "+TwitterToken)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return err, body
}

type APIResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	Id            string        `json:"id"`
	Name          string        `json:"name"`
	Username      string        `json:"username"`
	PublicMetrics PublicMetrics `json:"public_metrics"`
}

type PublicMetrics struct {
	FollowersCount int `json:"followers_count"`
	FollowingCount int `json:"following_count"`
	TweetCount     int `json:"tweet_count"`
	ListedCount    int `json:"listed_count"`
}

func getPublicMetrics(apiResponse []byte) PublicMetrics {
	var response APIResponse
	if err := json.Unmarshal(apiResponse, &response); err != nil {
		log.Println(err)
	}
	return response.Data.PublicMetrics

}
