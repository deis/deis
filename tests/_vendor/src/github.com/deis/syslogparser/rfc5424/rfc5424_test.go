package rfc5424

import (
	"fmt"
	"github.com/jeromer/syslogparser"
	. "launchpad.net/gocheck"
	"testing"
	"time"
)

// Hooks up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type Rfc5424TestSuite struct {
}

var _ = Suite(&Rfc5424TestSuite{})

func (s *Rfc5424TestSuite) TestParser_Valid(c *C) {
	fixtures := []string{
		// no STRUCTURED-DATA
		"<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su - ID47 - 'su root' failed for lonvick on /dev/pts/8",
		"<165>1 2003-08-24T05:14:15.000003-07:00 192.0.2.1 myproc 8710 - - %% It's time to make the do-nuts.",
		// with STRUCTURED-DATA
		`<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] An application event log entry...`,

		// STRUCTURED-DATA Only
		`<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3" eventSource= "Application" eventID="1011"][examplePriority@32473 class="high"]`,
	}

	tmpTs, err := time.Parse("-07:00", "-07:00")
	c.Assert(err, IsNil)

	expected := []syslogparser.LogParts{
		syslogparser.LogParts{
			"priority":        34,
			"facility":        4,
			"severity":        2,
			"version":         1,
			"timestamp":       time.Date(2003, time.October, 11, 22, 14, 15, 3*10e5, time.UTC),
			"hostname":        "mymachine.example.com",
			"app_name":        "su",
			"proc_id":         "-",
			"msg_id":          "ID47",
			"structured_data": "-",
			"message":         "'su root' failed for lonvick on /dev/pts/8",
		},
		syslogparser.LogParts{
			"priority":        165,
			"facility":        20,
			"severity":        5,
			"version":         1,
			"timestamp":       time.Date(2003, time.August, 24, 5, 14, 15, 3*10e2, tmpTs.Location()),
			"hostname":        "192.0.2.1",
			"app_name":        "myproc",
			"proc_id":         "8710",
			"msg_id":          "-",
			"structured_data": "-",
			"message":         "%% It's time to make the do-nuts.",
		},
		syslogparser.LogParts{
			"priority":        165,
			"facility":        20,
			"severity":        5,
			"version":         1,
			"timestamp":       time.Date(2003, time.October, 11, 22, 14, 15, 3*10e5, time.UTC),
			"hostname":        "mymachine.example.com",
			"app_name":        "evntslog",
			"proc_id":         "-",
			"msg_id":          "ID47",
			"structured_data": `[exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"]`,
			"message":         "An application event log entry...",
		},
		syslogparser.LogParts{
			"priority":        165,
			"facility":        20,
			"severity":        5,
			"version":         1,
			"timestamp":       time.Date(2003, time.October, 11, 22, 14, 15, 3*10e5, time.UTC),
			"hostname":        "mymachine.example.com",
			"app_name":        "evntslog",
			"proc_id":         "-",
			"msg_id":          "ID47",
			"structured_data": `[exampleSDID@32473 iut="3" eventSource= "Application" eventID="1011"][examplePriority@32473 class="high"]`,
			"message":         "",
		},
	}

	c.Assert(len(fixtures), Equals, len(expected))
	start := 0
	for i, buff := range fixtures {
		expectedP := &Parser{
			buff:   []byte(buff),
			cursor: start,
			l:      len(buff),
		}

		p := NewParser([]byte(buff))
		c.Assert(p, DeepEquals, expectedP)

		err := p.Parse()
		c.Assert(err, IsNil)

		obtained := p.Dump()
		for k, v := range obtained {
			c.Assert(v, DeepEquals, expected[i][k])
		}
	}
}

