package petsmart

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

const ScrapeUrl string = "https://www.petsmart.com"
const HeroSelector string = ".campaign-hero-container"

func RetrievePromoCode() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	c, _ := chromedp.NewContext(ctx)
	defer func() {
		if err := chromedp.Cancel(c); err != nil {
			panic("chromedp could not be cancelled")
		}
	}()

	fmt.Println("Scraping: " + ScrapeUrl)

	var heroText string
	err := chromedp.Run(c,
		chromedp.Navigate(ScrapeUrl),
		chromedp.WaitEnabled(HeroSelector, chromedp.ByQuery),
		chromedp.Text(HeroSelector, &heroText, chromedp.ByQuery),
	)
	if err != nil {
		return "", err
	}

	promoCode, err := parseHero(heroText)
	if err != nil {
		return "", err
	}

	return promoCode, nil
}

func parseHero(text string) (string, error) {
	if strings.Contains(text, "SAVE20") {
		return "SAVE20", nil
	}

	return "", nil
}
