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
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parsel (Number, Client, Status, Address, CreatedAt) VALUES (:Number, :Client, :Status, :Address, :CreatedAt)",
		sql.Named("Number", p.Number),
		sql.Named("Client", p.Client),
		sql.Named("Status", p.Status),
		sql.Named("Address", p.Address),
		sql.Named("CreatedAt", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	i, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(i), nil

}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	parcel := Parcel{}
	rows, err := s.db.Query("SELECT Number, Client, Status, Address, CreatedAt FROM parsel WHERE Number = :number", sql.Named("Number", number))
	if err != nil {
		return parcel, err
	}
	defer rows.Close()
	// заполните объект Parcel данными из таблицы

	scanErr := rows.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
	if scanErr != nil {
		return parcel, err
	}

	return parcel, scanErr
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	rows, err := s.db.Query("SELECT Number, Client, Status, Address, CreatedAt FROM parsel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		parcel := Parcel{}

		err := rows.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, parcel)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE Status SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	par, err := s.Get(number)
	if err == nil && par.Status == "registered" {
		_, err := s.db.Exec("UPDATE Address SET status = :registered WHERE number = :number",
			sql.Named("address", address),
			sql.Named("number", number))
		return err
	}
	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	par, err := s.Get(number)
	if err == nil && par.Status == "registered" {
		_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))
		return err
	}
	return err
}
