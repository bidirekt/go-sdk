package api

import "example.com/ticketing/internal/tickets"

type BookTicketRequest struct {
	EventID  int64   `json:"eventId" binding:"required" bidirekt:"field"`
	Quantity int     `json:"quantity" binding:"required,min=1,max=10" bidirekt:"field"`
	Email    string  `json:"email" binding:"required,email" bidirekt:"field"`
	Coupon   *string `json:"coupon,omitempty" binding:"omitempty,alphanum" bidirekt:"field"`
}

type TicketResponse struct {
	ID       int64  `json:"id" bidirekt:"field"`
	EventID  int64  `json:"eventId" bidirekt:"field"`
	Quantity int    `json:"quantity" bidirekt:"field"`
	Email    string `json:"email" bidirekt:"field"`
	Canceled bool   `json:"canceled" bidirekt:"field"`
}

type ErrorResponse struct {
	Error string `json:"error" bidirekt:"field"`
}

func ticketToResponse(ticket tickets.Ticket) TicketResponse {
	return TicketResponse{
		ID:       ticket.ID,
		EventID:  ticket.EventID,
		Quantity: ticket.Quantity,
		Email:    ticket.Email,
		Canceled: ticket.Canceled,
	}
}
