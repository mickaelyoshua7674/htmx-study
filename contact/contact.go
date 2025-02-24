package contact

import (
	"encoding/json"
	"os"
)

type contact struct {
	Id    int `json:"id"`
	Name  string `json:"name"`
	Phone int `json:"phone"`
	Email string `json:"email"`
}

type Contacts []contact

const fileName = "contacts.json"

func NewContact(id int, name string, phone int, email string) contact {
	return contact{
		Id:    id,
		Name:  name,
		Phone: phone,
		Email: email,
	}
}

func ReadJSON() (Contacts, error) {
	var cts Contacts
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return cts, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cts)
	if err != nil {
		return cts, err
	}
	return cts, nil
}

func (cts Contacts) WriteJSON() error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&cts)
	if err != nil {
		return err
	}
	return nil
}