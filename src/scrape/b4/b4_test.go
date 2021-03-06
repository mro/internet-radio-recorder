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

// http://golang.org/pkg/testing/
// http://blog.stretchr.com/2014/03/05/test-driven-development-specifically-in-golang/
// https://xivilization.net/~marek/blog/2015/05/04/go-1-dot-4-2-for-raspberry-pi/
package b4 // import "purl.mro.name/recorder/radio/scrape/b4"

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	r "purl.mro.name/recorder/radio/scrape"
)

func TestTimeZone(t *testing.T) {
	b4 := Station("b4")
	assert.Equal(t, "Europe/Berlin", b4.TimeZone.String(), "foo")
}

func TestDayURLForDate(t *testing.T) {
	s := Station("b4")
	u, err := s.calendarItemRangeURLForTime(time.Date(2015, 11, 30, 5, 6, 7, 8, s.TimeZone))
	assert.Nil(t, err, "ouch: err")
	assert.Equal(t, "https://www.br-klassik.de/programm/radio/radiosendungen-100~calendarItems.jsp?rows=800&from=2015-11-30T05:07:07&to=2015-11-30T06:07:07", u.Source.String(), "ouch")
	assert.Equal(t, "2015-11-30T05:07:07+01:00", u.Time.Format(time.RFC3339), "ouch")
}

func TestUnmarshalBuiltMyTimeJSON(t *testing.T) {
	res := Time{}
	err := json.Unmarshal([]byte(`"2015-11-30T05:02:03"`), &res)
	assert.Nil(t, err, "ouch")
	assert.Equal(t, "2015-11-30T05:02:03+01:00", time.Time(res).Format(time.RFC3339), "ouch3")
}

func TestUnmarshalCalendarItemJSON(t *testing.T) {
	data := []byte(`{
	"datetime": "2015-11-30T05:00:00",
	"html": "\r\n    \r\n\r\n\r\n\r\n\r\n\r\n\r\n<li class=\"br-entry\" data-datetime=\"2015-11-30T05:00:00\">\r\n    \r\n    <ul>\r\n        <li class=\"br-time\">\r\n            <a class=\"br-toggle\">05:00</a>\r\n        </li>\r\n        <li class=\"br-content\">\r\n            <a class=\"br-toggle\">\r\n                \r\n                    \r\n                    \r\n                        <p class=\"br-type\">radio</p>\r\n\r\n                        <p class=\"br-title\">Nachrichten, Wetter</p>\r\n\r\n                        <p class=\"br-text\"></p>\r\n                    \r\n                \r\n            </a>\r\n\r\n            <div class=\"br-detail\">\r\n                <a href=\"/programm/radio/ausstrahlung-512526.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    \r\n                        <img alt=\"Sendungsbild: Nachrichten, Wetter, Verkehr | Bild: BR\" title=\"Sendungsbild: Nachrichten, Wetter, Verkehr | Bild: BR\" src=\"/programm/radio/sendungsbild-nachrichten-wetter-verkehr100~_h-364_v-img__16__9__xl_w-648_-be6819cc57a5436fe2e22755fd9495d5c6ac08f6.jpg?version=50e7f\"/>\r\n                    \r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512526.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    \r\n                        \r\n                        \r\n                            \r\n                            \r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n                            <p>\r\n\r\n\r\n\r\n\r\n\r\n\r\n</p>\r\n                                <p>\r\n\r\n\r\n\r\n\r\n\r\n\r\n</p>\r\n                        \r\n                    \r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512526.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    <span class=\"br-more\">Mehr<span class=\"br-sprite br-sprite-arrow-link\"></span></span>\r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512526~exportICS.ics\" class=\"br-ics-download br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    <span class=\"br-more\">Zum Kalender hinzuf\u00FCgen<span class=\"br-sprite br-sprite-arrow-link\"></span></span>\r\n                </a>\r\n\r\n                \r\n                    \r\n                    \r\n                \r\n            </div>\r\n        </li>\r\n    </ul>\r\n</li>\r\n"
}`)
	res := calendarItem{}
	err := json.Unmarshal(data, &res)
	assert.Nil(t, err, "ouch")
	assert.Equal(t, "2015-11-30T05:00:00+01:00", time.Time(res.DateTime).Format(time.RFC3339), "ouch3")
	assert.Equal(t, "\r\n    \r\n\r\n\r\n\r\n\r\n\r\n\r\n<li class=\"br-entry\" data-datetime=\"2015-11-30T05:00:00\">\r\n    \r\n    <ul>\r\n        <li class=\"br-time\">\r\n            <a class=\"br-toggle\">05:00</a>\r\n        </li>\r\n        <li class=\"br-content\">\r\n            <a class=\"br-toggle\">\r\n                \r\n                    \r\n                    \r\n                        <p class=\"br-type\">radio</p>\r\n\r\n                        <p class=\"br-title\">Nachrichten, Wetter</p>\r\n\r\n                        <p class=\"br-text\"></p>\r\n                    \r\n                \r\n            </a>\r\n\r\n            <div class=\"br-detail\">\r\n                <a href=\"/programm/radio/ausstrahlung-512526.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    \r\n                        <img alt=\"Sendungsbild: Nachrichten, Wetter, Verkehr | Bild: BR\" title=\"Sendungsbild: Nachrichten, Wetter, Verkehr | Bild: BR\" src=\"/programm/radio/sendungsbild-nachrichten-wetter-verkehr100~_h-364_v-img__16__9__xl_w-648_-be6819cc57a5436fe2e22755fd9495d5c6ac08f6.jpg?version=50e7f\"/>\r\n                    \r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512526.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    \r\n                        \r\n                        \r\n                            \r\n                            \r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n                            <p>\r\n\r\n\r\n\r\n\r\n\r\n\r\n</p>\r\n                                <p>\r\n\r\n\r\n\r\n\r\n\r\n\r\n</p>\r\n                        \r\n                    \r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512526.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    <span class=\"br-more\">Mehr<span class=\"br-sprite br-sprite-arrow-link\"></span></span>\r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512526~exportICS.ics\" class=\"br-ics-download br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    <span class=\"br-more\">Zum Kalender hinzuf\u00FCgen<span class=\"br-sprite br-sprite-arrow-link\"></span></span>\r\n                </a>\r\n\r\n                \r\n                    \r\n                    \r\n                \r\n            </div>\r\n        </li>\r\n    </ul>\r\n</li>\r\n", res.Html, "ouch3")
}

