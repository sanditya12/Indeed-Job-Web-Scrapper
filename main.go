package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type jobsRes struct {
	id    string
	title string
	loc   string
}

var baseUrl string = "https://id.indeed.com/jobs?q=golang&limit=50"
var jobUrl string = "https://id.indeed.com/lihat-lowongan-kerja?jk="

func main() {
	var jobs []jobsRes
	pages := getPages()

	for i := 0; i < pages; i++ {
		jobsResult := getPage(i)
		jobs = append(jobs, jobsResult...)
	}
	writeJobs(jobs)
}

func writeJobs(jobs []jobsRes) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	w := csv.NewWriter(file)
	defer w.Flush()
	header := []string{"id", "title", "loc"}
	wrErr := w.Write(header)
	checkErr(wrErr)

	for _, job := range jobs {
		content := []string{jobUrl + job.id, job.title, job.loc}
		wErr := w.Write(content)
		checkErr(wErr)
	}
}

func getPage(page int) []jobsRes {
	var jobs []jobsRes
	pageUrl := baseUrl + "&start=" + strconv.Itoa(page*50)
	res, err := http.Get(pageUrl)
	checkErr(err)
	checkStatus(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)

	doc.Find(".jobsearch-SerpJobCard").Each(func(i int, s *goquery.Selection) {
		job := extractJobs(s)
		jobs = append(jobs, job)
	})
	return jobs
}

func extractJobs(s *goquery.Selection) jobsRes {
	id, _ := s.Attr("data-jk")
	id = cleanString(id)
	title := cleanString(s.Find(".title>a").Text())
	loc := cleanString(s.Find(".sjcl>.location").Text())
	return jobsRes{id: id, title: title, loc: loc}
}

//get the total number of pages
func getPages() int {
	pages := 0
	res, err := http.Get(baseUrl)
	checkErr(err)
	checkStatus(res)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})
	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkStatus(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln(res)
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
