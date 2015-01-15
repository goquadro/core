package core

import (
	"log"
	"os"
	"time"
)

const (
	ISO8601 = "2006-01-02T15:04:05Z"
)

var (
	// config is the struct in which are stored all the credentials
	// used throughout the package.
	gqConfig config
)

type config struct {
	mailgunDomain         string
	mailgunKey            string
	mailgunPubKey         string
	notificationAddress   string
	mongoDBHosts          string
	authDatabase          string
	authUserName          string
	authPassword          string
	jobDatabase           string
	usersCollection       string
	documentsCollection   string
	signupCodesCollection string
	filesCollection       string
}

// Fast error checking
func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func TimeToIso(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}

// Gets the variable from the environment. `def` is the default value
// that gets used if no env is found with that name.
func getenv(varName, def string) string {
	if newVar := os.Getenv(varName); newVar != "" {
		return newVar
	}
	return def
}

func init() {
	gqConfig = config{
		mailgunDomain:       "goquadro.com",
		mailgunKey:          getenv("QDOC_MAILGUN_PRIVATE_KEY", ""),
		mailgunPubKey:       getenv("QDOC_MAILGUN_PUBLIC_KEY", ""),
		notificationAddress: "qdoc <notify@goquadro.com>",
		mongoDBHosts:        getenv("QDOC_MONGO_HOST", "localhost"),
		authDatabase:        getenv("QDOC_MONGO_DB", "qdoc"),
		authUserName:        getenv("QDOC_MONGO_USER", "qdoc1"),
		authPassword:        getenv("QDOC_MONGO_PW", "test"),
		jobDatabase:         getenv("QDOC_MONGO_DB", "qdoc"),
	}
}
