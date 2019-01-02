// Copyright 2018 by xruida.com, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package xlsx

func numTostring(t int) string {

	switch t {
	case 1:
		return "A"
	case 2:
		return "B"
	case 3:
		return "C"
	case 4:
		return "D"
	case 5:
		return "E"
	case 6:
		return "F"
	case 7:
		return "G"
	case 8:
		return "H"
	case 9:
		return "I"
	case 10:
		return "J"
	case 11:
		return "K"
	case 12:
		return "L"
	case 13:
		return "M"
	case 14:
		return "N"
	case 15:
		return "O"
	case 16:
		return "P"
	case 17:
		return "Q"
	case 18:
		return "R"
	case 19:
		return "S"
	case 20:
		return "T"
	case 21:
		return "U"
	case 22:
		return "V"
	case 23:
		return "W"
	case 24:
		return "X"
	case 25:
		return "Y"
	case 26:
		return "Z"
	}
	return "出错了"
}

func transformation(t int) string {
	if t <= 26 {
		return numTostring(t)
	} else if t > 26 && t <= 52 {
		return "A" + numTostring((t - 26))
	} else if t > 52 && t <= 78 {
		return "B" + numTostring((t - 52))
	} else if t > 78 && t <= 104 {
		return "C" + numTostring((t - 78))
	}
	return "出错了"
}
