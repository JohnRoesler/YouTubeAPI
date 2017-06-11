package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var listOfVideos []string
var listOfCommenters []string

const apiKey = ""

func main() {

	gitYouTubeVideos()
	//fmt.Println(listOfVideos)
	gitYouTubeComments()

	fmt.Println("Writing list of commenters")
	filename := ""
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file")
	}
	w := bufio.NewWriter(f)
	for x := 0; x < len(listOfCommenters); x++ {
		w.WriteString(listOfCommenters[x])
	}
	fmt.Println("Wrote list of commenters")

}

func gitYouTubeVideos() {

	type ChannelVideos struct {
		NextPageToken string `json:"nextPageToken"`
		PageInfo      struct {
			TotalResults   int `json:"totalResults"`
			ResultsPerPage int `json:"resultsPerPage"`
		} `json:"pageInfo"`
		Items []struct {
			ID struct {
				VideoID string `json:"videoId"`
			} `json:"id"`
		}
	}

	url := "https://www.googleapis.com/youtube/v3/search?key=" + apiKey + "&channelId=UCrOtGhui_jdLdoQNI7PU4Pg&part=snippet,id&order=date"

	req, _ := http.NewRequest("GET", url, nil)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode > 299 {
		log.Println("Error getting youtube videos:", res.StatusCode, ",", res.Status)
	}

	var item ChannelVideos
	err := json.Unmarshal(body, &item)
	if err != nil {
		fmt.Println("Error")
	}

	for x := 0; x < len(item.Items); x++ {
		listOfVideos = append(listOfVideos, item.Items[x].ID.VideoID)
	}

	for x := 2; x <= item.PageInfo.TotalResults/item.PageInfo.ResultsPerPage; x++ {
		fmt.Println("Getting list of videos working on page:", x)
		url := "https://www.googleapis.com/youtube/v3/search?key=" + apiKey + "&channelId=UCrOtGhui_jdLdoQNI7PU4Pg&part=snippet,id&order=date&pageToken=" + item.NextPageToken

		req, _ := http.NewRequest("GET", url, nil)

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		if res.StatusCode > 299 {
			log.Println("Error getting youtube videos:", res.StatusCode, ",", res.Status)
		}

		var item ChannelVideos
		err := json.Unmarshal(body, &item)
		if err != nil {
			log.Println("Error:", err)
		}
		for x := 0; x < len(item.Items); x++ {
			listOfVideos = append(listOfVideos, item.Items[x].ID.VideoID)
		}
	}

}

func gitYouTubeComments() {

	type VideoCommenters struct {
		NextPageToken string `json:"nextPageToken"`
		PageInfo      struct {
			TotalResults   int `json:"totalResults"`
			ResultsPerPage int `json:"resultsPerPage"`
		} `json:"pageInfo"`
		Items []struct {
			Snippet struct {
				TopLevelComment struct {
					Snippet struct {
						AuthorDisplayName string `json:"authorDisplayName"`
					} `json:"snippet"`
				} `json:"topLevelComment"`
			} `json:"snippet"`
			Replies struct {
				Comments []struct {
					Snippet struct {
						AuthorDisplayName string `json:"authorDisplayName"`
					} `json:"snippet"`
				} `json:"comments"`
			} `json:"replies"`
		} `json:"items"`
	}
	for x := 0; x < len(listOfVideos); x++ {
		fmt.Println("Working on list of commenters, on #", x, ", and video:", listOfVideos[x])

		if listOfVideos[x] == "" {
			continue
		}

		url := "https://www.googleapis.com/youtube/v3/commentThreads?part=snippet%2Creplies&key=" + apiKey + "&videoId=" + listOfVideos[x]

		req, _ := http.NewRequest("GET", url, nil)

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		if res.StatusCode > 299 {
			log.Println("Error getting commenters:", res.StatusCode, ",", res.Status, ",", listOfVideos[x])
		}

		var item VideoCommenters
		err := json.Unmarshal(body, &item)
		if err != nil {
			log.Println("Error:", err)
		}
		for x := 0; x < len(item.Items); x++ {
			listOfCommenters = append(listOfCommenters, item.Items[x].Snippet.TopLevelComment.Snippet.AuthorDisplayName+"\n")
			for y := 0; y < len(item.Items[x].Replies.Comments); y++ {
				listOfCommenters = append(listOfCommenters, item.Items[x].Replies.Comments[y].Snippet.AuthorDisplayName+"\n")
			}
		}
	}

}
