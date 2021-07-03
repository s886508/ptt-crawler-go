package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

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

	var s *storage.ESStorage = nil
	if len(*esServer) > 0 {
		s = &storage.ESStorage{}
		err := s.Init(*esServer)
		if err != nil {
			return
		}
	}

	articles := crawler.GetArticles(*startPage, *endPage, *board)
	for _, a := range articles {
		// output as files
		saveToFile(*outputDir, *board, a)

		// output to elasticsearch
		saveToElastic(s, *board, a)

		// output to stdout if no ouput option is specified
		printToScreen(*outputDir, *esServer, a)
	}
}

func saveToFile(dir string, board string, article *article.Article) {
	if len(dir) == 0 {
		return
	}

	filePath := fmt.Sprintf("%s/%s/%s.json", dir, board, article.Id)
	err := article.Save(filePath, false)
	if err != nil {
		log.Fatal(err)
	}
}

func saveToElastic(s *storage.ESStorage, board string, article *article.Article) {
	if s == nil {
		return
	}
	content, err := article.Dump()
	if err != nil {
		log.Fatal(err)
	}

	err = s.AddDocument(article.Id, content, board)
	if err != nil {
		log.Fatal(err)
	}
}

func printToScreen(dir string, elastic string, article *article.Article) {
	if len(dir) > 0 || len(elastic) > 0 {
		return
	}

	fmt.Printf("%s\n%v\n%s\n\n", strings.Repeat("%", 80), article, strings.Repeat("%", 80))
}
