package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"log"
)


func main()  {

	var files [3] string
	files[0]="/home/asqad/Pictures/girl.jpg"
	files[1]="/home/asqad/Pictures/carton.jpg"
	files[2]="/home/asqad/Pictures/island.jpg"

	fch :=make(chan string,3)

	// go func ()  {
	// 	for i := 0; i < 3; i++ {
	// 		fch<-files[i]
	// 	}
	// 	close(fch)
	// }()

	for i := 0; i < 3; i++ {
		fch<-files[i]
	}
	close(fch)
	
	total := makeThumbnails(fch)

	fmt.Println(total)

}

func makeThumbnails(filenames <-chan string)int64  {
	sizes := make(chan int64)
	var wg sync.WaitGroup

	for f := range filenames {
		wg.Add(1)

		//worker
		go func(f string) {
			defer wg.Done()
			thumb,err := ImageFile(f)
			if err!=nil {
				log.Println(err)
				return
			}
			info,_ :=os.Stat(thumb) //ok to ignore this error

			sizes <- info.Size()
		}(f)
	}

	// go func() {
	// 	wg.Wait()
	// 	close(sizes)
	// }()

	wg.Wait()
	close(sizes)
	

	var total int64
	for size := range sizes	 {
		total+=size
	}
	return total
}



// Image returns a thumbnail-size version of src.
func Image(src image.Image) image.Image {
	// Compute thumbnail size, preserving aspect ratio.
	xs := src.Bounds().Size().X
	ys := src.Bounds().Size().Y
	width, height := 128, 128
	if aspect := float64(xs) / float64(ys); aspect < 1.0 {
		width = int(128 * aspect) // portrait
	} else {
		height = int(128 / aspect) // landscape
	}
	xscale := float64(xs) / float64(width)
	yscale := float64(ys) / float64(height)

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	// a very crude scaling algorithm
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			srcx := int(float64(x) * xscale)
			srcy := int(float64(y) * yscale)
			dst.Set(x, y, src.At(srcx, srcy))
		}
	}
	return dst
}

// ImageStream reads an image from r and
// writes a thumbnail-size version of it to w.
func ImageStream(w io.Writer, r io.Reader) error {
	src, _, err := image.Decode(r)
	if err != nil {
		return err
	}
	dst := Image(src)
	return jpeg.Encode(w, dst, nil)
}

// ImageFile2 reads an image from infile and writes
// a thumbnail-size version of it to outfile.
func ImageFile2(outfile, infile string) (err error) {
	in, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outfile)
	if err != nil {
		return err
	}

	if err := ImageStream(out, in); err != nil {
		out.Close()
		return fmt.Errorf("scaling %s to %s: %s", infile, outfile, err)
	}
	return out.Close()
}

// ImageFile reads an image from infile and writes
// a thumbnail-size version of it in the same directory.
// It returns the generated file name, e.g. "foo.thumb.jpeg".
func ImageFile(infile string) (string, error) {
	ext := filepath.Ext(infile) // e.g., ".jpg", ".JPEG"
	outfile := strings.TrimSuffix(infile, ext) + ".thumb" + ext
	return outfile, ImageFile2(outfile, infile)
}