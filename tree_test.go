package webkit

// import "fmt"

// func pretty(n *node, prefix string) {
// 	if n == nil {
// 		return
// 	}
// 	fmt.Printf("%s%s(%d)\n", prefix, n.segment, n.typ)
// 	for _, child := range n.children {
// 		pretty(child, prefix+"  ")
// 	}
// }

// func ExampleNodeInsert() {
// 	var root node
// 	root.insert([]string{":a", ":b", "c"}, nil, 0)
// 	root.insert([]string{"b", "*"}, nil, 0)
// 	root.insert([]string{"*c"}, nil, 0)
// 	pretty(&root, "")

// 	// Output:
// 	// (0)
// 	//   b(4)
// 	//     *(2)
// 	//   *c(2)
// 	//   :a(1)
// 	//     :b(1)
// 	//       c(4)
// }

// func ExampleNodeLookup() {
// 	var root node
// 	var ps Params
// 	var n node

// 	root.insert([]string{":a", ":b", "c"}, nil, 0)
// 	root.insert([]string{"b", "*"}, nil, 0)

// 	root.lookup([]string{"a", "b", "c"}, 0, &ps, &n)
// 	fmt.Println(n)
// 	fmt.Println(ps)

// 	// Output:
// 	// {4 c <nil> []}
// 	// map[a:a b:b]
// }
