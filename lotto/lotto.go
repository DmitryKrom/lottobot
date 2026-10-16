package lotto

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
)

const (
	ZAHLENCLASS     = "mod-simple-list__inline"
	ZUSATZZAHLCLASS = "extra-number"
	HEADINGCLASS    = "heading-h5"
	COL_MD_12       = "col-md-12"
	POLYGON_LABEL   = "polygon-label"
	POLYGON_ZAHLEN  = "polygon-outlined"
	POLYGON_EXTRA   = "polygon-big"
	SUPERZAHL       = "Superzahl"
)

type Play struct {
	Lotto_zahlen  []int
	Zusatz_zahlen []int
}

var Play_list []Play
var Datum string
var Wg sync.WaitGroup

func Get_lotto(l_art int, msg_text string) (Play, error) {
	var lotto_art string

	if l_art == 3 {
		lotto_art = "eurojackpot"
	}
	if l_art == 4 {
		lotto_art = "lotto-6aus49"
	}
	// lotto_dat := strings.TrimPrefix(msg_text, lotto_art)
	lottoZahlen, err := get_lotto_gewinnzahlen(lotto_art, msg_text)
	return lottoZahlen, err
}

func NewPlay() Play {
	return Play{
		Lotto_zahlen:  make([]int, 0),
		Zusatz_zahlen: make([]int, 0),
	}
}

func (p *Play) sortSpielZahlen() {
	sort.Ints(p.Lotto_zahlen)
	sort.Ints(p.Zusatz_zahlen)
}

func (p *Play) printSpielZahlen() {
	fmt.Println("####### Lotto ##########")
	for i := 0; i < len(p.Lotto_zahlen); i++ {
		fmt.Print(p.Lotto_zahlen[i])
		fmt.Print("  ")
	}
	fmt.Println()
	fmt.Println("####### Zusatz ##########")
	for i := 0; i < len(p.Zusatz_zahlen); i++ {
		fmt.Print(p.Zusatz_zahlen[i])
		fmt.Print("  ")
	}
	fmt.Println()
	fmt.Println()
}

