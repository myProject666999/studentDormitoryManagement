package models

import (
	"student-dormitory-management/utils"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	Role      string         `json:"role" gorm:"size:20;not null"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") {
		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}
	return nil
}

type Admin struct {
	User
	Name  string `json:"name" gorm:"size:50"`
	Phone string `json:"phone" gorm:"size:20"`
	Email string `json:"email" gorm:"size:100"`
	Avatar string `json:"avatar" gorm:"size:255"`
}

func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	a.User.Role = "admin"
	return a.User.BeforeCreate(tx)
}

type Student struct {
	User
	StudentNumber string     `json:"student_number" gorm:"uniqueIndex;size:50;not null"`
	Name          string     `json:"name" gorm:"size:50;not null"`
	Gender        string     `json:"gender" gorm:"size:10"`
	Phone         string     `json:"phone" gorm:"size:20"`
	Email         string     `json:"email" gorm:"size:100"`
	Class         string     `json:"class" gorm:"size:50"`
	Major         string     `json:"major" gorm:"size:100"`
	Avatar        string     `json:"avatar" gorm:"size:255"`
	DormitoryID   *uint      `json:"dormitory_id" gorm:"index"`
	Dormitory     *Dormitory `json:"dormitory" gorm:"foreignKey:DormitoryID"`
}

func (s *Student) BeforeCreate(tx *gorm.DB) error {
	s.User.Role = "student"
	return s.User.BeforeCreate(tx)
}

type MaintenanceWorker struct {
	User
	WorkerNumber string `json:"worker_number" gorm:"uniqueIndex;size:50;not null"`
	Name         string `json:"name" gorm:"size:50;not null"`
	Gender       string `json:"gender" gorm:"size:10"`
	Phone        string `json:"phone" gorm:"size:20"`
	Email        string `json:"email" gorm:"size:100"`
	Specialty    string `json:"specialty" gorm:"size:100"`
	Avatar       string `json:"avatar" gorm:"size:255"`
}

func (m *MaintenanceWorker) BeforeCreate(tx *gorm.DB) error {
	m.User.Role = "worker"
	return m.User.BeforeCreate(tx)
}

type Dormitory struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Building     string         `json:"building" gorm:"size:50;not null"`
	Floor        int            `json:"floor"`
	RoomNumber   string         `json:"room_number" gorm:"size:20;not null"`
	RoomType     string         `json:"room_type" gorm:"size:50"`
	TotalBeds    int            `json:"total_beds" gorm:"default:4"`
	UsedBeds     int            `json:"used_beds" gorm:"default:0"`
	Gender       string         `json:"gender" gorm:"size:10"`
	Status       int            `json:"status" gorm:"default:1"`
	Image        string         `json:"image" gorm:"size:255"`
	Description  string         `json:"description" gorm:"type:text"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type DormitoryAssignment struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	StudentID   uint           `json:"student_id" gorm:"uniqueIndex;not null"`
	Student     *Student       `json:"student" gorm:"foreignKey:StudentID"`
	DormitoryID uint           `json:"dormitory_id" gorm:"not null"`
	Dormitory   *Dormitory     `json:"dormitory" gorm:"foreignKey:DormitoryID"`
	AssignmentDate time.Time   `json:"assignment_date"`
	LeaveDate    *time.Time    `json:"leave_date"`
	Status       int            `json:"status" gorm:"default:1"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type Notice struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:200;not null"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	Image       string         `json:"image" gorm:"size:255"`
	AuthorID    uint           `json:"author_id"`
	AuthorName  string         `json:"author_name" gorm:"size:50"`
	Status      int            `json:"status" gorm:"default:1"`
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type RepairRequest struct {
	ID              uint                 `json:"id" gorm:"primaryKey"`
	Title           string               `json:"title" gorm:"size:200;not null"`
	Description     string               `json:"description" gorm:"type:text;not null"`
	Attachment      string               `json:"attachment" gorm:"size:255"`
	StudentID       uint                 `json:"student_id"`
	Student         *Student             `json:"student" gorm:"foreignKey:StudentID"`
	DormitoryID     uint                 `json:"dormitory_id"`
	Dormitory       *Dormitory           `json:"dormitory" gorm:"foreignKey:DormitoryID"`
	WorkerID        *uint                `json:"worker_id"`
	Worker          *MaintenanceWorker   `json:"worker" gorm:"foreignKey:WorkerID"`
	Status          string               `json:"status" gorm:"size:20;default:'pending'"`
	Priority        string               `json:"priority" gorm:"size:20;default:'normal'"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	DeletedAt       gorm.DeletedAt       `json:"-" gorm:"index"`
}

type RepairRecord struct {
	ID              uint                 `json:"id" gorm:"primaryKey"`
	RepairRequestID uint                `json:"repair_request_id" gorm:"not null"`
	RepairRequest   *RepairRequest       `json:"repair_request" gorm:"foreignKey:RepairRequestID"`
	WorkerID        uint                 `json:"worker_id"`
	Worker          *MaintenanceWorker   `json:"worker" gorm:"foreignKey:WorkerID"`
	Description     string               `json:"description" gorm:"type:text"`
	Status          string               `json:"status" gorm:"size:20;default:'in_progress'"`
	StartTime       *time.Time           `json:"start_time"`
	EndTime         *time.Time           `json:"end_time"`
	Cost            float64              `json:"cost" gorm:"default:0"`
	Remark          string               `json:"remark" gorm:"type:text"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	DeletedAt       gorm.DeletedAt       `json:"-" gorm:"index"`
}
