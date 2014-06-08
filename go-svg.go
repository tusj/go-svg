// smartSVG implements a subset of the SVG standard to write good, well-formed svg in a simple fashion.
// Tags which encloses groups returns a new *SVG object which is possible to write to.
// Write writes the coded svg to the writer provided.
package smartSVG

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
)

func fold(init float64, f func(float64, float64) float64, vals ...float64) float64 {
	for _, v := range vals {
		init = f(init, v)
	}
	return init
}

func reduce(f func(float64, float64) float64, vals ...float64) float64 {
	a := vals[0]
	for _, b := range vals {
		a = f(a, b)
	}

	return a
}

func max(vals ...float64) float64 {
	return reduce(func(a, b float64) float64 {
		if a > b {
			return a
		}
		return b
	}, vals...)
}

func min(vals ...float64) float64 {
	return reduce(func(a, b float64) float64 {
		if a < b {
			return a
		}
		return b
	})
}

func average(vals ...float64) float64 {
	return reduce(func(a, b float64) float64 {
		return a + b
	}, vals...) / float64(len(vals))
}

func round(a float64) int {
	if a > 0.0 {
		return int(a + 0.5)
	}
	return int(a - 0.5)
}

// Sort atts easily
type keySorter struct {
	keys []string
	vals []interface{}
}

func newKeySorter(a Att) *keySorter {
	ret := new(keySorter)
	ret.keys = make([]string, len(a))
	ret.vals = make([]interface{}, len(a))
	i := 0
	for k, v := range a {
		ret.keys[i] = k
		ret.vals[i] = v
		i++
	}
	return ret
}

func (vs *keySorter) Len() int           { return len(vs.vals) }
func (vs *keySorter) Less(i, j int) bool { return vs.keys[i] < vs.keys[j] }
func (vs *keySorter) Swap(i, j int) {
	vs.vals[i], vs.vals[j] = vs.vals[j], vs.vals[i]
	vs.keys[i], vs.keys[j] = vs.keys[j], vs.keys[i]
}

func (a Att) Sort() ([]string, []interface{}) {
	sortable := newKeySorter(a)
	sort.Sort(sortable)
	return sortable.keys, sortable.vals
}

type Data struct {
	X, Y []float64
}

type Att map[string]interface{}

func (a Att) String() string {
	str := ""
	for k, v := range a {
		str += k + "=\"" + fmt.Sprint(v) + `" `
	}
	return str[:len(str)]
}

func (a Att) SetPos(x, y int) {
	a["x"] = fmt.Sprint(x)
	a["y"] = fmt.Sprint(y)
}

func (a Att) SetSize(width, height int) {
	a["width"] = fmt.Sprint(width)
	a["height"] = fmt.Sprint(height)
}

// Holds one svg group with nodes of children inside group.
// data is used to hold data encapsuled within xml tag. If non-empty, then mids should be empty.
type SVG struct {
	tag         string
	a           Att
	data        string
	comments    []string
	mids        []*SVG
	parent      *SVG
	declaration string
}

func (s *SVG) String() string {
	//	return fmt.Sprint("tag: ", s.tag, "Att: ", s.a, "data: ", s.data, "comments: ", s.comments, "mids: ", s.mids, "parent: ", s.parent, "declaration: ", s.declaration)
	//	return fmt.Sprint("tag: ", s.tag, " atts: ", s.a)
	buf := bytes.NewBuffer(nil)
	s.Write(buf)
	return buf.String()
}

const (
	svgInit = `<?xml version="1.0"?>` + "\n" + `<!-- Generated by smartSVG -->`
)

var defaultNamespace map[string]string

func init() {
	defaultNamespace = map[string]string{"xmlns": "http://www.w3.org/2000/svg", "xmlns:xlink": "http://www.w3.org/1999/xlink"}
}

// Climb upwards xml-tree to see if p is parent of s
func (s *SVG) IsParent(p *SVG) bool {
	g := s
	for g.parent != p && g.parent != nil {
		g = g.parent
	}
	if g.parent == p {
		return true
	}
	return false
}

// Create new SVG object to write to
func New(width, height int) *SVG {
	return &SVG{tag: "svg", data: svgInit, mids: make([]*SVG, 0), comments: make([]string, 0), a: Att{"preserveAspectRatio": "xMinYmin meet", "viewBox": "0 0 " + fmt.Sprint(width, " ", height)}}
}

