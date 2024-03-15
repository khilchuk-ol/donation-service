package mono

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vtopc/go-monobank"

	"donation-service/internal/data"
)

const (
	MinTs = 65 * time.Second

	AccountId = "account"
)

type Service struct {
	Client  *monobank.PersonalClient
	Storage data.DonationRepository
	Logger  *log.Logger

	AccumulatedTotal float32
}

func NewService(client *monobank.PersonalClient, repository data.DonationRepository, logger *log.Logger) Service {
	return Service{
		Client:           client,
		Storage:          repository,
		Logger:           logger,
		AccumulatedTotal: 0,
	}
}

func (s *Service) PoolAccountInfo(ctx context.Context, waitTs time.Duration) {
	if waitTs < MinTs {
		waitTs = MinTs
	}

	now := time.Now()

	for {
		<-time.After(waitTs)

		newNow := time.Now()

		go s.extractInfo(ctx, AccountId, now, newNow)

		now = newNow
	}
}

func (s *Service) extractInfo(ctx context.Context, accountID string, from, to time.Time) {
	resp, err := s.Client.Transactions(ctx, accountID, from, to)

	if err != nil {
		s.Logger.Print(fmt.Sprintf("could not get transactions from mono: %w", err))
	}

	for _, transaction := range resp {
		d := data.Donation{
			Amount:  float32(transaction.Amount) / 100,
			Comment: transaction.Comment,
			Sender:  transaction.Description,
		}

		s.AccumulatedTotal += d.Amount

		d, err = s.Storage.InsertDonation(d)
		if err != nil {
			s.Logger.Print(fmt.Sprintf("could not create new donation: %w", err))
		}
	}
}
