package rakuten

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

const ScrapeUrl string = "https://www.rakuten.com/shop/petsmart"
const CashBackCssSelector string = "span[data-testid=\"online-cash-back\"]"

func RetrieveCashback() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	c, _ := chromedp.NewContext(ctx)
	defer func() {
		if err := chromedp.Cancel(c); err != nil {
			panic("chromedp could not be cancelled")
		}
	}()

	fmt.Println("Scraping: " + ScrapeUrl)

	var cashbackText string
	err := chromedp.Run(c,
		chromedp.Navigate(ScrapeUrl),
		chromedp.WaitEnabled(CashBackCssSelector, chromedp.ByQuery),
		chromedp.Text(CashBackCssSelector, &cashbackText, chromedp.ByQuery),
	)
	if err != nil {
		return 0, err
	}

	cashback, err := parseCashback(cashbackText)
	if err != nil {
		return 0, err
	}

	return cashback, nil
}

func parseCashback(text string) (int, error) {
	text = strings.ReplaceAll(text, "% Cash Back", "")
	cashback, err := strconv.Atoi(text)
	if err != nil {
		return 0, err
	}

	return cashback, nil
}
