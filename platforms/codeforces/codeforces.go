package codeforces

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mxwell/wac/model"
)

type Codeforces struct {
}

func InitCodeforces() model.Platform {
	return Codeforces{}
}

func (a Codeforces) ValidUrl(url string) bool {
	_, err := trimUrl(url)
	return err == nil
}

const CodeforcesHost = "http://codeforces.com"

func (a Codeforces) GetContest(url string, rootDirName string) (*model.Contest, error) {
	url, err := trimUrl(url)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}
	titleElement := doc.Find("#sidebar table.rtable th a")
	if titleElement.Length() != 1 {
		return nil, fmt.Errorf("unable to detect contest name")
	}
	title := titleElement.Text()
	tasks := make(map[string]model.Task)
	doc.Find("table.problems tr").Each(func(i int, s *goquery.Selection) {
		column := s.Find("td")
		// if no <td>, then it's a header, so ignore the row
		if column.Length() < 1 {
			return
		}
		column = column.First()
		tokenElement := column.Find("a")
		if tokenElement.Length() != 1 {
			log.Printf("unable to find token of task %d\n", i)
			return
		}
		token := tokenElement.Text()
		token = strings.TrimSpace(token)
		token = strings.ToLower(token)
		href, ok := tokenElement.Attr("href")
		if !ok {
			log.Printf("unable to extract link to task %d\n", i)
			return
		}
		column = column.Next()
		if column.Length() == 0 {
			log.Printf("unable to find column with name of task %d\n", i)
			return
		}
		nameElement := column.Find("a")
		if nameElement.Length() != 1 {
			log.Printf("unable to find name of task %d\n", i)
			return
		}
		name := strings.TrimSpace(nameElement.Text())
		tasks[token] = model.Task{CodeforcesHost + href, name, token, make([]string, 0)}
	})
	return &model.Contest{url, title, tasks, rootDirName}, nil
}

func (a Codeforces) GetTests(task *model.Task) ([]model.Test, error) {
	doc, err := goquery.NewDocument(task.Link)
	if err != nil {
		return nil, err
	}
	sampleTestsElement := doc.Find("div.sample-tests div.sample-test")
	if sampleTestsElement.Length() != 1 {
		return nil, fmt.Errorf("element with tests is not found")
	}

	inputs := make(map[int]string)
	outputs := make(map[int]string)
	var idOrder []int

	sampleTestsElement.Find("div.input").Each(func(i int, s *goquery.Selection) {
		id := i + 1
		element := s.Find("pre")
		if element.Length() != 1 {
			log.Printf("WARN input is not found in sample %d, ignoring the sample...\n", id)
			return
		}
		if html, err := element.Html(); err == nil {
			inputs[id] = processHtml(html)
		} else {
			log.Printf("WARN unable to retrieve input entry in sample %d, ignoring the sample...\n", id)
			return
		}
		idOrder = append(idOrder, id)
	})

	sampleTestsElement.Find("div.output").Each(func(i int, s *goquery.Selection) {
		id := i + 1
		element := s.Find("pre")
		if element.Length() != 1 {
			log.Printf("WARN output is not found in sample %d, ignoring the sample...\n", id)
			return
		}
		if html, err := element.Html(); err == nil {
			outputs[id] = processHtml(html)
		} else {
			log.Printf("WARN unable to retrieve output entry in sample %d, ignoring the sample...\n", id)
			return
		}
	})

	var result []model.Test

	for _, id := range idOrder {
		input := inputs[id]
		output, ok := outputs[id]
		if !ok {
			log.Printf("WARN no output for sample %d\n", id)
			continue
		}
		token := fmt.Sprintf("sample%d", id)
		result = append(result, model.Test{token, input, output})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no valid sample tests were found")
	}

	return result, nil
}

var BrReplacer = strings.NewReplacer("<br/>", "\n")

func processHtml(html string) string {
	return BrReplacer.Replace(html)
}

const SchemeHttp = "http://"
const PathContest = "/contest/"

func trimUrl(url string) (string, error) {
	var offset int
	if strings.HasPrefix(url, SchemeHttp) {
		offset = len(SchemeHttp)
	} else {
		return "", fmt.Errorf("no valid scheme is detected in url - %s", url)
	}
	slash := strings.Index(url[offset:], "/")
	if slash < 0 {
		return "", fmt.Errorf("url must have path - %s", url)
	}
	slash = offset + slash
	if !strings.Contains(url[offset:slash], "codeforces.com") {
		return "", fmt.Errorf("host must be codeforces.com")
	}
	contest := strings.Index(url[slash:], PathContest)
	if contest != 0 {
		return "", fmt.Errorf("path must start with %s", PathContest)
	}
	contest = slash + contest
	offset = contest + len(PathContest)
	idEnd := strings.IndexAny(url[offset:], "/?")
	if idEnd >= 0 {
		idEnd = offset + idEnd
	} else {
		idEnd = len(url)
	}
	if _, err := strconv.Atoi(url[offset:idEnd]); err != nil {
		return "", fmt.Errorf("can't parse contest ID from url")
	}
	/* Specify locale for reproducibility */
	return url[:idEnd] + "?locale=en", nil
}
