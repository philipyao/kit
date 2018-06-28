package util

import (
    "testing"
)

func TestCompress(t *testing.T) {
    src := []byte("dfjslfjdslfjd;af1090233sdfjdlsfds002r32r2r-2sfdslfj1jds*o1022&^%@22")
    t.Logf("src: len %v, %v", len(src), src)
    dst := Compress(src)
    t.Logf("after compress, dst: len %v, %v", len(dst), dst)

    tmp := Decompress(dst)
    t.Logf("after decompress, tmp: len %v, %v", len(tmp), tmp)
}
