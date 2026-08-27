package entity

type Point struct {
	X float32
	Y float32
}

type FilterType string

const (
	FilterSharpen FilterType = "sharpen"
)

type FilterConfig struct {
	Type       FilterType
	KernelSize int
}

type DenoiseConfig struct {
	KernelSize int
}


type PipelineStep struct {
	Denoise *DenoiseConfig
	Points  []Point
	MaxDim  int
	Target  PipelineImageSource
}

type SegType string

const (
	SegThreshold      SegType = "threshold"
	SegEdge           SegType = "edge"
	SegKMeans         SegType = "kmeans"
	SegWatershed      SegType = "watershed"
	SegGrabCut        SegType = "grabcut"
	SegConnectedComponents SegType = "connected_components"
)

type SegConfig struct {
	Type        SegType
	K           int
	Threshold   float64
	MaxVal      float64
	CannyLow    float64
	CannyHigh   float64
	Iter        int
	RectX, RectY, RectW, RectH int
}

type PipelineImageSource struct {
	Bytes []byte
	Path  string
}

type PipelineConfig struct {
	Steps []PipelineStep
}