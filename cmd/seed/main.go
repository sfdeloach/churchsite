package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/sfdeloach/churchsite/internal/config"
	"github.com/sfdeloach/churchsite/internal/database"
	"github.com/sfdeloach/churchsite/internal/models"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	now := time.Now()

	events := []models.Event{
		{
			Title:       "Lord's Day Morning Worship",
			Description: "Join us for our regular Lord's Day morning worship service with preaching from God's Word.",
			EventDate:   nextSunday(now).Add(10*time.Hour + 30*time.Minute),
			Location:    "Main Sanctuary",
			IsPublic:    true,
		},
		{
			Title:       "Men's Prayer Breakfast",
			Description: "Monthly men's fellowship and prayer breakfast. All men are welcome.",
			EventDate:   nextSaturday(now).Add(8 * time.Hour),
			Location:    "Fellowship Hall",
			IsPublic:    true,
		},
		{
			Title:       "Women's Bible Study",
			Description: "Weekly women's Bible study exploring the book of Ruth.",
			EventDate:   nextWeekday(now, time.Tuesday).Add(10 * time.Hour),
			Location:    "Room 204",
			IsPublic:    true,
		},
		{
			Title:       "Youth Group Game Night",
			Description: "Fun and fellowship for students in grades 6-12. Bring a friend!",
			EventDate:   nextWeekday(now, time.Friday).Add(18*time.Hour + 30*time.Minute),
			Location:    "Youth Room",
			IsPublic:    true,
		},
		{
			Title:       "Church Picnic",
			Description: "Annual church picnic at the park. Bring a dish to share. Hamburgers and hot dogs provided.",
			EventDate:   now.AddDate(0, 1, 0).Truncate(24*time.Hour).Add(11 * time.Hour),
			Location:    "Riverside Park — Shelter 3",
			IsPublic:    true,
		},
	}

	for _, event := range events {
		if err := db.Postgres.Create(&event).Error; err != nil {
			slog.Error("failed to seed event", "title", event.Title, "error", err)
			continue
		}
		slog.Info("seeded event", "title", event.Title, "date", event.EventDate.Format("Jan 2, 2006 3:04 PM"))
	}

	slog.Info("seeding complete", "events", len(events))
}

func nextSunday(from time.Time) time.Time {
	return nextWeekday(from, time.Sunday)
}

func nextSaturday(from time.Time) time.Time {
	return nextWeekday(from, time.Saturday)
}

func nextWeekday(from time.Time, day time.Weekday) time.Time {
	daysUntil := int(day-from.Weekday()+7) % 7
	if daysUntil == 0 {
		daysUntil = 7
	}
	return from.AddDate(0, 0, daysUntil).Truncate(24 * time.Hour)
}
