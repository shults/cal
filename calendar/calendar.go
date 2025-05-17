package calendar

import (
	"bytes"
	"flag"
	"fmt"
	"time"
)

const (
	monthWith = 20
)

type (
	wHelper struct {
		buf bytes.Buffer
	}
)

func Run(args []string) ([]byte, error) {
	cmd := args[0]
	args = args[1:]

	now := time.Now()
	fs := newFlagSet(cmd)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	helper := newHelper()

	if fs.help {
		fs.Usage()
		return nil, nil
	} else {
		helper.printMonth(time.Date(fs.Year(), fs.Month(), 1, 0, 0, 0, 0, now.Location()))
	}

	return helper.bytes(), nil
}

func newHelper() *wHelper {
	return &wHelper{}
}

func (w *wHelper) printTitle(ts time.Time) {
	title := fmt.Sprintf("%s %d\n", ts.Month(), ts.Year())
	offset := (monthWith - len(title)) / 2
	w.printByteTimes(' ', offset+1)
	w.buf.WriteString(title)
}

func (w *wHelper) printDayNames() {
	_, _ = w.buf.Write([]byte("Su Mo Tu We Th Fr Sa\n"))
}

func (w *wHelper) printMonth(ts time.Time) {
	now := time.Now()

	year := ts.Year()
	month := ts.Month()

	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()

	weekdayOfFirstDay := firstDay.Weekday()

	w.printTitle(ts)
	w.printDayNames()

	for i := 0; i < int(weekdayOfFirstDay); i++ {
		w.buf.WriteString("   ")
	}

	for day := 1; day <= daysInMonth; day++ {
		if day == now.Day() && year == now.Year() && month == now.Month() {
			_, _ = w.buf.WriteString(fmt.Sprintf("\033[7m%2d\033[0m ", day))
		} else {
			_, _ = w.buf.WriteString(fmt.Sprintf("%2d ", day))
		}

		if day == daysInMonth {
			break
		}

		if (int(weekdayOfFirstDay)+day-1)%7 == 6 {
			w.printLineBreak()
		}
	}

	w.printLineBreak()
}

func (w *wHelper) printLineBreak() {
	_ = w.buf.WriteByte('\n')
}

func (w *wHelper) printByteTimes(ch byte, times int) {
	w.buf.Grow(w.buf.Len() + times)

	for i := 0; i < times; i++ {
		_ = w.buf.WriteByte(ch)
	}
}

func (w *wHelper) bytes() []byte {
	defer w.buf.Reset()
	return w.buf.Bytes()
}

type flagSet struct {
	year  int
	month int
	help  bool
	flags *flag.FlagSet
}

func newFlagSet(cmd string) *flagSet {
	now := time.Now()

	var fs = new(flagSet)
	fs.flags = flag.NewFlagSet(cmd, flag.ExitOnError)

	fs.flags.IntVar(&fs.year, "y", now.Year(), "year")
	fs.flags.IntVar(&fs.month, "m", int(now.Month()), "month")
	fs.flags.BoolVar(&fs.help, "h", false, "print help")

	return fs
}

func (fs *flagSet) Parse(args []string) error {
	if err := fs.flags.Parse(args); err != nil {
		return err
	}

	return nil
}

func (fs *flagSet) Usage() {
	fs.flags.Usage()
}

func (fs *flagSet) Year() int {
	return fs.year
}

func (fs *flagSet) Month() time.Month {
	return time.Month(fs.month)
}
