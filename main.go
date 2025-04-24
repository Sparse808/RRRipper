package main

import (
	"fmt"
	"net/http"
	"log"
	"golang.org/x/net/html"
	"strings"
	"strconv"
	//"io"
	"os"
	"bytes"
	"time"
	"math/rand"
)

type Chapter struct {
	Name string
	Link string
}

func main() {
	fmt.Println("Starting Program")
	rand.Seed(time.Now().UnixNano())
	url := "https://www.lightnovelpub.com/novel/atticuss-odyssey-reincarnated-into-a-playground/chapters"
	domain := strings.SplitAfter(url,".com")[0]
	bookName := strings.Split(strings.SplitAfter(strings.SplitAfter(url,"novel/")[1],"/")[0],"/")[0]
	fmt.Println("filename: " + bookName)

	req, err := http.NewRequest("GET",url,nil)
	if err != nil {
		log.Panic(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://www.google.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
		return
	} 

	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Panic(err)
		return
	}

	chapterPages := []*html.Node{doc}
	allChapters := []Chapter{}

	callPages(doc, domain, &chapterPages, client)
	

	fmt.Println(chapterPages)
	//Gets every chapter page
	for _, page := range chapterPages {
		allChapters = append(allChapters, getChapters(page, domain)...)
	}

	printStuff(allChapters)
	
	callChapter(allChapters, client,bookName)
	
}

func callChapter(allChapters []Chapter, client *http.Client, bookName string) {
	fullHTML := ""

	for index := 0; index < len(allChapters); index++ {
	//for _, chapterPage := range allChapters {
		fmt.Println(allChapters[index].Link)
		pageReq, err := http.NewRequest("GET", allChapters[index].Link, nil)
		if err != nil {
			log.Panic(err)
			return
		}
		pageReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
		pageReq.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		pageReq.Header.Set("Accept-Language", "en-US,en;q=0.5")
		pageReq.Header.Set("Connection", "keep-alive")
		pageReq.Header.Set("Referer", "https://www.google.com")

		time.Sleep(time.Duration(rand.Intn(2000)+1000) * time.Millisecond)

		pageResp, err2 := client.Do(pageReq)
		if err2 != nil {
			//log.Panic(err2)
			//return
			
		}
		// pageResp.Body.Close()

		fmt.Println(pageResp.StatusCode)

		for pageResp.StatusCode != 200 || err2 != nil{
			pageReq2, err5 := http.NewRequest("GET", allChapters[index].Link, nil)
		if err5 != nil {
			log.Panic(err5)
			return
		}
		pageReq2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
		pageReq2.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		pageReq2.Header.Set("Accept-Language", "en-US,en;q=0.5")
		pageReq2.Header.Set("Connection", "keep-alive")
		pageReq2.Header.Set("Referer", "https://www.google.com")

		time.Sleep(time.Duration(rand.Intn(5000)+10000) * time.Millisecond)

		pageResp2, err4 := client.Do(pageReq2)
		if err4 != nil {
			//log.Panic(err4)
			//return
		}


		pageResp = pageResp2
		}
		

		tempDoc, err3 := html.Parse(pageResp.Body)
		if err3 != nil {
			log.Panic(err3)
			return
		}
		//fmt.Println(tempDoc)

		var chapContent *html.Node

		for node := range tempDoc.Descendants(){
			if node.Type != html.ElementNode || node.Data != "div" {
				continue
			}
			for _, ele := range node.Attr {
				if ele.Key != "id"  || ele.Val != "chapter-container"{
					continue
				}
				
				node.PrevSibling = nil
				node.NextSibling = nil
				fmt.Println(node)
				chapContent = node
				break
			} 
			
		}

		var buf bytes.Buffer
		if err := html.Render(&buf, chapContent); err != nil {
			log.Fatal(err)
		}

		partHTML := "<h2>" + allChapters[index].Name +"</h2>" + buf.String() 
		fullHTML = fullHTML + partHTML
	}


	finishedHTML :="<!DOCTYPE html>\n<html>\n<head><meta charset=\"UTF-8\"></head>\n<body>\n" + fullHTML + "\n</body>\n</html>"
	//fmt.Println(finishedHTML)
	os.WriteFile(bookName + ".html", []byte(finishedHTML), 0644)
}
























