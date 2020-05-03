package db

import (
  "time"
  "github.com/lib/pq"
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
  Email        string  `json:"email"`
  IsStuff      bool    `json:"is_stuff" gorm:"default:false"`
  ProjectCount int     `json:"project_count" gorm:"default:0"`
  SuccessRate  float32 `json:"success_rate" gorm:"default:0"`
}

// Product test model
type Product struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
  DeletedAt *time.Time `sql:"index" json:"-"`
  
  Code string `json:"code"`
  Price uint  `json:"price"`
}

// Category of project
type Category struct {
	ID        uint   `gorm:"primary_key" json:"id"`  
  Alias     string `gorm:"not null" json:"alias"`
  Name      string `gorm:"not null" json:"name"`
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
