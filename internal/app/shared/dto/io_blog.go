package dto

type GetArticlesIn struct {
	FromSlug *string
}

type GetArticlesOut struct {
	Articles []Article
	HasMore  bool
}
