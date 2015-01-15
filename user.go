package core

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var InvalidCredentials = errors.New("Invalid username or password.")
var InvalidUsernameError = errors.New("Username not valid.")
var InvalidEmailAddressError = errors.New("Email address not accepted.")
var UsernameAlreadyTakenError = errors.New("Username already taken.")
var UserAlreadyLoggedIn = errors.New("Current user is already logged in.")
var InvalidUidError = errors.New("No user with that ID.")

type User struct {
	ID               bson.ObjectId `bson:"_id,omitempty"    json:"userID"`
	Username         string        `bson:"username"         json:"username"`
	Name             string        `bson:"name"             json:"name"`
	Location         string        `bson:"location"         json:"location"`
	URL              string        `bson:"url"              json:"url"`
	Email            string        `bson:"email"            json:"email"`
	EmailVerified    bool          `bson:"email_verified"   json:"-"`
	IsRegistered     bool          `bson:"is_registered"    json:"-"`
	HasPassword      bool          `bson:"has_password"     json:"-"`
	IsActive         bool          `bson:"is_active"        json:"-"`
	Password         []byte        `bson:"password"         json:"-"`
	Salt             []byte        `bson:"salt"             json:"-"`
	GoogleOAuthSub   string        `bson:"google_oauth_sub" json:"-"`
	LastLogin        time.Time     `bson:"last_login"       json:"-"`
	EnteredPassword  string        `bson:"-"                json:"password"`
	CodeUsed         bson.ObjectId `bson:"signup_code"      json:"-"`
	VerificationCode string        `bson:"confirm_code"     json:"-"`
	Role             int           `bson:"role"             json:"-"`
	FailedLogins     int           `bson:"fails"            json:"-"`
	//ProfileImageUrl         string `json:"profile_image_url"`
	//ProfileImageUrlHttps    string `json:"profile_image_url_https"`
}

// Sync overwrites the provided user object with the information
// stored in the database, using the ID property to find it.
func (u *User) Sync() error {
	locSession := getSession()
	defer locSession.Close()
	return locSession.DB(JobDatabase).C(UsersCollection).FindId(u.ID).One(u)
}

// Returns a pointer to a User object, given its ID in the form of a string.
func GetUserById(uid string) (*User, error) {
	u := new(User)
	if !bson.IsObjectIdHex(uid) {
		return u, InvalidUidError
	}
	u.ID = bson.ObjectIdHex(uid)
	err := u.Sync()
	return u, err
}

// Returns a pointer to a User object, given its username.
func GetUserByName(username string) (*User, error) {
	u := new(User)
	locSession := getSession()
	defer locSession.Close()
	c := locSession.DB(JobDatabase).C(UsersCollection)
	err := c.Find(bson.M{"username": username}).One(u)
	if err != nil {
		log.Println("GetUserByName error:", err) // Debug code
	}
	return u, err
}

// Generates a random urlencoded string of the specified length
func RandomUrlencodedString(length int) string {
	rb := make([]byte, length)
	_, err := rand.Read(rb)
	if err != nil {
		fmt.Println(err)
	}
	return base64.URLEncoding.EncodeToString(rb)
}

