package avl

//
// This package is a GO implementation of AVL trees
//
// Credit to Eric Biggers for the "C" implementation I used as a base
//

//
// Differences from "C" implementation.  To save memory, the "C"
// implementation used the lower 2 bits of the parent node pointer
// to store the balance factor.  This worked because an AVL tree
// guarantees no 2 leaf nodes are ever more than 1 difference in
// depth in the tree, and any structure on a modern system will be
// aligned on at least a 4-byte boundary.  We can't take that liberty
// with GO, without playing dangerous (e.g. 'unsafe' package) games
// with pointers.  We also can't safely cast pointers back and forth
// to allow transparent embedding of the AVL head in the structure the
// intended application defines.  What we do instead: the AVL node has
// a field called 'owner', which is an empty interface.  The single
// function that provides a node pointer to the AVL system then has an
// additional parameter, which is also an empty interface.  That parameter
// is a pointer to the beginning of the containing (owning) node.  So a
// call to AvlTreeInsert might look like this:
//
// type struct myNode {
//      avlHdr        AvlNode
//      id            int64
//      xxx           float64
// }
//
// var exampleNode    myNode
//
// exampleNode.id = 12345
// exampleNode.xxx = 3.14159
//
// AvlTreeInsert(root, &myNode.avlHdr, &myNode, cmp)
//
// The AVL code will stash the myNode pointer in the owner interface field
// of the AVL node structure.  Any exported AVL functions that return a
// structure to the caller (e.g. AvlTreeLookup), will pull the owner field
// out of the structure and return that interface. which the caller will
// then assign to a pointer.  Like so:
//
// var nodep *myNode
//
// nodep = AvlTreeLookup(root, 12345, cmpint64)
//
// The client then has to use a type assertion to pull a usable structure
// pointer out of the interface.  Like so:
//
// if nodep != nil {
//     pp = nodep.(*myNode)
//     (do something here)
// } else {
//     (do something else)
// }
//

type AvlNode struct {
	left    *AvlNode
	right   *AvlNode
	parent  *AvlNode
	owner   any
	balance int8
	pad     [3]int8 //nolint:unused
}

type CmpFuncKey func(any, any) int
type CmpFuncNode func(any, any) int

// Internal functions - not visible outside this package

// Returns the left child (sign < 0) or the right child (sign > 0) of the
// specified AVL tree node
// Note: for all calls of this routine, 'sign' is constant at compilation
// time, so the compiler can remove the conditional

func avlGetChild(parent *AvlNode, sign int) *AvlNode {
	if sign < 0 {
		return parent.left
	} else {
		return parent.right
	}
}

// Set the left child or right child of the specified node

func avlSetChild(parent *AvlNode, sign int, child *AvlNode) {
	if sign < 0 {
		parent.left = child
	} else {
		parent.right = child
	}
}

func avlTreeFirstOrLastInOrder(root *AvlNode, sign int) *AvlNode {

	first := root

	if first != nil {
		for avlGetChild(first, +sign) != nil {
			first = avlGetChild(first, +sign)
		}
	}

	return first
}

func avlGetParent(node *AvlNode) *AvlNode {
	return node.parent
}

func avlTreeNodeSetUnlinked(node *AvlNode) { //nolint:unused
	node.parent = node
}

func avlTreeNodeIsUnlinked(node *AvlNode) bool { //nolint:unused
	return node.parent == node
}

// Set the parent and balance factor of the specified node

func avlSetParentBalance(node, parent *AvlNode, balance int) {
	node.parent = parent
	node.balance = int8(balance + 1)
}

// Set the parent of specified node

func avlSetParent(node, parent *AvlNode) {
	node.parent = parent
}

// Returns the balance factor of the specified AVL tree node ---
// that is, the height of its right subtree minus the height of its
// left subtree

func avlGetBalanceFactor(node *AvlNode) int {

	return int(node.balance) - 1
}

// Adjust the balance factor of the specified node

func avlAdjustBalanceFactor(node *AvlNode, amount int) {

	node.balance += int8(amount)
}

// Replace a child

func avlReplaceChild(root **AvlNode, parent, oldChild, newChild *AvlNode) {
	if parent != nil {
		if oldChild == parent.left {
			parent.left = newChild
		} else {
			parent.right = newChild
		}
	} else {
		*root = newChild
	}
}

