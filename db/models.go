package db

import (
    "time"

    "github.com/lib/pq"
)

const (
    // StatusDraft draft project
    StatusDraft string = "draft"
    // StatusSuccess success project
    StatusSuccess string = "success"
    // StatusFail fail project
    StatusFail string = "fail"
    // StatusHarvest harvest project
    StatusHarvest string = "harvest"
    // StatusSearch search project
    StatusSearch string = "search"
)

// User model
type User struct {
    ID        uint       `json:"id" gorm:"primary_key" `
    CreatedAt time.Time  `json:"-"`
    UpdatedAt time.Time  `json:"-"`

    Username     string  `json:"username"`
    FirstName    string  `json:"first_name"`
    LastName     string  `json:"last_name"`
    Avatar       string  `json:"avatar"`
    Email        string  `gorm:"type:varchar(100)" json:"email"`
    IsStaff      bool    `gorm:"default:false" json:"is_staff"`
    ProjectCount int     `gorm:"default:0" json:"project_count"`
    SuccessRate  float32 `gorm:"default:0" json:"success_rate"`
}

// Project model
type Project struct {
    ID            uint        `gorm:"primary_key" json:"id"`
    CreatedAt     time.Time   `json:"-"`
    UpdatedAt     time.Time   `json:"-"`
    Title         string      `gorm:"not null; size:100" json:"title"`
    SubTitle      string      `gorm:"size:100" json:"subtitle"`
    ReleaseDate   time.Time   `gorm:"not null" json:"release_date"`
    EventDate     time.Time   `json:"event_date"`
    GoalPeople    uint        `gorm:"not null; default:0" json:"goal_people"`
    GoalAmount    uint        `gorm:"not null; default:0" json:"goal_amount"`
    Total         uint        `gorm:"not null; default:0" json:"total"`
    Description   string      `gorm:"size:1000" json:"description"`
    ImageLink     string      `gorm:"not null; size:200" json:"image_link"`
    Instructions  string      `gorm:"size:500" json:"instructions"`
    Locked        bool        `gorm:"not null; default=false" json:"-"`
    Published     bool        `gorm:"not null; default=false" json:"-"`
    Closed        bool        `gorm:"not null; default=false" json:"-"`
    Owner         User        `json:"owner"`
    OwnerID       uint        `json:"-"`
    Category      Category    `json:"category"`
    CategoryID    uint        `json:"-"`
    ProjectType   ProjectType `json:"project_type"`
    ProjectTypeID uint        `json:"-"`
}

// Status of project
func (p *Project) Status() string {
    if !p.Published {
        return StatusDraft
    }
    if p.Closed {
        if p.Locked {
            return StatusSuccess
        }
        return StatusFail
    }
    if p.Locked {
        return StatusHarvest
    }

    return StatusSearch
}

// Percent of project
func (p *Project) Percent() uint {
    return 30
}

// Category of project
type Category struct {
    ID    uint   `gorm:"primary_key" json:"id"`
    Alias string `gorm:"not null" json:"alias"`
    Name  string `gorm:"not null" json:"name"`
}

// ProjectType of project
type ProjectType struct {
    ID            uint           `gorm:"primary_key" json:"id"`
    Alias         string         `gorm:"not null" json:"alias"`
    Name          string         `gorm:"not null" json:"name"`
    Options       pq.StringArray `gorm:"not null;type:varchar(500)[]" json:"options"`
    GoalByPeople  bool           `gorm:"not null" json:"goal_by_people"`
    GoalByAmount  bool           `gorm:"not null" json:"goal_by_amount"`
    EndByGoalGain bool           `gorm:"not null" json:"end_by_goal_gain"`
}

// Donation of project
type Donation struct {
    ID        uint    `gorm:"primary_key" json:"id"`
    User      User    `json:"-"`
    UserID    int     `json:"user"`
    Project   Project `json:"-"`
    ProjectID int     `json:"project"`
    Payment   int     `gorm:"not null; default:0" json:"payment"`
    Locked    bool    `gorm:"not null" json:"locked"`
    Paid      bool    `gorm:"not null" json:"paid"`
}