// Creates child group nested from s
func (s *SVG) newGroup(tag string, a Att) *SVG {
	g := SVG{tag: tag, a: make(Att, len(a))}

	// Copy attributes
	for k, v := range a {
		g.a[k] = v
	}
	g.parent = s

	s.mids = append(s.mids, &g)
	return &g
}

// Add svg group to other svg group. If svg is parent, adding will break tree structure and is therefore forbidden.
// Adding svg group to other svg group if svg is child of s is ok because it only creates forward links.
func (s *SVG) Add(svg *SVG) error {
	if s.IsParent(svg) {
		return errors.New("svg is parent of s. Cannot Add svg as subgroup of s in order to prevent creation of cycles")
	}
	s.mids = append(s.mids, svg)
	return nil
}

// Encapsulate children of s within svg.
func (s *SVG) Insert(svg *SVG) (*SVG, error) {

	if s.IsParent(svg) {
		return nil, errors.New("svg is parent of s. Cannot insert svg as subgroup of s in order to prevent creation of cycles")
	}
	if s.Find(svg) {
		return nil, errors.New("svg is child of s. Cannot insert svg as subgroup of s in order to prevent creation of cycles")
	}
	// Add children of s to children of svg
	svg.mids = append(svg.mids, s.mids...)

	// Substitute children of s with children of svg
	s.mids = s.mids[:0]
	s.mids = append(s.mids, svg)
	return svg, nil
}

// Write svg file to w
func (s *SVG) Write(w io.Writer) error {
	//	if s.tag != "" {
	//		return errors.New("I will only write from the outermost svg element")
	//	}
	for k, v := range defaultNamespace {
		s.a[k] = v
	}

	var writeGroupPtr func(*SVG, int)
	writeGroup := func(s *SVG, level int) {
		tabs := bytes.Repeat([]byte("\t"), level)

		// Encode attributes to form att1="val1" att2="val2"...
		buf := bytes.NewBuffer(nil)
		keys, vals := s.a.Sort()
		for i := range keys {
			fmt.Fprintf(buf, `%s="%v" `, keys[i], vals[i])
		}

		atts := buf.Bytes()
		if len(atts) > 0 { // Remove superfluous space
			atts = atts[:len(atts)-1]
		}

		for _, v := range s.comments {
			w.Write([]byte("<!--" + v + "-->\n"))
		}

		w.Write(tabs)
		w.Write([]byte("<" + s.tag))
		if len(atts) != 0 {
			w.Write([]byte(" "))
		}
		w.Write(atts)

		switch {
		case len(s.mids) == 0 && len(s.data) == 0:
			w.Write([]byte(" />\n"))
		case len(s.mids) != 0:
			w.Write([]byte(">\n"))
			for _, v := range s.mids {
				writeGroupPtr(v, level+1)
			}
			w.Write(tabs)
			w.Write([]byte("</" + s.tag + ">\n"))
		default:
			w.Write([]byte(">" + s.data + "</" + s.tag + ">\n"))
		}
	}
	if s.declaration != "" {
		w.Write([]byte(s.declaration))
	}
	writeGroupPtr = writeGroup
	writeGroup(s, 0)
	return nil
}

// Add attributes to group
func (s *SVG) AddAtt(override bool, a ...Att) {
	for _, att := range a {
		for k, v := range att {
			var tmp string
			if !override {
				if w, ok := s.a[k]; ok {
					tmp = fmt.Sprint(w)
				}
				if tmp != "" {
					tmp += " "
				}
			}
			s.a[k] = tmp + fmt.Sprint(v)
		}
	}
}

// Start svg group. Ref http://www.w3.org/TR/SVG11/struct.html#SVGElement
func (s *SVG) Start(width, height int, a Att) *SVG {
	g := s.newGroup("svg", a)
	g.a["width"] = width
	g.a["height"] = height
	return g
}

// Start svg group with viewbox.
func (s *SVG) StartView(width, height, minX, minY, viewWidth, viewHeight int, a Att) *SVG {
	g := s.Start(width, height, a)
	g.a["viewBox"] = fmt.Sprintf("%d %d %d %d", minX, minY, viewWidth, viewHeight)
	return g
}

