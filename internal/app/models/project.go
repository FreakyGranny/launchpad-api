package models

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

//go:generate mockgen -source=$GOFILE -destination=../mocks/model_project_mock.go -package=mocks .

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

	userProjectCountLimit = 20
)

// ProjectImpl ...
type ProjectImpl interface {
	Get(id int) (*Project, bool)
	GetProjectsWithPagination(f *ProjectListFilter) (ProjectPaginatorImpl, error)
	GetUserProjects(f *ProjectUserFilter) (*[]Project, error)
	GetActiveProjects() (*[]Project, error)
	Create(p *Project) error
	Update(p *Project) error
	DropEventDate(p *Project) error
	Delete(p *Project) error
	UpdateTotalByPayment(p *Project) error
	UpdateTotalByCount(p *Project) error
	Lock(p *Project) error
	Close(p *Project) error
	CheckForPaid(projectID int) (bool, error)
	SetEqualDonation(p *Project) error	
}

// Project model
type Project struct {
	tableName     struct{} `pg:"projects,alias:p"` //nolint
	ID            int
	Title         string
	SubTitle      string
	ReleaseDate   time.Time
	EventDate     time.Time
	GoalPeople    int `pg:",notnull"`
	GoalAmount    int `pg:",notnull"`
	Total         int
	Description   string
	ImageLink     string
	Instructions  string
	Locked        bool `pg:",notnull"`
	Published     bool `pg:",notnull"`
	Closed        bool `pg:",notnull"`
	Owner         User
	OwnerID       int
	Category      Category
	CategoryID    int
	ProjectType   ProjectType
	ProjectTypeID int
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

// ProjectListFilter ...
type ProjectListFilter struct {
	Category    int
	ProjectType int
	OnlyOpen    bool
	PageSize    int
	Page        int
}

// ProjectUserFilter ...
type ProjectUserFilter struct {
	UserID      int
	Contributed bool
	Owned       bool
}

// ProjectPaginatorImpl ...
type ProjectPaginatorImpl interface {
	NextPage() (int, bool)
	Retrieve() (*[]Project, error)
}

// ProjectPaginator ...
type ProjectPaginator struct {
	EntryCount int
	PageSize   int
	Query      *orm.Query
	Page       int
	Values     *[]Project
}

// NextPage ...
func (pp *ProjectPaginator) NextPage() (int, bool) {
	if pp.Page*pp.PageSize < pp.EntryCount {
		return pp.Page + 1, true
	}
	return 0, false
}

// Retrieve ...
func (pp *ProjectPaginator) Retrieve() (*[]Project, error) {
	var ofst int

	if pp.Page > 1 {
		ofst = pp.PageSize * (pp.Page - 1)
	}
	err := pp.Query.Offset(ofst).Limit(pp.PageSize).Order("id DESC").Select()
	if err != nil {
		return nil, err
	}

	return pp.Values, nil
}

// ProjectRepo ...
type ProjectRepo struct {
	db *pg.DB
}

// NewProjectModel ...
func NewProjectModel(db *pg.DB) *ProjectRepo {
	return &ProjectRepo{
		db: db,
	}
}

// Get ...
func (r *ProjectRepo) Get(id int) (*Project, bool) {
	project := &Project{}
	err := r.db.Model(project).Relation("Owner").Relation("Category").Relation("ProjectType").Where("p.id = ?", id).Select()
	if err != nil {
		return project, false
	}

	return project, true
}

// GetActiveProjects returns projects on search or harvest stage
func (r *ProjectRepo) GetActiveProjects() (*[]Project, error) {
	projects := &[]Project{}
	err := r.db.Model(projects).
		Relation("ProjectType").
		Where("p.published = ?", true).
		Where("p.closed = ?", false).
		Select()

	return projects, err
}

// GetProjectsWithPagination ...
func (r *ProjectRepo) GetProjectsWithPagination(f *ProjectListFilter) (ProjectPaginatorImpl, error) {
	projects := []Project{}
	q := r.db.Model(&projects).Relation("Category").Relation("ProjectType").Where("p.published = ?", true)
	if f.Category != 0 {
		q = q.Where("category_id = ?", f.Category)
	}
	if f.ProjectType != 0 {
		q = q.Where("project_type_id = ?", f.ProjectType)
	}
	if f.OnlyOpen {
		q = q.Where("closed = ?", false)
	}
	x, err := q.Count()
	if err != nil {
		return nil, err
	}
	return &ProjectPaginator{
		EntryCount: x,
		Query:      q,
		Values:     &projects,
		Page:       f.Page,
		PageSize:   f.PageSize,
	}, nil
}

// GetUserProjects ...
func (r *ProjectRepo) GetUserProjects(f *ProjectUserFilter) (*[]Project, error) {
	projects := []Project{}
	q := r.db.Model(&projects).Relation("Category").Relation("ProjectType").Where("owner_id = ?", f.UserID)
	if !f.Owned {
		q = q.Where("p.published = ?", true)
	}
	if f.Contributed {
		cProjects := r.db.Model((*Donation)(nil)).ColumnExpr("project_id").Where("user_id = ?", f.UserID)
		q = q.Where("p.id IN (?)", cProjects)
	}
	err := q.Limit(userProjectCountLimit).Order("id DESC").Select()
	if err != nil {
		return nil, err
	}

	return &projects, nil
}

// Create new project
func (r *ProjectRepo) Create(p *Project) error {
	count, err := r.db.Model((*User)(nil)).Where("u.id = ?", p.OwnerID).Count()
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrUserNotFound
	}
	_, err = r.db.Model(p).Insert()

	return err
}