func getSalt() []byte {
	b := make([]byte, 60)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

// Extracts the time of creation from the bson.ID of the User
func (u *User) CreatedAt() time.Time {
	return u.ID.Time()
}

// Merges the provided user to the calling User object.
// Chooses the calling user's data upon conflict.
func (u1 *User) MergeAndClean(u2 User) error {
	if u1.Username == u1.ID.Hex() {
		u1.Username = u2.Username
	}
	if u1.Name == "" {
		u1.Name = u2.Name
	}
	if u1.Location == "" {
		u1.Location = u2.Location
	}
	if u1.URL == "" {
		u1.URL = u2.URL
	}
	if u1.Email == "" || (!u1.EmailVerified && u2.EmailVerified) {
		u1.Email = u2.Email
		u1.EmailVerified = u2.EmailVerified
	}
	if !u1.HasPassword && u2.HasPassword {
		u1.HasPassword = true
		u1.Password = u2.Password
	}
	if u1.GoogleOAuthSub == "" {
		u1.GoogleOAuthSub = u2.GoogleOAuthSub
	}
	if u1.LastLogin.Before(u2.LastLogin) {
		u1.LastLogin = u2.LastLogin
	}
	if u1.IsRegistered || u2.IsRegistered {
		u1.IsRegistered = true
	}
	locSession := getSession()
	defer locSession.Close()
	uc := locSession.DB(JobDatabase).C(UsersCollection)
	if bson.IsObjectIdHex(u2.ID.String()) {
		docsFinder := bson.M{"owner": u2.ID}
		change := bson.M{"$set": bson.M{"owner": u1.ID}}
		err := locSession.DB(JobDatabase).C(DocumentsCollection).Update(docsFinder, change)
		if err != nil {
			return err
		}
		err = uc.RemoveId(u2.ID)
		if err != nil {
			return err
		}
	}
	return uc.UpdateId(u1.ID, u1)
}

// Sets the User's email to the provided address, after some checking.
func (u *User) SetEmail(address string) error {
	email, err := mail.ParseAddress(address)
	if err != nil {
		return err
	}
	u.Email = email.Address
	return nil
}

// Checks the username against some basic rules.
func validateUsername(username string) error {
	pattern := "^[a-z0-9_-]{3,20}$" //Include from 3 to 20 letters or numbers.
	username = strings.ToLower(username)
	ok, err := regexp.MatchString(pattern, username)
	if !ok {
		return InvalidUsernameError
	}
	return err
}

// Sets a new username for the user, after performing some basic checking.
// Uniqueness is checked through the unique index in mongo.
func (u *User) SetUsername(username string) error {
	err := validateUsername(username)
	if err != nil {
		return err
	}
	u.Username = username
	return nil
}

// Checks the chosen password against some basic rules.
// Presently, it only checks for the length (8+).
func validatePassword(password string) error {
	if len(password) >= 8 {
		return nil
	}
	return errors.New("Inserted password is too short")
}

// Sets a new password for the calling user.
// Doesn't check for authentication.
func (u *User) SetPassword(password string) error {
	if err := validatePassword(password); err != nil {
		return err
	}
	u.Salt = getSalt()
	hashedPasswd, err := bcrypt.GenerateFromPassword(append([]byte(password), u.Salt...), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = hashedPasswd
	return nil
}

// Registers a new user, provided as an argument, and transfers its properties
// to the calling User object after sanitization.
func (u *User) Register(candidate User) error {
	locSession := getSession()
	defer locSession.Close()
	c := locSession.DB(JobDatabase).C(UsersCollection)
	err := u.SetUsername(candidate.Username)
	if err != nil {
		return err
	}
	err = u.SetEmail(candidate.Email)
	if err != nil {
		return err
	}
	err = u.SetPassword(candidate.EnteredPassword)
	if err != nil {
		return err
	}
	u.IsRegistered = true
	u.IsActive = true
	u.HasPassword = true
	u.VerificationCode = RandomUrlencodedString(35)
	u.LastLogin = time.Now()

	err = c.Insert(u)
	if err != nil {
		return err
	}
	go u.SendConfirmationEmail()
	return c.Find(bson.M{"username": u.Username}).One(u)
}

func (u *User) UniqueId() interface{} {
	return u.ID.Hex()
}

// GetById will populate a user object from a database model with
// a matching id.
func (u *User) GetById(uid string) error {
	if !bson.IsObjectIdHex(uid) {
		return errors.New(fmt.Sprint("User ID not valid:", uid))
	}
	bsonid := bson.ObjectIdHex(uid)
	locSession := getSession()
	defer locSession.Close()
	err := locSession.DB(JobDatabase).C(UsersCollection).FindId(bsonid).One(&u)
	return err
}