//
// Template for performing a single rotation ---
//
// sign > 0:  Rotate clockwise (right) rooted at A:
//
//           P?            P?
//           |             |
//           A             B
//          / \           / \
//         B   C?  =>    D?  A
//        / \               / \
//       D?  E?            E?  C?
//
// (nodes marked with ? may not exist)
//
// sign < 0:  Rotate counterclockwise (left) rooted at A:
//
//           P?            P?
//           |             |
//           A             B
//          / \           / \
//         C?  B   =>    A   D?
//            / \       / \
//           E?  D?    C?  E?
//
// This updates pointers but not balance factors!
//

func avlRotate(root **AvlNode, A *AvlNode, sign int) {
	B := avlGetChild(A, -sign)
	E := avlGetChild(B, +sign)
	P := avlGetParent(A)

	avlSetChild(A, -sign, E)
	avlSetParent(A, B)

	avlSetChild(B, +sign, A)
	avlSetParent(B, P)

	if E != nil {
		avlSetParent(E, A)
	}

	avlReplaceChild(root, P, A, B)
}

//
// Template for performing a double rotation ---
//
// sign > 0:  Rotate counterclockwise (left) rooted at B, then
//                   clockwise (right) rooted at A:
//
//           P?            P?          P?
//           |             |           |
//           A             A           E
//          / \           / \        /   \
//         B   C?  =>    E   C? =>  B     A
//        / \           / \        / \   / \
//       D?  E         B   G?     D?  F?G?  C?
//          / \       / \
//         F?  G?    D?  F?
//
// (nodes marked with ? may not exist)
//
// sign < 0:  Rotate clockwise (right) rooted at B, then
//                   counterclockwise (left) rooted at A:
//
//         P?          P?              P?
//         |           |               |
//         A           A               E
//        / \         / \            /   \
//       C?  B   =>  C?  E    =>    A     B
//          / \         / \        / \   / \
//         E   D?      G?  B      C?  G?F?  D?
//        / \             / \
//       G?  F?          F?  D?
//
// Returns a pointer to E and updates balance factors.  Except for those
// two things, this function is equivalent to:
//      avlRotate(root, B, -sign)
//      avlRotate(root, A, +sign)
//
// See comment in avlHandleSubtreeGrowth() for explanation of balance
// factor updates.

func avlDoDoubleRotate(root **AvlNode, B, A *AvlNode, sign int) *AvlNode {

	E := avlGetChild(B, +sign)
	F := avlGetChild(E, -sign)
	G := avlGetChild(E, +sign)
	P := avlGetParent(A)
	e := avlGetBalanceFactor(E)

	avlSetChild(A, -sign, G)
	if sign*e >= 0 {
		avlSetParentBalance(A, E, 0)
	} else {
		avlSetParentBalance(A, E, -e)
	}

	avlSetChild(B, +sign, F)
	if sign*e <= 0 {
		avlSetParentBalance(B, E, 0)
	} else {
		avlSetParentBalance(B, E, -e)
	}

	avlSetChild(E, +sign, A)
	avlSetChild(E, -sign, B)
	avlSetParentBalance(E, P, 0)

	if G != nil {
		avlSetParent(G, A)
	}

	if F != nil {
		avlSetParent(F, B)
	}

	avlReplaceChild(root, P, A, E)

	return E
}

//
// This function handles the growth of a subtree due to an insertion.
//
// root
//      Location of the tree's root pointer.
//
// node
//      A subtree that has increased in height by 1 due to an insertion.
//
// parent
//      Parent of node; must not be nil.
//
// sign
//      -1 if node is the left child of parent
//      +1 if node is the right child of parent
//
// This function will adjust parent's balance factor, then do a (single
// or double) rotation if necessary.  The return value will be true if
// the full AVL tree is now adequately balanced, or false if the subtree
// rooted at parent is now adequately balanced but has increased in
// height by 1, so the caller should continue up the tree.
//
// Note that if false is returned, no rotation will have been done.
// Indeed, a single node insertion cannot require that more than one
// (single or double) rotation be done.
//