func TestUnmarshalCalendarItemsJSON(t *testing.T) {
	data := []byte(`[{
	"datetime": "2015-11-30T05:00:00",
	"html": "\r\n    \r\n\r\n\r\n\r\n\r\n\r\n\r\n<li class=\"br-entry\" data-datetime=\"2015-11-30T05:00:00\">\r\n    \r\n    <ul>\r\n        <li class=\"br-time\">\r\n            <a class=\"br-toggle\">05:00</a>\r\n        </li>\r\n        <li class=\"br-content\">\r\n            <a class=\"br-toggle\">\r\n                \r\n                    \r\n                    \r\n                        <p class=\"br-type\">radio</p>\r\n\r\n                        <p class=\"br-title\">Nachrichten, Wetter</p>\r\n\r\n                        <p class=\"br-text\"></p>\r\n                    \r\n                \r\n            </a>\r\n\r\n            <div class=\"br-detail\">\r\n                <a href=\"/programm/radio/ausstrahlung-512526.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    \r\n                        <img alt=\"Sendungsbild: Nachrichten, Wetter, Verkehr | Bild: BR\" title=\"Sendungsbild: Nachrichten, Wetter, Verkehr | Bild: BR\" src=\"/programm/radio/sendungsbild-nachrichten-wetter-verkehr100~_h-364_v-img__16__9__xl_w-648_-be6819cc57a5436fe2e22755fd9495d5c6ac08f6.jpg?version=50e7f\"/>\r\n                    \r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512526.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    \r\n                        \r\n                        \r\n                            \r\n                            \r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n                            <p>\r\n\r\n\r\n\r\n\r\n\r\n\r\n</p>\r\n                                <p>\r\n\r\n\r\n\r\n\r\n\r\n\r\n</p>\r\n                        \r\n                    \r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512526.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    <span class=\"br-more\">Mehr<span class=\"br-sprite br-sprite-arrow-link\"></span></span>\r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512526~exportICS.ics\" class=\"br-ics-download br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    <span class=\"br-more\">Zum Kalender hinzuf\u00FCgen<span class=\"br-sprite br-sprite-arrow-link\"></span></span>\r\n                </a>\r\n\r\n                \r\n                    \r\n                    \r\n                \r\n            </div>\r\n        </li>\r\n    </ul>\r\n</li>\r\n"
}, {
	"datetime": "2015-11-30T05:03:00",
	"html": "\r\n    \r\n\r\n\r\n\r\n\r\n\r\n\r\n<li class=\"br-entry\" data-datetime=\"2015-11-30T05:03:00\">\r\n    \r\n    <ul>\r\n        <li class=\"br-time\">\r\n            <a class=\"br-toggle\">05:03</a>\r\n        </li>\r\n        <li class=\"br-content\">\r\n            <a class=\"br-toggle\">\r\n                \r\n                    \r\n                    \r\n                        <p class=\"br-type\">radio</p>\r\n\r\n                        <p class=\"br-title\">Das ARD-Nachtkonzert (IV)</p>\r\n\r\n                        <p class=\"br-text\"></p>\r\n                    \r\n                \r\n            </a>\r\n\r\n            <div class=\"br-detail\">\r\n                <a href=\"/programm/radio/ausstrahlung-512528.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    \r\n                        <img alt=\"Dirigentenh\u00E4nde | Bild: Digital Vision\" title=\"Dirigentenh\u00E4nde | Bild: Digital Vision\" src=\"/programm/radio/dirigentenhaende102~_h-364_v-img__16__9__xl_w-648_-be6819cc57a5436fe2e22755fd9495d5c6ac08f6.jpg?version=f1b6f\"/>\r\n                    \r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512528.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    \r\n                        \r\n                        \r\n                            \r\n                            \r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n                            <p>\r\n\r\n\r\n\r\n\r\n\r\n\r\n    \r\n    \r\n        \r\n            \r\n            \r\n                Fr\u00E9d\u00E9ric Chopin: Polonaise A-Dur, op. 40, Nr. 1 (Rafal Blechacz, Klavier); Arcangelo Corelli: Concerto grosso B-Dur, op. 6, Nr. 11 (Alba Roca, Violine; Gli Incogniti, Violine und Leitung: Amandine Beyer); Arnold Bax: Sonate D-Dur (Michael Collins); Joseph Haydn: Symphonie Nr. 27 G-Dur (Heidelberger Sinfoniker: Thomas Fey); Ole Bull: \"Vision im Gebirge\" (Arve Tellefsen, Violine; Trondheim Symphony Orchestra: Eivind Aadland); Con Conrad: \"Singin' the Blues\" (George Gershwin, Klavier)\r\n            \r\n\r\n            \r\n            \r\n\r\n        \r\n    \r\n\r\n</p>\r\n                                <p>\r\n\r\n\r\n\r\n\r\n\r\n\r\n</p>\r\n                        \r\n                    \r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512528.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    <span class=\"br-more\">Mehr<span class=\"br-sprite br-sprite-arrow-link\"></span></span>\r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512528~exportICS.ics\" class=\"br-ics-download br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    <span class=\"br-more\">Zum Kalender hinzuf\u00FCgen<span class=\"br-sprite br-sprite-arrow-link\"></span></span>\r\n                </a>\r\n\r\n                \r\n                    \r\n                    \r\n                \r\n            </div>\r\n        </li>\r\n    </ul>\r\n</li>\r\n"
}]`)
	res := make([]calendarItem, 0)
	err := json.Unmarshal(data, &res)
	assert.Nil(t, err, "ouch")
	assert.Equal(t, 2, len(res), "ouch2")
	assert.Equal(t, "2015-11-30T05:00:00+01:00", time.Time(res[0].DateTime).Format(time.RFC3339), "ouch3")
	assert.Equal(t, "\r\n    \r\n\r\n\r\n\r\n\r\n\r\n\r\n<li class=\"br-entry\" data-datetime=\"2015-11-30T05:00:00\">\r\n    \r\n    <ul>\r\n        <li class=\"br-time\">\r\n            <a class=\"br-toggle\">05:00</a>\r\n        </li>\r\n        <li class=\"br-content\">\r\n            <a class=\"br-toggle\">\r\n                \r\n                    \r\n                    \r\n                        <p class=\"br-type\">radio</p>\r\n\r\n                        <p class=\"br-title\">Nachrichten, Wetter</p>\r\n\r\n                        <p class=\"br-text\"></p>\r\n                    \r\n                \r\n            </a>\r\n\r\n            <div class=\"br-detail\">\r\n                <a href=\"/programm/radio/ausstrahlung-512526.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    \r\n                        <img alt=\"Sendungsbild: Nachrichten, Wetter, Verkehr | Bild: BR\" title=\"Sendungsbild: Nachrichten, Wetter, Verkehr | Bild: BR\" src=\"/programm/radio/sendungsbild-nachrichten-wetter-verkehr100~_h-364_v-img__16__9__xl_w-648_-be6819cc57a5436fe2e22755fd9495d5c6ac08f6.jpg?version=50e7f\"/>\r\n                    \r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512526.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    \r\n                        \r\n                        \r\n                            \r\n                            \r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n                            <p>\r\n\r\n\r\n\r\n\r\n\r\n\r\n</p>\r\n                                <p>\r\n\r\n\r\n\r\n\r\n\r\n\r\n</p>\r\n                        \r\n                    \r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512526.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    <span class=\"br-more\">Mehr<span class=\"br-sprite br-sprite-arrow-link\"></span></span>\r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512526~exportICS.ics\" class=\"br-ics-download br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    <span class=\"br-more\">Zum Kalender hinzuf\u00FCgen<span class=\"br-sprite br-sprite-arrow-link\"></span></span>\r\n                </a>\r\n\r\n                \r\n                    \r\n                    \r\n                \r\n            </div>\r\n        </li>\r\n    </ul>\r\n</li>\r\n", res[0].Html, "ouch3")
	assert.Equal(t, "2015-11-30T05:03:00+01:00", time.Time(res[1].DateTime).Format(time.RFC3339), "ouch3")
	assert.Equal(t, "\r\n    \r\n\r\n\r\n\r\n\r\n\r\n\r\n<li class=\"br-entry\" data-datetime=\"2015-11-30T05:03:00\">\r\n    \r\n    <ul>\r\n        <li class=\"br-time\">\r\n            <a class=\"br-toggle\">05:03</a>\r\n        </li>\r\n        <li class=\"br-content\">\r\n            <a class=\"br-toggle\">\r\n                \r\n                    \r\n                    \r\n                        <p class=\"br-type\">radio</p>\r\n\r\n                        <p class=\"br-title\">Das ARD-Nachtkonzert (IV)</p>\r\n\r\n                        <p class=\"br-text\"></p>\r\n                    \r\n                \r\n            </a>\r\n\r\n            <div class=\"br-detail\">\r\n                <a href=\"/programm/radio/ausstrahlung-512528.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    \r\n                        <img alt=\"Dirigentenh\u00E4nde | Bild: Digital Vision\" title=\"Dirigentenh\u00E4nde | Bild: Digital Vision\" src=\"/programm/radio/dirigentenhaende102~_h-364_v-img__16__9__xl_w-648_-be6819cc57a5436fe2e22755fd9495d5c6ac08f6.jpg?version=f1b6f\"/>\r\n                    \r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512528.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    \r\n                        \r\n                        \r\n                            \r\n                            \r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n\r\n                            <p>\r\n\r\n\r\n\r\n\r\n\r\n\r\n    \r\n    \r\n        \r\n            \r\n            \r\n                Fr\u00E9d\u00E9ric Chopin: Polonaise A-Dur, op. 40, Nr. 1 (Rafal Blechacz, Klavier); Arcangelo Corelli: Concerto grosso B-Dur, op. 6, Nr. 11 (Alba Roca, Violine; Gli Incogniti, Violine und Leitung: Amandine Beyer); Arnold Bax: Sonate D-Dur (Michael Collins); Joseph Haydn: Symphonie Nr. 27 G-Dur (Heidelberger Sinfoniker: Thomas Fey); Ole Bull: \"Vision im Gebirge\" (Arve Tellefsen, Violine; Trondheim Symphony Orchestra: Eivind Aadland); Con Conrad: \"Singin' the Blues\" (George Gershwin, Klavier)\r\n            \r\n\r\n            \r\n            \r\n\r\n        \r\n    \r\n\r\n</p>\r\n                                <p>\r\n\r\n\r\n\r\n\r\n\r\n\r\n</p>\r\n                        \r\n                    \r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512528.html\" class=\"br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    <span class=\"br-more\">Mehr<span class=\"br-sprite br-sprite-arrow-link\"></span></span>\r\n                </a>\r\n                <a href=\"/programm/radio/ausstrahlung-512528~exportICS.ics\" class=\"br-ics-download br-internal\" title=\"BroadcastScheduleSlot\">\r\n                    <span class=\"br-more\">Zum Kalender hinzuf\u00FCgen<span class=\"br-sprite br-sprite-arrow-link\"></span></span>\r\n                </a>\r\n\r\n                \r\n                    \r\n                    \r\n                \r\n            </div>\r\n        </li>\r\n    </ul>\r\n</li>\r\n", res[1].Html, "ouch3")
}

