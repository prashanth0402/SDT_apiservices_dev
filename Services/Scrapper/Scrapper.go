// Package scrapper is a reusable, site-agnostic HTML scraper built on colly.
// It has no knowledge of any particular target site or data shape (job
// boards, listings, etc.) - callers describe what to extract via Request and
// FieldSelector - so the package can be imported unmodified across multiple
// repos/projects (e.g. job scraping in one project, other data scraping in
// another).
package scrapper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// FieldSelector describes how to extract one field from a scraped element.
// Attr empty means take the element's trimmed text content instead of an
// attribute value.
type FieldSelector struct {
	Selector string
	Attr     string
}

// Config controls crawler-wide behavior and is shared by every request made
// through a Scraper instance.
type Config struct {
	AllowedDomains []string
	UserAgent      string
	Parallelism    int
	RequestDelay   time.Duration
	RandomDelay    time.Duration
	Timeout        time.Duration
	MaxDepth       int
}

// Request describes a single scrape job: a page (or first page of a
// paginated listing) plus how to extract repeating items and, optionally,
// how to follow pagination links.
type Request struct {
	URL          string
	ItemSelector string
	Fields       map[string]FieldSelector
	NextSelector string
	NextAttr     string
	MaxPages     int
}

// Item is one scraped record: field name -> extracted value.
type Item map[string]string

// Scraper is a reusable, site-agnostic HTML scraper. Callers configure it
// once via New and reuse it across many Scrape/ScrapeOne calls.
type Scraper struct {
	cfg Config
}

// New creates a Scraper with the given crawler configuration.
func New(cfg Config) *Scraper {
	if cfg.UserAgent == "" {
		cfg.UserAgent = "Mozilla/5.0 (compatible; SDTScrapperBot/1.0)"
	}
	if cfg.Parallelism <= 0 {
		cfg.Parallelism = 1
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	return &Scraper{cfg: cfg}
}

func (s *Scraper) newCollector() (*colly.Collector, error) {
	opts := []colly.CollectorOption{colly.UserAgent(s.cfg.UserAgent)}
	if len(s.cfg.AllowedDomains) > 0 {
		opts = append(opts, colly.AllowedDomains(s.cfg.AllowedDomains...))
	}
	if s.cfg.MaxDepth > 0 {
		opts = append(opts, colly.MaxDepth(s.cfg.MaxDepth))
	}

	c := colly.NewCollector(opts...)
	c.SetRequestTimeout(s.cfg.Timeout)

	if err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: s.cfg.Parallelism,
		Delay:       s.cfg.RequestDelay,
		RandomDelay: s.cfg.RandomDelay,
	}); err != nil {
		return nil, fmt.Errorf("scrapper: configuring rate limit: %w", err)
	}
	return c, nil
}

// Scrape fetches req.URL (and, if req.NextSelector is set, up to
// req.MaxPages following pages) and extracts one Item per element matched by
// req.ItemSelector, using req.Fields to pull data out of each matched
// element.
func (s *Scraper) Scrape(ctx context.Context, req Request) ([]Item, error) {
	if req.URL == "" {
		return nil, fmt.Errorf("scrapper: request URL is required")
	}
	if req.ItemSelector == "" {
		return nil, fmt.Errorf("scrapper: item selector is required")
	}
	if req.MaxPages <= 0 {
		req.MaxPages = 1
	}

	c, err := s.newCollector()
	if err != nil {
		return nil, err
	}

	var items []Item
	var scrapeErr error
	nextURL := ""

	c.OnHTML(req.ItemSelector, func(e *colly.HTMLElement) {
		item := make(Item, len(req.Fields))
		for name, fs := range req.Fields {
			item[name] = extractField(e, fs)
		}
		items = append(items, item)
	})

	if req.NextSelector != "" {
		attr := req.NextAttr
		if attr == "" {
			attr = "href"
		}
		c.OnHTML(req.NextSelector, func(e *colly.HTMLElement) {
			if nextURL == "" {
				nextURL = e.Request.AbsoluteURL(e.Attr(attr))
			}
		})
	}

	c.OnError(func(r *colly.Response, err error) {
		scrapeErr = fmt.Errorf("scrapper: request to %s failed: %w", r.Request.URL, err)
	})

	url := req.URL
	for page := 0; page < req.MaxPages; page++ {
		select {
		case <-ctx.Done():
			return items, ctx.Err()
		default:
		}

		nextURL = ""
		if err := c.Visit(url); err != nil {
			return items, fmt.Errorf("scrapper: visiting %s: %w", url, err)
		}
		c.Wait()

		if scrapeErr != nil {
			return items, scrapeErr
		}
		if nextURL == "" {
			break
		}
		url = nextURL
	}

	return items, nil
}

// ScrapeOne fetches a single page and extracts one Item directly from it
// (no repeating item selector) - useful for detail pages, e.g. a single job
// posting.
func (s *Scraper) ScrapeOne(ctx context.Context, url string, fields map[string]FieldSelector) (Item, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	c, err := s.newCollector()
	if err != nil {
		return nil, err
	}

	item := make(Item, len(fields))
	var scrapeErr error

	c.OnHTML("html", func(e *colly.HTMLElement) {
		for name, fs := range fields {
			item[name] = extractField(e, fs)
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		scrapeErr = fmt.Errorf("scrapper: request to %s failed: %w", r.Request.URL, err)
	})

	if err := c.Visit(url); err != nil {
		return nil, fmt.Errorf("scrapper: visiting %s: %w", url, err)
	}
	c.Wait()
	if scrapeErr != nil {
		return nil, scrapeErr
	}
	return item, nil
}

func extractField(e *colly.HTMLElement, fs FieldSelector) string {
	sel := e.DOM
	if fs.Selector != "" {
		sel = e.DOM.Find(fs.Selector)
	}
	if fs.Attr != "" {
		v, _ := sel.Attr(fs.Attr)
		return strings.TrimSpace(v)
	}
	return strings.TrimSpace(sel.Text())
}
