// ------------------
// User: pei
// DateTime: 2020/10/29 11:20
// Description: 
// ------------------

package ims

import (
	"github.com/guzhi17/gzu"
	"github.com/guzhi17/xcon"
	"google.golang.org/protobuf/proto"
	"log"
	"sync"
	"time"
)


//========================================================================================================================



//========================================================================================================================
//type Request interface {
//	//get query
//	RequestMarshal(proto.Message) error
//	ResponseOK(int, proto.Message) error
//	ResponseError(int, string) error
//}
//========================================================================================================================

type ImRequest struct {
	ses *Session
	pkg Package
}
func (s *ImRequest)RequestMarshal(v proto.Message) error{
	return proto.Unmarshal(s.pkg.Body(), v)
}
func (s *ImRequest)ResponseOK(code int, pkg proto.Message)(err error){
	var b []byte
	b, err = proto.Marshal(pkg)
	if err != nil{
		return
	}
	s.pkg.SetRes(1)
	s.pkg.SetCode(uint32(code))
	_, err = s.ses.Write(s.pkg.Header(), b)
	return
}
func (s *ImRequest)ResponseError(code int, msg string)(err error){
	s.pkg.SetRes(1)
	s.pkg.SetCode(uint32(code))
	_, err = s.ses.Write(s.pkg.Header())
	return
}

//========================================================================================================================
type Handler interface {
	ServeXCon(req *ImRequest)error
}
type HandlerFunc func(req *ImRequest)error
// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeXCon(req *ImRequest)error {
	return f(req)
}
//========================================================================================================================
type HandlerManager struct {
	handlers map[uint32]Handler
}

func CreateHandlerManager() *HandlerManager {
	return &HandlerManager{
		handlers: map[uint32]Handler{},
	}
}

func (h *HandlerManager)RegisterHandler(id uint32, handler Handler)  {
	h.handlers[id] = handler
}
func (h *HandlerManager)RegisterHandlerFun(id uint32, handler HandlerFunc)  {
	h.handlers[id] = handler
}

func (h *HandlerManager)Handle(ses *Session, pkg Package) error {
	if h, ok := h.handlers[pkg.Fid()]; ok{
		return h.ServeXCon(&ImRequest{
			ses: ses,
			pkg: pkg,
		})
	}
	return ErrorPackage{pkg}
}

func (h *HandlerManager)OnConn(c *xcon.Conn)(xcon.Session, error){
	log.Printf("[%p,OnConn]\n", c)
	return h.NewSession(c), nil
}
func (h *HandlerManager)OnClose(c xcon.Session, err error)error{
	switch v := c.(type) {
	default:
		log.Printf("[%p,OnClose[Unknown]%v]\n", c, err)
	case *Session:
		log.Printf("[%p,OnClose]%v\n", v.conn, err)
		v.OnClose(err)
	}
	return nil
}
//========================================================================================================================
type QueryHandler func(pkg Package, err error)error
type QueryManager struct {
	header Package
	idx gzu.AtomUint32
	qid gzu.AtomUint32 //current id
	mu sync.Mutex
	hs map[uint32]QueryHandler
}

func CreateQueryManager() *QueryManager {
	return &QueryManager{
		header: MakeHeader(Header{
			Tag:  0x1615,
			Ver:  1,
			Code: 200,
		}),
		hs: map[uint32]QueryHandler{},
	}
}

func (s *QueryManager) Handle(v Package) error {
	if h, ok := s.take(v.Qid()); ok{
		return h(v, nil)
	}
	return ErrorPackage{v}
}
func (s *QueryManager) set(h QueryHandler) ( id uint32) {
	id = s.qid.Inc()
	s.mu.Lock()
	s.hs[id] = h
	s.mu.Unlock()
	return
}
func (s *QueryManager) take(v uint32) (r QueryHandler, ok bool) {
	s.mu.Lock()
	if r, ok = s.hs[v]; ok{
		delete(s.hs, v)
	}
	s.mu.Unlock()
	return
}
func (s *QueryManager) Reset() ( r map[uint32]QueryHandler ) {
	s.mu.Lock()
	//s.qid.Store(0)
	r, s.hs = s.hs, map[uint32]QueryHandler{}
	s.mu.Unlock()
	return
}
//========================================================================================================================
type Session struct {
	qh *QueryManager
	m *HandlerManager
	conn *xcon.Conn
}

func (h *HandlerManager)NewSession(c *xcon.Conn) *Session {
	return &Session{
		m: h,
		conn:c,
	}
}

//when connection closed, handle and then try?
func (s *Session)OnClose(err error)error{
	//closed?
	if s.qh != nil{
		//reconnect?
		var handlers = s.qh.Reset()
		for _, hi := range handlers{
			hi(nil, err)
		}
	}
	return err
}

func (s *Session) Close() error {
	return s.conn.Close()
}

func (s *Session)Write(bs...[]byte)(n int, err error){
	return s.conn.Write(bs...)
}

//-----------------------------------------------------------------------
func (s *Session)query(id uint32, q []byte, h QueryHandler) (qid uint32, err error) {
	//b, err := proto.Marshal(q)
	//if err != nil{
	//	return 0, err
	//}
	qid = s.qh.set(h)
	defer func() {
		if err != nil{
			s.qh.take(qid)
		}
	}()
	_, err = s.conn.Write(s.qh.header.CloneHeader().SetFid(id).SetQid(qid), q)
	if err != nil{
		return
	}
	return
}

func (s *Session)Query(id uint32, q []byte, h QueryHandler) (err error) {
	_, err = s.query(id, q, h)
	return
}

func (s *Session)QueryTimeout(id uint32, q []byte, t time.Duration) (r Package, err error) {
	type Response struct {
		Package
		Err error
	}
	var wait = make(chan Response)
	defer func() {
		close(wait)
	}()

	var qid uint32
	qid, err = s.query(id, q, func(pkg Package, err error) error {
		defer func() {recover()}()
		wait <- Response{Package: pkg, Err: err}
		return err
	})
	if err != nil{
		return nil, err
	}
	//wait for response
	select {
	case pkg, ok := <- wait:
		if !ok{
			return nil, ErrUnknown
		}
		return pkg.Package, pkg.Err
	case  <- time.After(t):
		s.qh.take(qid)
		return nil, ErrTimeout
	}
}
//-----------------------------------------------------------------------

//-define(IM_SESSION_SIGN, 16#1615).
//-define(HEADER_BODY(Ver, Tp, Res, Qid, Fid, Code, Body),    <<?IM_SESSION_SIGN:16, Ver:16, Tp:8, Res:8, Qid:32, Fid:32, Code:32, Body/binary>>).
func (s *Session)OnData(b []byte)( err error){
	pkg, err := PackageFromBytes(b)
	if err != nil{
		return err
	}
	switch pkg.Res() {
	default:
	case 0:
		if s.m != nil{
			var e1 = s.m.Handle(s, pkg)
			if e1 != nil{
				log.Println(e1)
			}
		}
	case 1:
		if s.qh != nil{
			var e1 = s.qh.Handle(pkg)
			if e1 != nil{
				log.Println(e1)
			}
		}
	}
	return
}
//========================================================================================================================