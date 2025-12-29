package capture

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

const tagStyle = "style"

type Inliner struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewInliner(baseURL string) (*Inliner, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	return &Inliner{
		baseURL:    parsed,
		httpClient: &http.Client{},
	}, nil
}

func (i *Inliner) InlineHTML(htmlContent []byte) ([]byte, error) {
	doc, err := html.Parse(bytes.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	i.processNode(doc)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return nil, fmt.Errorf("failed to render HTML: %w", err)
	}

	return buf.Bytes(), nil
}

func (i *Inliner) processNode(n *html.Node) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "link":
			i.processLink(n)
		case "script":
			i.processScript(n)
		case "img":
			i.processImage(n)
		case tagStyle:
			i.processStyleTag(n)
		}

		for j := range n.Attr {
			if n.Attr[j].Key == tagStyle {
				n.Attr[j].Val = i.inlineStyleURLs(n.Attr[j].Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		i.processNode(c)
	}
}

func (i *Inliner) processLink(n *html.Node) {
	var rel, href string
	var hrefIdx int

	for j, attr := range n.Attr {
		switch attr.Key {
		case "rel":
			rel = attr.Val
		case "href":
			href = attr.Val
			hrefIdx = j
		}
	}

	if rel == "stylesheet" && href != "" {
		content, contentType, err := i.fetchResource(href)
		if err != nil {
			log.Debug().Err(err).Str("href", href).Msg("Failed to fetch stylesheet")
			return
		}

		content = []byte(i.inlineStyleURLs(string(content)))

		n.Data = "style"
		n.Attr = nil

		n.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: string(content),
		})

		_ = contentType
		_ = hrefIdx
	}
}

func (i *Inliner) processScript(n *html.Node) {
	var src string
	var srcIdx int

	for j, attr := range n.Attr {
		if attr.Key == "src" {
			src = attr.Val
			srcIdx = j
			break
		}
	}

	if src == "" {
		return
	}

	content, _, err := i.fetchResource(src)
	if err != nil {
		log.Debug().Err(err).Str("src", src).Msg("Failed to fetch script")
		return
	}

	n.Attr = append(n.Attr[:srcIdx], n.Attr[srcIdx+1:]...)

	n.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: string(content),
	})
}

func (i *Inliner) processImage(n *html.Node) {
	for j, attr := range n.Attr {
		if attr.Key == "src" && !strings.HasPrefix(attr.Val, "data:") {
			dataURI, err := i.toDataURI(attr.Val)
			if err != nil {
				log.Debug().Err(err).Str("src", attr.Val).Msg("Failed to convert image to data URI")
				continue
			}
			n.Attr[j].Val = dataURI
		}
	}
}

func (i *Inliner) processStyleTag(n *html.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			c.Data = i.inlineStyleURLs(c.Data)
		}
	}
}

func (i *Inliner) inlineStyleURLs(css string) string {
	urlRegex := regexp.MustCompile(`url\(['"]?([^'")\s]+)['"]?\)`)

	return urlRegex.ReplaceAllStringFunc(css, func(match string) string {
		submatch := urlRegex.FindStringSubmatch(match)
		if len(submatch) < 2 {
			return match
		}

		resourceURL := submatch[1]
		if strings.HasPrefix(resourceURL, "data:") {
			return match
		}

		dataURI, err := i.toDataURI(resourceURL)
		if err != nil {
			return match
		}

		return fmt.Sprintf("url(%s)", dataURI)
	})
}

func (i *Inliner) toDataURI(resourceURL string) (string, error) {
	content, contentType, err := i.fetchResource(resourceURL)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	return fmt.Sprintf("data:%s;base64,%s", contentType, encoded), nil
}

func (i *Inliner) fetchResource(resourceURL string) (content []byte, contentType string, err error) {
	absURL, err := i.resolveURL(resourceURL)
	if err != nil {
		return nil, "", err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, absURL, http.NoBody)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := i.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch resource: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("resource returned status %d", resp.StatusCode)
	}

	content, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read resource: %w", err)
	}

	contentType = resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	if idx := strings.Index(contentType, ";"); idx != -1 {
		contentType = contentType[:idx]
	}

	return content, contentType, nil
}

func (i *Inliner) resolveURL(resourceURL string) (string, error) {
	parsed, err := url.Parse(resourceURL)
	if err != nil {
		return "", err
	}

	resolved := i.baseURL.ResolveReference(parsed)
	return resolved.String(), nil
}