func TestParseCalendarItems(t *testing.T) {
	s := Station("b4")
	u, err := s.calendarItemRangeURLForTime(time.Date(2015, 11, 30, 5, 0, 0, 0, localLoc))
	assert.Equal(t, "https://www.br-klassik.de/programm/radio/radiosendungen-100~calendarItems.jsp?rows=800&from=2015-11-30T05:01:00&to=2015-11-30T06:01:00", u.Source.String(), "Nov")
	assert.Equal(t, "2015-11-30T05:01:00+01:00", u.Time.Format(time.RFC3339), "Nov")
	assert.Equal(t, "b4", u.Station.Identifier, "Nov")
	assert.NotNil(t, u.Station.ProgramURL, "Nov")

	f, err := os.Open("testdata/2015-11-30T05-b4-program.json")
	assert.NotNil(t, f, "ouch")
	assert.Nil(t, err, "ouch")

	cis, _ := u.parseCalendarItemsReader(f, nil)
	assert.Equal(t, 2, len(cis), "Nov")
	{
		item := cis[0]
		assert.NotNil(t, item.Station.ProgramURL, "Nov")
		bc, err := item.parseBroadcastSeedString(&item.Html)
		assert.Nil(t, err, "ouch")
		assert.Equal(t, "2015-11-30T05:00:00+01:00", bc.Time.Format(time.RFC3339), "Nov")
		assert.Equal(t, "https://www.br-klassik.de/programm/radio/ausstrahlung-512526.html", bc.Source.String(), "Nov")
		assert.Equal(t, "https://www.br-klassik.de/programm/radio/sendungsbild-nachrichten-wetter-verkehr100~_h-364_v-img__16__9__xl_w-648_-be6819cc57a5436fe2e22755fd9495d5c6ac08f6.jpg?version=50e7f", bc.Image.String(), "Nov")
	}
	{
		item := cis[1]
		assert.NotNil(t, item.Station.ProgramURL, "Nov")
		bc, err := item.parseBroadcastSeedString(&item.Html)
		assert.Nil(t, err, "ouch")
		assert.Equal(t, "2015-11-30T05:03:00+01:00", bc.Time.Format(time.RFC3339), "Nov")
		assert.Equal(t, "https://www.br-klassik.de/programm/radio/ausstrahlung-512528.html", bc.Source.String(), "Nov")
		assert.Equal(t, "https://www.br-klassik.de/programm/radio/dirigentenhaende102~_h-364_v-img__16__9__xl_w-648_-be6819cc57a5436fe2e22755fd9495d5c6ac08f6.jpg?version=f1b6f", bc.Image.String(), "Nov")
	}
}

