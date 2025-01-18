package articles_get

import (
	"net/http"

	"github.com/art-es/yet-another-service/internal/app/shared/dto"
)

type request struct {
	FromSlug *string
}

type response struct {
	Articles []article `json:"articles"`
	HasMore  bool      `json:"hasMore"`
}

type article struct {
	Slug    string  `json:"slug"`
	Title   string  `json:"title"`
	Content string  `json:"content"`
	Author  *author `json:"author,omitempty"`
}

type author struct {
	NickName    string `json:"nickName"`
	DisplayName string `json:"displayName"`
}

func parseRequest(in *http.Request) request {
	var out request

	query := in.URL.Query()
	if fromSlug := query.Get("fromSlug"); fromSlug != "" {
		out.FromSlug = &fromSlug
	}

	return out
}

func convertResponse(out *dto.GetArticlesOut) response {
	articles := make([]article, 0, len(out.Articles))
	for _, a := range out.Articles {
		articles = append(articles, convertArticle(a))
	}

	return response{
		Articles: articles,
		HasMore:  out.HasMore,
	}
}

func convertArticle(in dto.Article) article {
	return article{
		Slug:    in.Slug,
		Title:   in.Title,
		Content: in.Content,
		Author:  convertAuthor(in.Author),
	}
}

func convertAuthor(in *dto.ArticleAuthor) *author {
	if in == nil {
		return nil
	}

	return &author{
		NickName:    in.NickName,
		DisplayName: in.DisplayName,
	}
}
