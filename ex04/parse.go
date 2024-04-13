package MyHtmlParser

import (
	"io"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	nodes := getAnchorNodes(doc)
	var links []Link

	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	return links, nil

}

func getAnchorNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, getAnchorNodes(c)...)
	}
	return ret
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, a := range n.Attr {
		if a.Key == "href" {
			ret.Href = a.Val
			break
		}
	}
	ret.Text = getTextFromAnchor(n)
	return ret
}

func getTextFromAnchor(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}

	var ret string

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += getTextFromAnchor(c) + " "
	}

	return ret
}
