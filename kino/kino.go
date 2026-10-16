package kino

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
)

type Film struct {
	FilmName  string
	FilmImg   string
	FilmPreis string
}
type ShowFilm struct {
	Type        string `json:"_type"`
	Id          int64  `json:"id"`
	Cinema      int64  `json:"cinema"`
	Blockbuster bool   `json:"blockbuster"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	HasTrailer  bool   `json:"hasTrailer"`
	// Attributes       []string      `json:"attributes,omitempty"`
	Showtimes        []Showtime    `json:"showtimes,omitempty"`
	Date             string        `json:"date"`
	Movie            int64         `json:"movie"`
	RelatedShows     []RelatedShow `json:"relatedShows,omitempty"`
	PosterPreload    string        `json:"poster_preload"`
	Poster           string        `json:"poster"`
	DetailLink       string        `json:"detailLink"`
	StartDate        string        `json:"startDate"`
	ScreeningWeek    int           `json:"screeningWeek"`
	ShowtimeSchedule Scheduler     `json:"showtimeSchedule"`
	Description      string
}
type Scheduler struct{}

type RelatedShow struct{}

type Showtime struct {
	Id       int64  `json:"id"`
	Datetime string `json:"datetime"`
}

func NewFilm(name string, img string, preis string) *Film {
	return &Film{FilmName: name,
		FilmImg:   img,
		FilmPreis: preis,
	}
}

func GetKino(filter int) []*ShowFilm {
	// initialize a controllable Chrome instance

	switch filter {
	case 7:
		f := filmeHeute()
		return f
	case 8:
		f := filmeWoche()
		return f
		// case "cinema vorverkauf":
		// 	// f := filmVorverkauf()
		// case "cinema events":
		// case "cinema vorschau":
	}
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	// to release the browser resources when
	// it is no longer needed
	defer cancel()
	fmt.Println("in GETKINO")

	var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		// chromedp.Navigate("https://www.cinestar.de"),
		// chromedp.Click(`//*[@title="Akzeptieren & Schließen"]`, chromedp.BySearch),
		chromedp.Navigate("https://www.cinestar.de/kino-siegen#kinoprogramm"),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.WaitVisible(`.ShowTile`, chromedp.BySearch),
		chromedp.Nodes(`.ShowTile`, &nodes, chromedp.ByQueryAll, chromedp.BySearch),
	)

	fmt.Println("NODES: ", len(nodes))
	if err != nil {
		fmt.Println("RUN CHROMEDP ERROR: ", err)
	}
	var ok = true
	var f_name string
	var f_image string
	// var f_preis string
	var filme []*Film
	for _, node := range nodes {
		// extract data from the product HTML node
		err = chromedp.Run(ctx,
			chromedp.Text("span.title", &f_name, chromedp.ByQuery, chromedp.FromNode(node)),
			chromedp.AttributeValue("img", "src", &f_image, &ok, chromedp.ByQuery, chromedp.FromNode(node)),
		)
		if err != nil {
			log.Fatal("Error:", err)
		}
		film := NewFilm(f_name, f_image, "0")
		filme = append(filme, film)
	}
	return nil
}

func filmVorverkauf() []*ShowFilm {
	url := "https://www.cinestar.de/api/cinema/46/show/"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error building Request -> ", err)
		return nil
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making Request -> ", err)
		return nil
	}
	defer resp.Body.Close()
	reader, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading Body to Bytes -> ", err)
		return nil
	}
	var showFilms = make([]ShowFilm, 0)
	if err := json.Unmarshal(reader, &showFilms); err != nil {
		fmt.Println("Error unmarshaling Data -> ", err)
	}
	var filme []*ShowFilm
	now := time.Now()
	wg := &sync.WaitGroup{}
	for i := range showFilms {
		filmTimesHeute := make([]Showtime, 0)

		for x, times := range showFilms[i].Showtimes {
			dateAndTime := strings.Fields(times.Datetime)
			dT := dateAndTime[0] + " " + dateAndTime[1] + ":00"
			datum, err := time.Parse(time.DateTime, dT)
			if err != nil {
				fmt.Println("		ERROR Parsing datum ", datum, "  ", dT, " ", err)
			}
			w := now.Weekday()
			daysToSunday := 7 - w
			sunDay := time.Date(now.Year(), now.Month(), now.Day()+int(daysToSunday)+7, 0, 0, 0, 0, time.UTC)
			// fmt.Println("Datum : ", datum.String())
			// fmt.Println("Next  : ", sunDay.String())

			if datum.After(now) && sunDay.After(datum) {
				// film := NewFilm(showFilms[i].Title, showFilms[i].Poster, "0")
				filmTimesHeute = append(filmTimesHeute, showFilms[i].Showtimes[x])
			}
		}
		if len(filmTimesHeute) > 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				description(showFilms[i].DetailLink, i, &showFilms[i].Description)
			}()
			// showFilms[i].Description = descr
			// showFilms[i].Showtimes = filmTimesHeute
			filme = append(filme, &showFilms[i])
		}
	}
	wg.Wait()
	return filme

}

