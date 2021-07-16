package article

const (
	CommentTypeLike = iota
	CommentTypeBoo
	CommentTypeNeutral
)

type Comment struct {
	Author  string `json:"author"`
	Content string `json:"content"`
	Date    string `json:"date"`
	SrcIp   string `json:"srcip"`
	Type    int    `json:"type"`
}
