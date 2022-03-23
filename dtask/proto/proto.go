package proto

const (
	Start = 1
	Stop  = 2
)

type Proto struct {
	Uid string
	Opt int
	Seq int64
}
