package crawler

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/s886508/ptt-crawler-go/pkg/article"
)

var pttUrl string = "https://www.ptt.cc"

const (
	MetaIndexAuthor = iota
	MetaIndexTitle
	MetaIndexData
	MetaIndexCount
)

// GetArticles retrieves ptt articles from web.
func GetArticles(firstPage int32, lastPage int32, board string) []*article.Article {
	url := fmt.Sprintf("%s/bbs/%s/index.html", pttUrl, board) // Last page

	articles, err := retrieveArticles(url)
	if err != nil {
		return nil
	}
	return articles
}

func httpReq(url string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil
	}
	req.AddCookie(&http.Cookie{Name: "over18", Value: "1"})

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	if resp.StatusCode != 200 {
		return nil
	}
	return resp
}

func retrieveArticles(url string) ([]*article.Article, error) {
	if len(url) == 0 {
		log.Println("retreive articels fail, empty url")
		return nil, fmt.Errorf("empty url")
	}
	resp := httpReq(url)
	if resp == nil {
		log.Println("retrieve articles fail, empty response")
		return nil, fmt.Errorf("empty response")
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var articles []*article.Article
	sels := doc.Find("div[class=r-ent]")
	for i := range sels.Nodes {
		s := sels.Eq(i)
		linkS := s.Find("div[class=title]>a[href]")
		link, _ := linkS.Attr("href")
		atlLink := fmt.Sprintf("%s/%s", pttUrl, link)
		if atl := retrieveSingleArticle(atlLink); atl != nil {
			if len(atl.Author) == 0 {
				atl.Author = s.Find("div[class=meta]>div[class=author]").Text()
			}
			if len(atl.Title) == 0 {
				atl.Title = linkS.Text()
			}
			articles = append(articles, atl)
		}
	}

	return articles, nil
}

func retrieveSingleArticle(url string) *article.Article {
	if len(url) == 0 {
		log.Println("retrieve single article fail, empty response")
		return nil
	}

	log.Println("retrieve single article from: %s", url)
	resp := httpReq(url)
	if resp == nil {
		log.Println("retrieve single article fail, empty response")
		return nil
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	tokens := strings.Split(url, "/")

	article := &article.Article{
		Id:    strings.TrimSuffix(tokens[len(tokens)-1], ".html"),
		Board: tokens[len(tokens)-2],
	}

	mainSel := doc.Find("div#main-content")
	mainSel.Find("div[class=article-metaline-right]").Remove()

	// process metadata
	metaData := mainSel.Find("div[class=article-metaline]")
	for i := 0; i < MetaIndexCount; i++ {
		m := metaData.Eq(i)
		if m == nil {
			continue
		}
		v := m.Find("span[class=article-meta-value]")
		if v == nil {
			continue
		}
		switch i {
		case MetaIndexAuthor:
			article.Author = v.Text()
		case MetaIndexTitle:
			article.Title = v.Text()
		case MetaIndexData:
			article.Date = v.Text()
		}
	}
	metaData.Remove()

	// process ip address
	f2 := mainSel.Find("span[class=f2]")
	for i := range f2.Nodes {
		s := f2.Eq(i)
		text := s.Text()
		if strings.Contains(text, "※ 發信站:") {
			re := regexp.MustCompile(`[0-9]*\.[0-9]*\.[0-9]*\.[0-9]* \(.*\)`)
			article.SrcIp = re.FindString(text)
		}
	}
	f2.Remove()

	// process push messages
	pushes := mainSel.Find("div[class=push]")
	//	for i := range pushes.Nodes {
	//	}
	pushes.Remove()

	// process content
	article.Content = mainSel.Text()

	return article
}