// Copyright (c) 2015-2017 Marcus Rohrmoser, http://purl.mro.name/recorder
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

package dlf // import "purl.mro.name/recorder/radio/scrape/dlf"

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	r "purl.mro.name/recorder/radio/scrape"
)

/////////////////////////////////////////////////////////////////////////////
/// Just wrap Station into a distinct, local type - a Scraper, naturally
type station r.Station

// Station Factory
func Station(identifier string) *station {
	tz, _ := time.LoadLocation("Europe/Berlin")
	switch identifier {
	case
		"dlf":
		s := station(r.Station{Name: "Deutschlandfunk", CloseDown: "00:00", ProgramURL: r.MustParseURL("http://www.deutschlandfunk.de/programmvorschau.281.de.html"), Identifier: identifier, TimeZone: tz})
		return &s
	case
		"drk":
		s := station(r.Station{Name: "Deutschlandradio Kultur", CloseDown: "00:00", ProgramURL: r.MustParseURL("http://www.deutschlandradiokultur.de/programmvorschau.282.de.html"), Identifier: identifier, TimeZone: tz})
		return &s
	}
	return nil
}

/// Stringer
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
		day, _ := s.dayURLForDate(t0)
		jobs = append(jobs, r.Scraper(*day))
	}
	return
}

///////////////////////////////////////////////////////////////////////
// http://www.deutschlandfunk.de/programmvorschau.281.de.html?drbm:date=19.11.2015

func (s *station) dayURLForDate(day time.Time) (ret *timeURL, err error) {
	r := timeURL(r.TimeURL{
		Time:    time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, s.TimeZone),
		Source:  *r.MustParseURL(s.ProgramURL.String() + day.Format("?drbm:date=02.01.2006")),
		Station: r.Station(*s),
	})
	ret = &r
	// err = errors.New("Not ümplemented yet.")
	return
}

/////////////////////////////////////////////////////////////////////////////
/// Just wrap TimeURL into a distinct, local type - a Scraper, naturally
type timeURL r.TimeURL

/// r.Scraper
func (day timeURL) Matches(nows []time.Time) (ok bool) {
	return true
}

// Scrape broadcasts from a day page.
func (day timeURL) Scrape() (jobs []r.Scraper, results []r.Broadcaster, err error) {
	bcs, err := day.parseBroadcastsFromURL()
	if nil == err {
		for _, bc := range bcs {
			results = append(results, bc)
		}
	}
	return
}

var (
	langDe string = "de"
)

func (day *timeURL) parseBroadcastsFromNode(root *html.Node) (ret []*r.Broadcast, err error) {
	// fmt.Fprintf(os.Stderr, "%s\n", day.Source.String())
	index := 0
	for _, at := range scrape.FindAll(root, func(n *html.Node) bool {
		return atom.A == n.DataAtom &&
			atom.Td == n.Parent.DataAtom &&
			atom.Tr == n.Parent.Parent.DataAtom &&
			"time" == scrape.Attr(n.Parent, "class")
	}) {
		// prepare response
		bc := r.Broadcast{
			BroadcastURL: r.BroadcastURL{
				TimeURL: r.TimeURL(*day),
			},
		}

		// some defaults
		bc.Language = &langDe
		{
			publisher := "http://www.deutschlandfunk.de/"
			if "drk" == day.Station.Identifier {
				publisher = "http://www.deutschlandradiokultur.de/"
			}
			bc.Publisher = &publisher
		}
		// set start time
		{
			aID := scrape.Attr(at, "name")
			if "" == aID {
				continue
			}
			bc.Source.Fragment = aID
			hour := r.MustParseInt(aID[0:2])
			minute := r.MustParseInt(aID[2:4])
			if 24 < hour || 60 < minute {
				continue
			}
			bc.Time = time.Date(day.Year(), day.Month(), day.Day(), hour, minute, 0, 0, day.TimeZone)
			if index > 0 {
				ret[index-1].DtEnd = &bc.Time
			}
		}
		// Title
		for idx, h3 := range scrape.FindAll(at.Parent.Parent, func(n *html.Node) bool {
			return atom.H3 == n.DataAtom &&
				atom.Td == n.Parent.DataAtom &&
				atom.Tr == n.Parent.Parent.DataAtom &&
				"description" == scrape.Attr(n.Parent, "class")
		}) {
			if idx != 0 {
				err = errors.New("There was more than 1 <tr><td class='description'><h3>")
				return
			}
			// purge 'aufnehmen' link:
			for _, chi := range scrape.FindAll(h3, func(n *html.Node) bool {
				return atom.A == n.DataAtom &&
					"psradio" == scrape.Attr(n, "class")
			}) {
				h3.RemoveChild(chi)
			}
			// fmt.Fprintf(os.Stderr, " '%s'\n", scrape.Text(h3))

			for idx, h3A := range scrape.FindAll(h3, func(n *html.Node) bool {
				return atom.A == n.DataAtom
			}) {
				if idx != 0 {
					err = errors.New("There was more than 1 <tr><td class='description'><h3><a>")
					return
				}
				bc.Title = scrape.Text(h3A)
				u, _ := url.Parse(scrape.Attr(h3A, "href"))
				bc.Subject = day.Source.ResolveReference(u)
			}
			bc.Title = strings.TrimSpace(bc.Title)
			if "" == bc.Title {
				bc.Title = r.TextChildrenNoClimb(h3)
			}
			// fmt.Fprintf(os.Stderr, " '%s'", bc.Title)
			{
				description := r.TextWithBrFromNodeSet(scrape.FindAll(h3.Parent, func(n *html.Node) bool { return atom.P == n.DataAtom }))
				bc.Description = &description
			}
		}
		// fmt.Fprintf(os.Stderr, "\n")
		ret = append(ret, &bc)
		index += 1
	}
	// fmt.Fprintf(os.Stderr, "len(ret) = %d '%s'\n", len(ret), day.Source.String())
	if index > 0 {
		midnight := time.Date(day.Year(), day.Month(), day.Day(), 24, 0, 0, 0, day.TimeZone)
		ret[index-1].DtEnd = &midnight
	}
	return
}

func (day *timeURL) parseBroadcastsFromReader(read io.Reader, cr0 *r.CountingReader) (ret []*r.Broadcast, err error) {
	cr := r.NewCountingReader(read)
	root, err := html.Parse(cr)
	r.ReportLoad("🐦", cr0, cr, day.Source)
	if nil != err {
		return
	}
	return day.parseBroadcastsFromNode(root)
}

func (day *timeURL) parseBroadcastsFromURL() (ret []*r.Broadcast, err error) {
	bo, cr, err := r.HttpGetBody(day.Source)
	if nil == bo {
		return nil, err
	}
	return day.parseBroadcastsFromReader(bo, cr)
}
