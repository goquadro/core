package core

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"
)

////////////////////////////////
// SignupCodes are keys meant to grant access to signup either to a person with a code,
// either to a specific email address without providing a code.
////////////////////////////////

type SignupCode struct {
	ID         bson.ObjectId `bson:"_id,omitempty"     json:"-"`
	EmailBound bool          `bson:"is_email_bound"    json:"is_email_bound"`
	Email      string        `bson:"email"             json:"email"`
	Code       string        `bson:"code"              json:"code"`
	Used       time.Time     `bson:"used_at,omitempty" json:"used_at"`
}

// Saves a new code to the database.
func (s *SignupCode) Persist() error {
	locSession := getSession()
	defer locSession.Close()
	c := locSession.DB(gqConfig.jobDatabase).C(SignupCodesCollection)
	return c.Insert(s)
}

// Check whether user is entitled to sign up, calling User.Register if OK
func (u *User) SignupWithCode(code string) error {
	errorCodeNotRecognized := errors.New("Code not recognized")
	signupCodes := []SignupCode{}
	locSession := getSession()
	defer locSession.Close()
	c := locSession.DB(gqConfig.jobDatabase).C(SignupCodesCollection)
	err := c.Find(bson.M{"code": code}).All(&signupCodes)
	if err != nil {
		return err
	}
	for _, sc := range signupCodes {
		// If not used AND (not bound OR bound to the right address)
		if (sc.Used == time.Time{}) && (sc.EmailBound && u.Email == sc.Email || !sc.EmailBound) {
			u.CodeUsed = sc.ID
			err = u.Register(*u)
			if err != nil {
				return err
			}
			sc.Used = time.Now()
			c.UpdateId(sc.ID, sc)
			return nil
		}
	}
	return errorCodeNotRecognized
}
