package main

import "snippetbox.ergz.com/internal/models"

type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
}