func (s *Rfc5424TestSuite) TestParseHeader_Valid(c *C) {
	ts := time.Date(2003, time.October, 11, 22, 14, 15, 3*10e5, time.UTC)
	tsString := "2003-10-11T22:14:15.003Z"
	hostname := "mymachine.example.com"
	appName := "su"
	procId := "123"
	msgId := "ID47"
	nilValue := string(NILVALUE)
	headerFmt := "<165>1 %s %s %s %s %s "

	fixtures := []string{
		// HEADER complete
		fmt.Sprintf(headerFmt, tsString, hostname, appName, procId, msgId),
		// TIMESTAMP as NILVALUE
		fmt.Sprintf(headerFmt, nilValue, hostname, appName, procId, msgId),
		// HOSTNAME as NILVALUE
		fmt.Sprintf(headerFmt, tsString, nilValue, appName, procId, msgId),
		// APP-NAME as NILVALUE
		fmt.Sprintf(headerFmt, tsString, hostname, nilValue, procId, msgId),
		// PROCID as NILVALUE
		fmt.Sprintf(headerFmt, tsString, hostname, appName, nilValue, msgId),
		// MSGID as NILVALUE
		fmt.Sprintf(headerFmt, tsString, hostname, appName, procId, nilValue),
	}

	pri := syslogparser.Priority{
		P: 165,
		F: syslogparser.Facility{Value: 20},
		S: syslogparser.Severity{Value: 5},
	}

	expected := []header{
		// HEADER complete
		header{
			priority:  pri,
			version:   1,
			timestamp: ts,
			hostname:  hostname,
			appName:   appName,
			procId:    procId,
			msgId:     msgId,
		},
		// TIMESTAMP as NILVALUE
		header{
			priority:  pri,
			version:   1,
			timestamp: *new(time.Time),
			hostname:  hostname,
			appName:   appName,
			procId:    procId,
			msgId:     msgId,
		},
		// HOSTNAME as NILVALUE
		header{
			priority:  pri,
			version:   1,
			timestamp: ts,
			hostname:  nilValue,
			appName:   appName,
			procId:    procId,
			msgId:     msgId,
		},
		// APP-NAME as NILVALUE
		header{
			priority:  pri,
			version:   1,
			timestamp: ts,
			hostname:  hostname,
			appName:   nilValue,
			procId:    procId,
			msgId:     msgId,
		},
		// PROCID as NILVALUE
		header{
			priority:  pri,
			version:   1,
			timestamp: ts,
			hostname:  hostname,
			appName:   appName,
			procId:    nilValue,
			msgId:     msgId,
		},
		// MSGID as NILVALUE
		header{
			priority:  pri,
			version:   1,
			timestamp: ts,
			hostname:  hostname,
			appName:   appName,
			procId:    procId,
			msgId:     nilValue,
		},
	}

	for i, f := range fixtures {
		p := NewParser([]byte(f))
		obtained, err := p.parseHeader()
		c.Assert(err, IsNil)
		c.Assert(obtained, Equals, expected[i])
		c.Assert(p.cursor, Equals, len(f))
	}
}

func (s *Rfc5424TestSuite) TestParseTimestamp_UTC(c *C) {
	buff := []byte("1985-04-12T23:20:50.52Z")
	ts := time.Date(1985, time.April, 12, 23, 20, 50, 52*10e6, time.UTC)

	s.assertTimestamp(c, ts, buff, 23, nil)
}

func (s *Rfc5424TestSuite) TestParseTimestamp_NumericTimezone(c *C) {
	tz := "-04:00"
	buff := []byte("1985-04-12T19:20:50.52" + tz)

	tmpTs, err := time.Parse("-07:00", tz)
	c.Assert(err, IsNil)

	ts := time.Date(1985, time.April, 12, 19, 20, 50, 52*10e6, tmpTs.Location())

	s.assertTimestamp(c, ts, buff, len(buff), nil)
}

func (s *Rfc5424TestSuite) TestParseTimestamp_MilliSeconds(c *C) {
	buff := []byte("2003-10-11T22:14:15.003Z")

	ts := time.Date(2003, time.October, 11, 22, 14, 15, 3*10e5, time.UTC)

	s.assertTimestamp(c, ts, buff, len(buff), nil)
}

