package contact

import (
	"encoding/json"
	"os"
)

type Contact struct {
	Id    int `json:"id"`
	Name  string `json:"name"`
	Phone int `json:"phone"`
	Email string `json:"email"`
}

type Contacts []Contact

const fileName = "contacts.json"

func NewContact(id int, name string, phone int, email string) Contact {
	return Contact{
		Id:    id,
		Name:  name,
		Phone: phone,
		Email: email,
	}
}

func ReadJSON() Contacts {
	_, err := os.Stat(fileName)
	if err == os.ErrNotExist {
		_, err = os.Create(fileName)
		if err != nil {
			panic(err)
		}
		return nil
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var cts Contacts
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cts)
	if err != nil {
		panic(err)
	}
	return cts
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

func (cts Contacts) HaveEmail(email string) Contact {
	for _, c := range cts {
		if c.Email == email {
			return c
		}
	}
	return Contact{}
}