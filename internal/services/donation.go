package services

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"donation-service/internal/data"
	"donation-service/internal/services/mono"
	"donation-service/internal/storage"
)

type DonationService struct {
	Storage        storage.Storage
	Logger         *log.Logger
	MonoService    *mono.Service
	DonationsCache *storage.Cache
}

func NewDonationService(s storage.Storage, logger *log.Logger, mono *mono.Service, cache *storage.Cache) DonationService {
	return DonationService{
		Storage:        s,
		Logger:         logger,
		MonoService:    mono,
		DonationsCache: cache,
	}
}

func (s DonationService) GetNewDonationsFromCache() []data.Donation {
	donations := s.DonationsCache.FlushAll()

	res := make([]data.Donation, 0, len(donations))

	for _, d := range donations {
		if d.ID > 0 {
			res = append(res, d)
		}
	}

	return res
}

func (s DonationService) GetMaxDonation() (data.Donation, bool) {
	d, err := s.Storage.GetMaxDonation()
	if err != nil {
		s.Logger.Print(fmt.Sprintf("could not get max donation: %w", err))

		return data.Donation{}, false
	}

	return d, true
}

func (s DonationService) GetTotalSum() float32 {
	return s.MonoService.AccumulatedTotal
}

func (s DonationService) GetNewDonationsTest() []data.Donation {
	res := make([]data.Donation, 6)

	for i := 0; i < 6; i++ {
		sender := "Bід: " + RandStringRunes(7)
		comment := ""
		amount := (i + 1) * (rand.Intn(i+1) + 10)

		if i%3 == 0 {
			comment = RandStringRunes(15)
		}

		res[i] = data.Donation{
			Sender:  sender,
			Comment: comment,
			Amount:  float32(amount),
		}
	}

	return res
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