func avlHandleSubtreeGrowth(root **AvlNode, node, parent *AvlNode, sign int) bool {
	oldBalanceFactor := avlGetBalanceFactor(parent)

	if oldBalanceFactor == 0 {
		avlAdjustBalanceFactor(parent, sign)

		// parent is still sufficiently balanced (-1 or +1
		// balance factor), but must have increased in height.
		// Continue up the tree

		return false
	}

	newBalanceFactor := oldBalanceFactor + sign
	if newBalanceFactor == 0 {
		avlAdjustBalanceFactor(parent, sign)

		// parent is now perfectly balanced (0 balance factor).
		// It cannot have increased in height, so there is
		// nothing more to do

		return true
	}

	// parent is too left-heavy (newBalanceFactor == -2) or
	// too right-heavy (newBalanceFactor == +2)

	// Test whether node is left-heavy (-1 balance factor) or
	// right-heavy (+1 balance factor).
	// Note that it cannot be perfectly balanced (0 balance factor)
	// because here we are under the invariant that node has
	// increased in height due to the insertion.  */

	if sign*avlGetBalanceFactor(node) > 0 {

		// node (B below) is heavy in the same direction parent
		// (A below) is heavy.
		//
		// ==============================================================
		// The comment, diagram, and equations below assume sign < 0
		// The other case is symmetric!
		// =============================================================
		//
		// Do a clockwise rotation rooted at parent (A below):
		//
		//           A              B
		//          / \           /   \
		//         B   C?  =>    D     A
		//        / \           / \   / \
		//       D   E?        F?  G?E?  C?
		//      / \
		//     F?  G?
		//
		// Before the rotation:
		//      balance(A) = -2
		//      balance(B) = -1
		// Let x = height(C).  Then:
		//      height(B) = x + 2
		//      height(D) = x + 1
		//      height(E) = x
		//      max(height(F), height(G)) = x.
		//
		// After the rotation:
		//      height(D) = max(height(F), height(G)) + 1
		//                = x + 1
		//      height(A) = max(height(E), height(C)) + 1
		//                = max(x, x) + 1 = x + 1
		//      balance(B) = 0
		//      balance(A) = 0
		//

		avlRotate(root, parent, -sign)

		// Equivalent to setting parent's balance factor to 0.
		avlAdjustBalanceFactor(parent, -sign) /* A */

		// Equivalent to setting node's balance factor to 0.
		avlAdjustBalanceFactor(node, -sign) /* B */

	} else {
		// node (B below) is heavy in the direction opposite
		// from the direction parent (A below) is heavy.
		//
		// =============================================================
		// The comment, diagram, and equations below assume sign < 0
		// The other case is symmetric!
		// ============================================================
		//
		// Do a counterblockwise rotation rooted at node (B below),
		// then a clockwise rotation rooted at parent (A below):
		//
		//           A             A           E
		//          / \           / \        /   \
		//         B   C?  =>    E   C? =>  B     A
		//        / \           / \        / \   / \
		//       D?  E         B   G?     D?  F?G?  C?
		//          / \       / \
		//         F?  G?    D?  F?
		//
		// Before the rotation:
		//      balance(A) = -2
		//      balance(B) = +1
		// Let x = height(C).  Then:
		//      height(B) = x + 2
		//      height(E) = x + 1
		//      height(D) = x
		//      max(height(F), height(G)) = x
		//
		// After both rotations:
		//      height(A) = max(height(G), height(C)) + 1
		//                = x + 1
		//      balance(A) = balance(E{orig}) >= 0 ? 0 : -balance(E{orig})
		//      height(B) = max(height(D), height(F)) + 1
		//                = x + 1
		//      balance(B) = balance(E{orig} <= 0) ? 0 : -balance(E{orig})
		//
		//      height(E) = x + 2
		//      balance(E) = 0
		//

		avlDoDoubleRotate(root, node, parent, -sign)
	}

	// Height after rotation is unchanged; nothing more to do

	return true
}

// Rebalance the tree after insertion of the specified node

func avlTreeRebalanceAfterInsert(root **AvlNode, inserted *AvlNode) {

	inserted.left = nil
	inserted.right = nil

	node := inserted

	// Adjust balance factor of new node's parent
	// No rotation will need to be done at this level

	parent := avlGetParent(node)
	if parent == nil {
		return
	}

	if node == parent.left {
		avlAdjustBalanceFactor(parent, -1)
	} else {
		avlAdjustBalanceFactor(parent, +1)
	}

	if avlGetBalanceFactor(parent) == 0 {
		// parent did not change in height - Nothing more to do
		return
	}

	// The subtree rooted at parent increased in height by 1

	for done := false; !done; {
		// Adjust balance factor of next ancestor

		node = parent
		parent = avlGetParent(node)
		if parent == nil {
			return
		}

		// The subtree rooted at node has increased in height by 1
		if node == parent.left {
			done = avlHandleSubtreeGrowth(root, node, parent, -1)
		} else {
			done = avlHandleSubtreeGrowth(root, node, parent, +1)
		}
	}
}

