package article

import (
	"encoding/json"
	"io/ioutil"
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
