// ------------------
// User: pei
// DateTime: 2020/11/7 16:44
// Description: 
// ------------------

package main

import (
	"apps/ims"
	"apps/proto/build/go/pb"
	"fmt"
	"github.com/guzhi17/xcon"
	"google.golang.org/protobuf/proto"
	"log"
	"testing"
	"time"
)

func TestClient_Dial(t *testing.T) {
	var c = &ims.Client{
		Host: "10.10.1.99:3721",
		ConnConfig: xcon.ConnConfig{
			PackageMaxLength: 1 << 10,
			PackageMode:      xcon.Pm32,
			ReadTimeout:      time.Second*100,
			WriteTimeout:     time.Second*10,
		},
	}

	var err = c.Dial()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		var i = 0
		for{
			i++
			time.Sleep(time.Second)
			pkg, err := c.QueryTimeout(uint32(pb.Fids_Fid_SysEchoQuery), &pb.SysEchoQuery{Word: fmt.Sprintf("times: (%d)", i)}, time.Second)
			if err != nil{
				//log.Fatal(err)
				continue
			}
			//c.Close()
			var r pb.SysEchoQueryResponse
			err = proto.Unmarshal(pkg.Body(), &r)
			log.Println(err, r)
		}
	}()

	time.Sleep(time.Second*10)
}


type Heapable interface {
	Len()int
	Swap(i, j int)
	Less(i, j int) bool
}
func Heap(h Heapable){
	sz := h.Len()
	for i := sz>>1; i >= 0; i--{
		heap(h, sz, i)
	}
}


func heap(h Heapable, hz, i int) {
	var(
		r = (i+1) << 1
		l = r - 1
		largest = i
	)
	if l < hz && h.Less(l, i){
		largest = l
	}
	if r < hz && h.Less(r, largest){
		largest = r
	}
	if largest != i{
		h.Swap(i, largest)
		heap(h, hz, largest)
	}
}

type IntHeap []int
func (s IntHeap)Len()int{
	return len(s)
}
func (s IntHeap)Swap(i, j int){
	s[i], s[j] = s[j], s[i]
}
func (s IntHeap)Less(i, j int) bool{
	return s[i] > s[j]
}
func (s IntHeap) Top() int  {
	return s[0]
}
func (s *IntHeap) Extract() (int, bool) {
	var sz = s.Len()
	if sz < 1{
		return 0, false
	}
	var m = (*s)[0]
	//the last to the first and then re-heap
	(*s)[0] = (*s)[sz-1]
	*s = (*s)[:sz-1]
	Heap(s)
	return m, true
}

func (s IntHeap) increase(i int, key int) int {
	s[i] = key
	for ; i > 0 && s[i>>1] < s[i];  {
		s.Swap(i, i>>1)
		i = i>>1
	}
	return 0
}

func (s *IntHeap) Insert(key int) {
	var l = s.Len()
	*s = append(*s, key)
	s.increase(l, key)
}


//type Heap []int
//func MaxHeap(heap Heap){
//	sz := len(heap)
//	for i := sz>>1; i >= 0; i--{
//		maxHeap(heap, sz, i)
//	}
//}
//func maxHeap(heap Heap, hz, i int) {
//	var(
//		r = (i+1) << 1
//		l = r - 1
//		largest = i
//	)
//	if l < hz && heap[l] > heap[i]{
//		largest = l
//	}
//	if r < hz && heap[r] > heap[largest]{
//		largest = r
//	}
//	if largest != i{
//		heap[i], heap[largest] = heap[largest], heap[i]
//		maxHeap(heap, hz, largest)
//	}
//}
//func MinHeap(heap Heap){
//	sz := len(heap)
//	for i := sz>>1; i >= 0; i--{
//		minHeap(heap, sz, i)
//	}
//}
//func minHeap(heap Heap, hz, i int) {
//	var(
//		r = (i+1) << 1
//		l = r - 1
//		largest = i
//	)
//	if l < hz && heap[l] < heap[i]{
//		largest = l
//	}
//	if r < hz && heap[r] < heap[largest]{
//		largest = r
//	}
//	if largest != i{
//		heap[i], heap[largest] = heap[largest], heap[i]
//		minHeap(heap, hz, largest)
//	}
//}

func TestHeap(t *testing.T) {
	var h = IntHeap{4,1,3,2,16,9,10,14,8,7}//{1,2,3,4,5,6,7,8,9}
	//MinHeap(heap)
	Heap(h)
	log.Println(h)
	log.Println(h.Extract())
	log.Println(h)
	h.Insert(16)
	log.Println(h)

	for v, has := h.Extract(); has; v, has = h.Extract() {
		log.Println(v)
	}
}