// Update project
func (r *ProjectRepo) Update(p *Project) error {
	_, err := r.db.Model(p).WherePK().UpdateNotZero()

	return err
}

// DropEventDate set event_date to null value
func (r *ProjectRepo) DropEventDate(p *Project) error {
	_, err := r.db.Model(p).Set("event_date = null").WherePK().Update()

	return err
}

// Delete project by id
func (r *ProjectRepo) Delete(p *Project) error {
	_, err := r.db.Model(p).WherePK().Delete()

	return err
}

// donationSum returns sum of donations payment for given project
func (r *ProjectRepo) donationSum(id int) (int, error) {
	sum := 0
	err := r.db.Model((*Donation)(nil)).
		ColumnExpr("sum(d.payment)").
		Where("d.project_id = ?", id).
		Select(&sum)
	if err != nil {
		return sum, err
	}

	return sum, nil
}

// donationCount returns count of donations for given project
func (r *ProjectRepo) donationCount(id int) (int, error) {
	return r.db.Model((*Donation)(nil)).Where("d.project_id = ?", id).Count()
}

func (r *ProjectRepo) saveTotal(p *Project, value int) error {
	p.Total = value
	_, err := r.db.Model(p).Set("total = ?total").WherePK().Update()

	return err

}

// UpdateTotalByPayment ...
func (r *ProjectRepo) UpdateTotalByPayment(p *Project) error {
	sum, err := r.donationSum(p.ID)
	if err != nil {
		return err
	}

	return r.saveTotal(p, sum)
}

// UpdateTotalByCount ...
func (r *ProjectRepo) UpdateTotalByCount(p *Project) error {
	count, err := r.donationCount(p.ID)
	if err != nil {
		return err
	}

	return r.saveTotal(p, count)
}

// Lock project with associated donations
func (r *ProjectRepo) Lock(p *Project) error {
	_, err := r.db.Model(p).Set("locked = TRUE").WherePK().Update()
	if err != nil {
		return err
	}
	_, err = r.db.Model((*Donation)(nil)).
		Set("locked = TRUE").
		Where("d.project_id = ?", p.ID).
		Update()

	return err
}

// Close project with associated donations
func (r *ProjectRepo) Close(p *Project) error {
	_, err := r.db.Model(p).Set("locked = TRUE").WherePK().Update()
	if err != nil {
		return err
	}
	if !p.Locked {
		_, err = r.db.Model((*Donation)(nil)).
			Set("locked = TRUE").
			Where("d.project_id = ?", p.ID).
			Update()
	}

	return err
}

// CheckForPaid checks if all donations are paid
func (r *ProjectRepo) CheckForPaid(projectID int) (bool, error) {
	allPaid := false
	err := r.db.Model((*Donation)(nil)).
		ColumnExpr("bool_and(d.paid)").
		Where("d.project_id = ?", projectID).
		Select(&allPaid)
	if err != nil {
		return allPaid, err
	}

	return allPaid, nil
}

// SetEqualDonation checks if all donations are paid
func (r *ProjectRepo) SetEqualDonation(p *Project) error {
	_, err := r.db.Model((*Donation)(nil)).
		Set("payment = ?", p.GoalAmount / p.GoalPeople).
		Where("d.project_id = ?", p.ID).
		Update()

	return err
}
