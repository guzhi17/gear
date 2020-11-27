// ------------------
// User: pei
// DateTime: 2020/11/6 15:23
// Description: 
// ------------------

package ims

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/guzhi17/gzu"
	"github.com/guzhi17/xcon"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"sync"
	"time"
)

const(
	StatusNone = 0
	StatusConnecting = 1
	StatusConnected = 2
	//StatusShutdown = 3
)

var(
	ErrorStatus = errors.New("error status")
	ErrorNoSession = errors.New("error no session")
	ErrorClosed = errors.New("error closed")
)

type Client struct {
	TlsConfig *tls.Config
	Host string
	xcon.ConnConfig

	status gzu.AtomUint32

	closed gzu.AtomUint32
	mu sync.RWMutex
	ses *Session //the connection
}

func (c *Client)Close() (err error) {
	var v = c.closed.Inc()
	if v > 1{
		return nil
	}
	ses := c.getSession()
	if ses == nil{
		return nil
	}
	err = ses.Close()
	return
}
func (c *Client)Closed() bool {
	return c.closed.Load() > 0
}

func (c *Client)setSession(ses *Session)  {
	var so *Session
	c.mu.Lock()
	so = c.ses
	c.ses = nil
	if !c.Closed(){
		c.ses = ses
	}
	c.mu.Unlock()

	//clear old
	if so != nil{
		so.Close()
	}
	//if closed, clear new
	if ses!=nil && c.Closed() {
		ses.Close()
	}
}
func (c *Client)getSession() (ses *Session) {
	c.mu.RLock()
	ses = c.ses
	c.mu.RUnlock()
	return
}

func (c *Client) Dial() error {
	if !c.status.CAS(StatusNone, StatusConnecting){
		return ErrorStatus
	}
	var (
		con net.Conn
		err error
	)
	//only StatusConnecting goes here
	if c.TlsConfig != nil{
		con, err = tls.Dial("tcp", c.Host, c.TlsConfig)
	}else{
		con, err = net.Dial("tcp", c.Host)
	}
	if err != nil{
		//dial failed, set status to StatusNone
		c.status.Store(StatusNone)
		return err
	}
	var conn = xcon.NewConn(con, c.ConnConfig)
	var s = &Session{
		qh: CreateQueryManager(),
		conn: conn,
	}
	//connected
	c.status.Store(StatusConnected)
	go c.serve(s)

	c.setSession(s)
	return nil
}

func (c *Client)serve(s *Session)  {
	//only StatusConnected goes here
	var err error
	defer func() {
		c.status.CAS(StatusConnected, StatusNone)
		//c.status.Store(StatusNone)
		c.setSession(nil)
		s.OnClose(err)
		//do reconnect ?
		var retry = time.Millisecond*100
		for{
			time.Sleep(retry)
			if c.Closed(){
				return
			}
			log.Println("Client retry to Dial...")
			err = c.Dial()
			if err == nil{
				//connected, exit this one
				return
			}
			if err == ErrorStatus{
				log.Println("other session is connected, exit this one")
				return
			}
			if retry < time.Second*10{
				retry = retry * 2
			}
			log.Println(err, " retry later ", retry)
		}
	}()

	if c.Closed(){
		s.Close()
		err = ErrorClosed
		return
	}
	//todo use a context to cancel sub sessions?
	err = s.conn.Serve(s, context.Background())
}

func (c *Client)Query(id uint32, q proto.Message, h QueryHandler) (err error) {
	b, err := proto.Marshal(q)
	if err != nil{
		return err
	}
	if c.status.Load() != StatusConnected{
		return ErrorStatus
	}
	ses := c.getSession()
	if ses == nil{
		return ErrorNoSession
	}
	err = ses.Query(id, b, h)
	return
}

func (c *Client)QueryTimeout(id uint32, q proto.Message, t time.Duration) (r Package, err error) {
	b, err := proto.Marshal(q)
	if err != nil{
		return nil, err
	}
	if c.status.Load() != StatusConnected{
		return nil, ErrorStatus
	}
	ses := c.getSession()
	if ses == nil{
		return nil, ErrorNoSession
	}
	r, err = ses.QueryTimeout(id, b, t)
	return
}