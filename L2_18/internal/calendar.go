// Package calendar provides simple calendar logic
package calendar

import (
	"fmt"
	"time"
)

// Calendar represents an event storage system
type Calendar struct {
	events []Event
}

// NewCalendar creates new calendar object
func NewCalendar() *Calendar {
	return &Calendar{
		events: []Event{},
	}
}

// Event represents a calendar event with user, date and description
type Event struct {
	userID string
	date   time.Time
	text   string
}

func newEvent(userID string, date string, text string) (*Event, error) {
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("Invalid date: %v", err)
	}

	if userID == "" {
		return nil, fmt.Errorf("UserID cant be empty")
	}

	if text == "" {
		return nil, fmt.Errorf("Event text cant be empty")
	}

	return &Event{
		userID: userID,
		date:   dateTime,
		text:   text,
	}, nil
}

// Add adds new event into caldenar
func (c *Calendar) Add(userID string, date string, text string) error {
	event, err := newEvent(userID, date, text)
	if err != nil {
		return fmt.Errorf("Error creating new event: %v", err)
	}
	c.events = append(c.events, *event)
	return nil
}

func (c *Calendar) findEvent(userID string, date time.Time, text string) (int, *Event, error) {
	for i, event := range c.events {
		if event.userID == userID && event.date == date && event.text == text {
			return i, &c.events[i], nil
		}
	}
	return 0, nil, fmt.Errorf("Event not found")
}

// Update provides ability to change event text
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

// Delete delets event from calendar
func (c *Calendar) Delete(userID string, date string, text string) error {
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

// GetEventsByDay returns events texts by day
func (c *Calendar) GetEventsByDay(userID string, day string) ([]string, error) {
	dayEvents := []string{}

	if dayDate, err := time.Parse("2006-01-02", day); err == nil {
		for _, event := range c.events {
			if event.userID == userID && event.date == dayDate {
				dayEvents = append(dayEvents, event.text)
			}
		}
		return dayEvents, nil
	}

	return nil, fmt.Errorf("Invalid format of date")
}

// GetEventsByWeek returns events texts by week
func (c *Calendar) GetEventsByWeek(userID string, week string) ([]string, error) {
	weekEvents := []string{}

	startDate, err := time.Parse("2006-01-02", week)
	if err != nil {
		return nil, fmt.Errorf("Invalid date format")
	}

	endDate := startDate.AddDate(0, 0, 6)

	for _, event := range c.events {
		if event.userID == userID &&
			(event.date.Equal(startDate) || event.date.After(startDate)) &&
			(event.date.Equal(endDate) || event.date.Before(endDate)) {
			weekEvents = append(weekEvents, event.text)
		}
	}

	return weekEvents, nil
}

// GetEventsByMonth returns events texts by month
func (c *Calendar) GetEventsByMonth(userID string, day string) ([]string, error) {
	monthEvents := []string{}

	if monthDate, err := time.Parse("2006-01-02", day); err == nil {
		for _, event := range c.events {
			if event.userID == userID &&
				event.date.Month() == monthDate.Month() &&
				event.date.Year() == monthDate.Year() {
				monthEvents = append(monthEvents, event.text)
			}
		}
		return monthEvents, nil
	}

	return nil, fmt.Errorf("Invalid format of date")
}
