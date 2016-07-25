package main

import (
  "errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const CONNECTIONSTRING = "mongodb://127.0.0.1"

type MongoConnection struct {
    originalSession *mgo.Session
}

func NewDBConnection() (conn *MongoConnection) {
	conn = new(MongoConnection)
	conn.createConnection()
	return
}

type mongoDocument struct {
  Id         bson.ObjectId `bson:"_id"`
	ShortUrl   string        `bson:"shorturl"`
	LongUrl    string        `bson:"longurl"`
  ClickCount int           `bson:"click_count"`
}

func (c *MongoConnection) createConnection()(err error){
  fmt.Println("Connecting to local mongo server....")
  c.originalSession, err = mgo.Dial(CONNECTIONSTRING)
  if err == nil {
    fmt.Println("Connection established to mongo server")
    urlCollection := c.originalSession.DB("Shorty").C("shorturls")
    fmt.Println(urlCollection)
    if urlCollection == nil {
			err = errors.New("Collection could not be created, maybe need to create it manually")
		}
    index := mgo.Index{
      Key: []string{"$text:shorturl"},
      Unique: true,
      DropDups: true,
    }
    urlCollection.EnsureIndex(index)
    fmt.Println("made it!")
  } else {
		fmt.Printf("Error occured while creating mongodb connection: %s", err.Error())
  }
  return
}

func (c *MongoConnection) getSessionAndCollection() (session *mgo.Session, urlCollection *mgo.Collection, err error) {
  if c.originalSession != nil {
    session = c.originalSession.Copy()
    urlCollection = session.DB("Shorty").C("shorturls")
  }else {
		err = errors.New("No original session found")
	}
	return
}

func (c *MongoConnection) FindDoc(shorturl string)(doc mongoDocument, session*mgo.Session, urlCollection *mgo.Collection, err error){
  result := mongoDocument{}
  session, urlCollection, err = c.getSessionAndCollection()
  if err != nil {
    return
  }
  err = urlCollection.Find(bson.M{"shorturl":shorturl}).One(&result)
  if err != nil {
    fmt.Println(err)
    return
  }
  return result, session, urlCollection, nil
}

func (c *MongoConnection) FindUrl(shorturl string)(lUrl string, err error){
  result, session, urlCollection, err := c.FindDoc(shorturl)
  defer session.Close()
  if err != nil {
    fmt.Println(err)
    return
  }
  new_clicks := result.ClickCount + 1
  info, err := urlCollection.Upsert(bson.M{"_id": result.Id}, bson.M{"$set": bson.M{"click_count": new_clicks}})
  fmt.Println(info)
  if err != nil {
    fmt.Println(err)
    return
  }
  return result.LongUrl, nil
}

func (c *MongoConnection) FindCount(shorturl string)(clickCount int, err error){
  result, session, _, err := c.FindDoc(shorturl)
  defer session.Close()
  if err != nil {
    fmt.Println(err)
    return
  }
  return result.ClickCount, nil
}

func (c *MongoConnection) InsertUrl(longUrl string, shortUrl string)(sUrl string, err error){
  session, urlCollection, err := c.getSessionAndCollection()
  if err == nil {
    defer session.Close()
    err = urlCollection.Insert(
      &mongoDocument{
        Id: bson.NewObjectId(),
        ShortUrl: shortUrl,
        LongUrl: longUrl,
      },
    )
    if err != nil {
      if mgo.IsDup(err) {
        err = errors.New("Duplicate name exists for the shorturl")
      }
    }
  }
  return
}
