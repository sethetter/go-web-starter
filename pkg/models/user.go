package models

import (
	"math/rand"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	valid "github.com/asaskevich/govalidator"
	"github.com/jmoiron/sqlx"
)

type User struct {
	NewUser
	ID        string    `db:"id" valid:"uuidv4"`
	CreatedAt time.Time `db:"created_at" valid:"time"`
	UpdatedAt time.Time `db:"updated_at" valid:"time"`
}

type NewUser struct {
	Email string `db:"email" valid:"email"`
}

func (u NewUser) Validate() (bool, []error) {
	if pass, err := valid.ValidateStruct(u); !pass {
		return false, err.(valid.Errors).Errors()
	}
	return true, []error{}
}

func (u NewUser) InsertOrGet(db *sqlx.DB) (User, error) {
	var user User

	// TODO: can I make this work without doing the fake email update?
	// currently, it causes updated_at to change needlessly
	q := psql().Insert("users").Columns("email").Values(u.Email).
		Suffix(`ON CONFLICT("email") DO UPDATE SET email = $1 RETURNING *`)

	query, args, err := q.ToSql()
	if err != nil {
		return user, err
	}

	if err = db.Get(&user, query, args...); err != nil {
		return user, err
	}

	return user, nil
}

type AuthChallenge struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Code      string    `db:"code"`
	Expires   time.Time `db:"expires"`
	CreatedAt time.Time `db:"created_at"`
}

func (ac AuthChallenge) Check(provided string) bool {
	return strings.ToLower(provided) == strings.ToLower(ac.Code)
}

func NewAuthChallenge(u User, db *sqlx.DB) (AuthChallenge, error) {
	var authChallenge AuthChallenge

	q := psql().Insert("auth_challenges").
		Columns("user_id", "code", "expires").
		Values(u.ID, generateCode(), time.Now().Add(10*time.Minute)).
		Suffix(`ON CONFLICT("user_id") DO UPDATE SET (code, expires) = ($2, $3) RETURNING *`)

	query, args, err := q.ToSql()
	if err != nil {
		return authChallenge, err
	}

	if err = db.Get(&authChallenge, query, args...); err != nil {
		return authChallenge, err
	}

	return authChallenge, nil
}

func AuthChallengeForUserID(userId string, db *sqlx.DB) (AuthChallenge, error) {
	var authChallenge AuthChallenge
	q := psql().Select("*").From("auth_challenges").
		Where(sq.Eq{"user_id": userId})

	query, args, err := q.ToSql()
	if err != nil {
		return authChallenge, err
	}

	if err = db.Get(&authChallenge, query, args...); err != nil {
		return authChallenge, err
	}

	return authChallenge, nil
}

func generateCode() string {
	charset := "0123456789"

	var str strings.Builder
	for i := 0; i < 6; i++ {
		c := rand.Intn(len(charset))
		str.WriteString(string(charset[c]))
	}

	return str.String()
}
