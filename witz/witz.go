package witz

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/gocolly/colly"
)

var quelle = "https://witze-paradies.de/"
var witzapi = "https://witzapi.de/api/joke"
var endPoints = []string{"geile-witze/", "busfahrer-witze/", "bier-witze/", "beamte-witze", "blondinenwitze/", "bodenlose-witze/", "bundeswehr-witze/", "dumme-witze/", "blode-spruche/", "arzt-witze/", "bose-witze/", "buro-witze/"}

type Witz struct {
	Thema string
	Text  string `json:"text"`
}
type W struct {
	data []map[string]string
}

func GetWitz() Witz {
	var witz Witz
	w := rand.Intn(2)
	idx := rand.Intn(len(endPoints))
	switch w {
	case 0:
		url := quelle + endPoints[idx]
		fmt.Println(url)
		witz.Thema = endPoints[idx]

		c := colly.NewCollector()
		c.OnHTML(`.blog_post_single`, func(h *colly.HTMLElement) {
			// fmt.Println("NAME ", h.Name)
			// fmt.Println("Attr class ", h.Attr("class"))
			h.ForEach(`.blog_content`, func(i int, b *colly.HTMLElement) {
				witz.Text = witz.Text + b.Text
				// witz.Text = witz.Text + "\n"
			})
		})
		c.Visit(url)
		if witz.Text == "" {
			witz.Text = "Kein Witz gefunden !!!"
		}
	case 1:
		url := witzapi
		resp, err := http.Get(url)
		if err != nil {
			witz.Text = err.Error()
			return witz
		}
		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			witz.Text = err.Error()
			return witz
		}
		fmt.Println(string(b))
		var w W
		w.data = make([]map[string]string, 0)
		if err = json.Unmarshal(b, &w.data); err != nil {
			witz.Text = err.Error() + " hallo form error"
			return witz
		}
		witz.Text = w.data[0]["text"]
	}
	return witz
}