func get_lotto_gewinnzahlen(art, datum string) (Play, error) {
	d := datum
	var lotto DBOperation
	if art == "lotto-6aus49" {
		datum = "date=" + datum
		lotto = New6aus49()
	}
	if art == "eurojackpot" {
		datum = "datum=" + datum
		lotto = NewJP()
	}

	db := NewDatabase()
	db.Connect("postgres://postgres:postgres@localhost:5432/lotto")
	defer db.Pool.Close()
	fmt.Println("  --- vor ---")
	err := lotto.Select(d, db.Pool)
	fmt.Println("  --- nach ---")

	if err == nil {
		switch art {
		case "lotto-6aus49":
			lotto6aus49 := lotto.(*Lotto6aus49Zahlen)
			gewinnZahlen := NewPlay()
			erst, _ := strconv.Atoi(lotto6aus49.firstNum)
			zw, _ := strconv.Atoi(lotto6aus49.secondNum)
			dr, _ := strconv.Atoi(lotto6aus49.thirdNum)
			vier, _ := strconv.Atoi(lotto6aus49.fourthNum)
			fuen, _ := strconv.Atoi(lotto6aus49.fifthNum)
			sechs, _ := strconv.Atoi(lotto6aus49.sixthNum)
			s1, _ := strconv.Atoi(lotto6aus49.super)

			gewinnZahlen.Lotto_zahlen = append(gewinnZahlen.Lotto_zahlen, erst)
			gewinnZahlen.Lotto_zahlen = append(gewinnZahlen.Lotto_zahlen, zw)
			gewinnZahlen.Lotto_zahlen = append(gewinnZahlen.Lotto_zahlen, dr)
			gewinnZahlen.Lotto_zahlen = append(gewinnZahlen.Lotto_zahlen, vier)
			gewinnZahlen.Lotto_zahlen = append(gewinnZahlen.Lotto_zahlen, fuen)
			gewinnZahlen.Lotto_zahlen = append(gewinnZahlen.Lotto_zahlen, sechs)
			gewinnZahlen.Zusatz_zahlen = append(gewinnZahlen.Zusatz_zahlen, s1)
			return gewinnZahlen, nil

		case "eurojackpot":
			eurojackpot := lotto.(*EurojackpotZahlen)
			gewinnZahlen := NewPlay()
			erst, _ := strconv.Atoi(eurojackpot.firstNum)
			zw, _ := strconv.Atoi(eurojackpot.secondNum)
			dr, _ := strconv.Atoi(eurojackpot.thirdNum)
			vier, _ := strconv.Atoi(eurojackpot.fourthNum)
			fuen, _ := strconv.Atoi(eurojackpot.fifthNum)
			s1, _ := strconv.Atoi(eurojackpot.supOne)
			s2, _ := strconv.Atoi(eurojackpot.supTwo)

			gewinnZahlen.Lotto_zahlen = append(gewinnZahlen.Lotto_zahlen, erst)
			gewinnZahlen.Lotto_zahlen = append(gewinnZahlen.Lotto_zahlen, zw)
			gewinnZahlen.Lotto_zahlen = append(gewinnZahlen.Lotto_zahlen, dr)
			gewinnZahlen.Lotto_zahlen = append(gewinnZahlen.Lotto_zahlen, vier)
			gewinnZahlen.Lotto_zahlen = append(gewinnZahlen.Lotto_zahlen, fuen)
			gewinnZahlen.Zusatz_zahlen = append(gewinnZahlen.Zusatz_zahlen, s1)
			gewinnZahlen.Zusatz_zahlen = append(gewinnZahlen.Zusatz_zahlen, s2)
			return gewinnZahlen, nil
		}
	}

	switch art {
	case "lotto-6aus49":
		gewinnzahlen := visitLotto(datum)
		if len(gewinnzahlen.Lotto_zahlen) > 0 && len(gewinnzahlen.Zusatz_zahlen) > 0 {
			gewinnzahlen.sortSpielZahlen()
			lotto6aus49 := lotto.(*Lotto6aus49Zahlen)
			lotto6aus49.firstNum = strconv.Itoa(gewinnzahlen.Lotto_zahlen[0])
			lotto6aus49.secondNum = strconv.Itoa(gewinnzahlen.Lotto_zahlen[1])
			lotto6aus49.thirdNum = strconv.Itoa(gewinnzahlen.Lotto_zahlen[2])
			lotto6aus49.fourthNum = strconv.Itoa(gewinnzahlen.Lotto_zahlen[3])
			lotto6aus49.fifthNum = strconv.Itoa(gewinnzahlen.Lotto_zahlen[4])
			lotto6aus49.sixthNum = strconv.Itoa(gewinnzahlen.Lotto_zahlen[5])
			lotto6aus49.super = strconv.Itoa(gewinnzahlen.Zusatz_zahlen[0])
			lotto6aus49.Insert(d, db.Pool)
			return gewinnzahlen, nil
		}
		return gewinnzahlen, nil
	case "eurojackpot":
		gewinnzahlen, err := visitJackpot(datum)
		if err != nil {
			// fmt.Println("[ERROR: ] Seite meldet: ", err)
			return Play{}, err
		}
		if len(gewinnzahlen.Lotto_zahlen) > 0 && len(gewinnzahlen.Zusatz_zahlen) > 0 {
			gewinnzahlen.sortSpielZahlen()
			eurojackpot := lotto.(*EurojackpotZahlen)
			eurojackpot.firstNum = strconv.Itoa(gewinnzahlen.Lotto_zahlen[0])
			eurojackpot.secondNum = strconv.Itoa(gewinnzahlen.Lotto_zahlen[1])
			eurojackpot.thirdNum = strconv.Itoa(gewinnzahlen.Lotto_zahlen[2])
			eurojackpot.fourthNum = strconv.Itoa(gewinnzahlen.Lotto_zahlen[3])
			eurojackpot.fifthNum = strconv.Itoa(gewinnzahlen.Lotto_zahlen[4])
			eurojackpot.supOne = strconv.Itoa(gewinnzahlen.Zusatz_zahlen[0])
			eurojackpot.supTwo = strconv.Itoa(gewinnzahlen.Zusatz_zahlen[1])
			eurojackpot.Insert(d, db.Pool)
			return gewinnzahlen, nil
		}
		return gewinnzahlen, nil
	}
	// fmt.Println("*** Gewinnzahlen ****")
	// gewinnzahlen.printSpielZahlen()

	return Play{}, nil
}

