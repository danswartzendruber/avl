package avl

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type myNode struct {
	avlHeader AvlNode
	hash      string
	id        int32
	deleted   bool
}

const maxNodes = 1000000
const hashLen = 32

var root *AvlNode

var p interface{}
var myp, mypp *myNode

var nodes [maxNodes]myNode

func cmpNameKey(key interface{}, node interface{}) int {

	myName1 := key.(string)

	myName2 := node.(*myNode).hash

	if myName1 < myName2 {
		return -1
	} else if myName1 > myName2 {
		return 1
	} else {
		return 0
	}
}

func cmpNameNode(node1 interface{}, node2 interface{}) int {

	myName1 := node1.(*myNode).hash

	myName2 := node2.(*myNode).hash

	if myName1 < myName2 {
		return -1
	} else if myName1 > myName2 {
		return 1
	} else {
		return 0
	}
}

func generateHash(i int32) string {

	hash := md5.Sum([]byte(fmt.Sprintf("%0d", i)))
	return hex.EncodeToString(hash[:])
}

func TestAvlInit(t *testing.T) {

	var i int32

	for i = 0; i < maxNodes; i++ {
		nodes[i].id = i
		nodes[i].hash = generateHash(i)
	}
}

func TestAvlTreeInsert(t *testing.T) {

	var it interface{}
	var i int32

	for i = 0; i < maxNodes; i++ {
		it = AvlTreeInsert(&root, &nodes[i].avlHeader, &nodes[i], cmpNameNode)
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
		p = AvlTreeLookup(root, nodes[i].hash, cmpNameKey)
		if p != nil {
			assert.False(t, nodes[i].deleted,
				fmt.Sprintf("Node %d is NOT in tree!", i))
		} else {
			assert.True(t, nodes[i].deleted,
				fmt.Sprintf("Node %d IS in tree!", i))
		}
	}
}