func (s *Rfc5424TestSuite) TestParseTimestamp_MicroSeconds(c *C) {
	tz := "-07:00"
	buff := []byte("2003-08-24T05:14:15.000003" + tz)

	tmpTs, err := time.Parse("-07:00", tz)
	c.Assert(err, IsNil)

	ts := time.Date(2003, time.August, 24, 5, 14, 15, 3*10e2, tmpTs.Location())

	s.assertTimestamp(c, ts, buff, len(buff), nil)
}

func (s *Rfc5424TestSuite) TestParseTimestamp_NanoSeconds(c *C) {
	buff := []byte("2003-08-24T05:14:15.000000003-07:00")
	ts := new(time.Time)

	s.assertTimestamp(c, *ts, buff, 26, syslogparser.ErrTimestampUnknownFormat)
}

func (s *Rfc5424TestSuite) TestParseTimestamp_NilValue(c *C) {
	buff := []byte("-")
	ts := new(time.Time)

	s.assertTimestamp(c, *ts, buff, 1, nil)
}

func (s *Rfc5424TestSuite) TestFindNextSpace_NoSpace(c *C) {
	buff := []byte("aaaaaa")

	s.assertFindNextSpace(c, 0, buff, syslogparser.ErrNoSpace)
}

func (s *Rfc5424TestSuite) TestFindNextSpace_SpaceFound(c *C) {
	buff := []byte("foo bar baz")

	s.assertFindNextSpace(c, 4, buff, nil)
}

func (s *Rfc5424TestSuite) TestParseYear_Invalid(c *C) {
	buff := []byte("1a2b")
	expected := 0

	s.assertParseYear(c, expected, buff, 4, ErrYearInvalid)
}

func (s *Rfc5424TestSuite) TestParseYear_TooShort(c *C) {
	buff := []byte("123")
	expected := 0

	s.assertParseYear(c, expected, buff, 0, syslogparser.ErrEOL)
}

func (s *Rfc5424TestSuite) TestParseYear_Valid(c *C) {
	buff := []byte("2013")
	expected := 2013

	s.assertParseYear(c, expected, buff, 4, nil)
}

func (s *Rfc5424TestSuite) TestParseMonth_InvalidString(c *C) {
	buff := []byte("ab")
	expected := 0

	s.assertParseMonth(c, expected, buff, 2, ErrMonthInvalid)
}

func (s *Rfc5424TestSuite) TestParseMonth_InvalidRange(c *C) {
	buff := []byte("00")
	expected := 0

	s.assertParseMonth(c, expected, buff, 2, ErrMonthInvalid)

	// ----

	buff = []byte("13")

	s.assertParseMonth(c, expected, buff, 2, ErrMonthInvalid)
}

func (s *Rfc5424TestSuite) TestParseMonth_TooShort(c *C) {
	buff := []byte("1")
	expected := 0

	s.assertParseMonth(c, expected, buff, 0, syslogparser.ErrEOL)
}

func (s *Rfc5424TestSuite) TestParseMonth_Valid(c *C) {
	buff := []byte("02")
	expected := 2

	s.assertParseMonth(c, expected, buff, 2, nil)
}

func (s *Rfc5424TestSuite) TestParseDay_InvalidString(c *C) {
	buff := []byte("ab")
	expected := 0

	s.assertParseDay(c, expected, buff, 2, ErrDayInvalid)
}

func (s *Rfc5424TestSuite) TestParseDay_TooShort(c *C) {
	buff := []byte("1")
	expected := 0

	s.assertParseDay(c, expected, buff, 0, syslogparser.ErrEOL)
}

func (s *Rfc5424TestSuite) TestParseDay_InvalidRange(c *C) {
	buff := []byte("00")
	expected := 0

	s.assertParseDay(c, expected, buff, 2, ErrDayInvalid)

	// ----

	buff = []byte("32")

	s.assertParseDay(c, expected, buff, 2, ErrDayInvalid)
}

func (s *Rfc5424TestSuite) TestParseDay_Valid(c *C) {
	buff := []byte("02")
	expected := 2

	s.assertParseDay(c, expected, buff, 2, nil)
}