// Set information about external style sheet
func (s *SVG) SetCSS(href string) {
	s.declaration = `<?xml-stylesheet type="text/css" href="` + href + "\" ?>\n"
}

// Create symbol group
func (s *SVG) Symbol(id string, a Att) *SVG {
	g := s.newGroup("symbol", a)
	g.a["id"] = id
	return g
}

// Create symbol group with viewbox attribute
func (s *SVG) SymbolWithViewbox(id string, minX, minY, viewWidth, viewHeight int, a Att) *SVG {
	g := s.Symbol(id, a)
	g.a["viewBox"] = fmt.Sprintf("%d %d %d %d", minX, minY, viewWidth, viewHeight)
	return g
}

// Use element given in id
func (s *SVG) Use(id string, a Att) *SVG {
	g := s.newGroup("use", a)
	g.a["xlink:href"] = "#" + id
	return g
}

// Add comment to group
func (s *SVG) Comment(comment string) {
	s.comments = append(s.comments, comment)
}

// Create title
func (s *SVG) Title(title string) *SVG {
	g := s.newGroup("title", nil)
	g.data = title
	return g
}

// Create description
func (s *SVG) Desc(description string) *SVG {
	g := s.newGroup("desc", nil)
	g.data = description
	return g
}

// Create definitions
func (s *SVG) Def() *SVG {
	return s.newGroup("defs", nil)
}

// Create marker
func (s *SVG) Marker(id string, a Att) *SVG {
	g := s.newGroup("marker", a)
	g.a["id"] = id
	return g
}

// Create group
func (s *SVG) G(a Att) *SVG {
	return s.newGroup("g", a)
}

// Create group with id
func (s *SVG) GID(id string, a Att) *SVG {
	g := s.G(a)
	g.a["id"] = id
	return g
}

// Set ID attribute of group
func (s *SVG) ID(id string) {
	s.a["id"] = id
}

// Create group with translation of coordinate system
func (s *SVG) Translate(x, y float64) *SVG {
	a := Att{"transform": fmt.Sprintf("translate(%g, %g)", x, y)}
	return s.newGroup("g", a)
}

// Create group with scale of coordinate system
func (s *SVG) Scale(x, y float64) *SVG {
	a := Att{"transform": fmt.Sprintf("scale(%g, %g)", x, y)}
	return s.newGroup("g", a)
}

// Make translation attribute of coordinate system
func Translate(x, y float64) Att {
	return Att{"transform": fmt.Sprintf("translate(%g, %g)", x, y)}
}

// Make scale attribute of coordinate system
func Scale(x, y float64) Att {
	return Att{"transform": fmt.Sprintf("scale(%g, %g)", x, y)}
}

func ViewBox(x, y, xMin, yMin int) Att {
	return Att{"viewBox": fmt.Sprintf("%d %d %d %d", x, y, xMin, yMin)}
}

func SumAtts(a ...Att) Att {
	ret := make(Att)
	for _, c := range a {
		for k, v := range c {
			ret[k] = v
		}
	}
	return ret
}

// Make object clickable
func (s *SVG) A(URL string) (*SVG, error) {
	g := SVG{tag: "a", a: Att{"xlink:href": URL, "xlink:show": "replace", "target": "_parent"}}
	return s.Insert(&g)
}

// Draw circle
func (s *SVG) Circle(x, y, r int, a Att) *SVG {
	g := s.newGroup("circle", a)
	g.a["cx"] = fmt.Sprint(x)
	g.a["cy"] = fmt.Sprint(y)
	g.a["r"] = fmt.Sprint(r)
	return g
}

// Draw rectangle
func (s *SVG) Rect(x, y, width, height int, a Att) *SVG {
	g := s.newGroup("rect", a)
	g.a["x"] = fmt.Sprint(x)
	g.a["y"] = fmt.Sprint(y)
	g.a["width"] = fmt.Sprint(width)
	g.a["height"] = fmt.Sprint(height)
	return g
}

// Draw line
func (s *SVG) Line(x1, y1, x2, y2 int, a Att) *SVG {
	g := s.newGroup("line", a)
	g.a["x1"] = fmt.Sprint(x1)
	g.a["x2"] = fmt.Sprint(x2)
	g.a["y1"] = fmt.Sprint(y1)
	g.a["y2"] = fmt.Sprint(y2)
	return g
}

