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
const channelID = ""

func main() {

	gitYouTubeVideos()
	gitYouTubeComments()

	fmt.Println("Writing list of commenters")
	filename := "listOfCommenters.txt"
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

	url := "https://www.googleapis.com/youtube/v3/search?key=" + apiKey + "&channelId=" + channelID + "&part=snippet&order=date"

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
	fmt.Println("Getting list of videos working on page:", 0)
	for x := 0; x < len(item.Items); x++ {
		listOfVideos = append(listOfVideos, item.Items[x].ID.VideoID)
	}
	x := 1
	for item.NextPageToken != "" {
		fmt.Println("Getting list of videos working on page:", x)
		x++
		url := "https://www.googleapis.com/youtube/v3/search?key=" + apiKey + "&channelId=" + channelID + "&part=snippet,id&order=date&pageToken=" + item.NextPageToken

		req, _ := http.NewRequest("GET", url, nil)

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		if res.StatusCode > 299 {
			log.Println("Error getting youtube videos:", res.StatusCode, ",", res.Status)
		}
		item = ChannelVideos{}
		err := json.Unmarshal(body, &item)
		if err != nil {
			log.Println("Error:", err)
		}
		for y := 0; y < len(item.Items); y++ {
			if item.Items[y].ID.VideoID == "" {
				continue
			}
			listOfVideos = append(listOfVideos, item.Items[y].ID.VideoID)
		}
	}

}

func gitYouTubeComments() {

	type VideoCommenters struct {
		NextPageToken string `json:"nextPageToken"`
		Items         []struct {
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
		fmt.Println("Working on list of commenters, on #", x, " page 0, and video:", listOfVideos[x])

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
			if item.Items[x].Snippet.TopLevelComment.Snippet.AuthorDisplayName != "One Stop Co-op Shop" {
				listOfCommenters = append(listOfCommenters, item.Items[x].Snippet.TopLevelComment.Snippet.AuthorDisplayName+"\n")
			}
			for y := 0; y < len(item.Items[x].Replies.Comments); y++ {
				if item.Items[x].Replies.Comments[y].Snippet.AuthorDisplayName != "One Stop Co-op Shop" {
					listOfCommenters = append(listOfCommenters, item.Items[x].Replies.Comments[y].Snippet.AuthorDisplayName+"\n")
				}
			}
		}
		y := 1
		for item.NextPageToken != "" {
			fmt.Println("Working on list of commenters, on #", x, " page", y, ", and video:", listOfVideos[x])
			y++
			url := "https://www.googleapis.com/youtube/v3/commentThreads?part=snippet%2Creplies&key=" + apiKey + "&videoId=" + listOfVideos[x] + "&pageToken=" + item.NextPageToken

			req, _ := http.NewRequest("GET", url, nil)

			res, _ := http.DefaultClient.Do(req)

			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode > 299 {
				log.Println("Error getting commenters:", res.StatusCode, ",", res.Status, ",", listOfVideos[x])
			}

			item = VideoCommenters{}
			err := json.Unmarshal(body, &item)
			if err != nil {
				log.Println("Error:", err)
			}
			for x := 0; x < len(item.Items); x++ {
				if item.Items[x].Snippet.TopLevelComment.Snippet.AuthorDisplayName != "One Stop Co-op Shop" {
					listOfCommenters = append(listOfCommenters, item.Items[x].Snippet.TopLevelComment.Snippet.AuthorDisplayName+"\n")
				}
				for y := 0; y < len(item.Items[x].Replies.Comments); y++ {
					if item.Items[x].Replies.Comments[y].Snippet.AuthorDisplayName != "One Stop Co-op Shop" {
						listOfCommenters = append(listOfCommenters, item.Items[x].Replies.Comments[y].Snippet.AuthorDisplayName+"\n")
					}
				}
			}
		}
	}

}
