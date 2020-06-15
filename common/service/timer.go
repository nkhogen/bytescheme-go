package service

import (
	"bytescheme/common/db"
	"bytescheme/common/util"
	"context"
	"fmt"
	"strings"
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
)

// EventCallback is the callback invoked by the timer
type EventCallback func(eventID string, data map[string]interface{}) error

// Timer is the timer
type Timer struct {
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
		store:         store,
		eventTimers:   map[string]*time.Timer{},
		eventCallback: eventCallback,
	}
	timer.ctx, timer.cancel = context.WithCancel(context.Background())
	go timer.watch()
	util.ShutdownHandler.RegisterCloseable(timer)
	return timer
}

// calculateEventTime calculates the event time based on the current event time.
// It returns the next event time and the delay.
func (timer *Timer) calculateNextEventTime(eventTime time.Time, recurMins int) (time.Time, time.Duration) {
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

func (timer *Timer) saveEvent(event *Event) error {
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
			eventTime, eventDelay := timer.calculateNextEventTime(event.Time, event.RecurMins)
			if event.Version == 0 {
				event.Time = eventTime
				err = timer.saveEvent(event)
				if err != nil {
					continue
				}
			}
			if eventDelay > EventScheduleWindow {
				// No need to schedule now
				fmt.Printf("Ingoring event %+v as the event time is too far way\n", event)
				continue
			}
			activeEventIDs[event.ID] = true
			fmt.Printf("Scheduling event %+v\n", event)
			eventTimer = time.AfterFunc(eventDelay, func() {
				if event.RecurMins != 0 {
					eventTime, _ := timer.calculateNextEventTime(event.Time, event.RecurMins)
					// now + just greater
					event.Time = eventTime
					fmt.Printf("Registering the recurring event %+v with next time %s\n", event, event.Time.String())
					timer.saveEvent(event)
					if err != nil {
						return
					}
				}
				// Callback
				fmt.Printf("Triggering event %+v\n", event)
				err = timer.eventCallback(event.ID, event.Data)
				if err != nil {
					fmt.Printf("Error in invoking the callback for event %+v. Error: %s\n", event, err.Error())
					return
				}
				// Delete the executed timer
				delete(timer.eventTimers, event.ID)
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
