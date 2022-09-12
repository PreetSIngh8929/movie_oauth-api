package rest

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	fmt.Println("about to start test cases....")
	rest.StartMockupServer()
	os.Exit(m.Run())
}
func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"passowrd"}`,
		RespHTTPCode: -1,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "invalid restclient response when trying to login user", err.Message())
}
func TestLoginUserInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"the-password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"invalid login credentials","status":"404","error":"not_found"}`,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "invalid user interface when trying to login user", err.Message())
}
func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"the-password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"invalid login credentials","status":404,"error":"not_found"}`,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status())
	assert.EqualValues(t, "invalid login credentials", err.Message())
}
func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"thepassword"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": "5","first_name": "NEERAJ","last_name": "CHAND","email": "xysssssz123@gmail.com"}`,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "thepassword")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "error when tryin to unmarshal users login response", err.Message())

}
func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com","password":"the-password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id":5 ,"first_name": "NEERAJ","last_name": "CHAND","email": "email@gmail.com"}`,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 5, user.Id)
	assert.EqualValues(t, "NEERAJ", user.FirstName)
	assert.EqualValues(t, "CHAND", user.LastName)
	assert.EqualValues(t, "xysssssz123@gmail.com", user.Email)
}
