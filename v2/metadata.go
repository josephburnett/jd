package jd

type Metadata struct {
	Version int
	Merge   bool
}

func (m Metadata) Options() []Option {
	if m.Merge {
		return []Option{MERGE}
	} else {
		return []Option{}
	}
}
