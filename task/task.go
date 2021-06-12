package task

import (
	"bytes"
	"container/list"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/cheolgyu/dbment/ticker/utils"
)

type Task struct {
	Debug bool

	TimeFormat string
	LogDir     string

	TickerPlanCycle time.Duration
	TickerPlan      *list.List
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

	o.TickerPlan = list.New()
	o.TickerCH = make(chan bool)
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
	last := o.TickerPlan.Back()
	text := ""
	t := time.Now()

	lastExecTime := t
	if last == nil {
		runTime := time.Date(t.Year(), t.Month(), t.Day(), 15, 30, 0, 0, time.Local)
		lastExecTime = runTime.Add(-1 * time.Hour * 24)
	} else {
		lastExecTime = last.Value.(time.Time)
	}

	n_day := t.Weekday()
	var weekday = n_day != 0 && n_day != 6

	nextExecTime := lastExecTime.Add(o.TickerPlanCycle)
	if !weekday {
		runTime := time.Date(t.Year(), t.Month(), t.Day(), 15, 30, 0, 0, time.Local)
		if n_day == 0 {
			runTime.Add(1 * time.Hour * 24)
		} else if n_day == 6 {
			runTime.Add(2 * time.Hour * 24)
		}
		nextExecTime = runTime
	}
	if o.Debug {
		waiting := nextExecTime.Sub(t)
		text += "   now         :" + fmt.Sprint(t.Format(o.TimeFormat)) + "\n"
		text += " lastExecTime  :" + fmt.Sprint(lastExecTime.Format(o.TimeFormat)) + "\n"
		text += " nextExecTime  :" + fmt.Sprint(nextExecTime.Format(o.TimeFormat)) + "\n"
		text += " waiting       :" + fmt.Sprintf("%v", waiting) + "\n"

		o.log("\n" + text)
	}

	return nextExecTime

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

					o.log("[디비작업 시작]]" + time.Now().String())

					if o.Debug {
						execCmd_test()
					} else {
						execCmd()
					}

					o.log("[디비작업 종료]]" + time.Now().String())
					o.TickerPlan.PushBack(time.Now())

				}

			}
		}
	}()

	wg.Wait()

}

func execCmd_test() {

	time.Sleep(time.Second * 5)

}

func execCmd() {

	cmd := exec.Command("./dbment", "-run=daily", "-prod")

	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("in all caps: %q\n", out.String())
}