func (s *Rfc5424TestSuite) TestParseFullDate_Invalid(c *C) {
	buff := []byte("2013+10-28")
	fd := fullDate{}

	s.assertParseFullDate(c, fd, buff, 4, syslogparser.ErrTimestampUnknownFormat)

	// ---

	buff = []byte("2013-10+28")
	s.assertParseFullDate(c, fd, buff, 7, syslogparser.ErrTimestampUnknownFormat)
}

func (s *Rfc5424TestSuite) TestParseFullDate_Valid(c *C) {
	buff := []byte("2013-10-28")
	fd := fullDate{
		year:  2013,
		month: 10,
		day:   28,
	}

	s.assertParseFullDate(c, fd, buff, len(buff), nil)
}

func (s *Rfc5424TestSuite) TestParseHour_InvalidString(c *C) {
	buff := []byte("azer")
	expected := 0

	s.assertParseHour(c, expected, buff, 2, ErrHourInvalid)
}

func (s *Rfc5424TestSuite) TestParseHour_TooShort(c *C) {
	buff := []byte("1")
	expected := 0

	s.assertParseHour(c, expected, buff, 0, syslogparser.ErrEOL)
}

func (s *Rfc5424TestSuite) TestParseHour_InvalidRange(c *C) {
	buff := []byte("-1")
	expected := 0

	s.assertParseHour(c, expected, buff, 2, ErrHourInvalid)

	// ----

	buff = []byte("24")

	s.assertParseHour(c, expected, buff, 2, ErrHourInvalid)
}

func (s *Rfc5424TestSuite) TestParseHour_Valid(c *C) {
	buff := []byte("12")
	expected := 12

	s.assertParseHour(c, expected, buff, 2, nil)
}

func (s *Rfc5424TestSuite) TestParseMinute_InvalidString(c *C) {
	buff := []byte("azer")
	expected := 0

	s.assertParseMinute(c, expected, buff, 2, ErrMinuteInvalid)
}

func (s *Rfc5424TestSuite) TestParseMinute_TooShort(c *C) {
	buff := []byte("1")
	expected := 0

	s.assertParseMinute(c, expected, buff, 0, syslogparser.ErrEOL)
}

func (s *Rfc5424TestSuite) TestParseMinute_InvalidRange(c *C) {
	buff := []byte("-1")
	expected := 0

	s.assertParseMinute(c, expected, buff, 2, ErrMinuteInvalid)

	// ----

	buff = []byte("60")

	s.assertParseMinute(c, expected, buff, 2, ErrMinuteInvalid)
}

func (s *Rfc5424TestSuite) TestParseMinute_Valid(c *C) {
	buff := []byte("12")
	expected := 12

	s.assertParseMinute(c, expected, buff, 2, nil)
}

func (s *Rfc5424TestSuite) TestParseSecond_InvalidString(c *C) {
	buff := []byte("azer")
	expected := 0

	s.assertParseSecond(c, expected, buff, 2, ErrSecondInvalid)
}

func (s *Rfc5424TestSuite) TestParseSecond_TooShort(c *C) {
	buff := []byte("1")
	expected := 0

	s.assertParseSecond(c, expected, buff, 0, syslogparser.ErrEOL)
}

func (s *Rfc5424TestSuite) TestParseSecond_InvalidRange(c *C) {
	buff := []byte("-1")
	expected := 0

	s.assertParseSecond(c, expected, buff, 2, ErrSecondInvalid)

	// ----

	buff = []byte("60")

	s.assertParseSecond(c, expected, buff, 2, ErrSecondInvalid)
}

func (s *Rfc5424TestSuite) TestParseSecond_Valid(c *C) {
	buff := []byte("12")
	expected := 12

	s.assertParseSecond(c, expected, buff, 2, nil)
}

func (s *Rfc5424TestSuite) TestParseSecFrac_InvalidString(c *C) {
	buff := []byte("azerty")
	expected := 0.0

	s.assertParseSecFrac(c, expected, buff, 0, ErrSecFracInvalid)
}

func (s *Rfc5424TestSuite) TestParseSecFrac_NanoSeconds(c *C) {
	buff := []byte("123456789")
	expected := 0.123456

	s.assertParseSecFrac(c, expected, buff, 6, nil)
}

