package main

import (
	"github.com/s886508/ptt-crawler-go/pkg/crawler"
)

func main() {
	_ := crawler.GetArticles(0, 0, "car")
}
