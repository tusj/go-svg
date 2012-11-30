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

	$ godoc github.com/ajstarks/svgo
	

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
		drawing.Diagram(posX, posY, width, height, svg.Data{x, y}, "SVG Diagram", svg.Continuous)
		drawing.Flush(os.Stdout) 
	}

### This produces ###

	`<svg preserveAspectRatio="xMinYmin meet" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 500 500" xmlns:xlink="http://www.w3.org/1999/xlink">
		<g height="500" width="500" transform="translate(0, 0)" id="diagram">
			<rect y="0" x="0" height="500" width="500" />
			<text y="13" x="0" id="title" fill="black">SVG Diagram</text>
			<g transform="translate(70, 13)" id="plot">
				<g id="label" text-anchor="end">
					<text y="477" x="0">1.00</text>
					<text y="431" x="0">3.40</text>
					<text y="384" x="0">5.80</text>
					<text y="337" x="0">8.20</text>
					<text y="290" x="0">10.60</text>
					<text y="243" x="0">13.00</text>
					<text y="197" x="0">15.40</text>
					<text y="150" x="0">17.80</text>
					<text y="103" x="0">20.20</text>
					<text y="56" x="0">22.60</text>
					<text y="9" x="0">25.00</text>
				</g>
				<g id="label" text-anchor="start">
					<text y="487" x="0">1.00</text>
					<text y="487" x="43">1.40</text>
					<text y="487" x="86">1.80</text>
					<text y="487" x="129">2.20</text>
					<text y="487" x="172">2.60</text>
					<text y="487" x="215">3.00</text>
					<text y="487" x="258">3.40</text>
					<text y="487" x="301">3.80</text>
					<text y="487" x="344">4.20</text>
					<text y="487" x="387">4.60</text>
					<text y="487" x="430">5.00</text>
				</g>
				<svg height="477" width="430" viewBox="0 0 430 477">
					<g transform="translate(0, 477) scale(1, -1)" fill="none">
						<g transform="translate(2, 2) scale(0.9906976744186047, 0.9916142557651991)">
							<g stroke-width="1" stroke="black" id="grid">
								<defs>
									<line id="vLine" x2="0" x1="0" y1="0" y2="477" />
									<line id="hLine" x2="430" x1="0" y1="0" y2="0" />
								</defs>
								<use xlink:href="#vLine" x="0" />
								<use xlink:href="#hLine" y="0" />
								<use x="43" xlink:href="#vLine" />
								<use xlink:href="#hLine" y="48" />
								<use xlink:href="#vLine" x="86" />
								<use xlink:href="#hLine" y="95" />
								<use xlink:href="#vLine" x="129" />
								<use xlink:href="#hLine" y="143" />
								<use x="172" xlink:href="#vLine" />
								<use xlink:href="#hLine" y="191" />
								<use xlink:href="#vLine" x="215" />
								<use xlink:href="#hLine" y="238" />
								<use xlink:href="#vLine" x="258" />
								<use y="286" xlink:href="#hLine" />
								<use x="301" xlink:href="#vLine" />
								<use y="334" xlink:href="#hLine" />
								<use xlink:href="#vLine" x="344" />
								<use y="382" xlink:href="#hLine" />
								<use xlink:href="#vLine" x="387" />
								<use xlink:href="#hLine" y="429" />
								<use xlink:href="#vLine" x="430" />
								<use y="477" xlink:href="#hLine" />
							</g>
							<g transform="translate(-107.5, -19.875) scale(107.5, 19.875)" id="data">
								<defs>
									<marker preserveAspectRatio="xMidYMid meet" stroke="blue" id="polyline-midmarker" fill="none" orient="auto" refX="5" refY="5" vector-effect="non-scaling-stroke" viewBox="0 0 10 10">
										<circle cx="0" cy="0" r="1" />
									</marker>
								</defs>
								<polyline vector-effect="non-scaling-stroke" stroke="blueviolet" fill="none" points="1.000000,1.000000 2.000000,4.000000 3.000000,9.000000 4.000000,16.000000 5.000000,25.000000 " />
							</g>
						</g>
						<rect y="0" x="0" stroke-width="3" height="477" stroke="grey" width="430" />
					</g>
				</svg>
			</g>
		</g>
	</svg>`


### Design principle ###

Every SVG element which can enclose other SVG groups returns a reference to this element. To add an element to the specific group, use the 
return value of the SVG group. This builds a tree of svg elements which composes the XML structure. Upon writing, the Flush function writes the whole tree to the writer provided.