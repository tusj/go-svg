SmartSVG
========
SmartSVG is a work-in-progress SVG library for Go designed with flexibility in mind. It currently implements only a subset of the SVG definition, but all the basics are there.


Files
-----
* smartSVG.go: Library implementation
* constants.go: Colour definition with colour helper functions

Building and Usage
------------------

Usage: (assuming GOPATH is set)

	go get github.com/tusj/smartSVG
	go install github.com/tusj/smartSVG/

You can use godoc to browse the documentation from the command line:

	$ godoc github.com/tusj/smartSVG
	

Example program
---------------
	package main
	
	import (
		svg "github.com/tusj/smartSVG"
		"os"
	)
	
	func main() {
		x := []float64{1, 2, 3, 4, 5}
		y := []float64{1, 4, 9, 16, 25}
			width, height := 500, 500
		drawing := svg.New(width, height)
		posX, posY :=  0, 0
		att := svg.Att{"fill": "none"}
		plot, _ := drawing.Diagram(posX, posY, width, height/2, svg.Data{x, y}, "SVG Example", svg.Continuous)
		
			secondY := []float64{25, 16, 9, 4, 1}
			plot.AddPlot(svg.Data{x, secondY}, svg.Att{"stroke" : svg.GetColour()})
	
			plot.Legend("first", "second")
			plot.AddAtt(false, att)
	
		g := drawing.GID("lowerHalf", svg.Translate(0, float64(height)/2.0))
		defs := g.Def()
			defs.Rect(0, 0, 10, 10, svg.Att{"id" : "blueRectangle", "fill" : "blue"})
			defs.Circle(0, 0, 5, svg.Att{"id" : "yellowCircle", "fill" : "yellow"})
	
		g.Text(width/2, height/8, "Simple SVG example with demonstration of subgroups", svg.Att{"text-anchor": "middle"})
		leftMid := svg.Att{"x" : width/10, "y" : height / 8}
		rightMid := svg.Att{"x" : width - width/10, "y":  height / 8}
		midLower := svg.Att{"x" : width/2, "y" : height / 4}
		g.Use("blueRectangle", leftMid)
		g.Use("blueRectangle", rightMid)
		g.Use("yellowCircle", midLower)	
	
		drawing.Flush(os.Stdout)
	}

