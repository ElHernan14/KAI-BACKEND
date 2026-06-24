package home

import (
	"context"
	"testing"

	repository "kai-back/internal/modules/home/repository"

	"github.com/google/uuid"
)

type homeRepositoryStub struct {
	evolutionMessage *string
	fallbackMessage  *string
}

func (s *homeRepositoryStub) FindKaiSummary(context.Context, uuid.UUID) (*repository.KaiSummaryRow, error) {
	return nil, nil
}
func (s *homeRepositoryStub) FindTotalXP(context.Context, uuid.UUID) (int, error) {
	return 0, nil
}
func (s *homeRepositoryStub) FindCurrentStreak(context.Context, uuid.UUID) (int, error) {
	return 0, nil
}
func (s *homeRepositoryStub) FindDailyHabits(context.Context, uuid.UUID) ([]repository.DailyHabitRow, error) {
	return nil, nil
}
func (s *homeRepositoryStub) FindFallbackMessage(context.Context) (*string, error) {
	return s.fallbackMessage, nil
}
func (s *homeRepositoryStub) FindRandomEvolutionMessage(context.Context) (*string, error) {
	return s.evolutionMessage, nil
}

func TestResolveMessagePrefersEvolutionMessageDuringEvent(t *testing.T) {
	evolutionMessage := "Kai evolucionó"
	normalMessage := "Mensaje normal"
	service := &Service{repository: &homeRepositoryStub{evolutionMessage: &evolutionMessage}}

	got, err := service.resolveMessage(context.Background(), &normalMessage, true)
	if err != nil {
		t.Fatal(err)
	}
	if got != evolutionMessage {
		t.Fatalf("message = %q; want %q", got, evolutionMessage)
	}
}

func TestResolveMessageUsesNormalMessageOutsideEvent(t *testing.T) {
	evolutionMessage := "Kai evolucionó"
	normalMessage := "Mensaje normal"
	service := &Service{repository: &homeRepositoryStub{evolutionMessage: &evolutionMessage}}

	got, err := service.resolveMessage(context.Background(), &normalMessage, false)
	if err != nil {
		t.Fatal(err)
	}
	if got != normalMessage {
		t.Fatalf("message = %q; want %q", got, normalMessage)
	}
}
