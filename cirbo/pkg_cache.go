package cirbo

import (
	"github.com/cirbo-lang/cirbo/cbty"
	"github.com/cirbo-lang/cirbo/projpath"
)

type pkgCache map[projpath.FilePath]pkgCacheEntry

type pkgCacheEntry struct {
	Value cbty.Value
}

var noPkgCacheEntry = pkgCacheEntry{}

func (c pkgCache) Has(fp projpath.FilePath) bool {
	_, has := c[fp]
	return has
}

func (c pkgCache) Get(fp projpath.FilePath) pkgCacheEntry {
	return c[fp]
}

func (c pkgCache) GetOk(fp projpath.FilePath) (pkgCacheEntry, bool) {
	v, ok := c[fp]
	return v, ok
}

func (c pkgCache) Put(fp projpath.FilePath, entry pkgCacheEntry) {
	c[fp] = entry
}