func (s *Rfc5424TestSuite) TestParseSecFrac_Valid(c *C) {
	buff := []byte("0")

	expected := 0.0
	s.assertParseSecFrac(c, expected, buff, 1, nil)

	buff = []byte("52")
	expected = 0.52
	s.assertParseSecFrac(c, expected, buff, 2, nil)

	buff = []byte("003")
	expected = 0.003
	s.assertParseSecFrac(c, expected, buff, 3, nil)

	buff = []byte("000003")
	expected = 0.000003
	s.assertParseSecFrac(c, expected, buff, 6, nil)
}

func (s *Rfc5424TestSuite) TestParseNumericalTimeOffset_Valid(c *C) {
	buff := []byte("+02:00")
	cursor := 0
	l := len(buff)
	tmpTs, err := time.Parse("-07:00", string(buff))
	c.Assert(err, IsNil)

	obtained, err := parseNumericalTimeOffset(buff, &cursor, l)
	c.Assert(err, IsNil)

	expected := tmpTs.Location()
	c.Assert(obtained, DeepEquals, expected)

	c.Assert(cursor, Equals, 6)
}

func (s *Rfc5424TestSuite) TestParseTimeOffset_Valid(c *C) {
	buff := []byte("Z")
	cursor := 0
	l := len(buff)

	obtained, err := parseTimeOffset(buff, &cursor, l)
	c.Assert(err, IsNil)
	c.Assert(obtained, DeepEquals, time.UTC)
	c.Assert(cursor, Equals, 1)
}

func (s *Rfc5424TestSuite) TestGetHourMin_Valid(c *C) {
	buff := []byte("12:34")
	cursor := 0
	l := len(buff)

	expectedHour := 12
	expectedMinute := 34

	obtainedHour, obtainedMinute, err := getHourMinute(buff, &cursor, l)
	c.Assert(err, IsNil)
	c.Assert(obtainedHour, Equals, expectedHour)
	c.Assert(obtainedMinute, Equals, expectedMinute)

	c.Assert(cursor, Equals, l)
}

func (s *Rfc5424TestSuite) TestParsePartialTime_Valid(c *C) {
	buff := []byte("05:14:15.000003")
	cursor := 0
	l := len(buff)

	obtained, err := parsePartialTime(buff, &cursor, l)
	expected := partialTime{
		hour:    5,
		minute:  14,
		seconds: 15,
		secFrac: 0.000003,
	}

	c.Assert(err, IsNil)
	c.Assert(obtained, DeepEquals, expected)
	c.Assert(cursor, Equals, l)
}

func (s *Rfc5424TestSuite) TestParseFullTime_Valid(c *C) {
	tz := "-02:00"
	buff := []byte("05:14:15.000003" + tz)
	cursor := 0
	l := len(buff)

	tmpTs, err := time.Parse("-07:00", string(tz))
	c.Assert(err, IsNil)

	obtainedFt, err := parseFullTime(buff, &cursor, l)
	expectedFt := fullTime{
		pt: partialTime{
			hour:    5,
			minute:  14,
			seconds: 15,
			secFrac: 0.000003,
		},
		loc: tmpTs.Location(),
	}

	c.Assert(err, IsNil)
	c.Assert(obtainedFt, DeepEquals, expectedFt)
	c.Assert(cursor, Equals, 21)
}

func (s *Rfc5424TestSuite) TestToNSec(c *C) {
	fixtures := []float64{
		0.52,
		0.003,
		0.000003,
	}

	expected := []int{
		520000000,
		3000000,
		3000,
	}

	c.Assert(len(fixtures), Equals, len(expected))
	for i, f := range fixtures {
		obtained, err := toNSec(f)
		c.Assert(err, IsNil)
		c.Assert(obtained, Equals, expected[i])
	}
}

func (s *Rfc5424TestSuite) TestParseAppName_Valid(c *C) {
	buff := []byte("su ")
	appName := "su"

	s.assertParseAppName(c, appName, buff, 2, nil)
}