### This produces ###

		<svg preserveAspectRatio="xMinYmin meet" viewBox="0 0 500 500" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
		<g fill="none" height="250" id="diagram" transform="translate(0, 0)" width="500">
			<rect height="250" width="500" x="0" y="0" />
			<text fill="black" id="title" text-anchor="middle" x="250" y="18">SVG Example</text>
			<g id="plot" transform="translate(70, 25)">
				<g fill="black" id="label" text-anchor="end">
					<text x="0" y="227">1.00</text>
					<text x="0" y="207">3.40</text>
					<text x="0" y="186">5.80</text>
					<text x="0" y="166">8.20</text>
					<text x="0" y="145">10.60</text>
					<text x="0" y="125">13.00</text>
					<text x="0" y="104">15.40</text>
					<text x="0" y="83">17.80</text>
					<text x="0" y="63">20.20</text>
					<text x="0" y="42">22.60</text>
					<text x="0" y="22">25.00</text>
				</g>
				<g fill="black" id="label" text-anchor="start">
					<text x="0" y="237">1.00</text>
					<text x="43" y="237">1.40</text>
					<text x="86" y="237">1.80</text>
					<text x="129" y="237">2.20</text>
					<text x="172" y="237">2.60</text>
					<text x="215" y="237">3.00</text>
					<text x="258" y="237">3.40</text>
					<text x="301" y="237">3.80</text>
					<text x="344" y="237">4.20</text>
					<text x="387" y="237">4.60</text>
					<text x="430" y="237">5.00</text>
				</g>
				<svg height="227" viewBox="0 0 430 227" width="430">
					<g fill="none" transform="translate(0, 227) scale(1, -1)">
						<g transform="translate(2, 2) scale(0.9906976744186047, 0.9823788546255506)">
							<g id="grid" stroke="black" stroke-width="1">
								<defs>
									<line id="vLine" x1="0" x2="0" y1="0" y2="227" />
									<line id="hLine" x1="0" x2="430" y1="0" y2="0" />
								</defs>
								<use x="0" xlink:href="#vLine" />
								<use xlink:href="#hLine" y="0" />
								<use x="43" xlink:href="#vLine" />
								<use xlink:href="#hLine" y="23" />
								<use x="86" xlink:href="#vLine" />
								<use xlink:href="#hLine" y="45" />
								<use x="129" xlink:href="#vLine" />
								<use xlink:href="#hLine" y="68" />
								<use x="172" xlink:href="#vLine" />
								<use xlink:href="#hLine" y="91" />
								<use x="215" xlink:href="#vLine" />
								<use xlink:href="#hLine" y="114" />
								<use x="258" xlink:href="#vLine" />
								<use xlink:href="#hLine" y="136" />
								<use x="301" xlink:href="#vLine" />
								<use xlink:href="#hLine" y="159" />
								<use x="344" xlink:href="#vLine" />
								<use xlink:href="#hLine" y="182" />
								<use x="387" xlink:href="#vLine" />
								<use xlink:href="#hLine" y="204" />
								<use x="430" xlink:href="#vLine" />
								<use xlink:href="#hLine" y="227" />
							</g>
							<g id="data" transform="translate(-107.5, -9.458333333333334) scale(107.5, 9.458333333333334)">
								<defs>
									<marker fill="none" id="polyline-midmarker" orient="auto" preserveAspectRatio="xMidYMid meet" refX="5" refY="5" stroke="blue" vector-effect="non-scaling-stroke" viewBox="0 0 10 10">
										<circle cx="0" cy="0" r="1" />
									</marker>
								</defs>
								<polyline fill="none" points="1.000000,1.000000 2.000000,4.000000 3.000000,9.000000 4.000000,16.000000 5.000000,25.000000 " stroke="blueviolet" vector-effect="non-scaling-stroke" />
								<polyline points="1.000000,25.000000 2.000000,16.000000 3.000000,9.000000 4.000000,4.000000 5.000000,1.000000 " stroke="brown" vector-effect="non-scaling-stroke" />
							</g>
						</g>
						<rect height="227" stroke="grey" stroke-width="3" width="430" x="0" y="0" />
					</g>
				</svg>
			</g>
			<g id="legend" transform="translate(5, 30)" viewBox="0 0 20 215">
				<defs>
					<rect height="107" id="legendRect" width="20" x="0" y="0" />
				</defs>
				<use fill="blueviolet" xlink:href="#legendRect" y="0" />
				<text fill="black" text-anchor="middle" transform="rotate(90, 5, 0)" x="58" y="0">first</text>
				<use fill="brown" xlink:href="#legendRect" y="107" />
				<text fill="black" text-anchor="middle" transform="rotate(90, 5, 107)" x="58" y="107">second</text>
			</g>
		</g>
		<g id="lowerHalf" transform="translate(0, 250)">
			<defs>
				<rect fill="blue" height="10" id="blueRectangle" width="10" x="0" y="0" />
				<circle cx="0" cy="0" fill="yellow" id="yellowCircle" r="5" />
			</defs>
			<text text-anchor="middle" x="250" y="62">Simple SVG example with demonstration of subgroups</text>
			<use x="50" xlink:href="#blueRectangle" y="62" />
			<use x="450" xlink:href="#blueRectangle" y="62" />
			<use x="250" xlink:href="#yellowCircle" y="125" />
		</g>
	</svg>



### Design principle ###

Every SVG element which can enclose other SVG groups returns a reference to this element. To add an element to the specific group, use the 
return value of the SVG group. This builds a tree of svg elements which composes the XML structure. Upon writing, the Flush function writes the whole tree to the writer provided.

The SVG Library is also written with the use of CSS in mind. Therefore, the wrapper functions can easily have their style modified if an external CSS sheet is provided.