func filmeWoche() []*ShowFilm {
	url := "https://www.cinestar.de/api/cinema/46/show/"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error building Request -> ", err)
		return nil
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making Request -> ", err)
		return nil
	}
	defer resp.Body.Close()
	reader, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading Body to Bytes -> ", err)
		return nil
	}
	var showFilms = make([]ShowFilm, 0)
	if err := json.Unmarshal(reader, &showFilms); err != nil {
		fmt.Println("Error unmarshaling Data -> ", err)
	}
	var filme []*ShowFilm
	now := time.Now()
	wg := &sync.WaitGroup{}
	for i := range showFilms {
		filmTimesHeute := make([]Showtime, 0)

		for x, times := range showFilms[i].Showtimes {
			dateAndTime := strings.Fields(times.Datetime)
			dT := dateAndTime[0] + " " + dateAndTime[1] + ":00"
			datum, err := time.Parse(time.DateTime, dT)
			if err != nil {
				fmt.Println("		ERROR Parsing datum ", datum, "  ", dT, " ", err)
			}
			w := now.Weekday()
			daysToSunday := 7 - w
			sunDay := time.Date(now.Year(), now.Month(), now.Day()+int(daysToSunday)+7, 0, 0, 0, 0, time.UTC)
			// fmt.Println("Datum : ", datum.String())
			// fmt.Println("Next  : ", sunDay.String())

			if datum.After(now) && sunDay.After(datum) {
				// film := NewFilm(showFilms[i].Title, showFilms[i].Poster, "0")
				filmTimesHeute = append(filmTimesHeute, showFilms[i].Showtimes[x])
			}
		}
		if len(filmTimesHeute) > 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				description(showFilms[i].DetailLink, i, &showFilms[i].Description)
			}()
			// showFilms[i].Description = descr
			// showFilms[i].Showtimes = filmTimesHeute
			filme = append(filme, &showFilms[i])
		}
	}
	wg.Wait()
	return filme
}

func filmeHeute() []*ShowFilm {
	url := "https://www.cinestar.de/api/cinema/46/show/"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("Error building Request -> ", err)
		return nil
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making Request -> ", err)
		return nil
	}
	defer resp.Body.Close()
	reader, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading Body to Bytes -> ", err)
		return nil
	}
	var showFilms = make([]ShowFilm, 0)
	if err := json.Unmarshal(reader, &showFilms); err != nil {
		fmt.Println("Error unmarshaling Data -> ", err)
	}
	var filme []*ShowFilm
	now := time.Now()
	wg := &sync.WaitGroup{}
	for i := range showFilms {
		filmTimesHeute := make([]Showtime, 0)
		if strings.Contains(showFilms[i].Title, "igital") {
			// fmt.Println(showFilms[i].Title)
			// fmt.Println(showFilms[i].Showtimes)
		}

		for x, times := range showFilms[i].Showtimes {
			dateAndTime := strings.Fields(times.Datetime)
			dT := dateAndTime[0] + " " + dateAndTime[1] + ":00"
			datum, err := time.Parse(time.DateTime, dT)
			if err != nil {
				fmt.Println("		ERROR Parsing datum ", datum, "  ", dT, " ", err)
			}
			nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
			// fmt.Println("Datum : ", datum.String())
			// fmt.Println("Next  : ", nextDay.String())

			if datum.After(now) && nextDay.After(datum) {
				// film := NewFilm(showFilms[i].Title, showFilms[i].Poster, "0")
				filmTimesHeute = append(filmTimesHeute, showFilms[i].Showtimes[x])
			}
		}
		if len(filmTimesHeute) > 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				description(showFilms[i].DetailLink, i, &showFilms[i].Description)
			}()
			showFilms[i].Showtimes = filmTimesHeute
			filme = append(filme, &showFilms[i])
		}
	}
	wg.Wait()
	// for y := range filme {
	// fmt.Println("Titel: ", filme[y].Title)
	// fmt.Println("Beschreibung: ", filme[y].Description)
	// fmt.Println("----------------------------------------------------")
	// }
	return filme
}

func description(path string, idx int, sf *string) {
	url := "https://www.cinestar.de" + path
	var description string
	c := colly.NewCollector()
	c.OnHTML(`.textElement`, func(h *colly.HTMLElement) {
		description = h.Text
	})
	c.Visit(url)
	*sf = description
}