func (s *Rfc5424TestSuite) TestParseAppName_TooLong(c *C) {
	// > 48chars
	buff := []byte("suuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuu ")
	appName := ""

	s.assertParseAppName(c, appName, buff, 48, ErrInvalidAppName)
}

func (s *Rfc5424TestSuite) TestParseProcId_Valid(c *C) {
	buff := []byte("123foo ")
	procId := "123foo"

	s.assertParseProcId(c, procId, buff, 6, nil)
}

func (s *Rfc5424TestSuite) TestParseProcId_TooLong(c *C) {
	// > 128chars
	buff := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab ")
	procId := ""

	s.assertParseProcId(c, procId, buff, 128, ErrInvalidProcId)
}

func (s *Rfc5424TestSuite) TestParseMsgId_Valid(c *C) {
	buff := []byte("123foo ")
	procId := "123foo"

	s.assertParseMsgId(c, procId, buff, 6, nil)
}

func (s *Rfc5424TestSuite) TestParseMsgId_TooLong(c *C) {
	// > 32chars
	buff := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa ")
	procId := ""

	s.assertParseMsgId(c, procId, buff, 32, ErrInvalidMsgId)
}

func (s *Rfc5424TestSuite) TestParseStructuredData_NilValue(c *C) {
	// > 32chars
	buff := []byte("-")
	sdData := "-"

	s.assertParseSdName(c, sdData, buff, 1, nil)
}

func (s *Rfc5424TestSuite) TestParseStructuredData_SingleStructuredData(c *C) {
	sdData := `[exampleSDID@32473 iut="3" eventSource="Application"eventID="1011"]`
	buff := []byte(sdData)

	s.assertParseSdName(c, sdData, buff, len(buff), nil)
}

func (s *Rfc5424TestSuite) TestParseStructuredData_MultipleStructuredData(c *C) {
	sdData := `[exampleSDID@32473 iut="3" eventSource="Application"eventID="1011"][examplePriority@32473 class="high"]`
	buff := []byte(sdData)

	s.assertParseSdName(c, sdData, buff, len(buff), nil)
}

func (s *Rfc5424TestSuite) TestParseStructuredData_MultipleStructuredDataInvalid(c *C) {
	a := `[exampleSDID@32473 iut="3" eventSource="Application"eventID="1011"]`
	sdData := a + ` [examplePriority@32473 class="high"]`
	buff := []byte(sdData)

	s.assertParseSdName(c, a, buff, len(a), nil)
}

// -------------

func (s *Rfc5424TestSuite) BenchmarkParseTimestamp(c *C) {
	buff := []byte("2003-08-24T05:14:15.000003-07:00")

	p := NewParser(buff)

	for i := 0; i < c.N; i++ {
		_, err := p.parseTimestamp()
		if err != nil {
			panic(err)
		}

		p.cursor = 0
	}
}

func (s *Rfc5424TestSuite) BenchmarkParseHeader(c *C) {
	buff := []byte("<165>1 2003-10-11T22:14:15.003Z mymachine.example.com su 123 ID47")

	p := NewParser(buff)

	for i := 0; i < c.N; i++ {
		_, err := p.parseHeader()
		if err != nil {
			panic(err)
		}

		p.cursor = 0
	}
}

// -------------

func (s *Rfc5424TestSuite) assertTimestamp(c *C, ts time.Time, b []byte, expC int, e error) {
	p := NewParser(b)
	obtained, err := p.parseTimestamp()
	c.Assert(err, Equals, e)

	tFmt := time.RFC3339Nano
	c.Assert(obtained.Format(tFmt), Equals, ts.Format(tFmt))

	c.Assert(p.cursor, Equals, expC)
}

func (s *Rfc5424TestSuite) assertFindNextSpace(c *C, nextSpace int, b []byte, e error) {
	obtained, err := syslogparser.FindNextSpace(b, 0, len(b))
	c.Assert(obtained, Equals, nextSpace)
	c.Assert(err, Equals, e)
}

