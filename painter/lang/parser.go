package lang

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"strconv"
	"strings"

	"github.com/sueta2016/labik-3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
	lastBgColor painter.Operation
	lastBgRect  *painter.BgRectangle
	figures     []painter.Figure
	moveOps     []painter.Operation
	updateOp    painter.Operation
}

// Parse reads and parses input from the provided io.Reader and returns the corresponding list of painter.Operation.
func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() { // loop through the input stream using the scanner
		commandLine := scanner.Text()

		err := p.parse(commandLine) // parse the command line into an operation
		if err != nil {
			return nil, err
		}
	}
	return p.finalResult(), nil
}

func (p *Parser) finalResult() []painter.Operation {
	var res []painter.Operation
	if p.lastBgColor != nil {
		res = append(res, p.lastBgColor)
	}
	if p.lastBgRect != nil {
		res = append(res, p.lastBgRect)
	}
	if len(p.moveOps) != 0 {
		res = append(res, p.moveOps...)
	}
	p.moveOps = nil
	if len(p.figures) != 0 {
		println(len(p.figures))
		for _, figure := range p.figures {
			res = append(res, &figure)
		}
	}
	if p.updateOp != nil {
		res = append(res, p.updateOp)
	}
	return res
}

func (p *Parser) resetState() {
	p.lastBgColor = nil
	p.lastBgRect = nil
	p.figures = nil
	p.moveOps = nil
	p.updateOp = nil
}

func (p *Parser) parse(commandLine string) error {
	parts := strings.Split(commandLine, " ")
	instruction := parts[0]
	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}
	var iArgs []int
	for _, arg := range args {
		i, err := strconv.Atoi(arg)
		if err == nil {
			iArgs = append(iArgs, i)
		}
	}

	var figureOps []painter.Figure

	switch instruction {
	case "white":
		p.lastBgColor = painter.OperationFunc(painter.WhiteFill)
	case "green":
		p.lastBgColor = painter.OperationFunc(painter.GreenFill)
	case "bgrect":
		p.lastBgRect = &painter.BgRectangle{X1: iArgs[0], Y1: iArgs[1], X2: iArgs[2], Y2: iArgs[3]}
	case "figure":
		clr := color.RGBA{R: 255, G: 255, B: 0, A: 1}
		figure := painter.Figure{X: iArgs[0], Y: iArgs[1], C: clr}
		p.figures = append(p.figures, figure)
	case "move":
		moveOp := painter.Move{X: iArgs[0], Y: iArgs[1], Figures: figureOps}
		p.moveOps = append(p.moveOps, &moveOp)
	case "reset":
		p.resetState()
		p.lastBgColor = painter.OperationFunc(painter.ResetScreen)
	case "update":
		p.updateOp = painter.UpdateOp
	default:
		return fmt.Errorf("could not parse command %v", commandLine)
	}
	return nil
}