//
// This function handles the shrinkage of a subtree due to a deletion.
//
// root
//      Location of the tree's root pointer.
//
// parent
//      A node in the tree, exactly one of whose subtrees has decreased
//      in height by 1 due to a deletion.  (This includes the case where
//      one of the child pointers has become nil, since we can consider
//      the "nil" subtree to have a height of 0.)
//
// sign
//      +1 if the left subtree of parent has decreased in height by 1
//      -1 if the right subtree of parent has decreased in height by 1
//
// leftDeletedRet
//      If the return value is not nil, this will be set to true if the
//      left subtree of the returned node has decreased in height by 1,
//      or false if the right subtree of the returned node has decreased
//      in height by 1.
//
// This function will adjust parent's balance factor, then do a (single
// or double) rotation if necessary.  The return value will be nil if
// the full AVL tree is now adequately balanced, or a pointer to the
// parent of parent if parent is now adequately balanced but has
// decreased in height by 1.  Also in the latter case, leftDeletedRet
// will be set.
//

func avlHandleSubtreeShrink(root **AvlNode, parent *AvlNode, sign int, leftDeletedRet *bool) *AvlNode {

	var node *AvlNode

	oldBalanceFactor := avlGetBalanceFactor(parent)

	if oldBalanceFactor == 0 {
		// Prior to the deletion, the subtree rooted at
		// parent was perfectly balanced.  It's now
		// unbalanced by 1, but that's okay and its height
		// hasn't changed.  Nothing more to do

		avlAdjustBalanceFactor(parent, sign)
		return nil
	}

	newBalanceFactor := oldBalanceFactor + sign
	if newBalanceFactor == 0 {
		// The subtree rooted at parent is now perfectly
		// balanced, whereas before the deletion it was
		// unbalanced by 1.  Its height must have decreased
		// by 1.  No rotation is needed at this location,
		// but continue up the tree

		avlAdjustBalanceFactor(parent, sign)
		node = parent
	} else {
		// parent is too left-heavy (newBalanceFactor == -2) or
		// too right-heavy (newBalanceFactor == +2)

		node = avlGetChild(parent, sign)

		/* The rotations below are similar to those done during
		 * insertion (see avlHandleSubtreeGrowth()), so full
		 * comments are not provided.  The only new case is the
		 * one where node has a balance factor of 0, and that is
		 * commented.  */

		if sign*avlGetBalanceFactor(node) >= 0 {

			avlRotate(root, parent, -sign)

			if avlGetBalanceFactor(node) == 0 {

				//
				// node (B below) is perfectly balanced.
				//
				// ================================================
				// The comment, diagram, and equations
				// below assume sign < 0.  The other case
				// is symmetric!
				// ================================================
				//
				// Do a clockwise rotation rooted at
				// parent (A below):
				//
				//           A              B
				//          / \           /   \
				//         B   C?  =>    D     A
				//        / \           / \   / \
				//       D   E         F?  G?E   C?
				//      / \
				//     F?  G?
				//
				// Before the rotation:
				//      balance(A) = -2
				//      balance(B) =  0
				// Let x = height(C).  Then:
				//      height(B) = x + 2
				//      height(D) = x + 1
				//      height(E) = x + 1
				//      max(height(F), height(G)) = x.
				//
				// After the rotation:
				//      height(D) = max(height(F), height(G)) +
				//
				//                = x + 1
				//      height(A) = max(height(E), height(C)) +
				//
				//                = max(x + 1, x) + 1 = x + 2
				//      balance(A) = -1
				//      balance(B) = +1
				//

				// A: -2 => -1 (sign < 0)
				// or +2 => +1 (sign > 0)
				// No change needed --- that's the same as
				// oldBalanceFactor.  */

				// B: 0 => +1 (sign < 0)
				// or 0 => -1 (sign > 0)  */

				avlAdjustBalanceFactor(node, -sign)

				/* Height is unchanged; nothing more to do */
				return nil
			} else {
				avlAdjustBalanceFactor(parent, -sign)
				avlAdjustBalanceFactor(node, -sign)
			}
		} else {
			node = avlDoDoubleRotate(root, node, parent, -sign)
		}
	}

	parent = avlGetParent(node)
	if parent != nil {
		*leftDeletedRet = (node == parent.left)
	}

	return parent
}

