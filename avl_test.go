package avl

import (
	"encoding/base64"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

type myNode struct {
	avlHeader AvlNode
	hash      string
	id        int32
	deleted   bool
}

const maxNodes = 1000
const hashLen = 32

var root *AvlNode

var p interface{}
var myp, mypp *myNode

var nodes [maxNodes]myNode

func cmpNumKey(key interface{}, node interface{}) int {

	myNum1 := key.(int32)

	myNum2 := node.(*myNode).id

	if myNum1 < myNum2 {
		return -1
	} else if myNum1 > myNum2 {
		return 1
	} else {
		return 0
	}
}

func cmpNumNode(node1 interface{}, node2 interface{}) int {

	myNum1 := node1.(*myNode).id

	myNum2 := node2.(*myNode).id

	if myNum1 < myNum2 {
		return -1
	} else if myNum1 > myNum2 {
		return 1
	} else {
		return 0
	}
}

func generateHash() string {

	hash := make([]byte, hashLen)
	_, _ = rand.Read(hash)
	rs := base64.StdEncoding.EncodeToString(hash)
	return rs
}

func TestAvlInit(t *testing.T) {

	rand.Seed(time.Now().Unix())

	for i := 0; i < maxNodes; i++ {
		nodes[i].id = rand.Int31()
		nodes[i].hash = generateHash()
	}
}

func TestAvlTreeInsert(t *testing.T) {

	var it interface{}

	for i := 0; i < maxNodes; i++ {
		it = AvlTreeInsert(&root, &nodes[i].avlHeader, &nodes[i], cmpNumNode)
		assert.Nil(t, it, "node already in tree!")
	}
}

func TestAvlTreeFirstInOrder_1st(t *testing.T) {

	p = AvlTreeFirstInOrder(root)
	if p != nil {
		myp = p.(*myNode)
	}
}

func TestAvlTreeNextInOrder_1st(t *testing.T) {

	for p != nil {
		p = AvlTreeNextInOrder(&myp.avlHeader)
		if p != nil {
			myp = p.(*myNode)
		}
	}
}

func TestAvlTreeLastInOrder(t *testing.T) {

	p = AvlTreeLastInOrder(root)
	if p != nil {
		myp = p.(*myNode)
	}
}

func TestAvlTreePrevInOrder(t *testing.T) {

	for p != nil {
		p = AvlTreePrevInOrder(&myp.avlHeader)
		if p != nil {
			myp = p.(*myNode)
		}
	}
}

func TestAvlTreeFirstInPostOrder(t *testing.T) {

	p = AvlTreeFirstInPostOrder(root)
	if p != nil {
		myp = p.(*myNode)
	}
}

func TestAvlTreeNextInPostOrder(t *testing.T) {

	for p != nil {
		p = AvlGetParent(&myp.avlHeader)
		if p != nil {
			mypp = p.(*myNode)
			p = AvlTreeNextInPostOrder(&myp.avlHeader, &mypp.avlHeader)
			myp = p.(*myNode)
		}
	}
}

func TestAvlTreeRemove(t *testing.T) {

	for i := 0; i < maxNodes; i += 2 {
		AvlTreeRemove(&root, &nodes[i].avlHeader)
		nodes[i].deleted = true
	}
}

func TestAvlTreeFirstInOrder_2nd(t *testing.T) {

	p = AvlTreeFirstInOrder(root)
	if p != nil {
		myp = p.(*myNode)
	}
}

func TestAvlTreeNextInOrder_2nd(t *testing.T) {

	for p != nil {
		p = AvlTreeNextInOrder(&myp.avlHeader)
		if p != nil {
			myp = p.(*myNode)
		}
	}
}

func TestAvlTreeLookup(t *testing.T) {

	for i := 0; i < maxNodes; i++ {
		p = AvlTreeLookup(root, nodes[i].id, cmpNumKey)
		if p != nil {
			assert.False(t, nodes[i].deleted,
				fmt.Sprintf("Node %d is NOT in tree!", i))
		} else {
			assert.True(t, nodes[i].deleted,
				fmt.Sprintf("Node %d IS in tree!", i))
		}
	}
}
