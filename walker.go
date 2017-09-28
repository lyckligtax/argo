package ar

type Walker struct {
	Bytes    []byte
	position int
}

func (w *Walker) Next(n int) (b []byte) {
	b = w.Bytes[w.position : w.position+n]
	w.position += n
	return
}

func NewWalker(bytes []byte) *Walker {
	return &Walker{
		Bytes:    bytes,
		position: 0,
	}
}
