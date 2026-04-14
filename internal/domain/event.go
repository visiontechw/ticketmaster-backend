package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Base
	Name        string
	Description string
	OccursAt    time.Time
	Location    string
	Capacity    int64
	OwnerId     uuid.UUID
}

func (e *Event) validate() error {

	if e.OwnerId == uuid.Nil {
		return errors.New("owner id is required")
	}
	if strings.TrimSpace(e.Name) == "" {
		return errors.New("event name is required")
	}

	// Compara com o início do dia atual
	if e.OccursAt.Before(time.Now().Truncate(24 * time.Hour)) {
		return errors.New("event date cannot be in the past")
	}

	if strings.TrimSpace(e.Location) == "" {
		return errors.New("location is required")
	}

	if e.Capacity <= 0 {
		return errors.New("capacity must be greater than zero")
	}

	return nil
}

func NewEvent(ownerId uuid.UUID, name, description, location string, capacity int64, occursAt time.Time) (*Event, error) {
	event := &Event{
		Base:        NewBase(),
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Location:    strings.TrimSpace(location),
		Capacity:    capacity,
		OccursAt:    occursAt,
		OwnerId:     ownerId,
	}

	if err := event.validate(); err != nil {
		return nil, err
	}

	return event, nil
}

func (e *Event) CanEdit(userId uuid.UUID) bool {
	return e.OwnerId == userId
}

func (e *Event) Reschedule(newDate time.Time) error {

	if newDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return errors.New("event date cannot be in the past")
	}

	e.OccursAt = newDate
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Event) CanBeCanceled() bool {
	limit := time.Now().Add(24 * time.Hour)
	return e.OccursAt.After(limit)
}

func (e *Event) UpdateDescription(newDescription string) error {
	newDescription = strings.TrimSpace(newDescription)
	if newDescription == "" {
		return errors.New("new description cannot be empty")
	}
	e.Description = newDescription
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Event) UpdateLocation(newLocation string) error {
	newLocation = strings.TrimSpace(newLocation)
	if newLocation == "" {
		return errors.New("new location cannot be empty")
	}
	e.Location = newLocation
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Event) UpdateCapacity(newCapacity int64) error {

	if newCapacity < 1 {
		return errors.New("capacity must be positive")
	}

	e.Capacity = newCapacity
	e.UpdatedAt = time.Now()
	return nil
}