func TestParseBroadcast_914548(t *testing.T) {
	{
		t0, _ := time.Parse(time.RFC3339, "2016-11-27T20:30:00+01:00")
		assert.Equal(t, "2016-11-27T20:30:00+01:00", t0.Format(time.RFC3339), "oha")
	}
	f, err := os.Open("testdata/2016-11-27T2030-b4-ausstrahlung-914548.html")
	assert.NotNil(t, f, "ouch")
	assert.Nil(t, err, "ouch")

	s := Station("b4")

	t0 := broadcastURL{
		BroadcastURL: r.BroadcastURL{
			TimeURL: r.TimeURL{
				Time:    time.Date(2016, time.November, 27, 20, 30, 0, 0, localLoc),
				Source:  *r.MustParseURL("https://www.br-klassik.de/programm/radio/ausstrahlung-914548.html"),
				Station: r.Station(*s),
			},
			Title: "Intermezzo",
		},
		Image: r.MustParseURL("https://www.br-klassik.de/programm/radio/logo106~_h-558_v-img__16__9__xl_w-994_-e1d284d92729d9396a907e303225e0f2d9fa53b4.jpg?version=78f3c"),
	}
	bc, err := t0.parseBroadcastReader(f, nil)
	assert.Nil(t, err, "ouch")
	assert.Equal(t, "b4", bc.Station.Identifier, "ouch: Station.Identifier")
	assert.Equal(t, "Intermezzo", bc.Title, "ouch: Title")
	assert.Equal(t, "https://www.br-klassik.de/programm/radio/ausstrahlung-914548.html", bc.Source.String(), "ouch: Source")
	assert.NotNil(t, bc.Language, "ouch: Language")
	assert.Equal(t, "de", *bc.Language, "ouch: Language")
	assert.Equal(t, t0.Title, bc.Title, "ouch: Title")
	assert.Equal(t, "", *bc.TitleSeries, "ouch: TitleSeries")
	assert.Nil(t, bc.TitleEpisode, "ouch: TitleEpisode")
	assert.Equal(t, "2016-11-27T20:30:00+01:00", bc.Time.Format(time.RFC3339), "ouch: Time")
	assert.Equal(t, "2016-11-27T21:00:00+01:00", bc.DtEnd.Format(time.RFC3339), "ouch: DtEnd")
	assert.Nil(t, bc.Subject, "ouch: Subject")
	assert.Equal(t, "2016-11-25T13:15:07+01:00", bc.Modified.Format(time.RFC3339), "ouch: Modified")
	assert.Equal(t, "Bayerischer Rundfunk", *bc.Author, "ouch: Author")
	assert.NotNil(t, bc.Description, "ouch: Description")
	assert.Equal(t, "Ernst von Gemmingen: Violinkonzert Nr. 1 A-Dur (Kolja Lessing, Violine; Münchner Rundfunkorchester: Ulf Schirmer)", *bc.Description, "ouch: Description")
	assert.NotNil(t, bc.Image, "ouch: Image")
	assert.Equal(t, "https://www.br-klassik.de/programm/radio/logo106~_h-558_v-img__16__9__xl_w-994_-e1d284d92729d9396a907e303225e0f2d9fa53b4.jpg?version=78f3c", bc.Image.String(), "ouch: Image")
}

