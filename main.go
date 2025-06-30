package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"

	//"io"
	"bytes"
	"math/rand"
	"os"
	"regexp"
	"time"
)

type Chapter struct {
	Name string `json:"title"`
	Link string `json:"url"`
}

func main() {
	fmt.Println("Starting Program")

	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments given")
	}
	rand.Seed(time.Now().UnixNano())
	url := os.Args[1]
	domain := strings.SplitAfter(url, ".com")[0]
	bookName := strings.SplitAfter(url, "/")[5]
	fmt.Println("filename: " + bookName)

	client := &http.Client{}
	allChapters := []Chapter{}

	doc := requestPageHTML(url, client)

	getAllChaptersInfo(doc, domain, &allChapters)

	fmt.Println("All chapters:")
	printStuff(allChapters)

	callChapter(allChapters, client, bookName)

}

// Send a get request to a url and returns the a note tree
func requestPageHTML(url string, client *http.Client) *html.Node {
	var response http.Response
	tryCount := 0

	for response.StatusCode != 200 {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("Making Request to ", url)

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Referer", "https://www.google.com")

		resp, err := client.Do(req)

		if err != nil {
			log.Fatal("Can not connect to url")
		}
		defer resp.Body.Close()

		fmt.Println("Status Code: ", resp.StatusCode)

		if resp.StatusCode != 200 {
			if tryCount > 5 {
				log.Fatal("Cant make 200 connection")
			}
			tryCount++
			continue
		}

		response = *resp
	}

	doc, err := html.Parse(response.Body)
	if err != nil {
		log.Fatal("Failed parsing html")
	}

	return doc

}

// Gets us all the pages that have 100 chapters each
func getAllChaptersInfo(doc *html.Node, domain string, allChapters *[]Chapter) {
	for item := range doc.Descendants() {
		if item.Type == html.TextNode {
			if len(item.Data) > 27 {
				if strings.TrimSpace(item.Data[:27]) == "window.fiction" {

					re := regexp.MustCompile(`window\.chapters\s*=\s*(\[[\s\S]*?\]);`)

					jsons := re.FindStringSubmatch(item.Data)

					if len(jsons) < 2 {
						log.Fatal("err")
					}

					jsonArray := jsons[1]

					err := json.Unmarshal([]byte(jsonArray), &allChapters)

					if err != nil {
						log.Fatal("error parsing chapters json")
					}

				}
			}

		}

	}

	for index, item := range *allChapters {
		(*allChapters)[index].Link = domain + item.Link
	}

}

// Goes to each chapter and calls its content
func callChapter(allChapters []Chapter, client *http.Client, bookName string) {
	fullHTML := ""

	for index := 0; index < len(allChapters); index++ {

		tempDoc := requestPageHTML(allChapters[index].Link, client)

		var chapContent *html.Node

		for node := range tempDoc.Descendants() {
			if node.Type != html.ElementNode || node.Data != "div" {
				continue
			}
			for _, ele := range node.Attr {
				if ele.Key != "class" || ele.Val != "chapter-inner chapter-content" {
					continue
				}

				node.PrevSibling = nil
				node.NextSibling = nil
				fmt.Println("Saving content: ", node)
				fmt.Println("-------------------")
				chapContent = node
				break
			}

		}

		var buf bytes.Buffer
		if err := html.Render(&buf, chapContent); err != nil {
			log.Fatal(err)
		}

		partHTML := "<h2>" + allChapters[index].Name + "</h2>" + buf.String()
		fullHTML = fullHTML + partHTML
	}

	finishedHTML := "<!DOCTYPE html>\n<html>\n<head><meta charset=\"UTF-8\"></head>\n<body>\n" + fullHTML + "\n</body>\n</html>"
	//fmt.Println(finishedHTML)
	os.WriteFile(bookName+".html", []byte(finishedHTML), 0644)
}

func printStuff(chaps []Chapter) {
	for _, chap := range chaps {
		fmt.Println(chap.Name)
	}
}
