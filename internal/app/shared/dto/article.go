package dto

type Article struct {
	Slug     string
	Title    string
	Content  string
	AuthorID int64

	Author *ArticleAuthor
}

type ArticleAuthor struct {
	DisplayName string
	NickName    string
}
