package comet

type Node struct {
	Current *Channel
	Next    *Node
	Prev    *Node
}
