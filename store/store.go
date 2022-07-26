package store

type CreatePayload struct {
	Username  string
	Firstname string
	Lastname  string
}

type UpdatePayload struct {
	Username  string
	Firstname string
	Lastname  string
}

type Store interface {
	List(params ...string)
	Get(id string, params ...string)
	Create(payload CreatePayload) (interface{}, error)
	Update(id string, payload UpdatePayload)
	Delete(id string)
}

// Models
type User struct {
	Id        string
	Username  string
	Firstname string
	Lastname  string
}