func visitLotto(datum string) Play {
	gewinnZahlen := NewPlay()
	c := colly.NewCollector()

	c.OnHTML("div."+COL_MD_12, func(h *colly.HTMLElement) {
		h.ForEach("span."+POLYGON_ZAHLEN, func(i int, span *colly.HTMLElement) {
			// zahl, _ := strings.CutPrefix(span.Text, "\n")
			// zahl, _ = strings.CutSuffix(zahl, "\n")
			zahl := strings.TrimSpace(span.Text)
			// fmt.Println(zahl)
			z, err := strconv.Atoi(zahl)
			if err != nil {
				fmt.Println(err)
			}
			gewinnZahlen.Lotto_zahlen = append(gewinnZahlen.Lotto_zahlen, z)
		})
		if h.ChildText("p."+HEADINGCLASS) == SUPERZAHL {
			s := strings.TrimSpace(h.ChildText("span." + POLYGON_EXTRA))
			// fmt.Println(s)
			superzahl, _ := strconv.Atoi(s)
			gewinnZahlen.Zusatz_zahlen = append(gewinnZahlen.Zusatz_zahlen, superzahl)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		fmt.Println("ProxyURL", r.ProxyURL)
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("OnError: ", err)
		fmt.Println("responsFunc Status : ", r.StatusCode)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("OnResponce Status: ", r.StatusCode)
	})
	err := c.Visit(fmt.Sprintf("https://www.westlotto.de/lotto-6aus49/gewinnzahlen/gewinnzahlen.html?%s", datum))
	if err != nil {
		fmt.Printf("[ERROR] visiting the site: %v\n", err)
	}
	return gewinnZahlen
}

func visitJackpot(datum string) (Play, error) {
	gewinnZahlen := NewPlay()
	c := colly.NewCollector()
	c.OnHTML("ul."+ZAHLENCLASS, func(h *colly.HTMLElement) {
		h.ForEach("li", func(_ int, li *colly.HTMLElement) {
			z, _ := strconv.Atoi(li.Text)
			fmt.Println("eurojackpot Zahl: ", z)
			if li.Attr("class") == ZUSATZZAHLCLASS {
				gewinnZahlen.Zusatz_zahlen = append(gewinnZahlen.Zusatz_zahlen, z)
			}
			if li.Attr("class") == "" {
				gewinnZahlen.Lotto_zahlen = append(gewinnZahlen.Lotto_zahlen, z)
			}

		})

	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("OnError: ", err)
		fmt.Println("responsFunc Status : ", r.StatusCode)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("OnResponce Status: ", r.StatusCode)
		// fmt.Println("OnResponce URL: ", string(r.Body))
	})
	err := c.Visit(fmt.Sprintf("https://www.westlotto.de/eurojackpot/gewinnzahlen/gewinnzahlen.html?%s", datum))
	if err != nil {
		fmt.Println("[ERROR] visiting site: ", err)
	}
	if len(gewinnZahlen.Lotto_zahlen) == 0 {
		ctx, cancel := chromedp.NewContext(
			context.Background(),
		)
		defer cancel()
		fmt.Println("--- chromedp ---")
		var screenshotBuffer []byte
		err := chromedp.Run(ctx,
			chromedp.Navigate(fmt.Sprintf("https://www.westlotto.de/eurojackpot/gewinnzahlen/gewinnzahlen.html?%s", datum)),
			chromedp.Sleep(500*time.Millisecond),
			// chromedp.FullScreenshot(&screenshotBuffer, 100),
			chromedp.Screenshot(`div.mod-content`, &screenshotBuffer),
			// chromedp.CaptureScreenshot(&screenshotBuffer),
			chromedp.Sleep(500*time.Millisecond),
		)
		if err != nil {
			fmt.Println(err.Error())
			return Play{}, err
		}
		err = os.WriteFile("jackpot.png", screenshotBuffer, 0644)
		if err != nil {
			return Play{}, errors.New("Konnte keine Datei jackpot.png schreiben !!!")
		}
		fmt.Println("jackpot.png geschrieben !!!")
		return Play{}, errors.New("screenshot")
	}
	return gewinnZahlen, nil
}
