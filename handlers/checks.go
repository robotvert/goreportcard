package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/k0kubun/pp"
	"github.com/robotvert/goreportcard/check"
)

type score struct {
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	FileSummaries []check.FileSummary `json:"file_summaries"`
	Percentage    float64             `json:"percentage"`
}

type ChecksResp struct {
	Checks      []score   `json:"checks"`
	Average     float64   `json:"average"`
	Grade       Grade     `json:"grade"`
	Files       int       `json:"files"`
	Issues      int       `json:"issues"`
	Repo        string    `json:"repo"`
	LastRefresh time.Time `json:"last_refresh"`
}

func CheckPackage(path string) (ChecksResp, error) {

	filenames, err := check.GoFiles(path)
	pp.Println(filenames)
	if err != nil {
		return ChecksResp{}, fmt.Errorf("Could not get filenames: %v", err)
	}
	if len(filenames) == 0 {
		return ChecksResp{}, fmt.Errorf("No .go files found")
	}
	checks := []check.Check{check.GoFmt{Dir: path, Filenames: filenames},
		check.GoVet{Dir: path, Filenames: filenames},
		check.GoLint{Dir: path, Filenames: filenames},
		check.GoCyclo{Dir: path, Filenames: filenames},
	}

	ch := make(chan score)
	for _, c := range checks {
		go func(c check.Check) {
			p, summaries, err := c.Percentage()
			if err != nil {
				log.Printf("ERROR: (%s) %v", c.Name(), err)
			}
			s := score{
				Name:          c.Name(),
				Description:   c.Description(),
				FileSummaries: summaries,
				Percentage:    p,
			}
			ch <- s
		}(c)
	}

	resp := ChecksResp{Repo: path,
		Files:       len(filenames),
		LastRefresh: time.Now().UTC()}
	var avg float64
	var issues = make(map[string]bool)
	for i := 0; i < len(checks); i++ {
		s := <-ch
		resp.Checks = append(resp.Checks, s)
		avg += s.Percentage
		for _, fs := range s.FileSummaries {
			issues[fs.Filename] = true
		}
	}

	resp.Average = avg / float64(len(checks))
	resp.Issues = len(issues)
	resp.Grade = grade(resp.Average * 100)

	return resp, nil
}
