package contact

import (
	"encoding/json"
	"os"
)

type Contact struct {
	Id    int `json:"id,omitempty"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Errors map[string]error `json:"error"`
}

type Contacts []Contact

const fileName = "contacts.json"

func NewContact(name, phone, email string, err map[string]error) Contact {
	return Contact{
		Name:  name,
		Phone: phone,
		Email: email,
		Errors: err,
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

func (cts Contacts) GetMaxId() int {
	maxId := 0
	for _, c := range cts {
		if c.Id > maxId {
			maxId = c.Id
		}
	}
	return maxId
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

func (cts Contacts) HaveEmail(email string) Contacts {
	for _, c := range cts {
		if c.Email == email {
			return Contacts{c}
		}
	}
	return Contacts{}
}