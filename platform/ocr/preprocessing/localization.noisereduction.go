package preprocessing

import (
	"sort"

	"gocv.io/x/gocv"
)

func (n *Storage) HybridDenoise(src gocv.Mat, maxWindowSize int) (gocv.Mat, error) {
	if maxWindowSize < 3 {
		maxWindowSize = 3
	}
	if maxWindowSize%2 == 0 {
		maxWindowSize++
	}

	stage1 := n.amfStage(src, maxWindowSize)
	defer stage1.Close()

	stage2 := n.mdbmfStage(stage1, maxWindowSize)

	return stage2, nil
}

func (n *Storage) amfStage(src gocv.Mat, maxWindowSize int) gocv.Mat {
	channels := gocv.Split(src)
	defer func() {
		for _, ch := range channels {
			ch.Close()
		}
	}()

	res := make([]gocv.Mat, len(channels))
	for i, ch := range channels {
		dst := gocv.NewMat()
		ch.CopyTo(&dst)
		res[i] = dst
	}

	for ci, ch := range channels {
		rows, cols := ch.Rows(), ch.Cols()

		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				cur := ch.GetUCharAt(y, x)
				if cur != 0 && cur != 255 {
					continue
				}

				win := 3
				var med byte
				ok := false

				for win <= maxWindowSize {
					half := win / 2
					vals := make([]byte, 0, win*win)

					for dy := -half; dy <= half; dy++ {
						for dx := -half; dx <= half; dx++ {
							ny, nx := y+dy, x+dx
							if ny >= 0 && ny < rows && nx >= 0 && nx < cols {
								vals = append(vals, ch.GetUCharAt(ny, nx))
							}
						}
					}

					sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
					med = vals[len(vals)/2]

					if med != 0 && med != 255 {
						ok = true
						break
					}
					win += 2
				}

				if ok {
					res[ci].SetUCharAt(y, x, med)
				}
			}
		}
	}

	var dst gocv.Mat
	gocv.Merge(res, &dst)

	for _, r := range res {
		r.Close()
	}

	return dst
}

func (n *Storage) mdbmfStage(src gocv.Mat, maxWindowSize int) gocv.Mat {
	channels := gocv.Split(src)
	defer func() {
		for _, ch := range channels {
			ch.Close()
		}
	}()

	res := make([]gocv.Mat, len(channels))
	for i, ch := range channels {
		dst := gocv.NewMat()
		ch.CopyTo(&dst)
		res[i] = dst
	}

	for ci, ch := range channels {
		rows, cols := ch.Rows(), ch.Cols()
		cy, cx := rows/2, cols/2

		var lastProc byte

		for _, p := range n.centerOut(rows, cols, cy, cx) {
			y, x := p.y, p.x
			cur := ch.GetUCharAt(y, x)
			if cur != 0 && cur != 255 {
				lastProc = cur
				continue
			}

			vals := n.collectValid(ch, rows, cols, y, x, maxWindowSize)
			if len(vals) > 0 {
				sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
				med := vals[len(vals)/2]
				res[ci].SetUCharAt(y, x, med)
				lastProc = med
				continue
			}

			half := maxWindowSize / 2
			if lastProc == 0 || lastProc == 255 {
				var zeroCount, ffCount int
				for dy := -half; dy <= half; dy++ {
					for dx := -half; dx <= half; dx++ {
						ny, nx := y+dy, x+dx
						if ny >= 0 && ny < rows && nx >= 0 && nx < cols {
							v := ch.GetUCharAt(ny, nx)
							if v == 0 {
								zeroCount++
							} else if v == 255 {
								ffCount++
							}
						}
					}
				}
				total := (2*half + 1) * (2*half + 1)
				if zeroCount > total/2 {
					lastProc = 0
				} else if ffCount > total/2 {
					lastProc = 255
				} else {
					lastProc = cur
				}
			}
			res[ci].SetUCharAt(y, x, lastProc)
		}
	}

	var dst gocv.Mat
	gocv.Merge(res, &dst)

	for _, r := range res {
		r.Close()
	}

	return dst
}

type coord struct {
	y, x int
}

func (n *Storage) centerOut(rows, cols, cy, cx int) []coord {
	total := rows * cols
	out := make([]coord, 0, total)

	if cy >= 0 && cy < rows && cx >= 0 && cx < cols {
		out = append(out, coord{cy, cx})
	}

	quadrants := []struct {
		ystart, yend, ystep int
		xstart, xend, xstep int
	}{
		{cy - 1, -1, -1, cx - 1, -1, -1},
		{cy - 1, -1, -1, cx + 1, cols, 1},
		{cy + 1, rows, 1, cx - 1, -1, -1},
		{cy + 1, rows, 1, cx + 1, cols, 1},
	}

	for _, q := range quadrants {
		for y := q.ystart; y != q.yend; y += q.ystep {
			for x := q.xstart; x != q.xend; x += q.xstep {
				out = append(out, coord{y, x})
			}
		}
	}

	return out
}

func (n *Storage) collectValid(ch gocv.Mat, rows, cols, y, x, maxWin int) []byte {
	for win := 3; win <= maxWin; win += 2 {
		half := win / 2
		var vals []byte

		for dy := -half; dy <= half; dy++ {
			for dx := -half; dx <= half; dx++ {
				ny, nx := y+dy, x+dx
				if ny >= 0 && ny < rows && nx >= 0 && nx < cols {
					v := ch.GetUCharAt(ny, nx)
					if v != 0 && v != 255 {
						vals = append(vals, v)
					}
				}
			}
		}

		if len(vals) > 0 {
			return vals
		}
	}

	return nil
}
