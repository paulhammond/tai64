package tai64

import (
	"encoding/binary"
	"strconv"
	"time"
)

// This is a list of all leap seconds added since 1972, in TAI seconds since
// the unix epoch. It is derived from
// http://www.ietf.org/timezones/data/leap-seconds.list
// http://hpiers.obspm.fr/eop-pc/earthor/utc/UTC.html
// http://maia.usno.navy.mil/leapsec.html
var leapSeconds = []int64{
	// subtract 2208988800 to convert from NTP datetime to unix seconds
	// then add number of previous leap seconds to get TAI-since-unix-epoch
	1341100834,
	1230768033,
	1136073632,
	915148831,
	867715230,
	820454429,
	773020828,
	741484827,
	709948826,
	662688025,
	631152024,
	567993623,
	489024022,
	425865621,
	394329620,
	362793619,
	315532818,
	283996817,
	252460816,
	220924815,
	189302414,
	157766413,
	126230412,
	94694411,
	78796810,
	63072009,
}

type Error struct {
	message string
}

func (e Error) Error() string {
	return e.message
}

var parseError = Error{"Parse Error"}

func ParseTai64(s string) (time.Time, error) {
	if len(s) != 17 || s[0] != '@' {
		return time.Time{}, parseError
	}
	sec, err := strconv.ParseUint(s[1:], 16, 64)
	if err != nil {
		return time.Time{}, parseError
	}
	if sec > 1<<63 {
		return time.Time{}, parseError
	}
	return TaiDate(int64(sec-(1<<62)), 0), nil
}

func ParseTai64n(s string) (time.Time, error) {
	// http://cr.yp.to/daemontools/tai64n.html
	// http://cr.yp.to/libtai/tai64.html
	// "A TAI64N label is normally stored or communicated in external TAI64N
	// format, consisting of twelve 8-bit bytes", which is 24 chars of hex
	if len(s) != 25 || s[0] != '@' {
		return time.Time{}, parseError
	}
	// "The first eight bytes are the TAI64 label"
	sec, err := strconv.ParseUint(s[1:17], 16, 64)
	if err != nil {
		return time.Time{}, parseError
	}
	// "The last four bytes are the nanosecond counter in big-endian format"
	nsec, err := strconv.ParseUint(s[17:25], 16, 32)
	if err != nil {
		return time.Time{}, parseError
	}
	if sec > 1<<63 {
		return time.Time{}, parseError
	}
	return TaiDate(int64(sec-(1<<62)), int64(nsec)), nil
}

func DecodeTai64(b []byte) (time.Time, error) {
	if len(b) != 8 {
		return time.Time{}, parseError
	}
	sec := binary.BigEndian.Uint64(b)
	if sec > 1<<63 {
		return time.Time{}, parseError
	}
	return TaiDate(int64(sec-(1<<62)), 0), nil
}

func DecodeTai64n(b []byte) (time.Time, error) {
	if len(b) != 12 {
		return time.Time{}, parseError
	}
	sec := binary.BigEndian.Uint64(b[0:8])
	nsec := binary.BigEndian.Uint32(b[8:12])
	if sec > 1<<63 {
		return time.Time{}, parseError
	}
	return TaiDate(int64(sec-(1<<62)), int64(nsec)), nil
}

func TaiDate(secs, nsecs int64) time.Time {
	offset := len(leapSeconds) + 10
	for _, l := range leapSeconds {
		offset--
		if secs > l {
			break
		}
	}
	return time.Unix(secs-int64(offset), nsecs)
}
