// Copyright (c) 2016-2017 Marcus Rohrmoser, http://purl.mro.name/recorder
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
// associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute,
// sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or
// substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT
// NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES
// OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// MIT License http://opensource.org/licenses/MIT

// Scrape http://br-klassik.de program schedule + broadcast pages.
//
// import "purl.mro.name/recorder/radio/scrape/b4"

package b4

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	r "purl.mro.name/recorder/radio/scrape"
)

/////////////////////////////////////////////////////////////////////////////
///

/////////////////////////////////////////////////////////////////////////////
/// Just wrap Station into a distinct, local type.
type station r.Station

// Station Factory
//
// Returns a instance conforming to 'scrape.Scraper'
func Station(identifier string) *station {
	switch identifier {
	case
		"b4":
		s := station(r.Station{Name: "Bayern 4", CloseDown: "06:00", ProgramURL: r.MustParseURL("https://www.br-klassik.de/programm/radio/index.html"), Identifier: identifier, TimeZone: localLoc})
		return &s
	}
	return nil
}

func (s *station) String() string {
	return fmt.Sprintf("Station '%s'", s.Name)
}

func (s *station) Matches(nows []time.Time) (ok bool) {
	return true
}

// Synthesise calItemRangeURLs for incremental scraping and queue them up
func (s *station) Scrape() (jobs []r.Scraper, results []r.Broadcaster, err error) {
	now := time.Now()
	for _, t0 := range r.IncrementalNows(now) {
		u, _ := s.calendarItemRangeURLForTime(t0)
		jobs = append(jobs, r.Scraper(u))
	}
	return
}

///////////////////////////////////////////////////////////////////////
// https://www.br-klassik.de/programm/radio/radiosendungen-100~calendar-detail-inner.jsp?date=2016-01-11
// https://www.br-klassik.de/programm/radio/radiosendungen-100~calendarItems.jsp?from=2016-01-10T23:59:59&to=2016-01-11T00:10:00&rows=800
func (s *station) calendarItemRangeURLForTime(t time.Time) (ret *calItemRangeURL, err error) {
	if nil == s {
		panic("aua")
	}
	t0 := t.Add(time.Minute)
	t1 := t0.Add(time.Hour)
	r := calItemRangeURL(r.TimeURL{
		Time:    t0,
		Source:  *r.MustParseURL("https://www.br-klassik.de/programm/radio/radiosendungen-100~calendarItems.jsp?rows=800" + t0.Format("&from=2006-01-02T15:04:05") + t1.Format("&to=2006-01-02T15:04:05")),
		Station: r.Station(*s),
	})
	ret = &r
	// err = errors.New("Not ümplemented yet.")
	return
}

/////////////////////////////////////////////////////////////////////////////
/// Just wrap TimeURL into a distinct, local type - a Scraper, naturally
type calItemRangeURL r.TimeURL

// Fetch calendarItems in given interval (via json)
func (rangeURL *calItemRangeURL) Scrape() (jobs []r.Scraper, results []r.Broadcaster, err error) {
	calendarItems, err := rangeURL.parseCalendarItems()
	if nil != err {
		return
	}
	for _, cis := range calendarItems {
		bc, err := cis.parseBroadcastSeedString(&cis.Html)
		if nil == err {
			jobs = append(jobs, bc)
		}
	}
	return
}

func (rangeURL *calItemRangeURL) Matches(nows []time.Time) (ok bool) {
	return true
}

func (rangeURL *calItemRangeURL) parseCalendarItems() (cis []calendarItem, err error) {
	// fmt.Fprintf(os.Stderr, "GET %s\n", rangeURL.Source.String())
	bo, cr, err := r.HttpGetBody(rangeURL.Source)
	if nil == bo {
		return nil, err
	}
	return rangeURL.parseCalendarItemsReader(bo, cr)
}

func (rangeURL *calItemRangeURL) parseCalendarItemsReader(read io.Reader, cr0 *r.CountingReader) (cis []calendarItem, err error) {
	cr := r.NewCountingReader(io.LimitReader(read, 1048576))
	cis = make([]calendarItem, 0)
	err = json.NewDecoder(cr).Decode(&cis)
	r.ReportLoad("🐦", cr0, cr, rangeURL.Source)
	if nil != err {
		return
	}
	for i := range cis {
		cis[i].Station = &rangeURL.Station
	}
	return
}

