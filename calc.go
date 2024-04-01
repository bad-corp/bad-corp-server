package main

import (
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	tws "github.com/muyu66/two-way-score"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"slices"
)

func calcScore(users *[]User, comments *[]SubjectComment) map[int64]float64 {
	dg := simple.NewDirectedGraph()
	for _, user := range *users {
		node := graph.Node(simple.Node(user.Id))
		dg.AddNode(node)
	}

	for _, comment := range *comments {
		edge := ScoreEdge{F: simple.Node(comment.UserId), T: simple.Node(comment.SubjectIdCreator), Score: comment.Score}
		dg.SetEdge(edge)
	}

	var d = toFullGraph(dg, 2)
	//fmt.Printf("%+v\n", d)

	var ddd []tws.Node
	for _, dd := range d {
		ddd = append(ddd, tws.Node{
			RaterId:  dd.FromId,
			TargetId: dd.ToId,
			Deep:     dd.Deep,
			Score:    int64(dd.Score),
		})
	}
	//fmt.Printf("%+v\n", ddd)
	ss, _ := tws.Calc(&ddd)
	fmt.Printf("%+v\n", ss)

	var a = make(map[int64]float64)
	for k, v := range ss {
		a[k.(int64)] = v
	}
	return a
}

func toFullGraph(
	dg *simple.DirectedGraph,
	id int64,
) []Asd {
	// 获取节点的所有邻居
	neighbors := dg.To(id)
	neighbors2 := dg.From(id)

	var asdd = make([]Asd, 0)

	filter1 := bloom.NewWithEstimates(1000, 0.01)
	iterator(filter1, false, neighbors2, dg, 0, id, &asdd)

	var deep2 int64 = 0
	if len(asdd) > 0 {
		deep2 = slices.MaxFunc(asdd, func(a, b Asd) int {
			if a.Deep > b.Deep {
				return 1
			} else if a.Deep < b.Deep {
				return -1
			}
			return 0
		}).Deep
	}

	filter2 := bloom.NewWithEstimates(1000, 0.01)
	iterator(filter2, true, neighbors, dg, deep2, id, &asdd)

	var deep3 int64 = 0
	if len(asdd) > 0 {
		deep3 = slices.MinFunc(asdd, func(a, b Asd) int {
			if a.Deep > b.Deep {
				return 1
			} else if a.Deep < b.Deep {
				return -1
			}
			return 0
		}).Deep
	}

	// deep补正
	for i, _ := range asdd {
		asdd[i].Deep += -deep3 + 1
	}

	return asdd
}

func iterator(
	filter *bloom.BloomFilter,
	to bool,
	neighbors graph.Nodes,
	dg *simple.DirectedGraph,
	deep int64,
	fromId int64,
	asdd *[]Asd,
) {
	if to {
		deep++
	} else {
		deep--
	}
	for neighbors.Next() {
		currNode := neighbors.Node()
		if filter.Test(uint64ToBytes(uint64(currNode.ID()))) {
			return
		} else {
			filter.Add(uint64ToBytes(uint64(currNode.ID())))
		}
		var nodes graph.Nodes
		if to {
			e := dg.Edge(currNode.ID(), fromId).(ScoreEdge)
			*asdd = append(*asdd, Asd{
				FromId: currNode.ID(),
				ToId:   fromId,
				Deep:   deep,
				Score:  e.Score,
			})
			nodes = dg.To(currNode.ID())
		} else {
			e := dg.Edge(fromId, currNode.ID()).(ScoreEdge)
			*asdd = append(*asdd, Asd{
				FromId: fromId,
				ToId:   currNode.ID(),
				Deep:   deep,
				Score:  e.Score,
			})
			nodes = dg.From(currNode.ID())
		}
		iterator(filter, to, nodes, dg, deep, currNode.ID(), asdd)
	}
}

func (s ScoreEdge) From() graph.Node {
	return s.F
}

func (s ScoreEdge) To() graph.Node {
	return s.T
}

func (s ScoreEdge) ReversedEdge() graph.Edge {
	return nil
}

type ScoreEdge struct {
	F     graph.Node
	T     graph.Node
	Score int8
}

type Asd struct {
	FromId int64
	ToId   int64
	Deep   int64
	Score  int8
}
