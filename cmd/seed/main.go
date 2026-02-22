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

	// Seed staff members
	staffMembers := []models.StaffMember{
		{
			Name:         "Rev. James McAllister",
			Title:        "Senior Pastor",
			Bio:          "Pastor McAllister has faithfully served Saint Andrew's Chapel since 2008. A graduate of Reformed Theological Seminary, he is passionate about expository preaching and shepherding the flock entrusted to his care. He and his wife Margaret have three children.",
			Email:        "pastor@sachapel.com",
			DisplayOrder: 1,
			IsActive:     true,
			Category:     models.CategoryPastor,
		},
		{
			Name:         "Rev. David Kim",
			Title:        "Associate Pastor",
			Bio:          "Pastor Kim joined the staff in 2015 after serving as a church planter in South Korea. He oversees adult education, small groups, and missions. He holds an M.Div. from Westminster Theological Seminary.",
			Email:        "dkim@sachapel.com",
			DisplayOrder: 2,
			IsActive:     true,
			Category:     models.CategoryPastor,
		},
		{
			Name:         "Sarah Mitchell",
			Title:        "Director of Music",
			Bio:          "Sarah has led our music ministry since 2012, bringing a deep love for traditional Reformed hymnody and psalmody. She holds a Master of Music from the University of Michigan and also directs the church choir.",
			Email:        "music@sachapel.com",
			DisplayOrder: 1,
			IsActive:     true,
			Category:     models.CategoryStaff,
		},
		{
			Name:         "Robert Chen",
			Title:        "Youth Director",
			Bio:          "Robert joined Saint Andrew's in 2019 to lead our youth ministry. He is committed to teaching young people the fundamentals of the Reformed faith. He graduated from Covenant College with a degree in Bible and Theology.",
			Email:        "youth@sachapel.com",
			DisplayOrder: 2,
			IsActive:     true,
			Category:     models.CategoryStaff,
		},
		{
			Name:         "Linda Patterson",
			Title:        "Office Administrator",
			Bio:          "Linda keeps everything running smoothly at Saint Andrew's. She manages church communications, coordinates facility use, and supports all ministry activities. She has been a faithful member of the congregation for over 20 years.",
			Email:        "office@sachapel.com",
			DisplayOrder: 3,
			IsActive:     true,
			Category:     models.CategoryStaff,
		},
	}

	for _, member := range staffMembers {
		if err := db.Postgres.Create(&member).Error; err != nil {
			slog.Error("failed to seed staff member", "name", member.Name, "error", err)
			continue
		}
		slog.Info("seeded staff member", "name", member.Name, "title", member.Title)
	}

	// Seed ministries (use FirstOrCreate because slug has UNIQUE constraint)
	ministries := []models.Ministry{
		{
			Name:         "Sunday School",
			Slug:         "sunday-school",
			Description:  "Biblical instruction for all ages, grounding our congregation in the Reformed faith every Lord's Day morning.",
			ContactEmail: "sundayschool@sachapel.com",
			MeetingTime:  "Sundays, 9:15 AM",
			Location:     "Education Wing",
			IsActive:     true,
			SortOrder:    1,
			PageContent:  "<h2>Rooted in the Word</h2><p>Sunday School at Saint Andrew's Chapel is a cornerstone of our educational ministry. Each Sunday morning at 9:15 AM, our congregation gathers by age group to study the Word together.</p>",
		},
		{
			Name:         "Women's Ministry",
			Slug:         "womens-ministry",
			Description:  "Encouraging women to grow in grace through Bible study, fellowship, and service to the body of Christ.",
			ContactEmail: "women@sachapel.com",
			MeetingTime:  "Tuesdays, 10:00 AM",
			Location:     "Room 204",
			IsActive:     true,
			SortOrder:    2,
			PageContent:  "<h2>Women Growing Together in Grace</h2><p>The Women's Ministry exists to encourage and equip women to know Christ more deeply, love one another more faithfully, and serve the body of Christ more joyfully.</p>",
		},
		{
			Name:         "Youth Ministry",
			Slug:         "youth-ministry",
			Description:  "Discipling students in grades 6–12 in the truth of God's Word and the fellowship of the Reformed faith.",
			ContactEmail: "youth@sachapel.com",
			MeetingTime:  "Sundays, 5:00 PM",
			Location:     "Youth Room",
			IsActive:     true,
			SortOrder:    3,
			PageContent:  "<h2>Raising Up the Next Generation</h2><p>The Youth Ministry is committed to the discipleship of students in grades 6 through 12, rooting them in the truth of God's Word and the community of the covenant people.</p>",
		},
		{
			Name:         "Music Ministry",
			Slug:         "music-ministry",
			Description:  "Offering excellence in sacred music to the glory of God through psalmody, hymnody, and choral worship.",
			ContactEmail: "music@sachapel.com",
			MeetingTime:  "Thursdays, 7:00 PM (Choir rehearsal)",
			Location:     "Sanctuary",
			IsActive:     true,
			SortOrder:    4,
			PageContent:  "<h2>Singing to the Glory of God</h2><p>The Music Ministry understands corporate song as an act of worship, rooted in the historic Reformed tradition's high view of congregational singing.</p>",
		},
		{
			Name:         "Mercy Ministry",
			Slug:         "mercy-ministry",
			Description:  "Serving those in need within our congregation and community, reflecting the compassion of Christ.",
			ContactEmail: "mercy@sachapel.com",
			MeetingTime:  "As needs arise",
			Location:     "Various locations",
			IsActive:     true,
			SortOrder:    5,
			PageContent:  "<h2>The Ministry of Mercy</h2><p>The Mercy Ministry is grounded in the conviction that the church is called not only to proclaim the gospel in word but to embody it in deed.</p>",
		},
		{
			Name:         "Men's Fellowship",
			Slug:         "mens-fellowship",
			Description:  "Building men of God through prayer, accountability, Scripture study, and sacrificial service.",
			ContactEmail: "men@sachapel.com",
			MeetingTime:  "Second Saturday of each month, 8:00 AM",
			Location:     "Fellowship Hall",
			IsActive:     true,
			SortOrder:    6,
			PageContent:  "<h2>Iron Sharpening Iron</h2><p>The Men's Fellowship exists to build men who are rooted in Christ, committed to their families, faithful to the church, and engaged in the world for the glory of God.</p>",
		},
	}

	for _, ministry := range ministries {
		result := db.Postgres.Where(models.Ministry{Slug: ministry.Slug}).FirstOrCreate(&ministry)
		if result.Error != nil {
			slog.Error("failed to seed ministry", "slug", ministry.Slug, "error", result.Error)
			continue
		}
		slog.Info("seeded ministry", "name", ministry.Name, "slug", ministry.Slug)
	}

	slog.Info("seeding complete", "events", len(events), "staff_members", len(staffMembers), "ministries", len(ministries))
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
