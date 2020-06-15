package service

import (
	"bytescheme/common/db"
	"bytescheme/common/util"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	// TimerKeyPrefix is the key for timer event
	TimerKeyPrefix string = "timer/"
)

var (
	// EventScheduleWindow is the window for keeping timers in memory
	EventScheduleWindow = time.Second * 90
	// EventScanInterval is the DB scan interval
	EventScanInterval = time.Second * 60

	// CallbackRetryLimit is the retry limit for event callback
	CallbackRetryLimit = 3
	// CallbackRetryDelay is the retry delay
	CallbackRetryDelay = time.Millisecond * 300
)

// EventCallback is the callback invoked by the timer
type EventCallback func(eventID string, data map[string]interface{}) error

// Timer is the timer
type Timer struct {
	lock          *sync.Mutex
	store         *db.Store
	ctx           context.Context
	cancel        context.CancelFunc
	eventTimers   map[string]*time.Timer
	eventCallback EventCallback
}

// Event represents the event time
type Event struct {
	ID          string                 `json:"id"`
	Description string                 `json:"description"`
	Time        time.Time              `json:"time"`
	RecurMins   int                    `json:"recurMins"`
	Data        map[string]interface{} `json:"data"`
	Version     int64                  `json:"version"`
}

// NewTimer returns a new timer instance
func NewTimer(store *db.Store, eventCallback EventCallback) *Timer {
	timer := &Timer{
		lock:          &sync.Mutex{},
		store:         store,
		eventTimers:   map[string]*time.Timer{},
		eventCallback: eventCallback,
	}
	timer.ctx, timer.cancel = context.WithCancel(context.Background())
	go timer.watch()
	util.ShutdownHandler.RegisterCloseable(timer)
	return timer
}

// NextEventTime finds the next event time based on the current event time.
// It returns the next event time and the delay.
func (timer *Timer) NextEventTime(eventTime time.Time, recurMins int) (time.Time, time.Duration) {
	now := time.Now()
	if recurMins == 0 {
		remainingDuration := eventTime.Sub(now)
		if remainingDuration < time.Second {
			remainingDuration = time.Second
		}
		return eventTime, remainingDuration
	}
	delayDuration := time.Minute * time.Duration(recurMins)
	elapsedDuration := now.Sub(eventTime)
	remainingDuration := delayDuration - (elapsedDuration % delayDuration)
	return now.Add(remainingDuration), remainingDuration
}

// SaveEvent saves an event to the DB
func (timer *Timer) SaveEvent(event *Event) error {
	event.Version = time.Now().UnixNano()
	ba, err := util.ConvertToJSON(event)
	if err != nil {
		fmt.Printf("Error in scheduling the next event %+v. Error: %s\n", event, err.Error())
		return err
	}
	err = timer.store.Set(&db.KeyValue{
		Key:   TimerKeyPrefix + event.ID,
		Value: string(ba),
	})
	if err != nil {
		fmt.Printf("Error in scheduling the next event %+v. Error: %s\n", event, err.Error())
		return err
	}
	return nil
}

func (timer *Timer) watch() {
	scheduler := func() {
		fmt.Printf("Running scheduler...\n")
		keyValues, err := timer.store.Gets(TimerKeyPrefix)
		if err != nil {
			fmt.Printf("Error in fetching timers. Error: %s\n", err.Error())
			return
		}
		timer.lock.Lock()
		defer timer.lock.Unlock()
		activeEventIDs := map[string]bool{}
		for idx := range keyValues {
			keyValue := keyValues[idx]
			event := &Event{}
			err := util.ConvertFromJSON([]byte(keyValue.Value), event)
			if err != nil {
				fmt.Printf("Error in conversion for event %+v\n", *keyValue)
				continue
			}
			event.ID = strings.TrimPrefix(keyValue.Key, TimerKeyPrefix)
			eventTimer, ok := timer.eventTimers[event.ID]
			if ok {
				if event.Version == 0 {
					// Reload
					fmt.Printf("Cancelling exiting timer %+v as there is an update\n", event)
					eventTimer.Stop()
					delete(timer.eventTimers, event.ID)
				} else {
					continue
				}
			}
			now := time.Now()
			eventElapedTime := now.Sub(event.Time)
			if event.RecurMins == 0 && eventElapedTime > time.Minute {
				timer.store.Delete(TimerKeyPrefix + event.ID)
				continue
			}
			if event.Version == 0 {
				event.Time, _ = timer.NextEventTime(event.Time, event.RecurMins)
				err = timer.SaveEvent(event)
				if err != nil {
					continue
				}
			}
			eventDelay := event.Time.Sub(now)
			if eventDelay > EventScheduleWindow {
				// No need to schedule now
				fmt.Printf("Ingoring event %+v as the event time is too far way\n", event)
				continue
			}
			activeEventIDs[event.ID] = true
			fmt.Printf("Scheduling event %+v\n", event)
			eventTimer = time.AfterFunc(eventDelay, func() {
				defer func() {
					timer.lock.Lock()
					defer timer.lock.Unlock()
					// Delete the executed timer
					delete(timer.eventTimers, event.ID)
				}()
				// Callback
				fmt.Printf("Triggering event %+v\n", event)
				for i := 0; i < CallbackRetryLimit; i++ {
					err = timer.eventCallback(event.ID, event.Data)
					if err != nil {
						continue
					}
					time.Sleep(CallbackRetryDelay)
				}
				if err != nil {
					fmt.Printf("Error in invoking the callback for event %+v. Error: %s\n", event, err.Error())
					return
				}
				// Persist the next event time
				if event.RecurMins != 0 {
					event.Time, _ = timer.NextEventTime(event.Time, event.RecurMins)
					fmt.Printf("Registering the recurring event %+v with next time %s\n", event, event.Time.String())
					timer.SaveEvent(event)
				}
			})
			timer.eventTimers[event.ID] = eventTimer
		}
		// Any event ID which is not deleted
		for eventID, eventTimer := range timer.eventTimers {
			if activeEventIDs[eventID] {
				continue
			}
			fmt.Printf("Found inactive event ID %s\n", eventID)
			eventTimer.Stop()
			delete(timer.eventTimers, eventID)
		}
	}
	// Initial load
	scheduler()
	ticker := time.NewTicker(EventScanInterval)
	for {
		select {
		case <-timer.ctx.Done():
			return
		case <-ticker.C:
			scheduler()
		}
	}
}

// Close implements Closeable
func (timer *Timer) Close() error {
	timer.ctx.Done()
	return nil
}

// IsClosed returns if if the timer is already closed
func (timer *Timer) IsClosed() bool {
	select {
	case <-timer.ctx.Done():
		return true
	default:
		return false
	}
}

// Submit submits an event
func (timer *Timer) Submit(event *Event) (string, error) {
	// TODO
	return "", nil
}

// Cancel cancles an event
func (timer *Timer) Cancel(eventID string) (*Event, error) {
	// TODO
	return nil, nil
}
