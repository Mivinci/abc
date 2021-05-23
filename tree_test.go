package webkit

// import (
// 	"fmt"
// 	"testing"
// )

// func pretty(n *node, prefix string) {
// 	if n == nil {
// 		return
// 	}
// 	fmt.Printf("%s%s(%d)\n", prefix, n.segment, n.weight)
// 	for _, child := range n.children {
// 		pretty(child, prefix+"  ")
// 	}
// }

// func ExampleNodeInsert() {
// 	root := node{}
// 	root.insert([]string{":a", ":b", "c"}, nil, 0)
// 	root.insert([]string{"b", "*"}, nil, 0)
// 	root.insert([]string{"*c"}, nil, 0)
// 	pretty(&root, "")

// 	// Output:
// 	// (0)
// 	//   b(3)
// 	//     *(2)
// 	//   *c(2)
// 	//   :a(1)
// 	//     :b(1)
// 	//       c(3)
// }

// func ExampleNodeLookup() {
// 	root := node{}
// 	root.insert([]string{":a", ":b", "c"}, nil, 0)
// 	root.insert([]string{"b", "*"}, nil, 0)
// 	// root.insert([]string{"*c"}, nil, 0)

// 	var ps Params
// 	n := root.lookup([]string{"a", "b", "c"}, 0, &ps)
// 	fmt.Println(n)
// 	fmt.Println(ps)

// 	// Output:
// 	// 3 c <nil> 0
// 	// map[a:a b:b]
// }

// func BenchmarkNodeLookup(b *testing.B) {
// 	root := node{}
// 	root.insert([]string{":a", ":b", "c"}, nil, 0)
// 	var ps Params
// 	for i := 0; i < b.N; i++ {
// 		root.lookup([]string{"a", "b", "c"}, 0, &ps)
// 	}
// }
