package db

import (
	"time"

	"gorm.io/gorm"
)

// User has many CreditCards, UserID is the foreign key
//  functioneaza codul
// One(User) to many(Nodes)
// un utilizator poate sa trimita mai multe mesaje
// fieacare nod va avea un counter pt nr de sms-uri trimise(pt a putea face round robbin load balancing)

// cheia privata/publica vor fi generate doar la set-up atunci cand placa este setata
// pe placa ramane cheia privata si cea publica, pe server si in baza de date ajunge cheia publica

// ..../api/auth/node
//daca nu exista atunci sa fie oferit utilizatorului un link de autentificare ..../

type User struct {
  gorm.Model
  // UUID string
  Name string
  Email string
  Password string
  Nodes []Node
  Messages []Message
}

type Node struct {
  gorm.Model
  // relatii 
  UserID uint  

  // date
  // ?? la ce se refera number ??
  Number string
  Status string 
  MacAddress string
  ApiKey string
  SentMessages int
  // UUID string
  PublicKey string
  MessageQueue []Message
}

type Message struct {
  gorm.Model
  // relatii
  UserID uint
  NodeID uint
  // date
  Title string
  SentAt time.Time
  Status bool
  Method string
}
