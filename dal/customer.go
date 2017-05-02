package dal

import (
	"github.com/golang-rest-api/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

type BaseDAL interface {
	Get(id string) (models.Customer, error)
	GetAll() ([]models.Customer, error)
	Create(customer models.Customer) error
	Update(id string, customer models.Customer) error
}

type CustomerDAL struct {
	s *mgo.Session
}

func NewCustomerDAL() BaseDAL {
	session, err := mgo.Dial(os.Getenv("DB_HOST"))
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB(os.Getenv("DB_NAME")).C("customers")
	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	d := new(CustomerDAL)
	d.s = session.Copy()
	return d
}

func (c *CustomerDAL) Get(id string) (models.Customer, error) {
	session := c.s.Copy()
	defer session.Close()

	col := session.DB(os.Getenv("DB_NAME")).C("customers")

	var customer models.Customer
	err := col.Find(bson.M{"id": id}).One(&customer)
	return customer, err
}

func (c *CustomerDAL) GetAll() ([]models.Customer, error) {
	session := c.s.Copy()
	defer session.Close()

	col := session.DB(os.Getenv("DB_NAME")).C("customers")

	var customers []models.Customer
	err := col.Find(nil).All(&customers)
	return customers, err
}

func (c *CustomerDAL) Create(customer models.Customer) error {
	session := c.s.Copy()
	defer session.Close()

	col := session.DB(os.Getenv("DB_NAME")).C("customers")
	return col.Insert(customer)
}

func (c *CustomerDAL) Update(id string, customer models.Customer) error {
	session := c.s.Copy()
	defer session.Close()

	col := session.DB(os.Getenv("DB_NAME")).C("customers")
	return col.Update(bson.M{"id": id}, &customer)
}
