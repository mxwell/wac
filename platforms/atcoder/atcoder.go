package atcoder

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mxwell/wac/model"
)

func FetchContest(url string) (*model.Contest, error) {
	url, err := trimUrl(url)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocument(url + "/assignments")
	if err != nil {
		return nil, err
	}
	titleElement := doc.Find("span.contest-name")
	if titleElement.Length() != 1 {
		return nil, fmt.Errorf("unable to detect contest name")
	}
	title := titleElement.Text()
	tasks := make(map[string]model.Task)
	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		column := s.Find("td.center")
		if column.Length() == 0 {
			fmt.Printf("unable to parse task %d\n", i)
			return
		}
		tokenElement := column.Find("a.linkwrapper")
		if tokenElement.Length() != 1 {
			fmt.Printf("unable to find token of task %d\n", i)
			return
		}
		token := tokenElement.Text()
		href, ok := tokenElement.Attr("href")
		if !ok {
			fmt.Printf("unable to extract link to task %d\n", i)
			return
		}
		column = column.Next()
		if column.Length() == 0 {
			fmt.Printf("unable to parse task %d\n", i)
			return
		}
		nameElement := column.Find("a.linkwrapper")
		if nameElement.Length() != 1 {
			fmt.Printf("unable to find name of task %d\n", i)
			return
		}
		name := nameElement.Text()
		tasks[token] = model.Task{url + href, name, token}
	})
	return &model.Contest{url, title, tasks}, nil
}

const SchemeHttp = "http://"
const SchemeHttps = "https://"

func trimUrl(url string) (string, error) {
	var offset int
	if strings.HasPrefix(url, SchemeHttp) {
		offset = len(SchemeHttp)
	} else if strings.HasPrefix(url, SchemeHttps) {
		offset = len(SchemeHttps)
	} else {
		return "", fmt.Errorf("no valid scheme is detected in url - %s", url)
	}
	slash := strings.Index(url[offset:], "/")
	if slash >= 0 {
		url = url[:offset+slash]
	}
	if !strings.HasSuffix(url[offset:], ".contest.atcoder.jp") {
		return "", fmt.Errorf("bad contest URL")
	}
	return url, nil
}
