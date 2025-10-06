package grep

type Options struct {
	After      int
	Before     int
	Context    int
	CountOnly  bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
	LineNums   bool
	Pattern    string
}
