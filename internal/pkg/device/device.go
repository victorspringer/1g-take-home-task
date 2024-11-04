package device

import "time"

type Device struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Brand        string    `json:"brand"`
	CreationTime time.Time `json:"creationTime"`
	UpdateTime   time.Time `json:"updateTime"`
}

type Repository interface {
	Store(device *Device) error
	FindByID(id string) (*Device, error)
	List() ([]Device, error)
	Update(device *Device) error
	Remove(id string) error
	FindByBrand(brand string) ([]Device, error)
}
