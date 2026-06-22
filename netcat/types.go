package netcat

import (
	_ "embed"
	"net"
	"sync"
	"time"
)

type User struct {
	Name       string
	Channel    string
	Connection net.Conn
}

type Message struct {
	User    *User
	Time    time.Time
	Message string
}

type Channel struct {
}

type Server struct {
	DefaultChannel string
	UserSet        map[string]*User
	UserLimit      int
	Channels       map[string]*Channel
	History        []Message
	Admin          *User
	mu             sync.Mutex
}

//go:embed welcome.txt
var WelcomeMessage string