// Draw polyline
func (s *SVG) Polyline(d Data, a Att) (*SVG, error) {
	switch {
	case len(d.X) != len(d.Y):
		return nil, errors.New("length of data pair is not equal")
	case len(d.X) == 0:
		return nil, errors.New("length of data is zero")
	}

	encloseData := func(x, y float64) string {
		return fmt.Sprintf("%f,%f ", x, y)
	}

	// Create formatted data plots
	data := make([]byte, 0, len(d.X))
	for i := range d.X {
		data = append(data, []byte(encloseData(d.X[i], d.Y[i]))...)
	}

	// Draw the polyline
	g := s.newGroup("polyline", a)
	g.a["points"] = string(data)
	return g, nil
}

// Draw text
func (s *SVG) Text(x, y int, text string, a Att) *SVG {
	g := s.newGroup("text", a)
	g.data = text
	g.a["x"] = fmt.Sprint(x)
	g.a["y"] = fmt.Sprint(y)

	return g
}

func (s *SVG) Image(x, y, width, height int, link string, a Att) *SVG {
	g := s.newGroup("image", a)
	g.a["xlink:href"] = link
	g.a.SetPos(x, y)
	g.a.SetSize(width, height)

	return g
}

// Write text on line from p1 to p2, with cntGrids values as given in vals.
// Prerequisites: vals[] is linear
func (s *SVG) Label(x1, y1, x2, y2 int, vals []float64, cntGrids int, a Att) {
	g := s.GID("label", a)
	g.AddAtt(false, Att{"fill": "black"})
	xDiff := float64(x2 - x1)
	yDiff := float64(y2 - y1)
	angle := math.Abs(math.Atan(yDiff / xDiff))

	// Find increments
	max := max(vals...)
	min := min(vals...)
	valIncr := (max - min) / float64(cntGrids)
	xIncr := float64(xDiff) / float64(cntGrids) * math.Cos(angle)
	yIncr := float64(yDiff) / float64(cntGrids) * math.Sin(angle)

	// Start values
	val := min
	x := float64(x1)
	y := float64(y1)

	// Draw text
	for i := 0; i <= cntGrids; i++ {
		g.Text(round(x), round(y), fmt.Sprintf("%.2f", val), nil)
		val += valIncr
		x += xIncr
		y += yIncr
	}
}

// Add possibility to view legend

// Draw a grid with cntGrids horizontal and vertical lines
func (s *SVG) Grid(x, y, width, height, cntGrids int, a Att) *SVG {
	vLine := "vLine"
	hLine := "hLine"

	// Create group with defs
	g := s.GID("grid", a)
	d := g.Def()
	d.Line(0, 0, 0, height, Att{"id": vLine})
	d.Line(0, 0, width, 0, Att{"id": hLine})

	ix, iy := float64(x), float64(y)
	gridSizeX := float64(width) / float64(cntGrids)
	gridSizeY := float64(height) / float64(cntGrids)

	// Draw vertical and horizontal lines using the defined lines
	for i := 0; i <= cntGrids; i++ {
		g.Use(vLine, Att{"x": fmt.Sprintf("%.0f", ix)})
		g.Use(hLine, Att{"y": fmt.Sprintf("%.0f", iy)})
		ix += gridSizeX
		iy += gridSizeY
	}
	return g
}

// Display modes
const (
	Column = iota
	Continuous
)

