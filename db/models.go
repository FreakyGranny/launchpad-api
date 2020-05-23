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
    ID           uint      `gorm:"primary_key" json:"id"`
    CreatedAt    time.Time `json:"-"`
    UpdatedAt    time.Time `json:"-"`

    Username     string    `json:"username"`
    FirstName    string    `json:"first_name"`
    LastName     string    `json:"last_name"`
    Avatar       string    `json:"avatar"`
    Email        string    `gorm:"type:varchar(100)" json:"-"`
    IsStaff      bool      `gorm:"default=false" json:"is_staff"`
    ProjectCount uint      `gorm:"default=0" json:"project_count"`
    SuccessRate  float32   `gorm:"default=0" json:"success_rate"`
}

// Project model
type Project struct {
    ID            uint        `gorm:"primary_key"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
    Title         string      `gorm:"not null; size:100"`
    SubTitle      string      `gorm:"size:100"`
    ReleaseDate   time.Time   `gorm:"not null"`
    EventDate     time.Time
    GoalPeople    uint        `gorm:"not null; default=0"`
    GoalAmount    uint        `gorm:"not null; default=0"`
    Total         uint        `gorm:"not null; default=0"`
    Description   string      `gorm:"size:1000"`
    ImageLink     string      `gorm:"not null; size:200"`
    Instructions  string      `gorm:"size:500"`
    Locked        bool        `gorm:"not null; default=false"`
    Published     bool        `gorm:"not null; default=false"`
    Closed        bool        `gorm:"not null; default=false"`
    Owner         User
    OwnerID       uint
    Category      Category
    CategoryID    uint
    ProjectType   ProjectType
    ProjectTypeID uint
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

// Lock project and donations
func (p *Project) Lock() {
    dbClient := GetDbClient()

    dbClient.Model(&p).Update("locked", true)
    dbClient.Table("donations").Where("project_id = ?", p.ID).Update("locked", true)    
}

// Close project and donations
func (p *Project) Close() {
    dbClient := GetDbClient()

    dbClient.Model(&p).Update("closed", true)
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
    UserID    uint    `json:"user"`
    Project   Project `json:"-"`
    ProjectID uint    `json:"project"`
    Payment   uint    `gorm:"not null; default=0" json:"payment"`
    Locked    bool    `gorm:"not null; default=false" json:"locked"`
    Paid      bool    `gorm:"not null; default=false" json:"paid"`
}
