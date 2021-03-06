package core

import (
	"log"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	UsersCollection       = "users"
	DocumentsCollection   = "docs"
	SignupCodesCollection = "signupcodes"
	FilesCollection       = "files"
)

var mgoSession *mgo.Session

// Gets a new session clone.
func getSession() *mgo.Session {
	return mgoSession.Copy()
}

func init() {
	var err error

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{gqConfig.mongoDBHosts},
		Timeout:  60 * time.Second,
		Database: gqConfig.authDatabase,
		Username: gqConfig.authUserName,
		Password: gqConfig.authPassword,
	}

	mgoSession, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	mgoSession.SetMode(mgo.Monotonic, true)

	usernameIndex := mgo.Index{
		Key:    []string{"username"},
		Unique: true,
	}
	tagsIndex := mgo.Index{
		Key:    []string{"owner", "tags"},
		Unique: false,
		Sparse: false,
	}
	signupCodesIndex := mgo.Index{
		Key:    []string{"code"},
		Unique: false,
	}
	err = mgoSession.DB(gqConfig.jobDatabase).C(UsersCollection).EnsureIndex(usernameIndex)
	if err != nil {
		log.Fatal("Error creating users index:", err)
	}
	err = mgoSession.DB(gqConfig.jobDatabase).C(DocumentsCollection).EnsureIndex(tagsIndex)
	if err != nil {
		log.Fatal("Error creating documents index:", err)
	}
	err = mgoSession.DB(gqConfig.jobDatabase).C(SignupCodesCollection).EnsureIndex(signupCodesIndex)
	if err != nil {
		log.Fatal("Error creating signup codes index:", err)
	}
}