func callPages(doc *html.Node, domain string, chapterPages *[]*html.Node, client *http.Client) {
	//Gets us all the pages that have 100 chapters each
	for _, fullPage := range getTotalPages(doc, domain) {
		fmt.Println(fullPage)
		pageReq, err := http.NewRequest("GET", fullPage, nil)
		if err != nil {
			log.Panic(err)
			return
		}
		pageReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
		pageReq.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		pageReq.Header.Set("Accept-Language", "en-US,en;q=0.5")
		pageReq.Header.Set("Connection", "keep-alive")
		pageReq.Header.Set("Referer", "https://www.google.com")

		time.Sleep(time.Duration(rand.Intn(2000)+1000) * time.Millisecond)

		pageResp, err2 := client.Do(pageReq)
		if err2 != nil {
			log.Panic(err2)
			return
		}
		
		//defer pageReq.Body.Close()

		fmt.Println(pageResp.StatusCode)

		tempDoc, err3 := html.Parse(pageResp.Body)
		if err3 != nil {
			log.Panic(err3)
			return
		}
		*chapterPages = append(*chapterPages, tempDoc)

	}
}



func getTotalPages(doc *html.Node, domain string) []string {
	halfURL := ""
	var lastPageNumber int
	for node := range doc.Descendants() {
		if node.Type != html.ElementNode || node.Data != "ul" {
			continue
		}
		for _, ele := range node.Attr {
			if ele.Key != "class" {
				continue
			}
			if ele.Val == "pagination" {
				//inside of page navigation list
				for lists := range node.Descendants(){
					//fmt.Println(lists)
					if lists.Type != html.ElementNode || lists.Data != "li" {
						continue
					}
					//mt.Println(lists)
					if lists.NextSibling == nil {
						//fmt.Println(lists)
						lastPage := lists.FirstChild
						for _, attribute := range lastPage.Attr {
							if attribute.Key == "href" {
								split := strings.SplitAfter(attribute.Val,"=")
								halfURL = split[0]
								tempNumber, _ := strconv.Atoi(split[1])
								lastPageNumber = tempNumber
							}
							//fmt.Println(lastPageNumber)
						}
					}
				}
			}
		}
	}

	pageURLs := []string{}
	for index := 2; index <= lastPageNumber; index++ {
		pageURLs = append(pageURLs, domain + halfURL +strconv.Itoa(index))
	}

	return pageURLs
}


func getChapters(doc *html.Node, domain string) []Chapter{
	chapters := []Chapter{}

	for node := range doc.Descendants() {
		if node.Type != html.ElementNode || node.Data != "ul"{
			continue
		}
		for _, att := range node.Attr {
			if att.Key != "class" || att.Val != "chapter-list" {
				continue
			}
			//fmt.Println("Chapter List Node")
			//fmt.Println(node)
			for nodes := range node.Descendants() {
				if nodes.Type != html.ElementNode || nodes.Data != "a"{ 
					continue
				}
				tempChapter := Chapter{}
				for _, attribute := range nodes.Attr {
					if attribute.Key == "title" {
						//fmt.Println(attribute.Val)
						tempChapter.Name = attribute.Val
					
					}
					if attribute.Key == "href" {
						//fmt.Println(attribute.Val)
						tempChapter.Link = domain + attribute.Val
					}
				}
				//fmt.Println(tempChapter)
				chapters = append(chapters, tempChapter)
			}
		}
	}
	//fmt.Println(chapters)
	//fmt.Println(doc)
	return chapters
}


func printStuff(chaps []Chapter) {
	for _, chap := range chaps {
		fmt.Println(chap.Name)
	}
}