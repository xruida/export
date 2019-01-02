// Copyright 2018 by xruida.com, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package xlsx

func numTostring(t int) string {
	if t < 1 || t > 26 {
		return "出错了"
	}
	return string(t - 1 + 'A')
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
