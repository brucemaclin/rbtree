package rbtree

type Tree struct {
	root  *Node
	count int
}

const (
	red   = false
	black = true
)

type Node struct {
	left   *Node
	right  *Node
	parent *Node
	color  bool
	Item
}

type Item interface {
	Less(than Item) bool
}

func Init() *Tree {
	return &Tree{}
}
func parentOf(node *Node) *Node {
	if node == nil {
		return nil
	}
	return node.parent
}

func grandparentOf(n *Node) *Node {
	parent := parentOf(n)

	return parentOf(parent)
}
func uncleOf(n *Node) *Node {

	parent := parentOf(n)
	if parent == nil {
		return nil
	}
	return siblingOf(parent)

}
func (n *Node) rotateLeft() {
	right := n.right
	if right == nil {
		panic("right node can't be nil")
	}
	p := parentOf(n)
	n.right = right.left
	right.left = n
	n.parent = right
	if n.right != nil {
		n.right.parent = n
	}
	if p != nil {
		if p.left == n {
			p.left = right
		} else {
			p.right = right
		}
	}
	right.parent = p
}
func (n *Node) rotateRight() {
	left := n.left
	if left == nil {
		panic("left node can't be nil")
	}
	p := parentOf(n)
	n.left = left.right
	left.right = n
	n.parent = left
	if n.left != nil {
		n.left.parent = n
	}
	if p != nil {
		if n == p.left {
			p.left = left
		} else {
			n.right = left
		}
	}
	left.parent = p
}
func siblingOf(n *Node) *Node {
	parent := parentOf(n)
	if parent == nil {
		return nil
	}
	if parent.left == n {
		return parent.right
	}
	return parent.left
}
func leftOf(n *Node) *Node {
	if n == nil {
		return nil
	}
	return n.left
}
func rightOf(n *Node) *Node {
	if n == nil {
		return nil
	}
	return n.right
}
func (t *Tree) Get(item Item) (*Node, bool) {
	node := t.root
	for {
		if node == nil {
			return nil, false
		}
		if node.Item.Less(item) {
			node = node.right

		} else if item.Less(node.Item) {
			node = node.left

		} else {
			break
		}
	}
	return node, true
}
func (t *Tree) InsertNoReplace(item Item) bool {
	node := t.root
	newNode := &Node{}
	newNode.Item = item
	var parent *Node
	for {
		if node == nil {
			break
		}
		parent = node
		if node.Item.Less(item) {
			node = node.right
			if node == nil {
				newNode.parent = parent
				parent.right = newNode
			}
		} else if item.Less(node.Item) {
			node = node.left
			if node == nil {
				parent.left = newNode
				newNode.parent = parent
			}
		} else {
			return false
		}
	}
	t.fixAfterInsert(newNode)
	root := newNode
	for root.parent != nil {
		root = root.parent
	}
	t.root = root
	t.count++
	return true

}
func (t *Tree) isBlack(n *Node) bool {
	if n == nil {
		return black
	}
	return n.color == black
}
func (t *Tree) setColor(n *Node, color bool) {
	if n == nil {
		return
	}
	n.color = color
}

//insert case
//   case 1: N is the root node, i.e., first node of red–black tree.
//   case 2: N's parent (P) is black all valid
//   case 3: P is red (so it can't be the root of the tree) and N's uncle (U) is red
//   case 4: P is red and U is black
func (t *Tree) fixAfterInsert(newNode *Node) {
	parent := parentOf(newNode)

	node := newNode
	//root node
	if parent == nil {
		newNode.color = black
		return
	}
	if t.isBlack(parent) { //case 2
		return
	}

	if !t.isBlack(uncleOf(node)) { //case 3
		t.fixCase3(node)
	} else { //case 4
		t.fixCase4Step1(node)
	}

}

/*
							case 3 -node's uncle is red
						      G              g
						     / \            / \
							p	u   -->    P   U
					       /              /
	                      n              n
*/
func (t *Tree) fixCase3(n *Node) {
	t.setColor(parentOf(n), black)
	t.setColor(uncleOf(n), black)
	t.setColor(grandparentOf(n), red)
	t.fixAfterInsert(grandparentOf(n))
}

