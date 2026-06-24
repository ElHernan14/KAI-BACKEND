package evolution

import (
	"testing"
	"time"
)

func TestStageForXP(t *testing.T) {
	tests := []struct {
		name string
		xp   int
		want string
	}{
		{name: "sin experiencia", xp: 0, want: StagePuppy},
		{name: "antes de joven", xp: 29, want: StagePuppy},
		{name: "joven", xp: 30, want: StageYoung},
		{name: "antes de adulto", xp: 59, want: StageYoung},
		{name: "adulto", xp: 60, want: StageAdult},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StageForXP(tt.xp); got != tt.want {
				t.Fatalf("StageForXP(%d) = %q; want %q", tt.xp, got, tt.want)
			}
		})
	}
}

func TestEventWindow(t *testing.T) {
	startedAt := time.Date(2026, 6, 24, 12, 0, 0, 0, time.UTC)

	active, expiresAt := EventWindow(&startedAt, startedAt.Add(23*time.Hour+59*time.Minute))
	if !active {
		t.Fatal("expected evolution event to be active before 24 hours")
	}
	wantExpiration := startedAt.Add(EventDuration)
	if expiresAt == nil || !expiresAt.Equal(wantExpiration) {
		t.Fatalf("expiration = %v; want %v", expiresAt, wantExpiration)
	}

	active, _ = EventWindow(&startedAt, startedAt.Add(EventDuration))
	if active {
		t.Fatal("expected evolution event to expire at 24 hours")
	}
}
