package entity

// PS: 使用shouldBin Query 要使用 form tag，不然序列化不出来
type CommonListQuery struct {
	Limit     int    `form:"limit" validate:"gte=1,lte=100"  `
	Offset    int    `form:"offset" validate:"gte=0"  `
	Direction string `form:"direction" validate:"omitempty,oneof=desc asc"`
}