// Paint a diagram
func (s *SVG) Diagram(x, y, width, height int, d Data, title string, display int) (*SVG, error) {
	switch display {
	case Column, Continuous:
		break
	default:
		return nil, errors.New("Got unknown display mode")
	}
	switch {
	case len(d.X) == 0:
	case len(d.Y) == 0:
		return nil, errors.New("Got empty data set")
	case len(d.X) != len(d.Y):
		return nil, errors.New("Got data pair with uneven length")
	}

	last := d.X[0]
	for _, v := range d.X {
		if v < last {
			return nil, errors.New("Xvals is not sorted.")
		}
		last = v
	}

	titleHeight := 25
	textHeight := 10
	textRoomX := 70
	textRoomY := textRoomX / 3
	plotMargin := 2
	cntGrids := 10
	dWidth, dHeight := width-textRoomX, height-textRoomY

	top := s.GID("diagram", Att{"width": width, "height": height})
	top.AddAtt(false, Translate(float64(x), float64(y)))

	// Draw background in order to make whole object clickable, and to set the background of the diagram
	top.Rect(0, 0, width, height, nil)
	top.Text(width/2, 3*titleHeight/4, title, Att{"text-anchor": "middle", "fill": "black", "id": "title"})

	// Draw outer frame
	//defer v.Rect(0, 0, width, height, map[string]string{"stroke": "grey", "stroke-width": "1", "fill": "none"})

	// New group with plot, move to upper right corner
	g := top.G(Translate(float64(textRoomX), float64(titleHeight)))
	g.ID("plot")

	// Mark data on the axes
	alignRight := Att{"text-anchor": "end"}
	alignLeft := Att{"text-anchor": "start"}
	// Vertical
	g.Label(0, height-textRoomY, 0, titleHeight-2*plotMargin, d.Y, cntGrids, alignRight)
	// Horizontal
	g.Label(0, height-textRoomY+textHeight, width-textRoomX, height-textRoomY+textHeight, d.X, cntGrids, alignLeft)

	dd := g.StartView(dWidth, dHeight, 0, 0, dWidth, dHeight, nil)
	cartesian := dd.Translate(0.0, float64(dHeight))
	cartesian.AddAtt(false, Scale(1, -1), Att{"fill": "none"})
	// Draw inner frame almost last
	defer cartesian.Rect(0, 0, dWidth, dHeight, Att{"stroke": "grey", "stroke-width": "3"})

	marginShift := cartesian.Translate(float64(plotMargin), float64(plotMargin))
	marginShift.AddAtt(false, Scale(float64(dWidth-2*plotMargin)/float64(dWidth), float64(dHeight-2*plotMargin)/float64(dHeight)))
	marginShift.Grid(0, 0, dWidth, dHeight, cntGrids, Att{"stroke": "black", "stroke-width": "1"})
	// Finds the scales and shift of data
	resize := func(vals []float64, length int) (scale, shift float64) {
		min := min(vals...) // ignoring err because data has already been tested
		max := max(vals...)
		scale = float64(length) / (max - min)
		shift = -min * scale
		return
	}
	xScale, xShift := resize(d.X, dWidth)
	yScale, yShift := resize(d.Y, dHeight)

	// Scales and shifts the plot
	plot := marginShift.Translate(xShift, yShift)
	plot.ID("data")
	plot.AddAtt(false, Scale(xScale, yScale))

	// Create marker inside defs to be used with plot
	def := plot.Def()
	def.Marker("polyline-midmarker", Att{"viewBox": "0 0 10 10",
		"preserveAspectRatio": "xMidYMid meet",
		"refX":                "5",
		"refY":                "5",
		"stroke":              GetColour(),
		"fill":                "none",
		"orient":              "auto",
		"vector-effect":       "non-scaling-stroke"}).Circle(0, 0, 1, nil)

	// Draws the plot
	att := Att{"fill": "none", "stroke": GetColour(), "vector-effect": "non-scaling-stroke"} //, "marker-mid": "url(#polyline-midmarker)"})
	switch display {
	case Continuous:
		break
	case Column:
		// Create marker which stands as columns
		att["stroke"] = "none"
		att["marker-mid"] = "url(#column-marker)"
		def.Marker("column-marker", Att{"viewBox": "0 0 10 10",
			"preserveAspectRatio": "xMidYMid meet",
			"refX":                "5",
			"refY":                "5",
			"stroke":              GetColour(),
			"fill":                "none",
			"orient":              "fixed",
			"vector-effect":       "non-scaling-stroke"}).Rect(0, 0, 1, 1000, nil)
	}
	_, err := plot.Polyline(d, att)
	return top, err
}

// Test to prevent adding plot to previous plot not working? // Messes up scale when used?
func (s *SVG) AddPlot(d Data, a Att) (*SVG, error) {
	if s.a["id"] != "diagram" || len(s.FindGroups("polyline")) == 0 {
		return nil, errors.New("Will only add plot to existing diagram: Could not find id with diagram nor data with polyline groups")
	}

	var plot *SVG
	if plot = s.FindID("data"); plot == nil {
		return nil, errors.New("Could not find any existing data to add plot with")
	}

	line, err := plot.Polyline(d, SumAtts(Att{"vector-effect": "non-scaling-stroke"}, a))
	if err != nil {
		return nil, err
	}
	line.AddAtt(true, a)

	return line, nil
}

