package parse

type LineScan struct {
	Line string
	Pos  int
}

/*
	LineScan
	========
	Given a command:
			sleep 1 & echo "sup"

	Parse a pipeline of commands:
	-----------------------------
	- Start from first non-whitespace
	- Read until space (aka parse-command)
	- Read args
		* Command ends at &, |, &&, ||, or \n

	How does a command pipeline get represented (tree of dependencies)

*/

func (l LineScan) Consume() {

}
