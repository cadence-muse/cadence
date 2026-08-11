package postgresql

import (
	"fmt"
	"net/url"
)

type DSN struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func (D DSN) String() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(D.User, D.Password),
		Host:   fmt.Sprintf("%s:%d", D.Host, D.Port),
		Path:   "/" + D.Database,
	}
	return u.String()
}
