package crawler

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/ilam072/wbtech-l2/16/internal/storage"
	"github.com/ilam072/wbtech-l2/16/internal/utils"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
)

const prefix = "sites/"

type Crawler struct {
	url          string
	depth        int
	workersLimit int
	visited      map[string]bool
	mu           sync.Mutex
}

func New(url string, depth int, workersLimit int) *Crawler {
	return &Crawler{
		url:          url,
		depth:        depth,
		workersLimit: workersLimit,
		visited:      make(map[string]bool),
	}
}

func (c *Crawler) Run() {
	col := c.newCollector()

	col.OnRequest(func(r *colly.Request) {
		norm := utils.Normalize(r.URL)

		c.mu.Lock()
		if c.visited[norm] {
			c.mu.Unlock()
			r.Abort()
			return
		}
		c.visited[norm] = true
		c.mu.Unlock()
		fmt.Println("Обработка сайта:", r.URL.String())
	})

	col.OnError(func(r *colly.Response, err error) {
		log.Printf("Ошибка запроса %s: %v", r.Request.URL, err)
	})

	col.OnResponse(func(r *colly.Response) {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
		if err != nil {
			log.Printf("Ошибка парсинга HTML %s: %v", r.Request.URL, err)
			return
		}

		doc.Find("[href], [src]").Each(func(i int, s *goquery.Selection) {
			attr := "href"
			if _, exists := s.Attr("src"); exists {
				attr = "src"
			}

			val, _ := s.Attr(attr)
			if len(val) > 0 && val[0] == '/' {
				s.SetAttr(attr, val[1:])
			}
		})

		html, err := doc.Html()
		if err != nil {
			log.Printf("Ошибка генерации HTML %s: %v", r.Request.URL, err)
			return
		}

		localPath := storage.LocalPath(prefix, r.Request.URL)
		if err := storage.Save(localPath, []byte(html)); err != nil {
			log.Printf("Ошибка сохранения файла %s: %v", localPath, err)
		}
	})

	visitURL := func(u string) {
		if err := col.Visit(u); err != nil {
			if !strings.Contains(err.Error(), "already visited") {
				log.Printf("Ошибка посещения %s: %v", u, err)
			}
		}
	}

	col.OnHTML("a[href]", func(e *colly.HTMLElement) { visitURL(e.Request.AbsoluteURL(e.Attr("href"))) })
	col.OnHTML("img[src]", func(e *colly.HTMLElement) { visitURL(e.Request.AbsoluteURL(e.Attr("src"))) })
	col.OnHTML("link[href]", func(e *colly.HTMLElement) { visitURL(e.Request.AbsoluteURL(e.Attr("href"))) })
	col.OnHTML("script[src]", func(e *colly.HTMLElement) { visitURL(e.Request.AbsoluteURL(e.Attr("src"))) })

	visitURL(c.url)
	col.Wait()
}

func (c *Crawler) newCollector() *colly.Collector {
	u, err := url.Parse(c.url)
	if err != nil {
		log.Fatalf("Неверный URL: %v", err)
	}

	col := colly.NewCollector(
		colly.MaxDepth(c.depth),
		colly.AllowedDomains(u.Host),
		colly.Async(true),
	)
	col.SetRequestTimeout(15 * time.Second)

	if err = col.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: c.workersLimit,
	}); err != nil {
		log.Fatalf("Ошибка установки лимитов: %v", err)
	}

	return col
}
