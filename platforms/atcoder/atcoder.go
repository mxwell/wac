package atcoder

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mxwell/wac/model"
)

type AtCoder struct {
}

func InitAtCoder() model.Platform {
	return AtCoder{}
}

func (a AtCoder) ValidUrl(url string) bool {
	_, err := trimUrl(url)
	return err == nil
}

func (a AtCoder) GetContest(url string, rootDirName string) (*model.Contest, error) {
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
			log.Printf("WARN unable to parse task %d\n", i)
			return
		}
		tokenElement := column.Find("a.linkwrapper")
		if tokenElement.Length() != 1 {
			log.Printf("WARN unable to find token of task %d\n", i)
			return
		}
		token := strings.ToLower(tokenElement.Text())
		href, ok := tokenElement.Attr("href")
		if !ok {
			log.Printf("WARN unable to extract link to task %d\n", i)
			return
		}
		column = column.Next()
		if column.Length() == 0 {
			log.Printf("WARN unable to parse task %d\n", i)
			return
		}
		nameElement := column.Find("a.linkwrapper")
		if nameElement.Length() != 1 {
			log.Printf("WARN unable to find name of task %d\n", i)
			return
		}
		name := nameElement.Text()
		tasks[token] = model.Task{url + href, name, token, make([]string, 0)}
	})
	return &model.Contest{url, title, tasks, rootDirName}, nil
}

func contains(arr *[]int, value int) bool {
	for _, e := range *arr {
		if e == value {
			return true
		}
	}
	return false
}

func (a AtCoder) GetTests(task *model.Task) ([]model.Test, error) {
	doc, err := goquery.NewDocument(task.Link)
	if err != nil {
		return nil, err
	}
	statementElement := doc.Find("#task-statement")
	if statementElement.Length() != 1 {
		return nil, fmt.Errorf("can't detect task-statement uniquely: %d item(s) found", statementElement.Length())
	}
	enSpanElement := statementElement.Find("span.lang-en")
	if enSpanElement.Length() != 1 {
		return nil, fmt.Errorf("can't detect span in English uniquely")
	}

	sampleInputs := make(map[int]string)
	sampleOutputs := make(map[int]string)
	/* add sample IDs in order of appearance to restore the order later */
	var idOrder []int

	enSpanElement.Find("div.part").Each(func(i int, s *goquery.Selection) {
		headerElement := s.Find("section h3")
		if headerElement.Length() != 1 {
			return
		}
		preElement := s.Find("section pre")
		if preElement.Length() != 1 {
			return
		}
		header := headerElement.Text()
		pre := preElement.Text()
		if strings.HasPrefix(header, "Sample Input ") {
			if id, err := strconv.Atoi(strings.TrimPrefix(header, "Sample Input ")); err == nil {
				sampleInputs[id] = pre
				if !contains(&idOrder, id) {
					idOrder = append(idOrder, id)
				}
			} else {
				log.Printf("WARN can't parse sample id from header '%s'\n", header)
			}
		} else if strings.HasPrefix(header, "Sample Output ") {
			if id, err := strconv.Atoi(strings.TrimPrefix(header, "Sample Output ")); err == nil {
				sampleOutputs[id] = pre
				if !contains(&idOrder, id) {
					idOrder = append(idOrder, id)
				}
			} else {
				log.Printf("WARN can't parse sample id from header '%s'\n", header)
			}
		}
	})

	var result []model.Test
	for _, id := range idOrder {
		input := sampleInputs[id]
		output, ok := sampleOutputs[id]
		if !ok {
			log.Printf("WARN no output for sample %d\n", id)
			continue
		}
		token := "sample" + strconv.Itoa(id)
		result = append(result, model.Test{token, input, output})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("tests are not found")
	}

	return result, nil
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
