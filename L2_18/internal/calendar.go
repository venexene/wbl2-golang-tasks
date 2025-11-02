package calendar

import (
	"time"
	"fmt"
)


type Calendar struct {
	events []Event
}

func NewCalendar() *Calendar {
    return &Calendar{
        events: []Event{},
    }
}


type Event struct {
	userID string
	date time.Time
	text string
}


func newEvent(userID string, date string, text string) (*Event, error) {
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("Invalid date: %v", err)
	}

	if userID == "" {
		return nil, fmt.Errorf("UserID can't be empty")
	}

	if text == "" {
		return nil, fmt.Errorf("Event text can't be empty")
	}

	return &Event{
		userID: userID,
		date: 	dateTime,
		text:	text,
	}, nil
}

func (c *Calendar) Add(userID string, date string, text string) error {
	event, err := newEvent(userID, date, text)
	if err != nil {
		return fmt.Errorf("Error creating new event: %v", err)
	}
	c.events = append(c.events, *event)
	return nil
}

func (c *Calendar) findEvent(userID string, date time.Time, text string) (int ,*Event, error) {
	for i, event := range c.events {
		if event.userID == userID && event.date == date && event.text == text {
			return i, &c.events[i], nil
		}
	}
	return 0, nil, fmt.Errorf("Event not found")
} 

func (c *Calendar) Update(userID string, date string, text string, newText string) error {
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("Invalid date: %v", err)
	}
	
	_, event, err := c.findEvent(userID, dateTime, text) 
	if err != nil {
		return fmt.Errorf("Error updating event: %v", err)
	}
	event.text = newText
	return nil
}

func (c *Calendar) Delete(userID string, date string, text string, newText string) error {
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("Invalid date: %v", err)
	}
	
	ind, _, err := c.findEvent(userID, dateTime, text) 
	if err != nil {
		return fmt.Errorf("Error deleting event: %v", err)
	}
	c.events = append(c.events[:ind], c.events[ind+1:]...)
	return nil
}

func (c *Calendar) GetEventByPeriod(period string) ([]string, error) {
	periodEvents := []string{}

	if periodDate, err := time.Parse("02-01-2006", period); err == nil {
		for _, event := range c.events {
			if event.date == periodDate {
				periodEvents = append(periodEvents, event.text)
			}
		}
		return periodEvents, nil
	}

	if periodDate, err := time.Parse("01-2006", period); err == nil {
		for _, event := range c.events {
			if event.date.Year() == periodDate.Year() && event.date.Month() == periodDate.Month() {
				periodEvents = append(periodEvents, event.text)
			}
		}
		return periodEvents, nil
	}

	if periodDate, err := time.Parse("2006", period); err == nil {
		for _, event := range c.events {
			if event.date.Year() == periodDate.Year() {
				periodEvents = append(periodEvents, event.text)
			}
		}
		return periodEvents, nil
	}

	return nil, fmt.Errorf("Not expected format of date")
}