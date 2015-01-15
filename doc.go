package core

import (
	"errors"
	"net/url"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var InvalidBsonIdError = errors.New("Provided an invalid bson.id object")

type Document struct {
	ID           bson.ObjectId   `bson:"_id"      json:"docID"`
	Owner        bson.ObjectId   `bson:"user"     json:"owner"`
	LinkedFile   string          `bson:"furl"     json:"-"`
	Url          string          `bson:"url"      json:"url"`
	Title        string          `bson:"title"    json:"title"`
	Children     []bson.ObjectId `bson:"children" json:"children"`
	Tags         []string        `bson:"tag"      json:"tags"`
	Thumb        string          `bson:"thumb"    json:"thumbnailUrl"`
	ThumbMobile  string          `bson:"thumbm"   json:"thumbnailUrl"`
	FavIconUrl   string          `bson:"iconurl"  json:"favIconUrl"`
	Parents      []bson.ObjectId `bson:"-"        json:"-"`
	LastModified time.Time       `bson:"lastmod"  json:"lastModified"`
	ToBeDeleted  bool            `bson:"rm"       json:"-"`
	EnteredId    string          `bson:"-"        json:"-"`
}

// Returns the time of creation, extracting the information from the bsonID
func (d *Document) CreatedAt() time.Time {
	return d.ID.Time()
}

// Get a user's documents from the proper collection
func (u User) Documents() (*[]Document, error) {
	docs := []Document{}
	locSession := getSession()
	defer locSession.Close()
	err := locSession.DB(gqConfig.jobDatabase).C(DocumentsCollection).Find(bson.M{"user": u.ID}).All(&docs)
	return &docs, err
}

/*
// This should be done directly through a mgo query.
func (u User) DocumentsByTopic() (*map[string][]Document, error) {
	docsByTopic := make(map[string][]Document)
	docs, err := u.Documents()
	if err != nil {
		return nil, err
	}
	for _, doc := range *docs {
		orphan := true
		for _, topic := range doc.Topics {
			docsByTopic[topic] = append(docsByTopic[topic], doc)
			orphan = false
		}
		if orphan {
			docsByTopic["uncategorized"] = append(docsByTopic["uncategorized"], doc)
		}
	}
	return &docsByTopic, nil
}
*/

// Returns zero doc & err==nil if not found
func (user *User) GetDocumentById(id string) (*Document, error) {
	doc := Document{}
	if !bson.IsObjectIdHex(id) {
		return &doc, InvalidBsonIdError
	}
	bsonId := bson.ObjectIdHex(id)
	locSession := getSession()
	defer locSession.Close()
	docFinder := bson.M{"user": user.ID, "_id": bsonId}
	err := locSession.DB(gqConfig.jobDatabase).C(DocumentsCollection).Find(docFinder).One(&doc)
	return &doc, err
}

// sanitizeUrl performs basic url checking. To be worked on.
func sanitizeUrl(docurl string) (string, error) {
	_, err := url.Parse(docurl)
	if err != nil {
		return "", err
	}
	return docurl, nil
}

// User.AddDocument persists a document belonging to the acting user.
func (u *User) AddDocument(doc *Document) error {
	doc.ID = bson.NewObjectId()
	doc.Owner = u.ID
	doc.Url, _ = sanitizeUrl(doc.Url)
	//doc.Name
	//doc.Children
	//doc.Tags
	//doc.Topics
	//doc.Thumb
	doc.LastModified = doc.CreatedAt()

	locSession := getSession()
	defer locSession.Close()
	c := locSession.DB(gqConfig.jobDatabase).C(DocumentsCollection)
	err := c.Insert(doc)

	if err == nil {
		for _, parentId := range doc.Parents {
			parentDoc := Document{
				ID:    parentId,
				Owner: u.ID,
			}
			parentDoc.AddChild(doc)
		}
	}
	return err
}

// AddChild adds a document's ID to the list of children of its parent.
// Doesn't check for document owner consistency against any specific user.
func (d Document) AddChild(child *Document) error {
	locSession := getSession()
	defer locSession.Close()
	c := locSession.DB(gqConfig.jobDatabase).C(DocumentsCollection)
	docFinder := bson.M{"_id": d.ID, "user": d.Owner}
	err := c.Find(docFinder).One(d)
	if err != nil {
		return err
	}
	d.Children = append(d.Children, child.ID)
	change := bson.M{"$set": bson.M{"children": d.Children, "last_modified": time.Now()}}
	return c.Update(docFinder, change)
}

/*
func (u *User) GetDocsByTags(tags []string) (*[]Document, error) {
	docs := []Document{}
	locSession := getSession()
	defer locSession.Close()
	c := locSession.DB(gqConfig.jobDatabase).C(DocumentsCollection)
	docFinder := bson.M{"user": u.ID, "tags": tags}
	err := c.Find(docFinder).All(&docs)
	return &docs, err
}
*/

// Change a document's ownerID to a selected userID.
// Doesn't perform any auth check.
func (d *Document) ChangeOwner(u *User) error {
	locSession := getSession()
	defer locSession.Close()
	c := locSession.DB(gqConfig.jobDatabase).C(DocumentsCollection)
	docFinder := bson.M{"_id": d.ID, "user": d.Owner}
	change := bson.M{"$set": bson.M{"user": u.ID, "last_modified": time.Now()}}
	return c.Update(docFinder, change)
}

// ToDo: Change behaviour to just mark a doc for deletion, in order to permit undo.
// Some worker can then actually delete them after some time.
func (u *User) DeleteDocument(d *Document) error {
	candidateDoc, err := u.GetDocumentById(d.EnteredId)
	if err != nil {
		candidateDoc = d
	}
	locSession := getSession()
	defer locSession.Close()
	return locSession.DB(gqConfig.jobDatabase).C(DocumentsCollection).Remove(bson.M{"_id": candidateDoc.ID, "user": u.ID})
}

// User.PutDocument is a PUT (full overwrite) scheme document modifier
func (u *User) PutDocument(d *Document) error {
	locSession := getSession()
	defer locSession.Close()
	c := locSession.DB(gqConfig.jobDatabase).C(DocumentsCollection)
	selector := bson.M{"_id": d.ID, "user": u.ID}
	return c.Update(selector, d)
}

// Document.NamePreview aims to provide a human-friendly name for the document,
// limited to the given length (int)
func (d *Document) NamePreview(length int) string {
	newName := ""
	if d.Title != "" {
		newName = d.Title
	} else {
		newName = d.Url
	}
	if length > 0 && len(newName) > length {
		return newName[:length] + "..."
	}
	return newName
}
