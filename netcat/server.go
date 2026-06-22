package netcat

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

func NewServer(userLimit int, defaultChannel string) *Server {
	admin := &User{
		"ADMIN",
		"",
		nil,
	}
	srv := &Server{
		defaultChannel,
		make(map[string]*User, userLimit+1),
		userLimit,
		make(map[string]*Channel, 1),
		[]Message{},
		admin,
		sync.Mutex{},
	}
	srv.AddUser(admin)
	return srv
}

func (srv *Server) Start(port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	defer listener.Close()

	messages := make(chan Message, 256)
	go srv.MessageSender(messages)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		user := User{
			Channel:    srv.DefaultChannel,
			Connection: conn,
		}
		srv.mu.Lock()
		errChan := make(chan error, 256)
		srv.mu.Unlock()
		go user.HandleUser(srv, messages, errChan)
	}
}

func (user *User) HandleUser(srv *Server, messageChannel chan<- Message, errChan chan<- error) {
	defer user.Connection.Close()
	userLimitMsg := "Can't start connection: server reached user limit\n"
	srv.mu.Lock()
	if len(srv.UserSet) > srv.UserLimit {
		srv.mu.Unlock()
		user.Connection.Write([]byte(userLimitMsg))
		return
	}
	srv.mu.Unlock()
	user.Connection.Write([]byte(WelcomeMessage))
	scanner := bufio.NewScanner(user.Connection)
	for scanner.Scan() {
		user.Name = scanner.Text()
		srv.mu.Lock()
		if len(srv.UserSet) > srv.UserLimit {
			srv.mu.Unlock()
			user.Connection.Write([]byte(userLimitMsg))
			return
		}
		srv.mu.Unlock()
		if trim := strings.TrimSpace(scanner.Text()); trim == "" {
			user.Connection.Write([]byte("Username cannot be empty\n" + "[ENTER YOUR NAME]: "))
			continue
		}
		if !utf8.ValidString(scanner.Text()) ||
			strings.ContainsFunc(scanner.Text(), func(r rune) bool { return !unicode.IsGraphic(r) }) {
			continue
		}
		if srv.AddUser(user) {
			break
		}
		user.Connection.Write([]byte(
			"User named `" + user.Name + "` already exists\n" +
				"[ENTER YOUR NAME]: "))
	}

	messageChannel <- Message{
		srv.Admin,
		time.Now(),
		user.Name + " has joined our chat...",
	}
	err := user.ReceiveHistory(srv)
	if err != nil {
		errChan <- err
	}
	for scanner.Scan() {
		msg := scanner.Text()
		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}

		messageChannel <- Message{
			user,
			time.Now(),
			msg,
		}
	}
	messageChannel <- Message{srv.Admin,
		time.Now(),
		user.Name + " has left our chat..."}

	srv.RemoveUser(user)
}

func (srv *Server) MessageSender(messageChannel <-chan Message) {
	for {
		msg := <-messageChannel
		msg.Message = strings.TrimSpace(msg.Message)
		if msg.Message == "" || !utf8.ValidString(msg.Message) ||
			strings.ContainsFunc(msg.Message, func(r rune) bool {
				return !unicode.IsGraphic(r)
			}) {
			continue
		}
		srv.AddMessage(msg)
		srv.mu.Lock()
		for _, user := range srv.UserSet {
			err := user.ReceiveFromServer(msg, srv)
			if err != nil {
				fmt.Println("Message sender ", err)
			}
		}
		srv.mu.Unlock()
	}
}

func (srv *Server) AddMessage(msg Message) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	srv.History = append(srv.History, msg)
}

func (user *User) ReceiveFromServer(msg Message, srv *Server) error {
	if user.Connection == nil {
		return nil
	}
	if user.Channel == msg.User.Channel || msg.User.Channel == "" {
		var err error
		if msg.User != srv.Admin {
			_, err = fmt.Fprintf(user.Connection,
				"[%s][%s]: %s\n", msg.Time.Format("2006-01-02 15:04:05"), msg.User.Name, msg.Message)
		} else {
			_, err = fmt.Fprintln(user.Connection,
				msg.Message)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (user *User) ReceiveHistory(srv *Server) error {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	for _, msg := range srv.History {
		err := user.ReceiveFromServer(msg, srv)
		if err != nil {
			fmt.Println("Receive from history ", err)
		}
	}
	return nil
}

func (srv *Server) AddUser(user *User) bool {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	user.Name = strings.TrimSpace(user.Name)
	lower := strings.ToLower(user.Name)
	_, ok := srv.UserSet[lower]
	if ok {
		return false
	}
	srv.UserSet[lower] = user
	return true
}

func (srv *Server) RemoveUser(user *User) {
	srv.mu.Lock()
	defer srv.mu.Unlock()

	delete(srv.UserSet, user.Name)
	user.Connection = nil
}
