package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Article struct {
	Type  string  `json:"type"`
	ID    string  `json:"harvesterId"`
	Score float64 `json:"cerebro-score"`
	URL   string  `json:"url"`
	Title string  `json:"title"`
	Img   string  `json:"cleanImage"`
}
type Articles []Article
type RespArticles struct {
	Status   int `json:"httpStatus"`
	Responce struct {
		Items Articles `json:"items"`
	} `json:"response"`
}

type ContentMarketing struct {
	Type              string  `json:"type"`
	ID                string  `json:"harvesterId"`
	CommercialPartner string  `json:"commercialPartner"`
	LogoURL           string  `json:"logoURL"`
	CerebroScore      float64 `json:"cerebro-score"`
	URL               string  `json:"url"`
}
type ContentMarketings []ContentMarketing
type RespContent struct {
	Status   int `json:"httpStatus"`
	Responce struct {
		Items ContentMarketings `json:"items"`
	} `json:"response"`
}

// response block
type ResponseObj struct {
	ArticlesArr         Articles
	ContentMarketingObj ContentMarketing
}
type Responses []ResponseObj

type ResponseStruct struct {
	Status   int `json:"httpStatus"`
	Response struct {
		Items Responses `json:"items"`
	} `json:"response"`
}

func (box *ResponseStruct) AddItem() {
	var temp ResponseObj
	box.Response.Items = append(box.Response.Items, temp)
}
func (box *ResponseStruct) AddArticle(article Article) {

	ref := box.Response.Items
	lastObj := len(ref) - 1
	ref[lastObj].ArticlesArr = append(ref[lastObj].ArticlesArr, article)
}
func (box *ResponseStruct) AddContentMarketing(c ContentMarketings, index int) {

	ref := box.Response.Items

	lastObj := len(ref) - 1
	lenC := len(c) - 1

	if lenC >= index {
		// replace by real content
		ref[lastObj].ContentMarketingObj = c[index]
	} else {
		// fake content
		ref[lastObj].ContentMarketingObj.Type = "Ad"
	}
}

// build response
func BuildResponse(a Articles, c ContentMarketings) ResponseStruct {

	var result ResponseStruct

	lenA := len(a)
	lenC := len(c)

	if lenA > 0 {

		fmt.Printf("len of Articles: %d \n", lenA)
		fmt.Printf("len of ContentMarketings: %d \n", lenC)

		each := 5
		count := -1

		// For each article
		for i, article := range a {

			// Create new item of response items for each new portion
			if i%each == 0 {
				fmt.Printf("each %d detected i == %d\n", each, i)
				result.AddItem()
				// incriment count of portion
				count++
			}

			// Process
			result.AddArticle(article)

			result.AddContentMarketing(c, count)
		}
	}

	return result

}

// get data from url
func GetArticles(str []byte) Articles {

	var data RespArticles
	json.Unmarshal(str, &data)

	return data.Responce.Items
}
func GetContents(str []byte) ContentMarketings {

	var data RespContent
	json.Unmarshal(str, &data)

	return data.Responce.Items
}
func parseBody(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		print(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return body
}

// server block
func homePage(w http.ResponseWriter, r *http.Request) {

	urlArticles := "https://storage.googleapis.com/aller-structure-task/articles.json"
	urlContents := "https://storage.googleapis.com/aller-structure-task/contentmarketing.json"
	articles := GetArticles(parseBody(urlArticles))
	contents := GetContents(parseBody(urlContents))

	result := BuildResponse(articles, contents)
	result.Status = 200
	json.NewEncoder(w).Encode(result)

}
func handleRequests() {
	fmt.Println("Hello Server")
	http.HandleFunc("/data", homePage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
func main() {
	handleRequests()
}
