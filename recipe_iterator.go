package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/antchfx/htmlquery"
)

const (
	recipeQuery = `//div[@itemtype="http://schema.org/Recipe"]`
)

type RecipeIterator struct {
	addr  string
	path  string
	query string

	currentPage int
	pageSize    int
	lastPage    bool
	recipes     []*Recipe
	recipe      *Recipe

	m   sync.Mutex
	err error
}

func NewRecipeIterator(addr, path, query string) *RecipeIterator {
	return &RecipeIterator{
		addr:     addr,
		path:     path,
		query:    query,
		pageSize: 6,
	}
}

func (ri *RecipeIterator) Next() bool {
	if ri.err != nil {
		return false
	}

	ri.m.Lock()
	defer ri.m.Unlock()
	if len(ri.recipes) > 0 {
		ri.recipes = ri.recipes[1:]
	}
	if len(ri.recipes) == 0 && !ri.lastPage {
		ri.err = ri.queryNextRecipes()
	}

	return ri.err == nil && len(ri.recipes) > 0
}

func (ri *RecipeIterator) Recipe() *Recipe {
	ri.m.Lock()
	defer ri.m.Unlock()

	if ri.err != nil || len(ri.recipes) == 0 {
		return nil
	}
	return ri.recipes[0]
}

func (ri *RecipeIterator) Err() error {
	return ri.err
}

func (ri *RecipeIterator) queryNextRecipes() error {
	req, err := http.NewRequest("GET", ri.addr+ri.path, nil)
	if err != nil {
		return err
	}
	req.URL.RawQuery = fmt.Sprintf("q=%s&page=%d&pageSize=%d", ri.query, ri.currentPage, ri.pageSize)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return err
	}

	var recipes []*Recipe
	for _, n := range htmlquery.Find(doc, recipeQuery) {
		recipe, err := ParseRecipe(ri.addr, n)
		if err != nil {
			return err
		}

		recipes = append(recipes, recipe)
	}

	ri.recipes = append(ri.recipes, recipes...)
	ri.lastPage = (len(recipes) != ri.pageSize)
	ri.currentPage++

	return nil
}