// Swaps node X, which must have 2 children, with its in-order successor,
// then unlinks node X.  Returns the parent of X just before unlinking,
// without its balance factor having been updated to account for the unlink

func avlTreeSwapWithSuccessor(root **AvlNode, X *AvlNode, leftDeletedRet *bool) *AvlNode {

	var Y, ret, Q *AvlNode

	Y = X.right
	if Y.left == nil {
		//
		//     P?           P?           P?
		//     |            |            |
		//     X            Y            Y
		//    / \          / \          / \
		//   A   Y    =>  A   X    =>  A   B?
		//      / \          / \
		//    (0)  B?      (0)  B?
		//
		// [ X unlinked, Y returned ]
		//

		ret = Y
		*leftDeletedRet = false
	} else {

		for {
			Q = Y
			Y = Y.left
			if Y.left == nil {
				break
			}
		}

		//
		//     P?           P?           P?
		//     |            |            |
		//     X            Y            Y
		//    / \          / \          / \
		//   A   ...  =>  A  ...   =>  A  ...
		//       |            |            |
		//       Q            Q            Q
		//      /            /            /
		//     Y            X            B?
		//    / \          / \
		//  (0)  B?      (0)  B?
		//
		//
		// [ X unlinked, Q returned ]
		//

		Q.left = Y.right
		if Q.left != nil {
			avlSetParent(Q.left, Q)
		}
		Y.right = X.right
		avlSetParent(X.right, Y)
		ret = Q
		*leftDeletedRet = true
	}

	Y.left = X.left
	avlSetParent(X.left, Y)

	Y.parent = X.parent
	Y.balance = X.balance
	avlReplaceChild(root, avlGetParent(X), X, Y)

	return ret
}

func avlTreeNextOrPrevInOrder(node *AvlNode, sign int) *AvlNode {

	var next *AvlNode

	if avlGetChild(node, +sign) != nil {
		for next = avlGetChild(node, +sign); avlGetChild(next, -sign) != nil; {
			next = avlGetChild(next, -sign)
		}
	} else {
		for next = avlGetParent(node); next != nil && node == avlGetChild(next, +sign); {
			node = next
			next = avlGetParent(next)
		}
	}

	return next
}

// Exported functions

// Look up a specified key.  nil if not present

func AvlTreeLookup(root *AvlNode, key any, cmp CmpFuncKey) any {

	cur := root

	for cur != nil {
		res := cmp(key, cur.owner)
		if res < 0 {
			cur = cur.left
		} else if res > 0 {
			cur = cur.right
		} else {
			break
		}
	}

	if cur != nil {
		return cur.owner
	} else {
		return nil
	}
}

// Insert a node into the tree.  Returns nil if not already present,
// and existing node address if already present

func AvlTreeInsert(root **AvlNode, item *AvlNode,
	owner any, cmp CmpFuncNode) any {

	curPtr := root
	var cur *AvlNode = nil

	for *curPtr != nil {
		cur = *curPtr

		res := cmp(owner, cur.owner)
		if res < 0 {
			curPtr = &cur.left
		} else if res > 0 {
			curPtr = &cur.right
		} else {
			return cur.owner
		}
	}

	*curPtr = item

	item.parent = cur
	item.balance = 1
	item.owner = owner

	avlTreeRebalanceAfterInsert(root, item)

	return nil
}

// Removes an item from the specified AVL tree.
//
// root
//      Location of the AVL tree's root pointer.  Indirection is needed
//      because the root node may change if the tree needed to be rebalanced
//      because of the deletion or if node was the root node.
//
// node
//      Pointer to the `AvlNode' embedded in the item to remove from the tree
//
// Note: This function *only* removes the node and rebalances the tree.
// It does not free any memory, nor does it do the equivalent of
// avl_TreeNodeSetUnlinked()

