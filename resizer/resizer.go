package resizer

import (
	"context"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nfnt/resize"

	"image/jpeg"
	// Init jpeg processor
	_ "image/jpeg"
)

// Resizer ...
type Resizer struct {
	threads    int
	scale      float64
	inputPath  string
	outputPath string
}

/*
New returns new instance of Resizer
*/
func New(threads int, scale float64, inputPath, outputPath string) *Resizer {
	input, err := os.Stat(inputPath)
	if err != nil || input.Mode().IsRegular() {
		exitProgram("Invalid input folder")
	}

	output, err := os.Stat(outputPath)
	if err != nil || output.Mode().IsRegular() {
		exitProgram("Invalid output folder")
	}

	return &Resizer{threads, scale, inputPath, outputPath}
}

func exitProgram(message string) {
	fmt.Println(message)
	os.Exit(0)
}

/*
StartProcessing starts processing images,
returns
	WaitGroup to wait for all threads to end the job
	Error when any error occurs
*/
func (r *Resizer) StartProcessing(ctx context.Context) (*sync.WaitGroup, error) {

	files, err := r.getFiles(ctx)
	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}
	r.processFiles(ctx, wg, files)

	return wg, nil
}

/*
StartProcessingWithCancel starts processing images
returns
	WaitGroup to wait for all threads to end the job
	CancelFunc to cancel processing
	Error when any error occurs
*/
func (r *Resizer) StartProcessingWithCancel() (*sync.WaitGroup, context.CancelFunc, error) {

	ctx, cancelFunc := context.WithCancel(context.Background())

	wg, err := r.StartProcessing(ctx)
	if err != nil {
		cancelFunc()
		return nil, nil, err
	}

	return wg, cancelFunc, nil
}

func (r *Resizer) getFiles(ctx context.Context) (<-chan string, error) {
	files, err := ioutil.ReadDir(r.inputPath)
	if err != nil {
		return nil, err
	}

	ch := make(chan string)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, f := range files {
			if f.Mode().IsRegular() && strings.Contains(f.Name(), "jpg") || strings.Contains(f.Name(), "jpeg") {
				fullPath := filepath.Join(r.inputPath, f.Name())
				select {
				case <-ctx.Done():
					return
				case ch <- fullPath:
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch, nil
}

func (r *Resizer) processFiles(ctx context.Context, wg *sync.WaitGroup, files <-chan string) {

	for i := 0; i < r.threads; i++ {
		wg.Add(1)
		go r.resizeFiles(ctx, wg, files)
	}

}

func (r *Resizer) resizeFiles(ctx context.Context, wg *sync.WaitGroup, files <-chan string) {

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case f, ok := <-files:
			if !ok {
				return
			}
			resizeFile(f, r.outputPath, r.scale)
		}
	}
}

func resizeFile(filepath string, outputPath string, scale float64) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file", err.Error())
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decodng image", err.Error())
		return
	}

	width := scale / 100 * float64(img.Bounds().Dx())
	height := scale / 100 * float64(img.Bounds().Dy())

	img = resize.Resize(uint(width), uint(height), img, resize.Bicubic)

	filename, _ := file.Stat()
	placeToSave, err := os.Create(path.Join(outputPath, filename.Name()))
	if err != nil {
		fmt.Println("Failed to create destination file!")
		return
	}

	jpeg.Encode(
		placeToSave,
		img,
		&jpeg.Options{
			Quality: jpeg.DefaultQuality,
		},
	)
}
