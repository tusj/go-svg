// smartSVG implements svg helper functions to write good, well-formed svg in a simple fashion.
// Tags which encloses groups returns a new *SVG object which is possible to write to.
// Flush writes the coded svg to the writer provided.
package smartSVG

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"marketExchanger/helpers"
	"math"
)

const (
	svgInit = `<?xml version="1.0"?>
<!-- Generated by S-martVGo -->`
)

var defaultNamespace map[string]string

func init() {
	defaultNamespace = map[string]string{"xmlns": "http://www.w3.org/2000/svg", "xmlns:xlink": "http://www.w3.org/1999/xlink"}
}

// Tags
const (
	svg       = "svg"
	sym       = "sym"
	use       = "use"
	title     = "title"
	desc      = "desc"
	marker    = "marker"
	defs      = "defs"
	g         = "g"
	transform = "transform"
	circle    = "circle"
	rect      = "rect"
	line      = "line"
	polyline  = "polyline"
	text      = "text"
)

// Atts
const (
	x      = "x"
	y      = "y"
	rad    = "r"
	width  = "width"
	height = "height"
	points = "points"
)

// Holds one svg group with nodes of children inside group.
type SVG struct {
	tag      string
	atts     map[string]string
	data     string
	comments []string
	mids     []*SVG
}

type Att struct {
	atts map[string]string
}

func (a *Att) Translate(x, y int) {
	a.atts["transform"] = a.atts["transform"] + fmt.Sprintf(" translate(%d, %d)", x, y)
}

func (a *Att) Rotate(angle float64) {
	a.atts["transform"] = a.atts["transform"] + fmt.Sprintf(" rotate(%g)", angle)
}

func (a *Att) Scale(x, y float64) {
	a.atts["transform"] = a.atts["transform"] + fmt.Sprintf(" scale(%g, %g)", x, y)
}

func (a *Att) ID (id string) {
	a.atts["id"] = id
}

// Creates child group nested from s
func (s *SVG) newGroup(tag string, attributes map[string]string) *SVG {
	g := SVG{tag: tag, atts: make(map[string]string, len(attributes))}

	// Copy map
	for k, v := range attributes {
		g.atts[k] = v
	}

	s.mids = append(s.mids, &g)
	return &g
}

// Create new SVG object to write to
func New(width, height int) *SVG {
	return &SVG{data: svgInit, mids: make([]*SVG, 0), comments: make([]string, 0), atts: map[string]string{"width": fmt.Sprint(width), "height": fmt.Sprint(height)}}
}

// Flush svg in order to write the SVG
func (s *SVG) Flush(w io.Writer) error {
	if s.tag != "" {
		return errors.New("I will only write from the outermost svg element")
	}
	s.tag = "svg"
	for k, v := range defaultNamespace {
		s.atts[k] = v
	}

	var writeGroupPtr func(*SVG, int)
	writeGroup := func(s *SVG, level int) {
		tabs := bytes.Repeat([]byte("\t"), level)

		// Encode attributes to form att1="val1" att2="val2"... 
		buf := bytes.NewBuffer(nil)
		for k, v := range s.atts {
			fmt.Fprintf(buf, `%s="%v" `, k, v)
		}

		b := buf.Bytes()
		if len(b) > 0 { // Remove superfluous space
			b = b[:len(b)-1]
		}

		wb := func(str string) {
			w.Write(tabs)
			w.Write([]byte(str))
			w.Write(b)
		}

		for _, v := range s.comments {
			w.Write([]byte("<!--" + v + "-->\n"))
		}

		switch {
		case len(s.mids) == 0 && len(s.data) == 0:
			wb("<" + s.tag + " ")
			w.Write([]byte(" />\n"))
		case len(s.mids) != 0:
			wb("<" + s.tag + " ")
			w.Write([]byte(">\n"))
			for _, v := range s.mids {
				writeGroupPtr(v, level+1)
			}
			w.Write(tabs)
			w.Write([]byte("</" + s.tag + ">\n"))
		default:
			wb("<" + s.tag + " ")
			w.Write([]byte(">" + s.data + "</" + s.tag + ">\n"))
		}
	}
	writeGroupPtr = writeGroup
	writeGroup(s, 0)
	return nil
}

func (s *SVG) AddAtt(atts map[string]string) {
	for k, v := range atts {
		s.atts[k] = v
	}
}

