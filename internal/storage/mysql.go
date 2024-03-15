package storage

import (
	"context"
	"database/sql"
	"donation-service/internal/data"
	"fmt"
)

type NewDonationListener interface {
	NotifyNewDonation(donation data.Donation)
}

type Storage struct {
	DB        *sql.DB
	listeners []NewDonationListener
}

func NewStorage(db *sql.DB, listeners ...NewDonationListener) Storage {
	return Storage{
		DB:        db,
		listeners: listeners,
	}
}

func (s Storage) InsertDonation(d data.Donation) (data.Donation, error) {
	var insertResult sql.Result
	var err error

	if d.Comment != "" {
		query := "INSERT INTO `donations` (`amount`, `sender`, `comment`, `time`) VALUES (?, ?, ?, NOW())"

		insertResult, err = s.DB.ExecContext(context.Background(), query, d.Amount, d.Sender, d.Comment)
	} else {
		query := "INSERT INTO `donations` (`amount`, `sender`, `time`) VALUES (?, ?, NOW())"

		insertResult, err = s.DB.ExecContext(context.Background(), query, d.Amount, d.Sender)
	}

	if err != nil {
		return d, fmt.Errorf("could not insert donation: %w", err)
	}

	id, err := insertResult.LastInsertId()
	if err != nil {
		return d, fmt.Errorf("could not retrieve donation id: %w", err)
	}

	d.ID = int(id)

	s.notifyNewDonation(d)

	return d, nil
}

func (s Storage) GetDonations() ([]data.Donation, error) {
	results, err := s.DB.Query("SELECT `id`, `amount`, `sender`, `comment` FROM `donations`")
	if err != nil {
		return nil, fmt.Errorf("could not get donations: %w", err)
	}

	res := make([]data.Donation, 0, 10)

	for results.Next() {
		var d data.Donation

		err = results.Scan(&d.ID, &d.Amount, &d.Sender, &d.Comment)
		if err != nil {
			return nil, fmt.Errorf("could not get info for donation: %w", err)
		}

		res = append(res, d)
	}

	return res, nil
}

func (s Storage) GetDonationsByID(id int) (data.Donation, error) {
	var d data.Donation

	err := s.DB.QueryRow("SELECT `id`, `amount`, `sender`, `comment` FROM `donations` where `id` = ?", id).
		Scan(&d.ID, &d.Amount, &d.Sender, &d.Comment)

	if err != nil {
		return data.Donation{}, fmt.Errorf("could not get donation: %w", err)
	}

	return d, nil
}

func (s Storage) GetMaxDonation() (data.Donation, error) {
	var d data.Donation

	err := s.DB.QueryRow("SELECT `id`, `amount`, `sender`, `comment` FROM `donations` where `amount` = (select max(`amount`) from donations) order by time limit 1").
		Scan(&d.ID, &d.Amount, &d.Sender, &d.Comment)

	if err != nil {
		return data.Donation{}, fmt.Errorf("could not get max donation: %w", err)
	}

	return d, nil
}

func (s Storage) notifyNewDonation(d data.Donation) {
	for _, listener := range s.listeners {
		listener.NotifyNewDonation(d)
	}
}
