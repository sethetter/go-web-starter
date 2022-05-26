package models

import (
	sq "github.com/Masterminds/squirrel"
	valid "github.com/asaskevich/govalidator"
)

func init() {
	valid.SetFieldsRequiredByDefault(true)
}

func psql() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
