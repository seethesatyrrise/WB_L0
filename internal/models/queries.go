package models

type Restoration struct {
	ID   string `sql:"id"`
	Data []byte `sql:"_data"`
}

type SelectData struct {
	Data []byte `sql:"select_data"`
}
