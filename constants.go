package smartSVG

import (
	"container/ring"
	"math/rand"
	"strconv"
)

const (
	// Matches REGEX-pattern ([A-Za-z]+) += ("[a-z]+") +// rgb\([0-9]+, [0-9]+, [0-9]+\)
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
	Limegreen            = "limegreen"            // rgb(50, 205, 50)
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

var c *ring.Ring

func init() {
	count := 146
	c = ring.New(count)

	// Set random start point
	defer func() {
		for i := 0; i < rand.Intn(count); i++ {
			c = c.Next()
		}
	}()
	c.Value = Aliceblue
	c = c.Next()
	c.Value = Antiquewhite
	c = c.Next()
	c.Value = Aqua
	c = c.Next()
	c.Value = Aquamarine
	c = c.Next()
	c.Value = Azure
	c = c.Next()
	c.Value = Beige
	c = c.Next()
	c.Value = Bisque
	c = c.Next()
	c.Value = Black
	c = c.Next()
	c.Value = Blanchedalmond
	c = c.Next()
	c.Value = Blue
	c = c.Next()
	c.Value = Blueviolet
	c = c.Next()
	c.Value = Brown
	c = c.Next()
	c.Value = Burlywood
	c = c.Next()
	c.Value = Cadetblue
	c = c.Next()
	c.Value = Chartreuse
	c = c.Next()
	c.Value = Chocolate
	c = c.Next()
	c.Value = Coral
	c = c.Next()
	c.Value = Cornflowerblue
	c = c.Next()
	c.Value = Cornsilk
	c = c.Next()
	c.Value = Crimson
	c = c.Next()
	c.Value = Cyan
	c = c.Next()
	c.Value = Darkblue
	c = c.Next()
	c.Value = Darkcyan
	c = c.Next()
	c.Value = Darkgoldenrod
	c = c.Next()
	c.Value = Darkgray
	c = c.Next()
	c.Value = Darkgreen
	c = c.Next()
	c.Value = Darkgrey
	c = c.Next()
	c.Value = Darkkhaki
	c = c.Next()
	c.Value = Darkmagenta
	c = c.Next()
	c.Value = Darkolivegreen
	c = c.Next()
	c.Value = Darkorange
	c = c.Next()
	c.Value = Darkorchid
	c = c.Next()
	c.Value = Darkred
	c = c.Next()
	c.Value = Darksalmon
	c = c.Next()
	c.Value = Darkseagreen
	c = c.Next()
	c.Value = Darkslateblue
	c = c.Next()
	c.Value = Darkslategray
	c = c.Next()
	c.Value = Darkslategrey
	c = c.Next()
	c.Value = Darkturquoise
	c = c.Next()
	c.Value = Darkviolet
	c = c.Next()
	c.Value = Deeppink
	c = c.Next()
	c.Value = Deepskyblue
	c = c.Next()
	c.Value = Dimgray
	c = c.Next()
	c.Value = Dodgerblue
	c = c.Next()
	c.Value = Firebrick
	c = c.Next()
	c.Value = Floralwhite
	c = c.Next()
	c.Value = Forestgreen
	c = c.Next()
	c.Value = Fuchsia
	c = c.Next()
	c.Value = Gainsboro
	c = c.Next()
	c.Value = Ghostwhite
	c = c.Next()
	c.Value = Gold
	c = c.Next()
	c.Value = Goldenrod
	c = c.Next()
	c.Value = Gray
	c = c.Next()
	c.Value = Grey
	c = c.Next()
	c.Value = Green
	c = c.Next()
	c.Value = Greenyellow
	c = c.Next()
	c.Value = Honeydew
	c = c.Next()
	c.Value = Hotpink
	c = c.Next()
	c.Value = Indianred
	c = c.Next()
	c.Value = Indigo
	c = c.Next()
	c.Value = Ivory
	c = c.Next()
	c.Value = Khaki
	c = c.Next()
	c.Value = Lavender
	c = c.Next()
	c.Value = Lavenderblush
	c = c.Next()
	c.Value = Lawngreen
	c = c.Next()
	c.Value = Lemonchiffon
	c = c.Next()
	c.Value = Lightblue
	c = c.Next()
	c.Value = Lightcoral
	c = c.Next()
	c.Value = Lightcyan
	c = c.Next()
	c.Value = Lightgoldenrodyellow
	c = c.Next()
	c.Value = Lightgray
	c = c.Next()
	c.Value = Lightgreen
	c = c.Next()
	c.Value = Lightgrey
	c = c.Next()
	c.Value = Lightpink
	c = c.Next()
	c.Value = Lightsalmon
	c = c.Next()
	c.Value = Lightseagreen
	c = c.Next()
	c.Value = Lightskyblue
	c = c.Next()
	c.Value = Lightslategray
	c = c.Next()
	c.Value = Lightslategrey
	c = c.Next()
	c.Value = Lightsteelblue
	c = c.Next()
	c.Value = Lightyellow
	c = c.Next()
	c.Value = Lime
	c = c.Next()
	c.Value = Limegreen
	c = c.Next()
	c.Value = Linen
	c = c.Next()
	c.Value = Magenta
	c = c.Next()
	c.Value = Maroon
	c = c.Next()
	c.Value = Mediumaquamarine
	c = c.Next()
	c.Value = Mediumblue
	c = c.Next()
	c.Value = Mediumorchid
	c = c.Next()
	c.Value = Mediumpurple
	c = c.Next()
	c.Value = Mediumseagreen
	c = c.Next()
	c.Value = Mediumslateblue
	c = c.Next()
	c.Value = Mediumspringgreen
	c = c.Next()
	c.Value = Mediumturquoise
	c = c.Next()
	c.Value = Mediumvioletred
	c = c.Next()
	c.Value = Midnightblue
	c = c.Next()
	c.Value = Mintcream
	c = c.Next()
	c.Value = Mistyrose
	c = c.Next()
	c.Value = Moccasin
	c = c.Next()
	c.Value = Navajowhite
	c = c.Next()
	c.Value = Navy
	c = c.Next()
	c.Value = Oldlace
	c = c.Next()
	c.Value = Olive
	c = c.Next()
	c.Value = Olivedrab
	c = c.Next()
	c.Value = Orange
	c = c.Next()
	c.Value = Orangered
	c = c.Next()
	c.Value = Orchid
	c = c.Next()
	c.Value = Palegoldenrod
	c = c.Next()
	c.Value = Palegreen
	c = c.Next()
	c.Value = Paleturquoise
	c = c.Next()
	c.Value = Palevioletred
	c = c.Next()
	c.Value = Papayawhip
	c = c.Next()
	c.Value = Peachpuff
	c = c.Next()
	c.Value = Peru
	c = c.Next()
	c.Value = Pink
	c = c.Next()
	c.Value = Plum
	c = c.Next()
	c.Value = Powderblue
	c = c.Next()
	c.Value = Purple
	c = c.Next()
	c.Value = Red
	c = c.Next()
	c.Value = Rosybrown
	c = c.Next()
	c.Value = Royalblue
	c = c.Next()
	c.Value = Saddlebrown
	c = c.Next()
	c.Value = Salmon
	c = c.Next()
	c.Value = Sandybrown
	c = c.Next()
	c.Value = Seagreen
	c = c.Next()
	c.Value = Seashell
	c = c.Next()
	c.Value = Sienna
	c = c.Next()
	c.Value = Silver
	c = c.Next()
	c.Value = Skyblue
	c = c.Next()
	c.Value = Slateblue
	c = c.Next()
	c.Value = Slategray
	c = c.Next()
	c.Value = Slategrey
	c = c.Next()
	c.Value = Snow
	c = c.Next()
	c.Value = Springgreen
	c = c.Next()
	c.Value = Steelblue
	c = c.Next()
	c.Value = Tan
	c = c.Next()
	c.Value = Teal
	c = c.Next()
	c.Value = Thistle
	c = c.Next()
	c.Value = Tomato
	c = c.Next()
	c.Value = Turquoise
	c = c.Next()
	c.Value = Violet
	c = c.Next()
	c.Value = Wheat
	c = c.Next()
	c.Value = White
	c = c.Next()
	c.Value = Whitesmoke
	c = c.Next()
	c.Value = Yellow
	c = c.Next()
	c.Value = Yellowgreen
	c = c.Next()
}

func GetColour() (ret string) {
	ret = c.Value.(string)
	c = c.Next()
	return
}

func GetRandomColour() string {
	return RGB(rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func RGB(r, g, b int) string {
	var ret string = "rgb("
	for _, v := range []int{r, g, b} {
		ret += strconv.Itoa(v) + ", "
	}
	ret = ret[:len(ret)-2] // Remove superfluous ",space"
	ret += ")"
	return ret
}
