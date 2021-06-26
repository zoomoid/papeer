package types

type GridOptions struct {
	Enabled bool
	Opacity float32
	Color   string
	Scale   float32
}

type DeltaOptions struct {
	X string
	Y string
}

type WrapperOptions struct {
	Format    string
	Landscape bool
	Scale     float32
	Pages     string
	Filename  string
	Delta     DeltaOptions
	Grid      GridOptions
}
