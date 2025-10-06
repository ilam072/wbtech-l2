package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const prefix = "sites/"

func main() {
	link := `https://wooordhunt.ru/word/word`
	u, err := url.Parse(link)
	if err != nil {
		panic(err)
	}

	var depth int
	fmt.Scan(&depth)
	col := colly.NewCollector(colly.MaxDepth(depth), colly.AllowedDomains(u.Host))

	col.OnRequest(func(r *colly.Request) {
		fmt.Printf("Обработка сайта: %v\n", r.URL.String())
	})

	col.OnResponse(func(r *colly.Response) {
		localPath := path(prefix, r.Request.URL)

		contentType := r.Headers.Get("Content-Type")

		// Если это HTML
		if strings.HasPrefix(contentType, "text/html") {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(r.Body)))
			if err != nil {
				panic(err)
			}

			// Переписываем ссылки
			doc.Find("a[href], link[href], script[src], img[src]").Each(func(i int, s *goquery.Selection) {
				if val, ok := s.Attr("href"); ok {
					s.SetAttr("href", localLink(r.Request.URL, val))
				}
				if val, ok := s.Attr("src"); ok {
					s.SetAttr("src", localLink(r.Request.URL, val))
				}
			})

			// Сохраняем изменённый HTML
			html, err := doc.Html()
			if err != nil {
				panic(err)
			}
			if err := os.WriteFile(localPath, []byte(html), 0644); err != nil {
				panic(err)
			}
		} else {
			// Остальные файлы (css, js, картинки)
			if err := r.Save(localPath); err != nil {
				panic(err)
			}
		}
	})

	col.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})
	col.OnHTML("img[src]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Request.AbsoluteURL(e.Attr("src")))
	})
	col.OnHTML("link[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})
	col.OnHTML("script[src]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Request.AbsoluteURL(e.Attr("src")))
	})

	// Старт
	err = col.Visit(link)
	if err != nil {
		log.Fatal("Ошибка загрузки:", err)
	}

	col.Wait()
}

func localLink(u *url.URL, link string) string {
	abs := u.ResolveReference(&url.URL{Path: link})
	return "/" + filepath.Join(abs.Host, abs.Path)
}

func path(baseDir string, u *url.URL) string {
	// Базовый путь: sites/hostname/...
	p := filepath.Join(baseDir, u.Host, u.Path)

	// Если путь пустой или заканчивается на "/", то добавляем index.html
	if u.Path == "" || u.Path[len(u.Path)-1] == '/' {
		p = filepath.Join(baseDir, u.Host, u.Path, "index.html")
	} else if filepath.Ext(u.Path) == "" {
		// Если путь без расширения (например, /catalogue/category)
		p = filepath.Join(baseDir, u.Host, u.Path, "index.html")
	}

	// Создаём папку под файл
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	return p
}