// Write Start. Ref http://www.w3.org/TR/SVG11/struct.html#SVGElement
func (s *SVG) Start(width, height int, atts map[string]string) *SVG {
	g := s.newGroup(svg, atts)
	g.atts["width"] = fmt.Sprint(width)
	g.atts["height"] = fmt.Sprint(height)
	return g
}

// Start with viewbox. 
func (s *SVG) StartView(width, height, minX, minY, viewWidth, viewHeight int, atts map[string]string) *SVG {
	g := s.Start(width, height, atts)
	g.atts["viewBox"] = fmt.Sprintf("%d %d %d %d", minX, minY, viewWidth, viewHeight)
	return g
}

func (s *SVG) Symbol(id string, atts map[string]string) *SVG {
	g := s.newGroup(sym, atts)
	g.atts["id"] = id
	return g
}

func (s *SVG) SymbolWithViewbox(id string, minX, minY, viewWidth, viewHeight int, atts map[string]string) *SVG {
	g := s.Symbol(id, atts)
	g.atts["viewBox"] = fmt.Sprintf("%d %d %d %d", minX, minY, viewWidth, viewHeight)
	return g
}

func (s *SVG) Use(id string, atts map[string]string) *SVG {
	g := s.newGroup(use, atts)
	g.atts["xlink:href"] = "#" + id
	return g
}

func (s *SVG) Comment(comment string) {
	s.comments = append(s.comments, comment)
}

// Create title
func (s *SVG) Title(t string) *SVG {
	g := s.newGroup(title, nil)
	g.data = t
	return g
}

// Create description
func (s *SVG) Desc(d string) *SVG {
	g := s.newGroup(desc, nil)
	g.data = d
	return g
}

// Create definitions
func (s *SVG) Def() *SVG {
	return s.newGroup(defs, nil)
}

// Create marker
func (s *SVG) Marker(id string, atts map[string]string) *SVG {
	g := s.newGroup(marker, atts)
	g.atts["id"] = id
	return g
}

// Create group
func (s *SVG) G(atts map[string]string) *SVG {
	return s.newGroup(g, atts)
}

// Create group with id
func (s *SVG) GID(id string, atts map[string]string) *SVG {
	g := s.G(atts)
	g.atts["id"] = id
	return g
}

// Translate coordinate system
func (s *SVG) Translate(x, y int) *SVG {
	atts := map[string]string{"transform": fmt.Sprintf("translate(%d, %d)", x, y)}
	return s.newGroup(g, atts)
}

// Scale coordinate system
func (s *SVG) Scale(x, y float64) *SVG {
	atts := map[string]string{"transform": fmt.Sprintf("scale(%g, %g)", x, y)}
	return s.newGroup(g, atts)
}

// Draw circle
func (s *SVG) Circle(x, y, r int, atts map[string]string) *SVG {
	g := s.newGroup(circle, atts)
	g.atts["cx"] = fmt.Sprint(x)
	g.atts["cy"] = fmt.Sprint(y)
	g.atts["r"] = fmt.Sprint(r)
	return g
}

// Draw rectangle
func (s *SVG) Rect(x, y, width, height int, atts map[string]string) *SVG {
	g := s.newGroup(rect, atts)
	g.atts["x"] = fmt.Sprint(x)
	g.atts["y"] = fmt.Sprint(y)
	g.atts["width"] = fmt.Sprint(width)
	g.atts["height"] = fmt.Sprint(height)
	return g
}

// Draw line
func (s *SVG) Line(x1, y1, x2, y2 int, atts map[string]string) *SVG {
	g := s.newGroup(line, atts)
	g.atts["x1"] = fmt.Sprint(x1)
	g.atts["x2"] = fmt.Sprint(x2)
	g.atts["y1"] = fmt.Sprint(y1)
	g.atts["y2"] = fmt.Sprint(y2)
	return g
}

// Draw polyline
func (s *SVG) Polyline(x, y []float64, atts map[string]string) (error, *SVG) {
	switch {
	case len(x) != len(y):
		return errors.New("length of x and y data is not equal"), nil
	case len(x) == 0:
		return errors.New("length of x is zero"), nil
	}

	encloseData := func(x, y float64) string {
		return fmt.Sprintf("%f, %f ", x, y)
	}

	data := make([]byte, 0, len(x)*4)
	for i := range x {
		data = append(data, []byte(encloseData(x[i], y[i]))...)
	}

	g := s.newGroup(polyline, atts)
	g.atts["points"] = string(data)
	return nil, g
}

