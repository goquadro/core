package core

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gopkg.in/mgo.v2/bson"
)

/*
type LoginAttempt struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Username   string        `bson:"username"`
	Date       time.Time     `bson:"last_login"`
	Ip         net.IP        `bson:"ip"`
	Successful bool          `bson:"success"`
	Trial      int           `bson:"trial"`
	TotalTrial int           `bson:"total_trial"`
}

func (l *LoginAttempt) Time() time.Time {
	return l.ID.Time()
}

func (l *LoginAttempt) Do() *User {
	return new(User)
}
*/

// Login will perform any actions that are required to make a user model
// officially authenticated.
func (u *User) Login() {
	locSession := getSession()
	defer locSession.Close()
	c := locSession.DB(JobDatabase).C(UsersCollection)
	update := bson.M{"$set": bson.M{"last_login": time.Now(), "fails": 0}}
	err := c.UpdateId(u.ID, update)
	if err != nil {
		log.Panic(err)
	}
}

// Function triggered on failed login. Counts failed attempts.
func (u *User) LoginFailed() {
	locSession := getSession()
	defer locSession.Close()
	c := locSession.DB(JobDatabase).C(UsersCollection)
	update := bson.M{"$inc": bson.M{"fails": 1}}
	c.UpdateId(u.ID, update)
}

/*
// Logout will preform any actions that are required to completely
// logout a user.
func (u *User) Logout() {
	_ = "Do nothing."
}
*/

// Check password against user.
// Returns nil if successful.
func (u *User) CheckPassword() error {
	password := u.EnteredPassword
	locSession := getSession()
	defer locSession.Close()
	c := locSession.DB(JobDatabase).C(UsersCollection)
	err := c.Find(bson.M{"username": u.Username}).One(&u)
	if err != nil {
		return err
	}
	salted := append([]byte(password), u.Salt...)
	err = bcrypt.CompareHashAndPassword(u.Password, salted)
	if err != nil {
		return err
	}
	return nil
}
