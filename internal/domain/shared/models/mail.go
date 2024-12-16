package models

type Mail struct {
	ID      string
	Address string
	Subject string
	Content string
	Mailed  bool
}

func (m *Mail) Stored() bool {
	return m.ID != ""
}

func (m *Mail) SetMailed() {
	m.Mailed = true
}
