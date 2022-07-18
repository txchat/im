// Copyright 2015, Hu Keping. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/oofpgDLD/dtask/pkg/rbtree"
)

func main() {
	rbt := rbtree.New()

	rbt.Insert(rbtree.String("Hello"))
	rbt.Insert(rbtree.String("World"))

	rbt.Ascend(rbt.Min(), show)
}

func show(item rbtree.Item) bool {
	i, ok := item.(rbtree.String)
	if !ok {
		return false
	}
	fmt.Println(i)
	return true
}
