package eplister

import (
	"log"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/meros/go-svtdownloader/libs/gqcommon"
)

const baseURLString = "https://www.svtplay.se"

// Episodes contain info on one specific episode available on svtplay.se
type Episode struct {
	Series  string
	Season  string
	Episode string
	Url     url.URL
}

// Get retrieves and returns the list of episodes for a specific series on svtplay.se
func Get(series string) (episodes []Episode, err error) {
	baseURL, err := url.Parse(baseURLString)
	if err != nil {
		log.Println("Could not parse base URL", err)
		return
	}

	url, err := url.Parse(series)

	// Request the HTML page.
	seriesURL := baseURL.ResolveReference(url)
	doc, err := gqcommon.GetDoc(*seriesURL)
	if err != nil {
		log.Println("Could not fetch series url", err)
		return
	}

	seriesTitle := doc.Find("[data-rt='title-page-title']").Text()

	// Find the review items
	var seasonHrefs []string
	for _, node := range doc.Find(".play_accordion__section-title").Nodes {
		// For each item found, get the band and title
		href, hasHref := goquery.NewDocumentFromNode(node).Attr("href")
		if !hasHref {
			log.Println("Could not find expected href", err)
			continue
		}

		seasonHrefs = append(seasonHrefs, href)
	}

	for _, seasonHref := range seasonHrefs {
		url, err = url.Parse(seasonHref)
		if err != nil {
			log.Println("Could not parse seasonHref", err)
			continue
		}

		seasonURL := seriesURL.ResolveReference(url)

		seasonDoc, err := gqcommon.GetDoc(*seasonURL)
		if err != nil {
			log.Println("Could not fetch season URL", err)
			continue
		}

		seasonNodes := seasonDoc.
			Find(".play_accordion__section-title[href='" + seasonHref + "']").
			Nodes

		if len(seasonNodes) != 1 {
			log.Println("Expected one node for season")
			continue
		}

		seasonNodeDoc := goquery.
			NewDocumentFromNode(seasonNodes[0]).
			Parent().
			Parent()

		seasonTitle := seasonNodeDoc.Find(".play_accordion__section-title-inner").Text()
		titleNodes := seasonNodeDoc.Find(".play_related-list__item").Nodes

		for _, titleNode := range titleNodes {
			titleDocument := goquery.NewDocumentFromNode(titleNode)
			href, hasHref := titleDocument.Find(".play_related-item__link").Attr("href")
			if !hasHref {
				continue
			}

			episodeTitle := titleDocument.Find(".play_related-item__title").Text()

			url, err := url.Parse(href)
			if err != nil {
				log.Println("Failed to parse titleHref URL", err)
				continue
			}

			episodeURL := baseURL.ResolveReference(url)
			episodes = append(episodes, Episode{
				seriesTitle,
				seasonTitle,
				episodeTitle,
				*episodeURL})
		}
	}

	return
}
