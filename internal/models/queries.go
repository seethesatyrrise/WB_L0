package models

type Restoration struct {
	ID   string `sql:"id"`
	Data Order  `sql:"_data"`
}

type SelectData struct {
	Data Order `sql:"select_data"`
}
