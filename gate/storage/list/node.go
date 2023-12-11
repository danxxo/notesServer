package list

type node struct {
	value interface{}
	index int64
	next  *node
}
