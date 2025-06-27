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

var tree *AvlTree

var p any
var myp, mypp *myNode

var nodes [maxNodes]myNode

func cmpNameKey(key any, node any) int {

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

func cmpNameNode(node1 any, node2 any) int {

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

	tree = NewAvlTree()

	for i = 0; i < maxNodes; i++ {
		nodes[i].id = i
		nodes[i].hash = generateHash(i)
	}
}

func TestAvlTreeInsert(t *testing.T) {

	var it any
	var i int32

	for i = 0; i < maxNodes; i++ {
		it = tree.AvlTreeInsert(&nodes[i].avlHeader, &nodes[i], cmpNameNode)
		assert.Nil(t, it, "node already in tree!")
		if it != nil {
			// if the insert failed, mark the node as deleted to ensure
			// we don't try to remove it later
			nodes[i].deleted = true
		}
	}
}

func TestAvlTreeFirstInOrder_1st(t *testing.T) {

	p = tree.AvlTreeFirstInOrder()
	if p != nil {
		myp = p.(*myNode)
	}
}

func TestAvlTreeNextInOrder_1st(t *testing.T) {

	for ap := &myp.avlHeader; p != nil; {
		p = ap.AvlTreeNextInOrder()
		if p != nil {
			curp := myp
			myp = p.(*myNode)
			assert.Less(t, curp.hash, myp.hash)
			ap = &myp.avlHeader
		}
	}
}

func TestAvlTreeLastInOrder(t *testing.T) {

	p = tree.AvlTreeLastInOrder()
	if p != nil {
		myp = p.(*myNode)
	}
}

func TestAvlTreePrevInOrder(t *testing.T) {

	for ap := &myp.avlHeader; p != nil; {
		p = ap.AvlTreePrevInOrder()
		if p != nil {
			curp := myp
			myp = p.(*myNode)
			assert.Greater(t, curp.hash, myp.hash)
			ap = &myp.avlHeader
		}
	}
}

func TestAvlTreeFirstInPostOrder(t *testing.T) {

	p = tree.AvlTreeFirstInPostOrder()
	if p != nil {
		myp = p.(*myNode)
	}
}

func TestAvlTreeNextInPostOrder(t *testing.T) {

	for ap := &myp.avlHeader; p != nil; {
		p = ap.AvlGetParent()
		if p != nil {
			mypp = p.(*myNode)
			p = ap.AvlTreeNextInPostOrder(&mypp.avlHeader)
			myp = p.(*myNode)
			ap = &myp.avlHeader
		}
	}
}

func TestAvlTreeRemove(t *testing.T) {

	for i := 0; i < maxNodes; i += 2 {
		tree.AvlTreeRemove(&nodes[i].avlHeader)
		nodes[i].deleted = true
	}
}

func TestAvlTreeFirstInOrder_2nd(t *testing.T) {

	p = tree.AvlTreeFirstInOrder()
	if p != nil {
		myp = p.(*myNode)
	}
}

func TestAvlTreeNextInOrder_2nd(t *testing.T) {

	for ap := &myp.avlHeader; p != nil; {
		p = ap.AvlTreeNextInOrder()
		if p != nil {
			myp = p.(*myNode)
			ap = &myp.avlHeader
		}
	}
}

func TestAvlTreeLookup(t *testing.T) {

	for i := 0; i < maxNodes; i++ {
		p = tree.AvlTreeLookup(nodes[i].hash, cmpNameKey)
		if p != nil {
			assert.False(t, nodes[i].deleted,
				fmt.Sprintf("Node %d is NOT in tree!", i))
		} else {
			assert.True(t, nodes[i].deleted,
				fmt.Sprintf("Node %d IS in tree!", i))
		}
	}
}