/// Time and what it's about.
///
/// datetime from JSON response
/// https://www.br-klassik.de/programm/radio/radiosendungen-100~calendarItems.jsp?rows=800&from=2015-11-30T04:59:59&to=2015-11-30T06:00:00
type Time time.Time

// http://stackoverflow.com/a/25088079
func (t *Time) UnmarshalJSON(b []byte) error {
	tmp, err := time.ParseInLocation(jsonTimeFmt, string(b[:]), localLoc)
	*t = Time(tmp)
	return err
}

var (
	jsonTimeFmt = "\"" + "2006-01-02T15:04:05" + "\""
	localLoc    *time.Location
)

func init() {
	var err error
	localLoc, err = time.LoadLocation("Europe/Berlin")
	if nil != err {
		panic(err)
	}
}

/////////////////////////////////////////////////////////////////////////////
/// item from JSON response
/// https://www.br-klassik.de/programm/radio/radiosendungen-100~calendarItems.jsp?rows=800&from=2015-11-30T04:59:59&to=2015-11-30T06:00:00
type calendarItem struct {
	DateTime Time
	Html     string
	Station  *r.Station
}

func (item *calendarItem) parseBroadcastSeedString(htm *string) (bc *broadcastURL, err error) {
	root, err := html.Parse(strings.NewReader(*htm))
	if nil != err {
		return
	}
	return item.parseBroadcastSeedNode(root)
}

// Get Time, Source and Image from json html snippet
func (item *calendarItem) parseBroadcastSeedNode(root *html.Node) (bc *broadcastURL, err error) {
	bc = &broadcastURL{}
	bc.Station = *item.Station
	bc.Time = time.Time(item.DateTime)
	for _, a := range scrape.FindAll(root, func(n *html.Node) bool {
		if atom.A != n.DataAtom {
			return false
		}
		href := scrape.Attr(n, "href")
		return strings.HasPrefix(href, "/programm/radio/ausstrahlung-") && strings.HasSuffix(href, ".html")
	}) {
		ru, _ := url.Parse(scrape.Attr(a, "href"))
		bc.Source = *item.Station.ProgramURL.ResolveReference(ru)
	}
	for _, img := range scrape.FindAll(root, func(n *html.Node) bool { return atom.Img == n.DataAtom }) {
		ru, _ := url.Parse(scrape.Attr(img, "src"))
		bc.Image = item.Station.ProgramURL.ResolveReference(ru)
	}
	return
}

/////////////////////////////////////////////////////////////////////////////
///
type broadcastURL struct {
	r.BroadcastURL
	Image *url.URL
}

func (bcu *broadcastURL) Matches(nows []time.Time) (ok bool) {
	return true
}

func (bcu *broadcastURL) Scrape() (jobs []r.Scraper, results []r.Broadcaster, err error) {
	bc, err := bcu.parseBroadcast()
	if nil == err {
		results = append(results, bc)
	}
	return
}

/////////////////////////////////////////////////////////////////////////////
/// Parse broadcast
/////////////////////////////////////////////////////////////////////////////
var (
	bcDateRegExp = regexp.MustCompile("(\\d{2}):(\\d{2})\\s+bis\\s+(\\d{2}):(\\d{2})\\s+Uhr")
)

