package app

type CreatePayload struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type UpdatePayload struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type Actions interface {
	List(params ...string)
	Get(id string, params ...string)
	Create(payload CreatePayload) (interface{}, error)
	Update(id string, payload UpdatePayload)
	Delete(id string)
}
