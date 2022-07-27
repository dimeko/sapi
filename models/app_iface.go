package models

type Actions interface {
	List(limit string, offset string) ([]*User, error)
	Get(id string) (*User, error)
	Create(payload UserPayload) (*User, error)
	Update(id string, payload UserPayload) (*User, error)
	Delete(id string) error
}