// Completely re-scrape everything and verify consistence at least of Time, evtl. Title
func (bcu *broadcastURL) parseBroadcastNode(root *html.Node) (bc r.Broadcast, err error) {
	bc.Station = bcu.Station
	if "" == bc.Station.Identifier {
		panic("How can the identifier miss?")
	}
	bc.Source = bcu.Source
	bc.Time = bcu.Time
	bc.Image = bcu.Image
	{
		s := "de"
		bc.Language = &s
	}

	for i, main := range scrape.FindAll(root, func(n *html.Node) bool { return atom.Div == n.DataAtom && "br-main-text" == scrape.Attr(n, "class") }) {
		if 1 < i {
			err = errors.New("unexpected 2nd <div class='br-main-text'> ")
			return
		}

		// Subject
		for idx, h3 := range scrape.FindAll(root, func(n *html.Node) bool {
			return atom.H3 == n.DataAtom && "Weitere Informationen" == scrape.Text(n)
		}) {
			// fmt.Fprintf(os.Stderr, "GET %s\n", "uhu")
			if idx != 0 {
				err = errors.New("There was more than 1 <h3>Weitere Informationen")
				return
			}
			for _, a := range scrape.FindAll(h3.Parent, func(n *html.Node) bool {
				return atom.A == n.DataAtom
			}) {
				u, _ := url.Parse(scrape.Attr(a, "href"))
				bc.Subject = bc.Source.ResolveReference(u)
			}

			h3.Parent.Parent.RemoveChild(h3.Parent)
		}

		for i1, h2 := range scrape.FindAll(main, func(n *html.Node) bool { return atom.H2 == n.DataAtom }) {
			if 1 < i1 {
				err = errors.New("unexpected 2nd <h2> ")
				return
			}
			for i4, em := range scrape.FindAll(h2, func(n *html.Node) bool { return atom.Em == n.DataAtom }) {
				if 1 < i4 {
					err = errors.New("unexpected 2nd <em> ")
					return
				}
				bc.Title = scrape.Text(em)
				em.Parent.RemoveChild(em)
			}
			s := scrape.Text(h2)
			bc.TitleSeries = &s

			for i2, h3 := range scrape.FindAll(main, func(n *html.Node) bool { return atom.H3 == n.DataAtom }) {
				if 1 < i2 {
					err = errors.New("unexpected 2nd <h3> ")
					return
				}
				s := scrape.Text(h3)
				bc.TitleEpisode = &s
				h3.Parent.RemoveChild(h3)
			}

			inner := h2.Parent.Parent.Parent
			h2.Parent.RemoveChild(h2)

			for ch := inner.FirstChild; ch != nil; ch = ch.NextSibling {
				if atom.Div == ch.DataAtom {
					inner.RemoveChild(ch) // once removed NextSibling returns nil
				}
			}

			// Description
			description := r.TextWithBrFromNodeSet(scrape.FindAll(inner, func(n *html.Node) bool { return atom.P == n.DataAtom || atom.Div == n.DataAtom }))
			bc.Description = &description
		}
	}

	// DtEnd
	for _, p := range scrape.FindAll(root, func(n *html.Node) bool { return atom.P == n.DataAtom && "br-time" == scrape.Attr(n, "class") }) {
		m := bcDateRegExp.FindStringSubmatch(scrape.Text(p))
		if nil == m {
			err = errors.New("There was no date match")
			return
		}
		i := r.MustParseInt
		// bc.Time = time.Date(i(m[3]), time.Month(i(m[2])), i(m[1]), i(m[4]), i(m[5]), 0, 0, localLoc)
		t := time.Date(bc.Time.Year(), bc.Time.Month(), bc.Time.Day(), i(m[3]), i(m[4]), 0, 0, localLoc)
		if bc.Time.Hour() > t.Hour() || (bc.Time.Hour() == t.Hour() && bc.Time.Minute() > t.Minute()) { // after midnight
			t = t.AddDate(0, 0, 1)
		}
		bc.DtEnd = &t
	}

	// Modified
	for idx, meta := range scrape.FindAll(root, func(n *html.Node) bool {
		return atom.Meta == n.DataAtom && "og:article:modified_time" == scrape.Attr(n, "property")
	}) {
		if idx != 0 {
			err = errors.New("There was more than 1 <meta property='og:article:modified_time'/>")
			return
		}
		v, _ := time.Parse(time.RFC3339, scrape.Attr(meta, "content"))
		bc.Modified = &v
	}

	// Author
	for idx, meta := range scrape.FindAll(root, func(n *html.Node) bool {
		return atom.Meta == n.DataAtom && "author" == scrape.Attr(n, "name")
	}) {
		if idx != 0 {
			err = errors.New("There was more than 1 <meta name='author'/>")
			return
		}
		s := scrape.Attr(meta, "content")
		bc.Author = &s
	}

	return
}

func (bcu *broadcastURL) parseBroadcastReader(read io.Reader, cr0 *r.CountingReader) (bc r.Broadcast, err error) {
	cr := r.NewCountingReader(read)
	root, err := html.Parse(cr)
	r.ReportLoad("🐠", cr0, cr, bcu.Source)
	if nil != err {
		return
	}
	return bcu.parseBroadcastNode(root)
}

func (bcu *broadcastURL) parseBroadcast() (bc r.Broadcast, err error) {
	bo, cr, err := r.HttpGetBody(bcu.Source)
	if nil == bo {
		return bc, err
	}
	return bcu.parseBroadcastReader(bo, cr)
}
