package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/s886508/ptt-crawler-go/pkg/article"
	"github.com/s886508/ptt-crawler-go/pkg/crawler"
	"github.com/s886508/ptt-crawler-go/pkg/storage"
)

func main() {
	board := flag.String("board", "", "board name to retrieve articles")
	startPage := flag.Int64("start", 1, "start page to retrieve articles")
	endPage := flag.Int64("end", 10, "end page to retrieve articles")
	outputDir := flag.String("directory", "", "[optional] save retrieved articles to directory, "+
		"print to stdout if none of output option is specified")
	esServer := flag.String("elastic", "", "[optional] save retrieved articles to elasticsearch server, "+
		"print to stdout if none of output option is specified")

	flag.Parse()

	if len(*board) == 0 {
		flag.PrintDefaults()
		return
	}

	if *startPage <= 0 || *endPage <= 0 || *endPage < *startPage {
		log.Fatal("invalid page numbers")
	}

	articles := crawler.GetArticles(*startPage, *endPage, *board)

	// output as files
	if len(*outputDir) > 0 {
		saveToFile(*outputDir, *board, articles)
	}

	if len(*esServer) > 0 {
		// output to elasticsearch
		saveToElastic(*esServer, *board, articles)
	}

	// output to stdout if no ouput option is specified
	if len(*outputDir) == 0 && len(*esServer) == 0 {
		printToScreen(articles)
	}
}

func saveToFile(dir string, board string, articles []*article.Article) {
	if len(dir) == 0 {
		return
	}

	for _, article := range articles {
		filePath := fmt.Sprintf("%s/%s/%s.json", dir, board, article.Id)
		err := article.Save(filePath, false)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func saveToElastic(eshost string, board string, articles []*article.Article) {
	s := &storage.ESStorage{}
	err := s.Init(eshost)
	if err != nil {
		log.Fatal(err)
	}

	for _, article := range articles {
		content, err := article.Dump()
		if err != nil {
			continue
		}

		err = s.AddDocument(article.Id, content, board)
		if err != nil {
			continue
		}
	}
}

func printToScreen(articles []*article.Article) {
	b, err := json.MarshalIndent(articles, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
