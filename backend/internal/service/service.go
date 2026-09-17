package service

import (
	"crypto/rand"
	"math/big"

	"travelcrm/internal/repository"
)

// Service holds business logic and orchestrates the repository.
type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

// code alphabet excludes ambiguous characters (0/O, 1/I) for readability.
const codeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

// randCode returns a random code of length n using a crypto-strong source.
func randCode(n int) string {
	b := make([]byte, n)
	max := big.NewInt(int64(len(codeAlphabet)))
	for i := range b {
		idx, _ := rand.Int(rand.Reader, max)
		b[i] = codeAlphabet[idx.Int64()]
	}
	return string(b)
}

// newBookingCode / newTicketCode produce human-friendly, prefixed codes.
func newBookingCode() string { return "BK-" + randCode(6) }
func newTicketCode() string  { return "TK-" + randCode(8) }