func (s *Rfc5424TestSuite) assertParseYear(c *C, year int, b []byte, expC int, e error) {
	cursor := 0
	obtained, err := parseYear(b, &cursor, len(b))
	c.Assert(obtained, Equals, year)
	c.Assert(err, Equals, e)
	c.Assert(cursor, Equals, expC)
}

func (s *Rfc5424TestSuite) assertParseMonth(c *C, month int, b []byte, expC int, e error) {
	cursor := 0
	obtained, err := parseMonth(b, &cursor, len(b))
	c.Assert(obtained, Equals, month)
	c.Assert(err, Equals, e)
	c.Assert(cursor, Equals, expC)
}

func (s *Rfc5424TestSuite) assertParseDay(c *C, day int, b []byte, expC int, e error) {
	cursor := 0
	obtained, err := parseDay(b, &cursor, len(b))
	c.Assert(obtained, Equals, day)
	c.Assert(err, Equals, e)
	c.Assert(cursor, Equals, expC)
}

func (s *Rfc5424TestSuite) assertParseFullDate(c *C, fd fullDate, b []byte, expC int, e error) {
	cursor := 0
	obtained, err := parseFullDate(b, &cursor, len(b))
	c.Assert(err, Equals, e)
	c.Assert(obtained, Equals, fd)
	c.Assert(cursor, Equals, expC)
}

func (s *Rfc5424TestSuite) assertParseHour(c *C, hour int, b []byte, expC int, e error) {
	cursor := 0
	obtained, err := parseHour(b, &cursor, len(b))
	c.Assert(obtained, Equals, hour)
	c.Assert(err, Equals, e)
	c.Assert(cursor, Equals, expC)
}

func (s *Rfc5424TestSuite) assertParseMinute(c *C, minute int, b []byte, expC int, e error) {
	cursor := 0
	obtained, err := parseMinute(b, &cursor, len(b))
	c.Assert(obtained, Equals, minute)
	c.Assert(err, Equals, e)
	c.Assert(cursor, Equals, expC)
}

func (s *Rfc5424TestSuite) assertParseSecond(c *C, second int, b []byte, expC int, e error) {
	cursor := 0
	obtained, err := parseSecond(b, &cursor, len(b))
	c.Assert(obtained, Equals, second)
	c.Assert(err, Equals, e)
	c.Assert(cursor, Equals, expC)
}

func (s *Rfc5424TestSuite) assertParseSecFrac(c *C, secFrac float64, b []byte, expC int, e error) {
	cursor := 0
	obtained, err := parseSecFrac(b, &cursor, len(b))
	c.Assert(obtained, Equals, secFrac)
	c.Assert(err, Equals, e)
	c.Assert(cursor, Equals, expC)
}

func (s *Rfc5424TestSuite) assertParseAppName(c *C, appName string, b []byte, expC int, e error) {
	p := NewParser(b)
	obtained, err := p.parseAppName()

	c.Assert(err, Equals, e)
	c.Assert(obtained, Equals, appName)
	c.Assert(p.cursor, Equals, expC)
}

func (s *Rfc5424TestSuite) assertParseProcId(c *C, procId string, b []byte, expC int, e error) {
	p := NewParser(b)
	obtained, err := p.parseProcId()

	c.Assert(err, Equals, e)
	c.Assert(obtained, Equals, procId)
	c.Assert(p.cursor, Equals, expC)
}

func (s *Rfc5424TestSuite) assertParseMsgId(c *C, msgId string, b []byte, expC int, e error) {
	p := NewParser(b)
	obtained, err := p.parseMsgId()

	c.Assert(err, Equals, e)
	c.Assert(obtained, Equals, msgId)
	c.Assert(p.cursor, Equals, expC)
}

func (s *Rfc5424TestSuite) assertParseSdName(c *C, sdData string, b []byte, expC int, e error) {
	cursor := 0
	obtained, err := parseStructuredData(b, &cursor, len(b))

	c.Assert(err, Equals, e)
	c.Assert(obtained, Equals, sdData)
	c.Assert(cursor, Equals, expC)
}
