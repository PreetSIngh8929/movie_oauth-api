package db

import (
	"errors"

	"github.com/PreetSIngh8929/movie_oauth-api/src/clients/cassandra"
	"github.com/PreetSIngh8929/movie_oauth-api/src/domain/access_token"
	"github.com/PreetSIngh8929/movie_utils-go/rest_errors"
	"github.com/gocql/gocql"
)

const (
	queryGetACcessToken    = "SELECT access_token , user_id, client_id,expires FROM access_tokens WHERE access_token=?;"
	queryCreateAccessToken = "INSERT INTO access_tokens (access_token , user_id, client_id,expires) VALUES (?, ?, ?,?);"
	queryUpdateExpires     = "UPDATE access_tokens SET expires = ? WHERE access_token=?;"
)

func NewRepository() DbRepository {
	return &dbRepository{}
}

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(access_token.AccessToken) rest_errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RestErr
}
type dbRepository struct {
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, rest_errors.RestErr) {
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetACcessToken, id).Scan(&result.AccessToken, &result.UserId, &result.ClientId, &result.Expires); err != nil {
		if err == gocql.ErrNotFound {
			return nil, rest_errors.NewNotFoundError("no access_token found with given id")
		}
		return nil, rest_errors.NewInternalServerError(err.Error(), errors.New("database error"))
	}
	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) rest_errors.RestErr {

	if err := cassandra.GetSession().Query(queryCreateAccessToken,
		at.AccessToken,
		at.UserId,
		at.ClientId,
		at.Expires,
	).Exec(); err != nil {
		return rest_errors.NewInternalServerError(err.Error(), errors.New("database error"))
	}
	return nil
}
func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) rest_errors.RestErr {

	if err := cassandra.GetSession().Query(queryUpdateExpires,
		at.Expires,
		at.AccessToken,
	).Exec(); err != nil {
		return rest_errors.NewInternalServerError(err.Error(), errors.New("database error"))
	}
	return nil
}
