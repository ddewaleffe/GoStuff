package main

// Sort of like the Erlang ring message benchmark

import (
	//	"net/http"
	"container/ring"
	"fmt"
	"time"
)

// A Node has a channel from which it will receive a message from a neighbour and the send it to the next node.
type Node struct {
	self             *ring.Ring //to allow navigation along the ring
	id               int
	ch               chan string // channel via which we receive messages
	sent             int         //counters
	rcv              int
	neighbourChannel chan string // where do I send?
	neighbourid      int
}

// INode defines the interface methods that we want for our nodes.
type INode interface {
	Init(r *ring.Ring, i int)
	RcvThenSend(r *ring.Ring)
	GetChan() chan string
	Link(next *ring.Ring)
}

// how many nodes?
const nnodes = 100000
const nloops = 1

// GetChan ...
// returns the Channel of the current node.
func (node *Node) GetChan() chan string {
	return node.ch
}

// Init ...
//  initializes current node of ring
func (node *Node) Init(rp *ring.Ring, index int) {
	//fmt.Printf("Init() called (index=:%v)\n", index)
	node.self = rp
	node.id = index
	node.ch = make(chan string)
	node.sent = 0
	node.rcv = 0
	//fmt.Printf("Init:value of n at end: %v\n", node)
}

// Link ..
// connects this node to the next one in the ring.
func (node *Node) Link(next *ring.Ring) {
	// Node that I understand this better I think liking by the id's should not be necessary.
	// the only essential thing is the channel of the neighbring node. And even then, this could be
	// fetched while running
	nextINode := next.Value.(INode)
	node.neighbourChannel = nextINode.GetChan()
	node.neighbourid = next.Value.(*Node).id
}

// RcvThenSend ...
// Waits for message on this node's channel, then sends it onto the node's chanel.
func (node *Node) RcvThenSend(r *ring.Ring) {
	//	nn:=r.Value
	//fmt.Printf("Starting RcvThenSend on node.id:%v\n", node.id)
	for {
		msg := <-node.ch
		node.rcv++
		//fmt.Printf("Node %v got '%v'\n", node.id, msg)
		// do this, but if node is 0 and has already sent nloop msgs we stop
		if !(node.id == 0 && node.sent >= nloops) {
			// send message to next in ring
			node.neighbourChannel <- msg + fmt.Sprintf("[%v->%v]", node.id, node.neighbourid)
			node.sent++
		} else {
			now := time.Now()
			fmt.Printf("Done. Time: %v\n",
				(now.UTC().Format("Mon, 2 Jan 2006 15:04:05 MST")))
		}
	}
}

func main() {
	now := time.Now()
	fmt.Printf("Starting. Time: %v\n",
		(now.UTC().Format("Mon, 2 Jan 2006 15:04:05 MST")))

	// create a ring of nodes
	r := ring.New(nnodes)
	fmt.Printf("r at begin of run: %v\n", r)

	for i := 0; i < nnodes; i++ {
		nn := new(Node)
		//fmt.Printf("type of nn: %v\n", fmt.Sprintf("%T", nn))
		//fmt.Printf("value of nn: %v\n", nn)
		//inn := INode(nn)
		//fmt.Printf("type of inn: %v\n", fmt.Sprintf("%T", inn))
		nn.Init(r, i)
		r.Value = nn
		r = r.Next()
	}
	now = time.Now()
	fmt.Printf("Nodes created. Time: %v\n",
		(now.UTC().Format("Mon, 2 Jan 2006 15:04:05 MST")))

	for i := 0; i < nnodes; i++ {
		//var current INode
		//fmt.Printf("type of r.Value: %T\n", r.Value)

		icurrent := INode(r.Value.(*Node))
		//fmt.Printf("type of current: %T\n", icurrent)
		//fmt.Printf("value of current: %v\n", icurrent)
		icurrent.Link(r.Next())
		go icurrent.RcvThenSend(r)
		r = r.Next()
	}
	now = time.Now()
	fmt.Printf("Finished linking nodes\nSending init msg. Time: %v\n",
		(now.UTC().Format("Mon, 2 Jan 2006 15:04:05 MST")))

	r.Value.(*Node).ch <- "hello"

	select {
	case <-time.After(150000 * time.Millisecond):
		fmt.Printf("Timout. Quitting.\n")
	}
	fmt.Printf("r at end of run: %v\n", r)

}
