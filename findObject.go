package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
)

// this also a lil comment. super cute.333
// this comment is also cute
// new line

var (
	rcolor = color.RGBA{R: 255, A: 255}
	lcolor = color.RGBA{R: 255, A: 255}

	// current hue values set to find red colors
	lhsv = gocv.Scalar{Val1: 10, Val2: 145, Val3: 135}
	hhsv = gocv.Scalar{Val1: 255, Val2: 245, Val3: 240}

	size = image.Point{X: 600, Y: 600}
	blur = image.Point{X: 11, Y: 11}

	wt     = gocv.NewWindow("thersholded")
	wi     = gocv.NewWindow("images")
	img    = gocv.NewMat()
	mask   = gocv.NewMat()
	frame  = gocv.NewMat()
	hsv    = gocv.NewMat()
	kernel = gocv.NewMat()
)

func main() {
	defer close()

	// windows for user
	wt.ResizeWindow(600, 600)
	wt.MoveWindow(0, 0)
	wi.MoveWindow(600, 0)
	wi.ResizeWindow(600, 600)

	// laptop camera - maybe connects to camera driver?
	video, _ := gocv.OpenVideoCapture(0)
	defer video.Close()

	for {

		// break if camera not connected
		if !video.Read(&img) {
			fmt.Println("yoyo")
			break
		}

		// filters
		gocv.Flip(img, &img, 1)
		gocv.Resize(img, &img, size, 0, 0, gocv.InterpolationLinear)
		gocv.GaussianBlur(img, &frame, blur, 0, 0, gocv.BorderReflect101) // blurr img
		gocv.CvtColor(frame, &hsv, gocv.ColorBGRToHSV)                    // convert blurr to HSV color

		// look for object
		gocv.InRangeWithScalar(hsv, lhsv, hhsv, &mask) // find pixels based on hhsv and lhsv
		gocv.Erode(mask, &mask, kernel)
		gocv.Dilate(mask, &mask, kernel)

		contour := bestContour(mask, 2000)
		if len(contour) == 0 { // if no contour present; loop will continue
			if imShow() {
				break
			}
			continue
		}
		rect := gocv.BoundingRect(contour)
		gocv.Rectangle(&img, rect, color.RGBA{0, 255, 0, 0}, 2)
		fmt.Println(rect.Max)
		fmt.Println(rect.Dx())

		if imShow() {
			break
		}

	} // end of for-ever

}

func getRectDim(rect image.Rectangle) (int, int) {
	return 0, 0
}

func bestContour(frame gocv.Mat, minArea float64) []image.Point {
	cntr := gocv.FindContours(frame, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	var (
		bestCntr []image.Point
		bestArea = minArea
	)
	for _, cnt := range cntr {
		if area := gocv.ContourArea(cnt); area > bestArea {
			bestArea = area
			bestCntr = cnt
		}
	}
	return bestCntr
}

// controls delay
func imShow() bool {
	wi.IMShow(img)
	wt.IMShow(mask)
	return wi.WaitKey(1) == 27 || wt.WaitKey(1) == 27
}

func close() {
	defer wi.Close()
	defer wt.Close()
	defer img.Close()
	defer mask.Close()
	defer frame.Close()
	defer hsv.Close()
	defer kernel.Close()
}
