package crawler

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

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

func sleepMoment() {
	rand.Seed(time.Now().UnixNano())
	tfs := rand.Int63n(20)
	log.Printf("sleep for %d seconds to get next page\n", tfs)
	time.Sleep(time.Duration(tfs) * time.Second)
}

// GetArticles retrieves ptt articles from web.
func GetArticles(firstPage int64, lastPage int64, board string) []*article.Article {
	var articles []*article.Article
	for page := firstPage; page <= lastPage; page++ {
		log.Printf("retrieve articles from board: %s, page: %d", board, page)
		url := fmt.Sprintf("%s/bbs/%s/index%d.html", pttUrl, board, page) // Last page

		as, err := retrieveArticles(url)
		if err != nil {
			return nil
		}

		articles = append(articles, as...)

		if page != lastPage {
			sleepMoment()
		}
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
		if len(link) == 0 {
			continue
		}
		atlLink := fmt.Sprintf("%s%s", pttUrl, link)
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

	log.Println("retrieve single article from: ", url)
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
