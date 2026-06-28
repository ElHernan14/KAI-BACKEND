package home

import (
	"context"
	"testing"
	"time"

	repository "kai-back/internal/modules/home/repository"

	"github.com/google/uuid"
)

type homeRepositoryStub struct {
	evolutionMessage *string
	fallbackMessage  *string
}

func TestIsFirstInteractionToday(t *testing.T) {
	now := time.Date(2026, time.June, 26, 18, 0, 0, 0, time.Local)
	sameDay := time.Date(2026, time.June, 26, 8, 0, 0, 0, time.Local)
	previousDay := time.Date(2026, time.June, 25, 23, 0, 0, 0, time.Local)

	if isFirstInteractionToday(&sameDay, now) {
		t.Fatal("no debe generar otro saludo durante el mismo dia")
	}
	if !isFirstInteractionToday(&previousDay, now) {
		t.Fatal("debe generar saludo en el primer ingreso de un nuevo dia")
	}
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
