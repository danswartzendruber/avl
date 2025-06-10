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

Features

Briefly, the supported operations are:

- Insertion
- Deletion
- Search
- In-order traversal (forwards and backwards)
- Post-order traversal

See avl_tree.h for details.

Files

- avl_tree.go  Interface functions.  We follow the GO convention that
               "internal" functions begin with lower-case letters, and
               "exported" functions with upper-case letters

- test.go:     A test program.

License

This code and its accompanying files have been released into the
public domain.  There is NO WARRANTY, to the extent permitted by law.
See the CC0 Public Domain Dedication in the COPYING file for details

*/

package avl
