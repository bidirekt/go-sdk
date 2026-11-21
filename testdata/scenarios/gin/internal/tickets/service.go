package tickets

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("ticket not found")

var ErrSoldOut = errors.New("event sold out")

type Ticket struct {
	ID       int64
	EventID  int64
	Quantity int
	Email    string
	Coupon   *string
	Canceled bool
}

type Service struct {
	mutex     sync.Mutex
	tickets   map[int64]Ticket
	remaining map[int64]int
	nextID    int64
}

func NewService() *Service {
	return &Service{tickets: map[int64]Ticket{}, remaining: map[int64]int{1: 100, 2: 20}, nextID: 1}
}

func (this *Service) Book(eventID int64, quantity int, email string, coupon *string) (Ticket, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.remaining[eventID] < quantity {
		return Ticket{}, ErrSoldOut
	}
	this.remaining[eventID] -= quantity
	ticket := Ticket{ID: this.nextID, EventID: eventID, Quantity: quantity, Email: email, Coupon: coupon}
	this.tickets[ticket.ID] = ticket
	this.nextID++
	return ticket, nil
}

func (this *Service) Find(id int64) (Ticket, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	ticket, found := this.tickets[id]
	if !found {
		return Ticket{}, ErrNotFound
	}
	return ticket, nil
}

func (this *Service) Cancel(id int64) (Ticket, error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	ticket, found := this.tickets[id]
	if !found {
		return Ticket{}, ErrNotFound
	}
	ticket.Canceled = true
	this.tickets[id] = ticket
	this.remaining[ticket.EventID] += ticket.Quantity
	return ticket, nil
}

func (this *Service) Counts() (active int, canceled int) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	for _, ticket := range this.tickets {
		if ticket.Canceled {
			canceled++
		} else {
			active++
		}
	}
	return active, canceled
}
