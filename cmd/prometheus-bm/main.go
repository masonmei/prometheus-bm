package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/ddliu/go-httpclient"
	"github.com/jessevdk/go-flags"
	"github.com/tidwall/gjson"
)

// Option define option info
type Option struct {
	Endpoints   []string `short:"e" long:"endpoints" description:"prometheus endpoints"`
	Concurrency int      `short:"c" long:"concurrency" description:"request concurrency"`
}

func main() {

	var opt Option
	_, err := flags.Parse(&opt)
	if err != nil {
		log.Fatal("Parse error:", err)
	}

	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: "prometheus benchmark client",
		"Accept-Language":        "zh-ch",
	})

	endpoints := opt.Endpoints

	var queue []string

	for i := range endpoints {
		ep := endpoints[i]
		targets := genLabelValuesRequests(ep)
		for j := range targets {
			target := targets[j]
			queue = append(queue, target)
		}
		targets = genSeriesRequests(ep)
		for j := range targets {
			target := targets[j]
			queue = append(queue, target)
		}

		targets = genQueryRequests(ep, []time.Duration{
			5 * time.Minute,
			30 * time.Minute,
			time.Hour,
			2 * time.Hour,
			3 * time.Hour,
			6 * time.Hour,
		})
		for j := range targets {
			target := targets[j]
			queue = append(queue, target)
		}

		targets = genQueryRangeRequests(ep, []time.Duration{
			5 * time.Minute,
			30 * time.Minute,
			time.Hour,
			2 * time.Hour,
			3 * time.Hour,
			6 * time.Hour,
		})
		for j := range targets {
			target := targets[j]
			queue = append(queue, target)
		}
	}

	var wg = sync.WaitGroup{}

	g := NewLimiter(opt.Concurrency)
	for i, _ := range queue {
		wg.Add(1)
		task := queue[i]

		goFunc := func() {
			// 做一些业务逻辑处理
			fmt.Printf("go func: %s\n", task)
			res, err := httpclient.Get(task)
			if err != nil {
				log.Printf("Failed to get task: %s", task)
			} else {
				_, _ = ioutil.ReadAll(res.Body)
			}
			wg.Done()
		}
		g.Run(goFunc)
	}

	wg.Wait()
}

func genLabelValuesRequests(endpoint string) []string {
	res, err := httpclient.Get(endpoint + "/api/v1/labels")
	if err != nil {
		return []string{}
	}
	body, _ := ioutil.ReadAll(res.Body)
	bodystring := string(body)

	var result []string
	get := gjson.Get(bodystring, "data")
	array := get.Array()
	for i := range array {
		result = append(result, fmt.Sprintf("%s/api/v1/label/%s/values", endpoint, array[i].Str))
	}
	return result
}

func genSeriesRequests(endpoint string) []string {
	// res, err := httpclient.Get(endpoint + "/api/v1/label/__name__/values")
	// if err != nil {
	// 	return []string{}
	// }
	// body, _ := ioutil.ReadAll(res.Body)
	// bodystring := string(body)

	var result []string
	// get := gjson.Get(bodystring, "data")
	// array := get.Array()
	// for i := range array {
	// 	result = append(result, fmt.Sprintf("%s/api/v1/series?match[]=%s", endpoint, array[i].Str))
	// }
	return result
}

func genQueryRequests(endpoint string, tr []time.Duration) []string {
	return nil
}

func genQueryRangeRequests(endpoint string, tr []time.Duration) []string {

	return nil
}

type ConcurrentLimiter struct {
	n int
	c chan struct{}
}

// initialization ConcurrentLimiter struct
func NewLimiter(concurrency int) *ConcurrentLimiter {
	return &ConcurrentLimiter{
		n: concurrency,
		c: make(chan struct{}, concurrency),
	}
}

// Run f in a new goroutine but with limit.
func (g *ConcurrentLimiter) Run(f func()) {
	g.c <- struct{}{}
	go func() {
		f()
		<-g.c
	}()
}
