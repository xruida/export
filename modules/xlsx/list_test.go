// Copyright 2018 by xruida.com, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package xlsx

import (
	"testing"

	"github.com/issue9/assert"
)

func TestNumTostring(t *testing.T) {
	a := assert.New(t)

	a.Equal("A", numTostring(1))
	a.Equal("Z", numTostring(26))
}
