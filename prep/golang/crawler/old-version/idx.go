package idx

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

var ErrMismatchedHostname = fmt.Errorf("mismatched host name")

type Writer interface {
	Write(string, []byte) error
}

type Index struct {
	Writer
	URL  *url.URL
	Seen map[string]struct{}
}

func New(w Writer) *Index {
	return &Index{Writer: w, Seen: map[string]struct{}{}}
}

func (i *Index) EnsureURL(u *url.URL) {
	fmt.Println("Ensuring " + u.String())
	if i.URL == nil {
		i.URL = u
	}
}

func (i *Index) MarkSeen(str string) {
	fmt.Println("Marking: " + str)
	i.Seen[str] = struct{}{}
}

func (i *Index) Index(u *url.URL) error {
	i.EnsureURL(u)
	str := u.String()
	if _, seen := i.Seen[str]; seen {
		return nil
	}
	i.MarkSeen(str)
	resp, err := i.fetch(u)
	if err != nil {
		if err == ErrMismatchedHostname {
			return nil
		}
		return err
	}
	fart := bufio.NewReader(resp.Body)
	raw, err := ioutil.ReadAll(fart)
	n, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}
	if err := i.Writer.Write(u.Path, raw); err != nil {
		return err
	}
	i.Crawl(n)
	return nil
}

func (i *Index) Write(u string, data []byte) error {
	parsed, err := url.Parse(u)
	if err != nil {
		return err
	}
	return i.Writer.Write(parsed.Path, data)
}

func (i *Index) fetch(u *url.URL) (*http.Response, error) {
	fail := func(err error) (*http.Response, error) { return nil, err }

	if i.URL == nil {
		return fail(
			fmt.Errorf("unexpected error: index's URL should not be blank %#v", i),
		)
	}
	if u.Host == "" {
		u.Host = i.URL.Host
	}
	if !sameHostname(u, i.URL) {
		return fail(ErrMismatchedHostname)
	}
	fmt.Println("Fetching: " + u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		return fail(err)
	}
	if resp.StatusCode > 299 {
		return fail(fmt.Errorf("bad status for %+v: %d", u, resp.StatusCode))
	}
	return resp, nil
}

func sameHostname(url1, url2 *url.URL) bool {
	// fmt.Printf("Comparing %+v with %+v\n", url1, url2)
	return url1.Hostname() == url2.Hostname()
}

func (i *Index) Crawl(n *html.Node) error {
	if u := href(n); u != "" {
		fmt.Printf("Crawling: %#v \n", n)
		parsed, err := url.Parse(u)
		if err != nil {
			return err
		}
		if err = i.Index(parsed); err != nil {
			return err
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if err := i.Crawl(c); err != nil {
			return err
		}
	}
	return nil
}

func href(n *html.Node) string {
	fmt.Printf("Examining node for href: %#v\n\tType: %v\tData: %s\n", *n, n.Type, n.Data)
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				return attr.Val
			}
		}
	}
	return ""
}