func AvlTreeRemove(root **AvlNode, node *AvlNode) {
	var parent *AvlNode
	leftDeleted := false

	if node.left != nil && node.right != nil {
		// node is fully internal, with two children.  Swap it
		// with its in-order successor (which must exist in the
		// right subtree of node and can have, at most, a right
		// child), then unlink node

		parent = avlTreeSwapWithSuccessor(root, node, &leftDeleted)

		// parent is now the parent of what was node's in-order
		// successor.  It cannot be NULL, since node itself was
		// an ancestor of its in-order successor.
		// leftDeleted has been set to true if node's
		// in-order successor was the left child of parent,
		// otherwise false

	} else {
		var child *AvlNode

		// node is missing at least one child.  Unlink it.
		// Set parent to node's parent, and set leftDeleted
		// to reflect which child of parent node was.
		// Or if node was the root node, simply update the
		// root node and return

		if node.left != nil {
			child = node.left
		} else {
			child = node.right
		}
		parent = avlGetParent(node)
		if parent != nil {
			if node == parent.left {
				parent.left = child
				leftDeleted = true
			} else {
				parent.right = child
				leftDeleted = false
			}
			if child != nil {
				avlSetParent(child, parent)
			}
		} else {
			if child != nil {
				avlSetParent(child, parent)
			} else if *root != node {
				//
				// If no children and not the root, this node is not
				// in the tree!
				//
				panic("Node not in tree!")
			}
			*root = child
			return
		}
	}

	// Rebalance the tree

	for {
		if leftDeleted {
			parent = avlHandleSubtreeShrink(root, parent, +1, &leftDeleted)
		} else {
			parent = avlHandleSubtreeShrink(root, parent, -1, &leftDeleted)
		}
		if parent == nil {
			break
		}
	}
}

// Starts an in-order traversal of the tree: returns the
// least-valued node, or nil if the tree is empty

func AvlTreeFirstInOrder(root *AvlNode) any {
	rp := avlTreeFirstOrLastInOrder(root, -1)
	if rp != nil {
		return rp.owner
	} else {
		return nil
	}
}

// Starts an reverse in-order traversal of the tree: returns the
// greatest-valued node, or nil if the tree is empty

func AvlTreeLastInOrder(root *AvlNode) any {
	rp := avlTreeFirstOrLastInOrder(root, 1)
	if rp != nil {
		return rp.owner
	} else {
		return nil
	}
}

// Continues an in-order traversal of the tree

func AvlTreeNextInOrder(node *AvlNode) any {
	rp := avlTreeNextOrPrevInOrder(node, 1)
	if rp != nil {
		return rp.owner
	} else {
		return nil
	}
}

// Continues a reverse in-order traversal of the tree

func AvlTreePrevInOrder(node *AvlNode) any {
	rp := avlTreeNextOrPrevInOrder(node, -1)
	if rp != nil {
		return rp.owner
	} else {
		return nil
	}
}

// Starts a postorder traversal of the tree

func AvlTreeFirstInPostOrder(root *AvlNode) any {

	var first *AvlNode

	if root != nil {
		for first = root; first.left != nil || first.right != nil; {
			if first.left != nil {
				first = first.left
			} else {
				first = first.right
			}
		}
	}

	rp := first
	if rp != nil {
		return rp.owner
	} else {
		return nil
	}
}

// Continues a postorder traversal of the tree

func AvlTreeNextInPostOrder(prev, prevParent *AvlNode) any {

	next := prevParent

	if next != nil && prev == next.left && next.right != nil {
		for next = next.right; next.left != nil || next.right != nil; {
			if next.left != nil {
				next = next.left
			} else {
				next = next.right
			}
		}
	}

	rp := next
	if rp != nil {
		return rp.owner
	} else {
		return nil
	}
}

// Return the parent of a node

func AvlGetParent(node *AvlNode) any {
	rp := avlGetParent(node)
	if rp != nil {
		return rp.owner
	} else {
		return nil
	}
}

// Return the left child of a node

func AvlLeftChild(node *AvlNode) any {
	rp := node.left
	if rp != nil {
		return rp.owner
	} else {
		return nil
	}
}

// Return the right child of a node

func AvlRightChild(node *AvlNode) any {
	rp := node.right
	if rp != nil {
		return rp.owner
	} else {
		return nil
	}
}

// Return the balance factor of a node

func AvlGetBalanceFactor(node *AvlNode) int {

	return avlGetBalanceFactor(node)

}
