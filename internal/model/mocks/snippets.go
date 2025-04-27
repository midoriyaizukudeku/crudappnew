package mocks

import (
	"time"

	"dream.website/internal/model"
)

var mockSnippet = &model.Snippet{
	Id:      1,
	Title:   "Last Time",
	Content: "Of A Begiener",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (*model.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, model.ErrnoRecord
	}
}

func (m *SnippetModel) Latest() ([]*model.Snippet, error) {
	return []*model.Snippet{mockSnippet}, nil
}
