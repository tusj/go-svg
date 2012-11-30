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

	package main
	
	import (
		svg "github.com/tusj/smartSVG"
		os
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

### Design principle ###

Every SVG element which can enclose other SVG groups returns a reference to this element. To add an element to the specific group, use the 
return value of the SVG group. This builds a tree of svg elements which composes the XML structure. Upon writing, the Flush function writes the whole tree to the writer provided.