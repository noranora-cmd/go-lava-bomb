package volcano

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type volcano struct {
	secCh, minCh, hourCh chan string
	stopTime             time.Duration
	confCheckTime        int64
	doneCh               chan struct{}
	errorCh              chan error
	noiseCh              chan string

	mu sync.WaitGroup
}

type noise struct {
	Second string `json:"second"`
	Minute string `json:"minute"`
	Hour   string `json:"hour"`
}

// New is a constructor for a new volcano
func New(minutes time.Duration) (*volcano, error) {
	if minutes == 0 {
		return nil, errors.New("No stop time provided")
	}

	return &volcano{
		secCh:         make(chan string, 1),
		minCh:         make(chan string, 1),
		hourCh:        make(chan string, 1),
		stopTime:      minutes,
		confCheckTime: time.Now().Unix(),

		doneCh:  make(chan struct{}, 1),
		errorCh: make(chan error, 1),
		noiseCh: make(chan string),
	}, nil
}

// Erupt starts the ticking
func (c *volcano) Erupt(configFile string) {

	c.mu.Add(3)
	go c.rumble()
	go c.noiseMonitor(configFile)
	go c.print()

	c.mu.Wait()
	fmt.Println("volcano finished rumbling ...")
}

// rumble implements printing noises at requested intervals. It receives
// dynamically updated noises via channels, sent by noiseMonitor.
func (c *volcano) rumble() {
	defer c.mu.Done()
	defer func() {
		close(c.noiseCh)
	}()

	secTick := time.NewTicker(1 * time.Second)
	minTick := time.NewTicker(1 * time.Minute)
	hourTick := time.NewTicker(1 * time.Hour)

	secM := "rumble"
	minM := "RUMBLE"
	hourM := "LAVAOVERFLOW"

	start := time.Now()

	i := 0
	for {

		if time.Since(start) > c.stopTime {
			c.doneCh <- struct{}{}
			return
		}

		select {
		case m := <-c.secCh:
			secM = m
		case m := <-c.minCh:
			minM = m
		case m := <-c.hourCh:
			hourM = m
		default:
		}

		select {
		case err := <-c.errorCh:
			log.Println(err)
			return

		case <-secTick.C:
			select {
			case m := <-c.secCh:
				secM = m
			default:
			}
			i++
			if i%60 != 0 {
				c.noiseCh <- formatPrintNoise(i, secM)
			}

		case <-minTick.C:
			select {
			case m := <-c.minCh:
				minM = m
			default:
			}

			if i%3600 != 0 {
				c.noiseCh <- formatPrintNoise(i, minM)
			}

		case <-hourTick.C:
			select {
			case m := <-c.hourCh:
				hourM = m
			default:
			}
			c.noiseCh <- formatPrintNoise(i, hourM)
			i = 0
		}
	}
}

// noiseMonitor is responsible for delivering dynamically changed noises.
// This implementation uses a file based approach. Second, minute and hour interval
// noises are stored in json format in a file, regularly monitored by the noiseMonitor.
func (c *volcano) noiseMonitor(configFile string) {
	defer c.mu.Done()

	for {
		select {
		case <-c.doneCh:
			return
		case <-c.errorCh:
			return

		default:
			err := c.monitor(c.confCheckTime, configFile)
			if err != nil {
				c.errorCh <- err
			}
		}

		time.Sleep(1000 * time.Millisecond)
	}
}

func (c *volcano) monitor(confCheckTime int64, configFile string) error {
	noise, err := checkFile(c.confCheckTime, configFile)

	if err != nil {
		return err
	}
	c.confCheckTime = time.Now().Unix()

	if noise != nil {
		if noise.Second != "" {
			c.secCh <- noise.Second
		}
		if noise.Minute != "" {
			c.minCh <- noise.Minute
		}
		if noise.Hour != "" {
			c.hourCh <- noise.Hour
		}
	}
	return nil
}

func checkFile(confCheckTime int64, configFile string) (*noise, error) {

	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	timeDiff := stat.ModTime().Unix() - confCheckTime
	if timeDiff >= 0 {

		content, err := ioutil.ReadFile(configFile)
		if err != nil {
			return nil, err
		}

		var noise noise
		err = json.Unmarshal(content, &noise)
		if err != nil {
			return nil, err
		}
		return &noise, nil
	}
	return nil, nil
}

func formatPrintNoise(i int, m string) string {
	return fmt.Sprintf("%d  ...  %s", i, m)
}

func (c *volcano) print() {
	defer c.mu.Done()
	for {
		select {
		case msg, ok := <-c.noiseCh:
			println(msg)
			if !ok {
				return
			}
		default:
		}
	}
}
