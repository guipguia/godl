package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/guipguia/godl/internal/util"
	"github.com/urfave/cli/v2"
)

type progressWriter struct {
	total            int64
	downloaded       int64
	lock             sync.Mutex
	startTime        time.Time
	endTime          time.Time
	currentSpeed     float64
	currentSpeedUnit string
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.lock.Lock()
	pw.downloaded += int64(n)
	pw.calculateSpeed()
	pw.lock.Unlock()
	return n, nil
}

func (pw *progressWriter) calculateSpeed() {
	elapsed := time.Since(pw.startTime).Seconds()
	downloadedBytes := pw.downloaded
	pw.currentSpeed = float64(downloadedBytes) / elapsed / 1024 // Speed in KB/s
}

func downloadRange(client *http.Client, url string, file *os.File, start, end int64, pw *progressWriter, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %s\n", err)
		return
	}

	rangeHeader := fmt.Sprintf("bytes=%d-%d", start, end)
	req.Header.Set("Range", rangeHeader)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error downloading chunk: %s\n", err)
		return
	}
	defer resp.Body.Close()

	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		fmt.Printf("Error seeking file: %s\n", err)
		return
	}

	_, err = io.Copy(file, io.TeeReader(resp.Body, pw))
	if err != nil {
		fmt.Printf("Error writing chunk to file: %s\n", err)
		return
	}
}

// func updateProgress(progress *progressWriter, totalSize int64) {
// 	progress.startTime = time.Now()

// 	for {
// 		progress.lock.Lock()
// 		downloaded := progress.downloaded
// 		currentSpeed := progress.currentSpeed
// 		progress.lock.Unlock()

// 		percentage := float64(downloaded) / float64(totalSize) * 100.0
// 		fmt.Printf("\rProgress: %.2f%% | Speed: %.2f KB/s", percentage, currentSpeed)

// 		if downloaded >= totalSize {
// 			break
// 		}

// 		time.Sleep(500 * time.Millisecond) // Update interval
// 	}
// }

func updateProgress(progress *progressWriter, totalSize int64) {
	progress.startTime = time.Now()
	prevProgressLen := 0

	for {
		progress.lock.Lock()
		downloaded := progress.downloaded
		currentSpeed := progress.currentSpeed
		currentSpeedUnit := progress.currentSpeedUnit
		progress.lock.Unlock()

		percentage := float64(downloaded) / float64(totalSize) * 100.0

		if currentSpeed >= 1024 {
			currentSpeedUnit = "MB/s"
		}

		if currentSpeedUnit == "MB/s" {
			currentSpeed /= 1024 // Convert average speed to MB/s if necessary
		}

		remainingBytes := totalSize - downloaded
		remainingTime := time.Duration(float64(remainingBytes)/(progress.currentSpeed*1024)) * time.Second

		progressMsg := fmt.Sprintf("\rProgress: %.2f%% | Speed: %.2f %s | Remaining Time: %s",
			percentage, currentSpeed, currentSpeedUnit, remainingTime.Round(time.Second))

		// Clear the previous progress output
		if len(progressMsg) < prevProgressLen {
			clearProgress := strings.Repeat(" ", prevProgressLen)
			fmt.Print("\r" + clearProgress + "\r")
		}

		// Print the current progress
		fmt.Print(progressMsg)
		prevProgressLen = len(progressMsg)

		if downloaded >= totalSize {
			break
		}

		time.Sleep(200 * time.Millisecond) // Update interval
	}

	progress.endTime = time.Now()
}

// doDowwnload start downloading the file and save it to specified location
func doDownload(downloadUrl string, directory string, filename string, concurrency int) {
	res, err := http.Get(downloadUrl)
	if err != nil {
		fmt.Printf("Faield to download file: %s\n", err)
		return
	}
	defer res.Body.Close()

	fileSize := res.ContentLength
	if fileSize <= 0 {
		fmt.Println("Invalid file size")
		return
	}

	file, err := os.Create(directory + "/" + filename)
	if err != nil {
		fmt.Printf("Failed to create file: %s\n", err)
		return
	}
	defer file.Close()

	// Create a HTTP client with timeout
	client := &http.Client{}

	// Create a progressWriter to track the download progress
	progress := &progressWriter{
		total:            fileSize,
		downloaded:       0,
		lock:             sync.Mutex{},
		currentSpeedUnit: "KB/s",
	}

	go updateProgress(progress, fileSize)

	chunk := fileSize / int64(concurrency)

	var wg sync.WaitGroup
	wg.Add(concurrency)

	// Start the goroutines for concurrent downloading
	for i := 0; i < concurrency; i++ {
		start := int64(i) * chunk
		end := start + chunk - 1

		// For the last goroutine, download the remaining bytes
		if i == concurrency-1 {
			end = fileSize - 1
		}

		go downloadRange(client, downloadUrl, file, start, end, progress, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("\nFile downloaded successfully.")
}

// Download is the command function where the downloading begin.
func Download() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		var outputFileName string = util.GetNameBasedOnUrl(ctx.Args().Get(0))
		concurrency := ctx.Int("concurrency")

		if ctx.Args().Len() == 0 {
			fmt.Println("Please provide a URL to download something.")
			os.Exit(1)
		}

		if len(ctx.String("filename")) != 0 && ctx.String("filename") != util.GetNameBasedOnUrl(ctx.Args().Get(0)) {
			outputFileName = ctx.String("filename")
		}

		doDownload(ctx.Args().Get(0), ctx.String("dir"), outputFileName, concurrency)
		return nil
	}
}
