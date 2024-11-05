package users

type UserRepository interface {
}

type PostgresUserRepository struct{}

func NewPostgresUserRepository() *PostgresUserRepository {
	return &PostgresUserRepository{}
}

