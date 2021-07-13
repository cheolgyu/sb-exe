package task

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/cheolgyu/stock-write-ticker/ticker/utils"
)

type Task struct {
	Debug bool
	Print bool

	TimeFormat string
	LogDir     string

	TickerPlanCycle time.Duration
	WorkList        map[int]bool
	TickerCH        chan bool
	Ticker          *time.Ticker

	logFile *os.File
	mylog   utils.Log
}

func (o *Task) init() {
	//o.TimeFormat = "2006-01-02_15_04_05"
	//o.LogDir = "./logs/"

	o.mylog = utils.Log{}
	o.mylog.LogDir = o.LogDir
	o.mylog.TimeFormat = o.TimeFormat

	o.logFile = o.mylog.CreateFile(time.Now().Format(o.TimeFormat))

	o.WorkList = make(map[int]bool)
	o.TickerCH = make(chan bool)

	o.Print = true
	//o.Ticker = time.NewTicker(1 * time.Second)
	o.log("hello world")

}

func (o *Task) log(text string) {
	o.mylog.Write(o.logFile, text)
}

func (o *Task) Run() {

	if len(os.Args) > 1 {
		o.Debug = true
		arg := os.Args[1]
		if arg == "test" {

		} else {
			panic("누구냐")
		}
	} else {
		o.Debug = false

	}
	if o.Debug {
		o.Ticker = time.NewTicker(1 * time.Second)
		o.TickerPlanCycle = time.Second * 10 //time.Hour * 24 //time.Second * 30
	} else {
		o.Ticker = time.NewTicker(1 * time.Minute)
		o.TickerPlanCycle = time.Hour * 24
	}
	o.init()
	o.ticker_exec()

}

func (o *Task) GetExecTime() time.Time {

	now := time.Now()
	cur_key := planRole(now)
	if _, working := o.WorkList[cur_key]; working {
		return not_work(now)
	}

	text := "===========================\n"

	weekday := fmt.Sprintf("%v", now.Weekday())
	var is_weekday = weekday != "Saturday" && weekday != "Sunday"

	nextExecTime := time.Date(now.Year(), now.Month(), now.Day(), 15, 30, 0, 0, time.Local)
	if !is_weekday {
		return not_work(now)
	}

	if o.Print {
		waiting := nextExecTime.Sub(now)
		text += "   now         :" + fmt.Sprint(now.Format(o.TimeFormat)) + "\n"
		text += " nextExecTime  :" + fmt.Sprint(nextExecTime.Format(o.TimeFormat)) + "\n"
		text += " waiting       :" + fmt.Sprintf("%v", waiting) + "\n"
		o.log("\n" + text)
	}

	return nextExecTime

}

func not_work(t time.Time) time.Time {

	return time.Date(t.Year()+1, t.Month(), t.Day(), 15, 30, 0, 0, time.Local)
}

func (o *Task) ticker_exec() {

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {

		for {
			select {
			case <-o.TickerCH:
				o.Ticker.Stop()
				wg.Done()
				return
			case <-o.Ticker.C:

				now := time.Now()

				if now.After(o.GetExecTime()) {
					key := planRole(now)
					o.WorkList[key] = false

					o.log("[디비작업 시작]]" + now.String())

					if o.Debug {
						execCmd_test()
					} else {
						execCmd()
					}

					o.log("[디비작업 종료]]" + now.String())
					o.WorkList[key] = true

				}

			}
		}
	}()

	wg.Wait()

}

func planRole(t time.Time) int {
	s := t.Format("20060102")
	i, e := strconv.ParseInt(s, 64, 0)
	if e != nil {
		log.Panicln(e)
	}

	return int(i)
}

func execCmd_test() {

	time.Sleep(time.Second * 5)

}

func execCmd() {

	cmd := exec.Command("/stock/stock-write")

	//cmd.Stdin = strings.NewReader("some input")
	//var out bytes.Buffer
	//cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("in all caps: %q\n", out.String())
}
