package flags

type Flags struct {
	Column   int
	Numeric  bool
	Reverse  bool
	Unique   bool
	Month    bool
	IgnoreTB bool
	Check    bool
	Human    bool
}

func New() Flags {
	return Flags{}
}
