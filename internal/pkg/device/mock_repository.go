package device

import "errors"

// MockRepository implements device.Repository for testing.
type MockRepository struct {
	Devices []Device
	Err     error
}

func (m *MockRepository) List() ([]Device, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Devices, nil
}

func (m *MockRepository) FindByID(id string) (*Device, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	for _, d := range m.Devices {
		if d.ID == id {
			return &d, nil
		}
	}
	return nil, nil
}

func (m *MockRepository) FindByBrand(brand string) ([]Device, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	var results []Device
	for _, d := range m.Devices {
		if d.Brand == brand {
			results = append(results, d)
		}
	}
	return results, nil
}

func (m *MockRepository) Store(device *Device) error {
	if m.Err != nil {
		return m.Err
	}
	m.Devices = append(m.Devices, *device)
	return nil
}

func (m *MockRepository) Update(device *Device) error {
	if m.Err != nil {
		return m.Err
	}
	for i, d := range m.Devices {
		if d.ID == device.ID {
			m.Devices[i] = *device
			return nil
		}
	}
	return errors.New("device not found")
}

func (m *MockRepository) Remove(id string) error {
	if m.Err != nil {
		return m.Err
	}
	for i, d := range m.Devices {
		if d.ID == id {
			m.Devices = append(m.Devices[:i], m.Devices[i+1:]...)
			return nil
		}
	}
	return errors.New("device not found")
}