/*
		case 4-1 -node's uncle is gparent's right child and black and node's parent's right child(left rotate at parent)
		case 4-2 -node's uncle is gparent's left child and  black and node's parent's left child(right rotate at parent)
		 G             G              G              G
		/ \           / \            / \            / \
	   p   U  --->   n   U  or      U   p  --->    U   n
		\           /                  /                \
		 n         p                  n                  p
*/
func (t *Tree) fixCase4Step1(n *Node) {

	p := parentOf(n)
	gp := grandparentOf(n)
	if p == gp.left && n == p.right {
		p.rotateLeft()
		n = n.left
	} else if p == gp.right && n == p.left {
		p.rotateRight()
		n = n.right
	}
	t.fixCase4Step2(n)
}

/*
       G             P                 G                    P
      / \           / \               / \                  / \
     p   U  --->   n   g      or     U   p       --->     g   n
    /                   \                 \              /
   n                     U                 n            U
*/
func (t *Tree) fixCase4Step2(n *Node) {

	p := parentOf(n)
	gp := grandparentOf(n)
	if n == p.left {
		gp.rotateRight()
	} else {
		gp.rotateLeft()
	}
	t.setColor(p, black)
	t.setColor(gp, red)

}

func (t *Tree) Delete(item Item) bool {
	node := t.root
	var result *Node
	for {
		if node == nil {
			return false
		}
		if node.Less(item) {
			node = node.right
		} else if item.Less(node.Item) {
			node = node.left
		} else {
			result = node
			break
		}
	}
	if result.left != nil && result.right != nil {
		s := t.successor(result)
		result.Item = s.Item
		result = s
	}
	var replace *Node
	if result.left != nil {
		replace = result.left
	} else {
		replace = result.right
	}

	if replace != nil {
		t.replaceNode(result, replace)
		if t.isBlack(result) {
			t.fixAfterDelete(replace)
		}

	} else if parentOf(result) == nil {
		t.root = nil
	} else {
		if t.isBlack(result) {
			t.fixAfterDelete(result)
		}

		t.replaceNode(result, replace)
	}

	t.count--

	return true

}

//from java sdk1.8 treemap
func (t *Tree) fixAfterDelete(node *Node) {

	for node != t.root && t.isBlack(node) {
		if node == parentOf(node).left {
			sib := siblingOf(node)
			/*
						   black(p.left) = black(p.right)-1  --->  black(S.left) = black(S.right)-1 not balance should continue
								P               S
			                   / \    --->     / \
							  N   s           p   SR
			                                 / \
											N   SL
			*/
			if !t.isBlack(sib) {
				t.setColor(sib, black)
				t.setColor(parentOf(node), red)
				parentOf(node).rotateLeft()
				sib = parentOf(node).right
			}
			/*
						       black(P.left) = black(P.right)-1 --> black(P.left) = black(P.right) path through p lose one black node,should check node P
								P                  P
							   / \	              / \
			                  N   S    --->      N   s
							     / \                / \
								SL  SR             SL  SR

								or                 p                   p
			                                      / \                 / \
								                 N   S      -->      N   s      also need to  check node p(set node p black)
												    / \                 / \
												   SL  SR	           SL  SR
			*/
			if t.isBlack(leftOf(sib)) && t.isBlack(rightOf(sib)) {
				t.setColor(sib, red)
				node = parentOf(node)
			} else {
				/*
				                            p?              p?
				                           / \             / \
										  N   S    --->   N   SL
										     / \               \
											sl  SR              s
											                     \
																  SR

				*/
				if t.isBlack(rightOf(sib)) {
					t.setColor(sib, red)
					t.setColor(leftOf(sib), black)
					sib.rotateRight()
					sib = rightOf(parentOf(node))
				}
				/*
										p?                    s?
				                       / \                   / \
									  N   S        --->     P   SR
									     / \               / \
										SL  sr            N   SL

				*/
				t.setColor(sib, parentOf(node).color)
				t.setColor(rightOf(sib), black)
				t.setColor(parentOf(node), black)
				parentOf(node).rotateLeft()
				node = t.root

			}
		} else {
			sib := siblingOf(node)
			if !t.isBlack(sib) {
				t.setColor(sib, black)
				t.setColor(parentOf(node), red)
				parentOf(node).rotateRight()
				sib = parentOf(node).left
			}
			if t.isBlack(leftOf(sib)) && t.isBlack(rightOf(sib)) {
				t.setColor(sib, red)
				node = parentOf(node)
			} else {
				if t.isBlack(leftOf(sib)) {
					t.setColor(sib, red)
					t.setColor(rightOf(sib), black)
					sib.rotateLeft()
					sib = leftOf(parentOf(node))
				}
				t.setColor(sib, parentOf(node).color)
				t.setColor(leftOf(sib), black)
				t.setColor(parentOf(node), black)
				parentOf(node).rotateRight()
				node = t.root

			}
		}
	}

	t.setColor(node, black)
}
func (t *Tree) deleteCase1(node *Node) {
	if parentOf(node) != nil {
		t.deleteCase2(node)
	} // else it's root do nothing
}