func TestParseBroadcast_866264(t *testing.T) {
	{
		t0, _ := time.Parse(time.RFC3339, "2016-11-27T23:05:00+01:00")
		assert.Equal(t, "2016-11-27T23:05:00+01:00", t0.Format(time.RFC3339), "oha")
	}
	f, err := os.Open("testdata/2016-11-27T2305-b4-ausstrahlung-866264.html")
	assert.NotNil(t, f, "ouch")
	assert.Nil(t, err, "ouch")

	s := Station("b4")

	t0 := broadcastURL{
		BroadcastURL: r.BroadcastURL{
			TimeURL: r.TimeURL{
				Time:    time.Date(2016, time.November, 27, 23, 5, 0, 0, localLoc),
				Source:  *r.MustParseURL("https://www.br-klassik.de/programm/radio/ausstrahlung-866264.html"),
				Station: r.Station(*s),
			},
			Title: "Musik der Welt",
		},
		Image: r.MustParseURL("https://www.br-klassik.de/pierre-dominique-ponnelle-102~_h-558_v-img__16__9__xl_w-994_-e1d284d92729d9396a907e303225e0f2d9fa53b4.jpg?version=1c52e"),
	}
	bc, err := t0.parseBroadcastReader(f, nil)
	assert.Nil(t, err, "ouch")
	assert.Equal(t, "b4", bc.Station.Identifier, "ouch: Station.Identifier")
	assert.Equal(t, "Musik der Welt", bc.Title, "ouch: Title")
	assert.Equal(t, "https://www.br-klassik.de/programm/radio/ausstrahlung-866264.html", bc.Source.String(), "ouch: Source")
	assert.NotNil(t, bc.Language, "ouch: Language")
	assert.Equal(t, "de", *bc.Language, "ouch: Language")
	assert.Equal(t, t0.Title, bc.Title, "ouch: Title")
	assert.Equal(t, "Aus dem Studio Franken:", *bc.TitleSeries, "ouch: TitleSeries")
	assert.Equal(t, "Minsk, Ufa, Taschkent", *bc.TitleEpisode, "ouch: TitleEpisode")
	assert.Equal(t, "2016-11-27T23:05:00+01:00", bc.Time.Format(time.RFC3339), "ouch: Time")
	assert.Equal(t, "2016-11-28T00:00:00+01:00", bc.DtEnd.Format(time.RFC3339), "ouch: DtEnd")
	assert.Equal(t, "https://www.br-klassik.de/themen/jazz-und-weltmusik/musik-der-welt-102.html", bc.Subject.String(), "ouch: Subject")
	assert.Equal(t, "2016-11-24T09:15:20+01:00", bc.Modified.Format(time.RFC3339), "ouch: Modified")
	assert.Equal(t, "Bayerischer Rundfunk", *bc.Author, "ouch: Author")
	assert.NotNil(t, bc.Description, "ouch: Description")
	assert.Equal(t, "Fremde Länder, fremde Klänge in Features, Reportagen und Interviews. Ein internationales Autoren-Team präsentiert Musik rund um den Globus, vom Porträt bis zum Konzertmitschnitt.\n\nDie Stationen seiner Karriere haben einen exotischen Klang. Er war Chefdirigent in Minsk. Er kennt die Opernhäuser von Bischkek, Almaty oder Baku wie seine Westentasche. Er hat die Staatsorchester von Armenien, Aserbaidschan und Georgien geleitet. Dabei hat es Pierre-Dominique Ponnelle, Sohn des berühmten Opernregisseurs Jean-Pierre Ponnelle, nach Lehrjahren bei Herbert von Karajan und Otmar Suitner eher durch Zufall gen Osten verschlagen. Grau in Grau hatte er sich diese Welt dort vorgestellt - und war dann fasziniert, gepackt von der Buntheit, der Vielfalt, der musikalischen Qualität, die er in den Ländern der ehemaligen Sowjetunion vorfand. Seine Reisen zu fremden Ländern und Menschen haben inzwischen auch in einer Reihe von Kompositionen ihren Niederschlag gefunden: Da gibt es eine liebevolle Hommage an den armenischen Maler Minas Avetisyan oder ein Streichquartett, das zentralasiatische Melodien verarbeitet. Thorsten Preuß sprach mit dem weitgereisten Komponisten und Dirigenten.\n\nDie musikalischen Reisen des Komponisten und Dirigenten Pierre-Dominique Ponnelle\nVon Thorsten Preuß", *bc.Description, "ouch: Description")
	assert.NotNil(t, bc.Image, "ouch: Image")
	assert.Equal(t, "https://www.br-klassik.de/pierre-dominique-ponnelle-102~_h-558_v-img__16__9__xl_w-994_-e1d284d92729d9396a907e303225e0f2d9fa53b4.jpg?version=1c52e", bc.Image.String(), "ouch: Image")
}
