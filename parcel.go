package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	//  rows to parcel table

	res, err := s.db.Exec(
		`INSERT INTO parcel (client, status, address, created_at) 
		 VALUES (:client, :status, :address, :created_at)`,
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt),
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {

	//  read line by given number

	row := s.db.QueryRow(
		`SELECT number, client, status, address, created_at 
		 FROM parcel 
		 WHERE number = :number`,
		sql.Named("number", number),
	)

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// read rows from parcel table by given client

	rows, err := s.db.Query(
		`SELECT number, client, status, address, created_at 
		 FROM parcel 
		 WHERE client = :client`,
		sql.Named("client", client),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Parcel
	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// update status in parcel table

	_, err := s.db.Exec(
		`UPDATE parcel 
		 SET status = :status 
		 WHERE number = :number`,
		sql.Named("status", status),
		sql.Named("number", number),
	)
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// update the address only if the status is "registered"
	_, err := s.db.Exec(
		`UPDATE parcel 
		 SET address = :address 
		 WHERE number = :number AND status = 'registered'`,
		sql.Named("address", address),
		sql.Named("number", number),
	)
	return err
}

func (s ParcelStore) Delete(number int) error {
	// delete only if status is 'registered'
	_, err := s.db.Exec(
		`DELETE FROM parcel 
		 WHERE number = :number AND status = 'registered'`,
		sql.Named("number", number),
	)
	return err
}
