package warehouse

import (
	"database/sql"

	"github.com/garycarr/book_club/common"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// Warehouse ...
type Warehouse struct {
	DB     *sql.DB
	logrus *logrus.Logger
}

// NewWarehouse ...
func NewWarehouse(dbConnection string, logger *logrus.Logger) (*Warehouse, error) {
	db, err := sql.Open("postgres", dbConnection)
	defer func() {
		if err != nil {
			db.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Warehouse{
		logrus: logger,
		DB:     db,
	}, nil
}

// CreateUser ...
func (w *Warehouse) CreateUser(rr common.RegisterRequest) (*common.User, error) {
	var id string
	sqlStatement := `INSERT INTO user_data (display_name, password, email)
		VALUES ($1, $2, $3)
		RETURNING id`
	if err := w.DB.QueryRow(sqlStatement, rr.DisplayName, rr.Password, rr.Email).Scan(&id); err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return nil, common.ErrLoginUserAlreadyExists
		}
		return nil, err
	}
	return &common.User{
		ID:          id,
		Email:       rr.Email,
		DisplayName: rr.DisplayName,
	}, nil
}

// GetUserWithEmail ...
func (w *Warehouse) GetUserWithEmail(email string) (*common.User, error) {
	u := common.User{}
	sqlStatement := `SELECT id, email, password, display_name
		FROM user_data
		WHERE email = $1`
	err := w.DB.QueryRow(sqlStatement, email).Scan(&u.ID, &u.Email, &u.Password, &u.DisplayName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrLoginUserNotFound
		}
		return nil, err
	}
	return &u, nil
}