// Draw text
func (s *SVG) Text(x, y int, text string, atts map[string]string) *SVG {
	g := s.newGroup("text", atts)
	g.data = text
	g.atts["x"] = fmt.Sprint(x)
	g.atts["y"] = fmt.Sprint(y)

	return g
}

const (
	verticalShift = iota
	horizontalShift
)

func findShiftAndScale(i []float64, length int) (scale, shift float64, err error) {
	switch {
	case i == nil:
		err = errors.New("Received invalid argument: Argument is nil")
	case len(i) == 0:
		err = errors.New("Received invalid argument: Argument is zero")
	}
	if err != nil {
		return
	}
	Max, Min := helpers.Max(i), helpers.Min(i)

	scale = float64(length*90/100) / (Max[0] - Min[0])
	shift = -Min[0]
	return
}

func shiftAndScale(i []float64, shift, scale float64) (scaled []int) {
	scaled = make([]int, len(i))

	for ind, v := range i {
		scaled[ind] = int(helpers.Round((v+shift)*scale, 0))
	}

	return
}

func (s *SVG) Label(x1, y1, x2, y2 int, vals []float64, cntGrids int, atts map[string]string) {
	xDiff := (x2 - x1)
	yDiff := (y2 - y1)
	length := math.Sqrt(float64(xDiff * xDiff + yDiff * yDiff))
	angle := math.Atan(float64(xDiff) / float64(yDiff))
	g := s.GID("label", atts)
	g.AddAtt(map[string]string{"transform": fmt.Sprintf("rotate(%g)", 180 * angle / math.Pi), "font-family": "sans-serif", "font-size": "14pt", "fill": "black"})
	pos := 0.0
	val := vals[0]
	
	posIncr := length / float64(cntGrids)
	valIncr := (vals[len(vals)-1] - vals[0]) / float64(cntGrids)
	
	for i := 0; i < cntGrids; i++ {
		g.Text(int(helpers.Round(pos, 0)), 0, fmt.Sprintf("%.2f", val), nil)
		val += valIncr
		pos += posIncr
	}
}

// Paint a diagram. Prerequisites: X must be sorted.
func (s *SVG) Diagram(width, height int, xVals, yVals []float64) (error, *SVG) {
	
	if len(xVals) != len(yVals) {
		return errors.New("Got data pair with inconsistent length. Len x: " + fmt.Sprint(len(xVals)) + "\tlen y: " + fmt.Sprint(len(xVals))), nil
	}

	minWidth, minHeight := 100, 100

	switch {
	case width < 100:
		return errors.New("Width is below minimum width, which is " + fmt.Sprint(minWidth)), nil
	case height < 100:
		return errors.New("Height is below minimum height, which is " + fmt.Sprint(minHeight)), nil
	}
	
	textRoom := 70
	cntGrids := 10
	dWidth, dHeight := width - textRoom, height - textRoom
	v := s.Start(width, height, nil)
	v.AddAtt(map[string]string{"id": "diagram"})
	v.Translate(textRoom, height - textRoom / 2).Label(textRoom / 2, 0, textRoom / 2, height, xVals, cntGrids, nil)
	v.Translate(textRoom / 2, 0).Label(textRoom, height - textRoom / 2, width, height - textRoom / 2, yVals, cntGrids, nil)
	defer v.Rect(0, 0, width, height, map[string]string{"stroke": "grey", "stroke-width": "2", "fill": "none"})
	g := v.GID("plot", map[string]string{"transform": fmt.Sprintf("translate(%d, %d)", textRoom, 0)})
	d := g.StartView(dWidth, dHeight, 0, 0, dWidth, dHeight, map[string]string{"id": "plot", "preserveAspectRatio": "xMinyMax meet"})
	defer d.Rect(0, 0, dWidth, dHeight, map[string]string{"stroke": "grey", "stroke-width": "1", "fill": "none"})
	
	d.Grid(0, 0, dWidth, dHeight, cntGrids, map[string]string{"stroke": "grey", "stroke-width": "1"})
	d.Def().Marker("polyline-midmarker", map[string]string{"style": "overflow:visible", "fill": "black"}).Circle(0, 0, 2, map[string]string{"fill": "none", "stroke": "black"})
	
	f := d.Translate(textRoom, 0)
	e := f.StartView(dWidth, dHeight, 0, 0, dWidth, dHeight, nil)
	e.AddAtt(map[string]string{"presereveAspectRatio": "xMaxYMin middle"})
	i := e.Scale(1, -1)
	j := i.G(map[string]string{"fill": "none", "stroke": "red", "stroke-width": "1"})
	err, k := j.Polyline(xVals, yVals, nil)
	return err, k
	//xShifts := make([]float64, len(vals))
	//yShifts := make([]float64, len(vals))
	//xScales := make([]float64, len(vals))
	//yScales := make([]float64, len(vals))

	//for i, xyPair := range vals {
	//xShifts[i], xScales[i], err = findShiftAndScale(xyPair[0], width)
	//if err != nil {
	//return
	//}

	//yShifts[i], yScales[i], err = findShiftAndScale(xyPair[1], height)
	//if err != nil {
	//return
	//}
	//}

	//scales := helpers.Min(xScales, yScales)
	//shifts := helpers.Min(xShifts, yShifts)

	//l := d.Scale(1, -1).G(`fill="none"`)
	//for _, xyPair := range vals {
	//xScales := shiftAndScale(xyPair[0], scales[0],shifts[0])
	//yScales := shiftAndScale(xyPair[1], scales[1], shifts[1])
	//l.PolyLine(xScales, yScales)
	//}
}

