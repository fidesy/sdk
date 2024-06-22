package postgres

type Model interface {
	TableName() string
}
