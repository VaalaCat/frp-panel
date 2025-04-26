package dao

import "github.com/VaalaCat/frp-panel/app"

type Query interface{}

type queryImpl struct {
	ctx *app.Context
}

func NewQuery(ctx *app.Context) *queryImpl {
	return &queryImpl{
		ctx: ctx,
	}
}
