package postgres

import (
	"errors"
	"github.com/go-pg/pg/v10"
	"todo/domain"
)

// This will hold our postgres database pointer.
type UserRepository struct {
	DB *pg.DB
}

// This is where the injection happens, we fill the UserRepository struct with our
// pointer to the postgres database.
func NewUserRepository(DB *pg.DB) *UserRepository {
	return &UserRepository{DB: DB}
}

/** domain/domain.go UserRepository interface implementations. **/
/*****************************************************************/
func (u *UserRepository) GetByEmail(email string) (*domain.User, error) {
	// New User struct is made like this.
	user := new(domain.User)

	// Get the first User with this email.
	err := u.DB.Model(user).Where("email = ?", email).First()

	if err != nil {
		// Check if error err is equal to postgres error ErrNoRows.
		if errors.Is(err, pg.ErrNoRows) {
			// Return our custom error which represent that user is not found.
			return nil, domain.ErrNoResult
		}

		return nil, err
	}

	return user, nil
}

func (u *UserRepository) GetByUsername(username string) (*domain.User, error) {

	user := new(domain.User)
	err := u.DB.Model(user).Where("username = ?", username).First()

	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, domain.ErrNoResult
		}

		return nil, err
	}

	return user, nil
}

func (u *UserRepository) GetById(id int64) (*domain.User, error) {
	user := new(domain.User)
	err := u.DB.Model(user).Where("id = ?", id).First()

	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, domain.ErrNoResult
		}
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) Create(user *domain.User) (*domain.User, error) {
	// Insert user data in the database, returning back all the data.
	// Because user is of type *domain.User (a pointer) we don't need to reassign the
	// user variable as this function will use the pointer to just change the data.
	_, err := u.DB.Model(user).Returning("*").Insert()

	if err != nil {
		return nil, err
	}

	return user, nil
}
/*****************************************************************/
/*****************************************************************/