func (s *SVG) Grid(x, y, width, height, cntGrids int, atts map[string]string) {
	g := s.GID("grid", atts)
	d := g.Def()

	vLine := "vLine"
	hLine := "hLine"
	d.GID(vLine, nil).Line(0, 0, 0, height, nil)
	d.GID(hLine, nil).Line(0, 0, width, 0, nil)
	
	gridSizeX := float64(width) / float64(cntGrids)
	gridSizeY := float64(height) / float64(cntGrids)
	
	for ix := float64(x); ix <= float64(x + width); ix += gridSizeX {
		g.Use(vLine, map[string]string{"x": fmt.Sprintf("%.0f", ix)})
	}

	for iy := float64(y); iy <= float64(y + height); iy += gridSizeY {
		g.Use(hLine, map[string]string{"y": fmt.Sprintf("%.0f", iy)})
	}
}

const (
	Aliceblue            = "aliceblue"            // rgb(240, 248, 255)
	Antiquewhite         = "antiquewhite"         // rgb(250, 235, 215)
	Aqua                 = "aqua"                 // rgb(0, 255, 255)
	Aquamarine           = "aquamarine"           // rgb(127, 255, 212)
	Azure                = "azure"                // rgb(240, 255, 255)
	Beige                = "beige"                // rgb(245, 245, 220)
	Bisque               = "bisque"               // rgb(255, 228, 196)
	Black                = "black"                // rgb(0, 0, 0)
	Blanchedalmond       = "blanchedalmond"       // rgb(255, 235, 205)
	Blue                 = "blue"                 // rgb(0, 0, 255)
	Blueviolet           = "blueviolet"           // rgb(138, 43, 226)
	Brown                = "brown"                // rgb(165, 42, 42)
	Burlywood            = "burlywood"            // rgb(222, 184, 135)
	Cadetblue            = "cadetblue"            // rgb(95, 158, 160)
	Chartreuse           = "chartreuse"           // rgb(127, 255, 0)
	Chocolate            = "chocolate"            // rgb(210, 105, 30)
	Coral                = "coral"                // rgb(255, 127, 80)
	Cornflowerblue       = "cornflowerblue"       // rgb(100, 149, 237)
	Cornsilk             = "cornsilk"             // rgb(255, 248, 220)
	Crimson              = "crimson"              // rgb(220, 20, 60)
	Cyan                 = "cyan"                 // rgb(0, 255, 255)
	Darkblue             = "darkblue"             // rgb(0, 0, 139)
	Darkcyan             = "darkcyan"             // rgb(0, 139, 139)
	Darkgoldenrod        = "darkgoldenrod"        // rgb(184, 134, 11)
	Darkgray             = "darkgray"             // rgb(169, 169, 169)
	Darkgreen            = "darkgreen"            // rgb(0, 100, 0)
	Darkgrey             = "darkgrey"             // rgb(169, 169, 169)
	Darkkhaki            = "darkkhaki"            // rgb(189, 183, 107)
	Darkmagenta          = "darkmagenta"          // rgb(139, 0, 139)
	Darkolivegreen       = "darkolivegreen"       // rgb(85, 107, 47)
	Darkorange           = "darkorange"           // rgb(255, 140, 0)
	Darkorchid           = "darkorchid"           // rgb(153, 50, 204)
	Darkred              = "darkred"              // rgb(139, 0, 0)
	Darksalmon           = "darksalmon"           // rgb(233, 150, 122)
	Darkseagreen         = "darkseagreen"         // rgb(143, 188, 143)
	Darkslateblue        = "darkslateblue"        // rgb(72, 61, 139)
	Darkslategray        = "darkslategray"        // rgb(47, 79, 79)
	Darkslategrey        = "darkslategrey"        // rgb(47, 79, 79)
	Darkturquoise        = "darkturquoise"        // rgb(0, 206, 209)
	Darkviolet           = "darkviolet"           // rgb(148, 0, 211)
	Deeppink             = "deeppink"             // rgb(255, 20, 147)
	Deepskyblue          = "deepskyblue"          // rgb(0, 191, 255)
	Dimgray              = "dimgray"              // rgb(105, 105, 105)
	Dodgerblue           = "dodgerblue"           // rgb(30, 144, 255)
	Firebrick            = "firebrick"            // rgb(178, 34, 34)
	Floralwhite          = "floralwhite"          // rgb(255, 250, 240)
	Forestgreen          = "forestgreen"          // rgb(34, 139, 34)
	Fuchsia              = "fuchsia"              // rgb(255, 0, 255)
	Gainsboro            = "gainsboro"            // rgb(220, 220, 220)
	Ghostwhite           = "ghostwhite"           // rgb(248, 248, 255)
	Gold                 = "gold"                 // rgb(255, 215, 0)
	Goldenrod            = "goldenrod"            // rgb(218, 165, 32)
	Gray                 = "gray"                 // rgb(128, 128, 128)
	Grey                 = "grey"                 // rgb(128, 128, 128)
	Green                = "green"                // rgb(0, 128, 0)
	Greenyellow          = "greenyellow"          // rgb(173, 255, 47)
	Honeydew             = "honeydew"             // rgb(240, 255, 240)
	Hotpink              = "hotpink"              // rgb(255, 105, 180)
	Indianred            = "indianred"            // rgb(205, 92, 92)
	Indigo               = "indigo"               // rgb(75, 0, 130)
	Ivory                = "ivory"                // rgb(255, 255, 240)
	Khaki                = "khaki"                // rgb(240, 230, 140)
	Lavender             = "lavender"             // rgb(230, 230, 250)
	Lavenderblush        = "lavenderblush"        // rgb(255, 240, 245)
	Lawngreen            = "lawngreen"            // rgb(124, 252, 0)
	Lemonchiffon         = "lemonchiffon"         // rgb(255, 250, 205)
	Lightblue            = "lightblue"            // rgb(173, 216, 230)
	Lightcoral           = "lightcoral"           // rgb(240, 128, 128)
	Lightcyan            = "lightcyan"            // rgb(224, 255, 255)
	Lightgoldenrodyellow = "lightgoldenrodyellow" // rgb(250, 250, 210)
	Lightgray            = "lightgray"            // rgb(211, 211, 211)
	Lightgreen           = "lightgreen"           // rgb(144, 238, 144)
	Lightgrey            = "lightgrey"            // rgb(211, 211, 211)
	Lightpink            = "lightpink"            // rgb(255, 182, 193)
	Lightsalmon          = "lightsalmon"          // rgb(255, 160, 122)
	Lightseagreen        = "lightseagreen"        // rgb(32, 178, 170)
	Lightskyblue         = "lightskyblue"         // rgb(135, 206, 250)
	Lightslategray       = "lightslategray"       // rgb(119, 136, 153)
	Lightslategrey       = "lightslategrey"       // rgb(119, 136, 153)
	Lightsteelblue       = "lightsteelblue"       // rgb(176, 196, 222)
	Lightyellow          = "lightyellow"          // rgb(255, 255, 224)
	Lime                 = "lime"                 // rgb(0, 255, 0)
	Limegreen            = "limegreen"            // rgb( 50, 205, 50)
	Linen                = "linen"                // rgb(250, 240, 230)
	Magenta              = "magenta"              // rgb(255, 0, 255)
	Maroon               = "maroon"               // rgb(128, 0, 0)
	Mediumaquamarine     = "mediumaquamarine"     // rgb(102, 205, 170)
	Mediumblue           = "mediumblue"           // rgb(0, 0, 205)
	Mediumorchid         = "mediumorchid"         // rgb(186, 85, 211)
	Mediumpurple         = "mediumpurple"         // rgb(147, 112, 219)
	Mediumseagreen       = "mediumseagreen"       // rgb(60, 179, 113)
	Mediumslateblue      = "mediumslateblue"      // rgb(123, 104, 238)
	Mediumspringgreen    = "mediumspringgreen"    // rgb(0, 250, 154)
	Mediumturquoise      = "mediumturquoise"      // rgb(72, 209, 204)
	Mediumvioletred      = "mediumvioletred"      // rgb(199, 21, 133)
	Midnightblue         = "midnightblue"         // rgb(25, 25, 112)
	Mintcream            = "mintcream"            // rgb(245, 255, 250)
	Mistyrose            = "mistyrose"            // rgb(255, 228, 225)
	Moccasin             = "moccasin"             // rgb(255, 228, 181)
	Navajowhite          = "navajowhite"          // rgb(255, 222, 173)
	Navy                 = "navy"                 // rgb(0, 0, 128)
	Oldlace              = "oldlace"              // rgb(253, 245, 230)
	Olive                = "olive"                // rgb(128, 128, 0)
	Olivedrab            = "olivedrab"            // rgb(107, 142, 35)
	Orange               = "orange"               // rgb(255, 165, 0)
	Orangered            = "orangered"            // rgb(255, 69, 0)
	Orchid               = "orchid"               // rgb(218, 112, 214)
	Palegoldenrod        = "palegoldenrod"        // rgb(238, 232, 170)
	Palegreen            = "palegreen"            // rgb(152, 251, 152)
	Paleturquoise        = "paleturquoise"        // rgb(175, 238, 238)
	Palevioletred        = "palevioletred"        // rgb(219, 112, 147)
	Papayawhip           = "papayawhip"           // rgb(255, 239, 213)
	Peachpuff            = "peachpuff"            // rgb(255, 218, 185)
	Peru                 = "peru"                 // rgb(205, 133, 63)
	Pink                 = "pink"                 // rgb(255, 192, 203)
	Plum                 = "plum"                 // rgb(221, 160, 221)
	Powderblue           = "powderblue"           // rgb(176, 224, 230)
	Purple               = "purple"               // rgb(128, 0, 128)
	Red                  = "red"                  // rgb(255, 0, 0)
	Rosybrown            = "rosybrown"            // rgb(188, 143, 143)
	Royalblue            = "royalblue"            // rgb(65, 105, 225)
	Saddlebrown          = "saddlebrown"          // rgb(139, 69, 19)
	Salmon               = "salmon"               // rgb(250, 128, 114)
	Sandybrown           = "sandybrown"           // rgb(244, 164, 96)
	Seagreen             = "seagreen"             // rgb(46, 139, 87)
	Seashell             = "seashell"             // rgb(255, 245, 238)
	Sienna               = "sienna"               // rgb(160, 82, 45)
	Silver               = "silver"               // rgb(192, 192, 192)
	Skyblue              = "skyblue"              // rgb(135, 206, 235)
	Slateblue            = "slateblue"            // rgb(106, 90, 205)
	Slategray            = "slategray"            // rgb(112, 128, 144)
	Slategrey            = "slategrey"            // rgb(112, 128, 144)
	Snow                 = "snow"                 // rgb(255, 250, 250)
	Springgreen          = "springgreen"          // rgb(0, 255, 127)
	Steelblue            = "steelblue"            // rgb(70, 130, 180)
	Tan                  = "tan"                  // rgb(210, 180, 140)
	Teal                 = "teal"                 // rgb(0, 128, 128)
	Thistle              = "thistle"              // rgb(216, 191, 216)
	Tomato               = "tomato"               // rgb(255, 99, 71)
	Turquoise            = "turquoise"            // rgb(64, 224, 208)
	Violet               = "violet"               // rgb(238, 130, 238)
	Wheat                = "wheat"                // rgb(245, 222, 179)
	White                = "white"                // rgb(255, 255, 255)
	Whitesmoke           = "whitesmoke"           // rgb(245, 245, 245)
	Yellow               = "yellow"               // rgb(255, 255, 0)
	Yellowgreen          = "yellowgreen"          // rgb(154, 205, 50)
)