package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gophercises/quiet_hn/hn"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.html"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Println("Running on Port:", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	sc := storyCache{
		numStories: numStories,
		duration:   6 * time.Second,
	}

	go func() {
		tk := time.NewTicker(3 * time.Second)
		for {
			newStoryCache := storyCache{
				numStories: numStories,
				duration:   6 * time.Second,
			}

			newStoryCache.stories()
			sc.mutex.Lock()
			sc.cache = newStoryCache.cache
			sc.expiration = newStoryCache.expiration
			sc.mutex.Unlock()
			<-tk.C
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := sc.stories()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := templateData{
			Stories: stories,
			Time:    time.Since(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

type storyCache struct {
	numStories int
	cache      []item
	expiration time.Time
	duration   time.Duration
	mutex      sync.Mutex
}

func (sc *storyCache) stories() ([]item, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	if time.Now().Before(sc.expiration) {
		return sc.cache, nil
	}

	stories, err := getTopStories(sc.numStories)
	if err != nil {
		return nil, err
	}

	sc.expiration = time.Now().Add(sc.duration)
	sc.cache = stories

	return sc.cache, nil
}

/* var (
	cache           []item
	cacheExpiration time.Time
	cacheMutex      sync.Mutex
) */

/* func getCachedStories(numStories int) ([]item, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	if time.Now().Before(cacheExpiration) {
		return cache, nil
	}

	stories, err := getTopStories(numStories)
	if err != nil {
		return nil, err
	}

	cache = stories
	cacheExpiration = time.Now().Add(5 * time.Minute)

	return stories, nil
} */

func getTopStories(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("failed to load top stories")
	}
	var stories []item
	start := 0
	for len(stories) < numStories {
		missing := (numStories - len(stories)) * 5 / 4
		stories = append(stories, getStories(ids[start:start+missing])...)
		start += missing
	}

	return stories[:numStories], nil
}

func getStories(ids []int) []item {
	type result struct {
		idx  int
		item item
		err  error
	}

	resultCh := make(chan result)
	for i := 0; i < len(ids); i++ {
		var client hn.Client
		go func(idx, id int) {
			hnItem, err := client.GetItem(id)
			if err != nil {
				resultCh <- result{idx: idx, err: err}
			}
			resultCh <- result{idx: idx, item: parseHNItem(hnItem)}
		}(i, ids[i])
	}

	var results []result
	for i := 0; i < len(ids); i++ {
		results = append(results, <-resultCh)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].idx < results[j].idx
	})

	var stories []item
	for _, res := range results {
		if res.err != nil {
			continue
		}
		if isStoryLink(res.item) {
			stories = append(stories, res.item)
		}
	}

	return stories
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