func (s *SVG) Legend(desc ...string) (*SVG, error) {
	if s.a["id"] != "diagram" {
		return nil, errors.New("Will only add legend to diagram")
	}
	data := s.FindGroups("polyline")
	if len(data) != len(desc) {
		return nil, errors.New("Amount of plots found is not the same as the amount of descriptors given. #Data: " + fmt.Sprint(len(data)) + " #Desc: " + fmt.Sprint(len(desc)) + ". Desc is " + fmt.Sprint(desc))
	}

	var (
		transform   string
		pageHeight  int
		lW, lH      int
		p, n        int
		titleHeight int
	)
	if plot := s.FindID("plot"); plot != nil {
		if tr, ok := plot.a["transform"]; ok {
			if transform, ok = tr.(string); !ok {
				return nil, errors.New("Could not decode transform attribute to string")
			}
		} else {
			return nil, errors.New("Could not find transform attribute of plot group")
		}
	} else {
		return nil, errors.New("Could not find plot group of diagram")
	}
	if h := s.a["height"]; h != nil {
		var ok bool
		if pageHeight, ok = h.(int); !ok {
			return nil, errors.New("Could not fetch page height")
		}
	} else {
		return nil, errors.New("Could not find pageHeight attribute")
	}

	p = strings.Index(transform, "(")
	n = strings.Index(transform, ",")
	if p == -1 && n == -1 {
		return nil, errors.New("Decode error: Could not find textWidth for translation of diagram. String searched: " + transform + ". Got p & n: " + fmt.Sprint(p, n))
	} else {
		var err error
		lW, err = strconv.Atoi(transform[p+1 : n])
		if err != nil {
			return nil, err
		}
	}
	p = strings.Index(transform, ")")
	if p != -1 {
		var err error
		titleHeight, err = strconv.Atoi(transform[n+2 : p])
		if err != nil {
			return nil, err
		}
		lH = pageHeight - titleHeight
	} else {
		return nil, errors.New("Decode error: Could not find textHeight for translation of diagram. String searched: " + transform + "Got p & n: " + fmt.Sprint(p, n))
	}

	textHeight := 10
	legendMargin := 5
	lW, lH = lW-2*legendMargin, lH-2*legendMargin
	lW /= 3
	legend := s.GID("legend", SumAtts(Translate(float64(legendMargin), float64(titleHeight+legendMargin)), ViewBox(0, 0, lW, lH)))
	def := legend.Def()

	def.Rect(0, 0, lW, lH/len(data), Att{"id": "legendRect"})
	yDiff := lH / len(data)
	for i := 0; i < len(data); i++ {
		if stroke, ok := data[i].a["stroke"]; ok {
			if colour, ok := stroke.(string); ok {
				legend.Use("legendRect", Att{"fill": colour, "y": yDiff * i})
				t := legend.Text(textHeight/2+yDiff/2, yDiff*i, desc[i], Att{"text-anchor": "middle", "fill": "black"})
				t.AddAtt(false, Att{"transform": fmt.Sprintf("rotate(90, %d, %d)", textHeight/2, yDiff*i)})
			} else {
				return nil, errors.New("Decode error: Could not find stroke colour for data at data element " + fmt.Sprint(i))
			}
		} else {
			return nil, errors.New("Decode error: Could not find stroke attribute for data at data element " + fmt.Sprint(i))
		}
	}
	return legend, nil
}

// Search downwards, return first group with correct id
func (s *SVG) FindID(id string) *SVG {
	if s.a["id"] == id {
		return s
	}
	for _, c := range s.mids {
		if found := c.FindID(id); found != nil {
			return found
		}
	}
	return nil
}

// Search downwards for groups with matching tag
func (s *SVG) FindGroups(tag string) (groups []*SVG) {
	if s.tag == tag {
		groups = append(groups, s)
	}
	for _, c := range s.mids {
		groups = append(groups, c.FindGroups(tag)...)
	}
	return
}

// Search downwards for svg
func (s *SVG) Find(svg *SVG) bool {
	if s == svg {
		return true
	}
	for _, c := range s.mids {
		if c.Find(svg) {
			return true
		}
	}
	return false
}

//func (s *SVG)BodePlot(d []complex128)
