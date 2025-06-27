//
// Copyright as per Creative Commons Legal Code license, which can
// be found in the file COPYING
//

/*

Overview

This is a zero-dependency, high performance GO implementation of AVL trees.
It is intended to be incorporated into GO programming projects that need to
use self-balancing binary search trees.

This implementation is "intrusive", meaning that the tree node structure
must be embedded inside the data structure to be indexed in the tree.
This is the style commonly used in kernel data structures.  This is
actually the more general style of implementation; a void pointer and
comparison callback-based implementation can (but does not have to be)
be built on top of it.

This implementation is non-recursive, so it does not suffer from stack
overflows.

Structures

   AvlTree - this contains only a pointer (*AvlNode) to the root node,
   nil if empty.  The sole purpose is because the root node pointer can
   change when we add or delete nodes, and we can't use a pointer to a
   pointer as the pointer receiver.

   AvlNode - described in avl.go.

Features

Briefly, the supported operations are:

- New

func NewAvl() *TreeAvlTree

- Insertion

func (tree *AvlTree) AvlTreeInsert(item *AvlNode,
    owner any, cmp CmpFuncNode) any

- Deletion

func (tree *AvlTree) AvlTreeRemove(node *AvlNode)

- Search

func (tree *AvlTree) AvlTreeLookup(key any, cmp CmpFuncKey) any

- In-order traversal (forwards and backwards)

func (tree *AvlTree) AvlTreeFirstInOrder() any
func (tree *AvlTree) AvlTreeLastInOrder() any
func (node *AvlNode) AvlTreeNextInOrder() any
func (node *AvlNode) AvlTreePrevInOrder() any

- Post-order traversal

func (tree *AvlTree) AvlTreeFirstInPostOrder() any
func (prev *AvlNode) AvlTreeNextInPostOrder(prevParent *AvlNode) any

- Miscellaneous

   func (node *AvlNode) AvlGetParent() any {
   func (node *AvlNode) AvlLeftChild() any {
   func (node *AvlNode) AvlRightChild() any {
   func (node *AvlNode) AvlGetBalanceFactor() int {

Files

- avl.go  Interface functions.  We follow the GO convention that
          "internal" functions begin with lower-case letters, and
          "exported" functions with upper-case letters

- avl_test.go: GO test code (invoked by 'go test')

License

This code and its accompanying files have been released into the
public domain.  There is NO WARRANTY, to the extent permitted by law.
See the CC0 Public Domain Dedication in the COPYING file for details

*/

package avl
