package fixtures

import "elle-go/history"

func Serializable() history.History {
	return history.History{
		{Process: 1, Type: history.Ok, Mops: []history.Mop{
			{Type: history.Append, Key: "x", App: 1},
		}},
		{Process: 2, Type: history.Ok, Mops: []history.Mop{
			{Type: history.Read, Key: "x", Read: []int{1}},
			{Type: history.Append, Key: "x", App: 2},
		}},
		{Process: 3, Type: history.Ok, Mops: []history.Mop{
			{Type: history.Read, Key: "x", Read: []int{1, 2}},
			{Type: history.Read, Key: "y", Read: []int{}},
		}},
	}
}
