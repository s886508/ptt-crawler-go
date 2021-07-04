package storage

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type ESStorage struct {
	client *elasticsearch.Client
}

func (s *ESStorage) Init(url string) error {
	if s.client != nil {
		log.Println("elastic search client already initialized")
		return fmt.Errorf("client inited")
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			url,
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	s.client = client

	resp, err := s.client.Info()
	if err != nil {
		log.Println("fail to get elasticsearch info: ", err.Error())
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (s *ESStorage) AddDocument(docId string, doc string, board string) error {
	if s.client == nil {
		fmt.Println("elasticsearch client has not initialized")
		return fmt.Errorf("elasticsearch client not inited")
	}

	req := esapi.IndexRequest{
		Index:      fmt.Sprintf("ptt-board-%s", strings.ToLower(board)),
		DocumentID: docId,
		Body:       strings.NewReader(doc),
		Refresh:    "true",
	}

	resp, err := req.Do(context.Background(), s.client)
	if err != nil {
		log.Printf("index request fail, doc id: %s, err: %s\n", docId, err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		log.Printf("index document error, doc id: %s, err: %s ", docId, resp.Status())
		return fmt.Errorf("index document error")
	}

	log.Println("index document successfully, doc ID: ", docId)
	return nil
}