/*
	P              S              P            S
   / \            / \            / \          / \
  N   s   -->    p   SR  or     s   N  -->   SL  p
 / \ / \        / \            / \ / \      /\  / \
               N   SL                          SR  N
not finished because path through origin P->N now S->N get  one less black node
*/
func (t *Tree) deleteCase2(node *Node) {
	sib := siblingOf(node)
	if !t.isBlack(sib) {
		p := parentOf(node)
		t.setColor(p, red)
		t.setColor(sib, black)
		if node == p.left {
			p.rotateLeft()
		} else {
			p.rotateRight()
		}

	}
	t.deleteCase3(node)
}

/*
        P           P
       / \   -->   / \
      N   S       N   s
         / \         / \
        SL  SR      SL  SR
   path through N get one less black so just set s red.
   this will make path through P get one less black we should rebalance at P
*/
func (t *Tree) deleteCase3(node *Node) {
	p := parentOf(node)
	sib := siblingOf(node)
	//if n!= nil and is black sib can't be nil
	if t.isBlack(p) && t.isBlack(sib) && t.isBlack(sib.left) && t.isBlack(sib.right) {
		t.setColor(sib, red)
		t.deleteCase1(p)
	} else {
		t.deleteCase4(node)
	}
}

/*
     p                    P
    / \                  / \
   N   S    --->        N   s
      / \                  / \
     SL  SR               SL  SR


*/

func (t *Tree) deleteCase4(node *Node) {
	p := parentOf(node)
	sib := siblingOf(node)
	if !t.isBlack(p) && t.isBlack(sib) && t.isBlack(sib.left) && t.isBlack(sib.right) {
		t.setColor(sib, red)
		t.setColor(p, black)
	} else {
		t.deleteCase5(node)
	}
}

/*
		S              SL                   S             SR
       / \    --->     / \       or        / \    --->    /
	  sl  SR              s               SL  sr         s
                           \                            /
                            SR                         SL
*/
func (t *Tree) deleteCase5(node *Node) {
	sib := siblingOf(node)
	if t.isBlack(sib) {
		if node == leftOf(parentOf(node)) {
			if !t.isBlack(leftOf(sib)) && t.isBlack(rightOf(sib)) {
				t.setColor(leftOf(sib), black)
				t.setColor(sib, red)
				sib.rotateRight()
			}
		} else {
			if !t.isBlack(rightOf(sib)) && t.isBlack(leftOf(sib)) {
				t.setColor(leftOf(sib), black)
				t.setColor(sib, red)
				sib.rotateLeft()
			}
		}
	}
	t.deleteCase6(node)
}

/*
	(p)           (s)             (p)                    (s)
    / \           / \             / \                    / \
   N   S   --->  P   SR    or    S   N      ----->      SL  P
        \       /               /                            \
         sr    N               sl                             N
*/
func (t *Tree) deleteCase6(node *Node) {
	sib := siblingOf(node)
	parent := parentOf(node)
	t.setColor(sib, t.isBlack(parent))
	t.setColor(parent, black)
	if node == leftOf(parent) {
		t.setColor(rightOf(sib), black)
		parent.rotateLeft()
	} else {
		t.setColor(leftOf(sib), black)
		parent.rotateRight()
	}

}
func (t *Tree) replaceNode(node, replace *Node) {
	if parentOf(node) == nil {
		t.root = replace
	} else if node == parentOf(node).left {
		parentOf(node).left = replace
	} else {
		parentOf(node).right = replace
	}
	if replace != nil {
		replace.parent = node.parent
	}

}
func (t *Tree) successor(node *Node) *Node {
	if node.right != nil {
		return minNodeOfRight(node)
	}
	p := parentOf(node)
	for p != nil && node == p.right {
		node = p
		p = parentOf(p)
	}
	return p
}

func minNodeOfRight(node *Node) *Node {
	n := node.right
	for n.left != nil {
		n = n.left
	}
	return n
}
func (t *Tree) delete(node *Node) {

}
