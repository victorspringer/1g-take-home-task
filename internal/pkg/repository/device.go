package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/victorspringer/1g-take-home-task/internal/pkg/device"
)

type Client struct {
	db *pgxpool.Pool
}

func New(connString string) (*Client, error) {
	ctx := context.Background()

	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	createSchema := `
	CREATE TABLE IF NOT EXISTS devices (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		brand TEXT NOT NULL,
		creation_time TIMESTAMPTZ NOT NULL,
		update_time TIMESTAMPTZ NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_brand ON devices(brand);
	`
	_, err = db.Exec(ctx, createSchema)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Client{db: db}, nil
}

func (c *Client) Close() {
	c.db.Close()
}

func (c *Client) Store(device *device.Device) error {
	device.ID = uuid.New().String()
	device.CreationTime = time.Now()
	device.UpdateTime = device.CreationTime

	_, err := c.db.Exec(
		context.Background(),
		"INSERT INTO devices (id, name, brand, creation_time, update_time) VALUES ($1, $2, $3, $4, $5)",
		device.ID, device.Name, device.Brand, device.CreationTime, device.UpdateTime,
	)

	return err
}

func (c *Client) FindByID(id string) (*device.Device, error) {
	row := c.db.QueryRow(context.Background(), "SELECT id, name, brand, creation_time, update_time FROM devices WHERE id=$1", id)

	device := &device.Device{}

	err := row.Scan(&device.ID, &device.Name, &device.Brand, &device.CreationTime, &device.UpdateTime)
	if err != nil && err.Error() == pgx.ErrNoRows.Error() {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return device, nil
}

func (c *Client) List() ([]device.Device, error) {
	rows, err := c.db.Query(context.Background(), "SELECT id, name, brand, creation_time, update_time FROM devices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []device.Device
	for rows.Next() {
		var device device.Device
		if err := rows.Scan(&device.ID, &device.Name, &device.Brand, &device.CreationTime, &device.UpdateTime); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func (c *Client) Update(device *device.Device) error {
	if device.ID == "" {
		return errors.New("device id is required")
	}

	if device.Name != "" && device.Brand != "" {
		_, err := c.db.Exec(
			context.Background(),
			"UPDATE devices SET name=$2, brand=$3, update_time=$4 WHERE id=$1",
			device.ID, device.Name, device.Brand, time.Now(),
		)
		return err
	} else if device.Name != "" {
		_, err := c.db.Exec(
			context.Background(),
			"UPDATE devices SET name=$2, update_time=$3 WHERE id=$1",
			device.ID, device.Name, time.Now(),
		)
		return err
	} else if device.Brand != "" {
		_, err := c.db.Exec(
			context.Background(),
			"UPDATE devices SET brand=$2, update_time=$3 WHERE id=$1",
			device.ID, device.Brand, time.Now(),
		)
		return err
	}

	return errors.New("invalid update request")
}

func (c *Client) Remove(id string) error {
	_, err := c.db.Exec(context.Background(), "DELETE FROM devices WHERE id=$1", id)
	return err
}

func (c *Client) FindByBrand(brand string) ([]device.Device, error) {
	rows, err := c.db.Query(context.Background(), "SELECT id, name, brand, creation_time, update_time FROM devices WHERE brand=$1", brand)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []device.Device
	for rows.Next() {
		var device device.Device
		if err := rows.Scan(&device.ID, &device.Name, &device.Brand, &device.CreationTime, &device.UpdateTime); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}
