package scrapper

import (
	"encoding/csv"
	"fmt"
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

var jobUrl string = "https://id.indeed.com/lihat-lowongan-kerja?jk="

func Scrap(key string) {
	var baseUrl string = "https://id.indeed.com/jobs?q=" + key + "&limit=50"
	c := make(chan []jobsRes)
	var jobs []jobsRes
	pages := getPages(baseUrl)

	for i := 0; i < pages; i++ {
		go getPage(i, baseUrl, c)
	}
	for i := 0; i < pages; i++ {
		jobsResult := <-c
		jobs = append(jobs, jobsResult...)
	}

	writeJobs(jobs) //make csv file
	fmt.Println(jobs)
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

func getPage(page int, baseUrl string, c chan<- []jobsRes) {
	var jobs []jobsRes
	d := make(chan jobsRes)
	pageUrl := baseUrl + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting ", pageUrl)
	res, err := http.Get(pageUrl)

	checkErr(err)
	checkStatus(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchedJob := doc.Find(".jobsearch-SerpJobCard")

	searchedJob.Each(func(i int, s *goquery.Selection) {
		go extractJobs(s, d)
	})

	for i := 0; i < searchedJob.Length(); i++ {
		job := <-d
		jobs = append(jobs, job)
	}
	c <- jobs
}

func extractJobs(s *goquery.Selection, d chan<- jobsRes) {
	id, _ := s.Attr("data-jk")
	id = CleanString(id)
	title := CleanString(s.Find(".title>a").Text())
	loc := CleanString(s.Find(".sjcl>.location").Text())
	d <- jobsRes{id: id, title: title, loc: loc}
}

//get the total number of pages
func getPages(baseUrl string) int {
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

//CleanString Clean the String from spaces around
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
