package preprocessing

import (
	"image"
	"math"
	"sort"

	"aidoc/internal/entity"

	"gocv.io/x/gocv"
)

func (e *Storage) DewarpDocument(img gocv.Mat, points []entity.Point) gocv.Mat {
	if len(points) != 4 {
		return img
	}

	ordered := e.orderPoints(points)

	src := gocv.NewPointVector()
	defer src.Close()
	for _, p := range ordered {
		src.Append(image.Point{X: int(p.X), Y: int(p.Y)})
	}

	w := int(math.Round(distance(ordered[0], ordered[1])))
	h := int(math.Round(distance(ordered[0], ordered[3])))

	if w <= 0 || h <= 0 {
		return img
	}

	dst := gocv.NewPointVector()
	defer dst.Close()
	dst.Append(image.Point{X: 0, Y: 0})
	dst.Append(image.Point{X: w, Y: 0})
	dst.Append(image.Point{X: w, Y: h})
	dst.Append(image.Point{X: 0, Y: h})

	matrix := gocv.GetPerspectiveTransform(src, dst)
	defer matrix.Close()

	result := gocv.NewMat()
	gocv.WarpPerspective(img, &result, matrix, image.Point{X: w, Y: h})
	return result
}

func (e *Storage) orderPoints(points []entity.Point) []entity.Point {
	cp := make([]entity.Point, len(points))
	copy(cp, points)

	sort.Slice(cp,
		func(i, j int) bool {
			return cp[i].Y < cp[j].Y
		},
	)

	top := cp[:2]
	bottom := cp[2:]

	sort.Slice(top,
		func(i, j int) bool {
			return top[i].X < top[j].X
		},
	)

	sort.Slice(bottom,
		func(i, j int) bool {
			return bottom[i].X < bottom[j].X
		},
	)

	return []entity.Point{top[0], top[1], bottom[1], bottom[0]}
}

func distance(a, b entity.Point) float64 {
	dx := float64(a.X - b.X)
	dy := float64(a.Y - b.Y)
	return math.Sqrt(dx*dx + dy*dy)
}
