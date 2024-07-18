package repository

// Repository определяет контракт для работы с базой данных
type Repository interface {
	SaveMessage(content string) (int, error)
	MarkMessageAsProcessed(id int) error
}

// PostgreSQLRepository реализует интерфейс Repository для работы с PostgreSQL
type PostgreSQLRepository struct {
	// Здесь могут быть добавлены необходимые поля для работы с PostgreSQL
}

func NewPostgreSQLRepository() *PostgreSQLRepository {
	// Инициализация экземпляра PostgreSQLRepository, если нужно
	return &PostgreSQLRepository{}
}

func (r *PostgreSQLRepository) SaveMessage(content string) (int, error) {
	// Логика сохранения сообщения в PostgreSQL
	return 0, nil
}

func (r *PostgreSQLRepository) MarkMessageAsProcessed(id int) error {
	// Логика пометки сообщения как обработанного в PostgreSQL
	return nil
}
