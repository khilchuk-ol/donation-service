package data

type Donation struct {
	ID       int     `json:"id"`
	Amount   float32 `json:"amount"`
	Sender   string  `json:"sender"`
	Comment  string  `json:"comment"`
	TimeUnix int64   `json:"time"`
}

type DonationRepository interface {
	InsertDonation(d Donation) (Donation, error)
	GetDonations() ([]Donation, error)
	GetDonationsByID(id int) (Donation, error)
}
