package article

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Article struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Board   string `json:"board"`
	Content string `json:"content"`
	Date    string `json:"date"`
	SrcIp   string `json:"srcip"`
}

func (a *Article) Save(filePath string, format bool) error {
	err := os.MkdirAll(path.Dir(filePath), 0755)
	if err != nil {
		return err
	}
	if format {
		b, err := json.Marshal(a)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filePath, b, 0644)
		if err != nil {
			return err
		}
	} else {
		b, err := json.MarshalIndent(a, "", " ")
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filePath, b, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Article) Dump() (string, error) {
	b, err := json.MarshalIndent(a, "", " ")
	if err != nil {
		log.Println("fail to dump article to json: ", err.Error())
		return "", err
	}

	return string(b), nil
}
