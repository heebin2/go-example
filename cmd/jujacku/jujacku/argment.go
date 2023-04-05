package jujacku

import (
	"flag"
	"fmt"
	"time"
)

type Config struct {
	Start time.Time
	End time.Time
	InterID int
	Flow int
	Percentage float32
}

func (c Config) String() string {
	return fmt.Sprintf(`Start      : %s
 End        : %s
 Inter      : %d
 Flow       : %d
 Percentage : %f`, c.Start, c.End, c.InterID, c.Flow, c.Percentage)
}

func NewConfig() (Config, error) {
	start := flag.String("s", "", "start time 2006-01-02_15:04:05")
	end := flag.String("e", "", "end time 2006-01-02_15:04:05")
	inter := flag.Int("i", 0, "intersection 1 ~")
	flow := flag.Int("f", 0, "flow 1,2,3,4")
	per := flag.Float64("p", 0, "add percentage 0 ~ 100")
	flag.Parse()

	var err error
	if start == nil || end == nil {
		return Config{}, fmt.Errorf("start time or end time is null")
	}

	ret := Config{}

	ret.Start, err = time.ParseInLocation("2006-01-02_15:04:05", *start, time.Local)
	if err != nil {
		return Config{}, fmt.Errorf("start time convert error : %s (%s)", err, *start)
	}

	ret.End, err = time.ParseInLocation("2006-01-02_15:04:05", *end, time.Local)
	if err != nil {
		return Config{}, fmt.Errorf("end time convert error : %s (%s)", err, *end)
	}

	if inter == nil || *inter == 0 {
		return Config{}, fmt.Errorf("inter error")
	}
	ret.InterID = *inter

	if flow == nil || *flow < 0 || *flow > 4 {
		return Config{}, fmt.Errorf("flow error")
	}
	ret.Flow = *flow

	if per == nil || *per < 0 || *per > 100 {
		return Config{}, fmt.Errorf("percentage error")
	}
	ret.Percentage = float32(*per)

	return ret, nil
}