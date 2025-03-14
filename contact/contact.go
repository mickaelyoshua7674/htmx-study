package contact

import (
	"encoding/json"
	"os"
	"slices"
)

type Contact struct {
	Id    int `json:"id,omitempty"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func (c *Contact) Update(name, phone, email string) {
	c.Name = name
	c.Phone = phone
	c.Email = email
}

type Contacts []Contact

const fileName = "contacts.json"

func NewContact(id int, name, phone, email string) Contact {
	return Contact{
		Id: id,
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

func (cts Contacts) GetMaxId() int {
	maxId := 0
	for _, c := range cts {
		if c.Id > maxId {
			maxId = c.Id
		}
	}
	return maxId
}

func (cts Contacts) GetIndexById(id int) int {
	for i, c := range cts {
		if c.Id == id {
			return i
		}
	}
	return -1
}

func (cts Contacts) GetContactById(id int) Contact {
	for _, c := range cts {
		if c.Id == id {
			return c
		}
	}
	return Contact{}
}

func (cts Contacts) GetIdByEmail(email string) int {
	for _, c := range cts {
		if c.Email == email {
			return c.Id
		}
	}
	return -1
}

func (cts *Contacts) DeleteById(id int) int {
	index := cts.GetIndexById(id)
	//*cts = append((*cts)[:index], (*cts)[index+1:]...)
	*cts = slices.Delete(*cts, index, index+1)
	return index
}

func (cts *Contacts) WriteJSON() error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cts)
	if err != nil {
		return err
	}
	return nil
}