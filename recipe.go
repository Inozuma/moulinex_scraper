package main

import (
	"net/http"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const (
	imageQuery = "//img"
	linkQuery  = "//a"

	infoQuery         = "//div[@class='recipe-info']"
	ingredientsQuery  = "//div[@class='recipe-ingredients']"
	instructionsQuery = "//div[@class='recipe-instructions']"
)

type Recipe struct {
	Title string
	Image string
	Link  string

	Category        string
	PreparationTime string
	CookTime        string

	Ingredients  []string
	Instructions []string
}

func ParseRecipe(addr string, node *html.Node) (*Recipe, error) {
	recipe := &Recipe{}

	recipe.Title = htmlquery.SelectAttr(htmlquery.FindOne(node, imageQuery), "title")
	recipe.Image = htmlquery.SelectAttr(htmlquery.FindOne(node, imageQuery), "data-src")
	recipe.Link = htmlquery.SelectAttr(htmlquery.FindOne(node, linkQuery), "href")

	if recipe.Link != "" {
		resp, err := http.Get(addr + recipe.Link)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		doc, err := htmlquery.Parse(resp.Body)
		if err != nil {
			return nil, err
		}

		infoNode := htmlquery.FindOne(doc, infoQuery)
		recipe.Category = htmlquery.InnerText(htmlquery.FindOne(infoNode, "//p[@itemprop='recipeCategory']"))
		recipe.PreparationTime = htmlquery.SelectAttr(htmlquery.FindOne(infoNode, "//time[@itemprop='prepTime']"), "datetime")
		recipe.CookTime = htmlquery.SelectAttr(htmlquery.FindOne(infoNode, "//time[@itemprop='cookTime']"), "datetime")

		ingredientsNode := htmlquery.FindOne(doc, ingredientsQuery)
		for _, n := range htmlquery.Find(ingredientsNode, "//li[@itemprop='ingredients']") {
			recipe.Ingredients = append(recipe.Ingredients, htmlquery.InnerText(n))
		}

		instructionsNode := htmlquery.FindOne(doc, instructionsQuery)
		for _, n := range htmlquery.Find(instructionsNode, "//p") {
			recipe.Instructions = append(recipe.Instructions, htmlquery.InnerText(n))
		}
	}

	return recipe, nil
}
