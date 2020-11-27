// ------------------
// User: pei
// DateTime: 2020/10/28 10:26
// Description: 
// ------------------

package main

import (
	"apps/ims"
	"apps/proto/build/go/pb"
	"github.com/guzhi17/gzu"
	"github.com/guzhi17/xcon"
	"log"
	"time"
)

func main() {
	log.SetFlags(11)
	m := ims.CreateHandlerManager()
	m.RegisterHandlerFun(uint32(pb.Fids_Fid_SysEchoQuery), func(req *ims.ImRequest) error {
		var q pb.SysEchoQuery
		err := req.RequestMarshal(&q)
		if err != nil {
			log.Println(err)
			return nil
		}
		req.ResponseOK(200, &pb.SysEchoQueryResponse{
			Word: "hello:" + q.Word,
		})
		return nil
	})
	s := xcon.CreateServer(xcon.Config{
		ConnConfig: xcon.ConnConfig{
			PackageMaxLength: 1 << 10,
			PackageMode:      xcon.Pm32,
			ReadTimeout:      time.Second*100,
			WriteTimeout:     time.Second*10,
		},
		Handler:          m,
		Addr:             ":3721",
	})
	gzu.AppRun(func() {
		err := s.ListenAndServe()
		log.Println(err)
		gzu.AppClose()
	})
}

