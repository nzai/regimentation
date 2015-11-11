package trading

type TurtleSimulate struct {
	*Simulate
	Holding int //	头寸
	N       int //	实波动幅度区间
	Enter   int //	入市区间
	Exit    int //	退市区间
	Stop    int //	止损区间
}

