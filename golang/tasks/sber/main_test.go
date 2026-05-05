package main

import "testing"

func Test_NewObject(t *testing.T) {
	tests := []struct {
		title  string
		str    string
		result Sayer
	}{
		{
			title:  "Success Base",
			str:    "Base",
			result: Base{name: "NewObjectBase"},
		},
		{
			title: "Success Child",
			str:   "Child",
			result: Child{
				Base: Base{
					name: "NewObjectChild",
				},
				lastName: "NewObjectLastNameChild",
			},
		},
		{
			title:  "Fail NewObject",
			str:    "tasdfasdfas",
			result: nil,
		},
	}

	for _, v := range tests {
		t.Run(v.title, func(t *testing.T) {
			newObj := NewObject(v.str)
			if v.result != newObj {
				t.Fatal()
			}
		})
	}
}
