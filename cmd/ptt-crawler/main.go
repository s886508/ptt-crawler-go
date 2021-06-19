package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/s886508/ptt-crawler-go/pkg/crawler"
)

func main() {
	board := flag.String("board", "", "board name to retrieve articles")
	startPage := flag.Int64("start", 1, "start page to retrieve articles")
	endPage := flag.Int64("end", 10, "end page to retrieve articles")
	outputDir := flag.String("directory", "", "directory to save retrieved articles")

	flag.Parse()

	if len(*board) == 0 {
		flag.PrintDefaults()
		return
	}

	if *startPage <= 0 || *endPage <= 0 || *endPage > *startPage {
		return
	}

	articles := crawler.GetArticles(*startPage, *endPage, *board)
	for _, a := range articles {
		if len(*outputDir) > 0 {
			filePath := fmt.Sprintf("%s/%s/%s.json", *outputDir, *board, a.Id)
			err := a.Save(filePath, false)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
