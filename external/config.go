package external

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"os"
)

type Host struct {
	ConnString string
	ServerPort string
}

func NewHost(connString string, serverPort string) *Host {
	return &Host{ConnString: connString, ServerPort: serverPort}
}

// HostProvider Provide the dependency of Host
func HostProvider() (*Host, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	connString := os.Getenv("CONNSTRING")
	serverPort := os.Getenv("SERVERPORT")
	return NewHost(connString, serverPort), nil
}
