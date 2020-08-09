package subscriber

import (
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/data"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const timeout = 5 * time.Second

type worker struct {
	sub      data.Subscription
	dao      data.HistoryDAO
	l        *zap.SugaredLogger
	stopChan chan command
}

func (w *worker) run() {
	select {
	case cmd := <-w.stopChan:
		w.handleStopCommand(cmd)
		return
	default:
		content, duration := w.fetch()
		// Check again before adding. Content is already created and has timestamp.
		// Therefore this ensures that we won't add to database something that was created after the update.
		// Moreover we DO NOT schedule next fetch before this check.
		// Therefore no other goroutine will try to read from the stopChan.
		select {
		case cmd := <-w.stopChan:
			w.handleStopCommand(cmd)
			return
		default:
			// Schedule next before adding to database. Adding might take some time.
			w.scheduleNextFetch(duration)
			if err := w.dao.AddToHistory(content); err != nil {
				w.l.Infof("worker id=%d: failed to add content %+v", w.sub.ID, *content)
			} else {
				w.l.Debugf("worker id=%d: added %+v", w.sub.ID, *content)
			}
			return
		}
	}
}

// Returns fetched content and duration of the whole process
func (w *worker) fetch() (*data.Content, time.Duration) {
	start := time.Now()

	c := &http.Client{
		Timeout: timeout,
	}
	resp, err := c.Get(w.sub.URL)
	if err != nil {
		urlErr := err.(*url.Error)
		if strings.Contains(urlErr.Err.Error(), "Client.Timeout") {
			return data.NewContent(w.sub.ID, nil, timeout.Seconds()), timeout
		} else {
			duration := time.Now().Sub(start)
			return data.NewContent(w.sub.ID, nil, duration.Seconds()), duration
		}
	}
	defer resp.Body.Close()

	responseRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.l.Infof("worker id=%d: failed to read response %s", w.sub.ID, err)
		duration := time.Now().Sub(start)
		return data.NewContent(w.sub.ID, nil, duration.Seconds()), duration
	}

	duration := time.Now().Sub(start)
	response := string(responseRaw)
	return data.NewContent(w.sub.ID, &response, duration.Seconds()), duration
}

func (w *worker) cleanHistory() {
	if deletedCount, err := w.dao.DeleteAllHistory(w.sub.ID); err != nil {
		w.l.Errorf("worker id=%d: failed to clean history %s", w.sub.ID, err)
	} else {
		w.l.Infof("worker id=%d: cleaned history, deleted %d entries", w.sub.ID, deletedCount)
	}
}

// Assumption: We want interval to be the time from the beginning of the previous download.
func (w *worker) scheduleNextFetch(previousFetchDuration time.Duration) {
	// Case 1
	// Previous fetch took more time and one or more intervals have passed.
	// For instance, interval = 3, previousFetchDuration = 5. Then we should wait (3 - 5 % 3) seconds for next fetch.
	// 5 % 3 represents time elapsed in current interval
	// Case 2
	// Previous fetch took less time than single interval
	// Then we should wait (interval - previousFetchDuration) == (interval - previousFetchDuration % interval)
	interval := time.Second * time.Duration(w.sub.Interval)
	timeToNext := interval - previousFetchDuration%interval
	time.AfterFunc(timeToNext, w.run)
}

func (w *worker) handleStopCommand(cmd command) {
	switch cmd {
	case stop:
		w.l.Infof("worker id=%d: received stop command", w.sub.ID)
		return
	case stopAndClean:
		w.l.Infof("worker id=%d: received stopAndClean command", w.sub.ID)
		w.cleanHistory()
		return
	default:
		w.l.Infof("worker id=%d: received unknown command", w.sub.ID)
		return
	}
}
