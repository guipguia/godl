package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/guipguia/godl/internal/util"
	"github.com/urfave/cli/v2"
)

type progressWriter struct {
	total        int64
	downloaded   int64
	lock         sync.Mutex
	startTime    time.Time
	currentSpeed float64
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

func updateProgress(progress *progressWriter, totalSize int64) {
	progress.startTime = time.Now()

	for {
		progress.lock.Lock()
		downloaded := progress.downloaded
		currentSpeed := progress.currentSpeed
		progress.lock.Unlock()

		percentage := float64(downloaded) / float64(totalSize) * 100.0
		fmt.Printf("\rProgress: %.2f%% | Speed: %.2f KB/s", percentage, currentSpeed)

		if downloaded >= totalSize {
			break
		}

		time.Sleep(500 * time.Millisecond) // Update interval
	}
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
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	// Create a progressWriter to track the download progress
	progress := &progressWriter{
		total:      fileSize,
		downloaded: 0,
		lock:       sync.Mutex{},
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

	fmt.Println("File downloaded successfully.")
}

func Download() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		var outputFileName string = util.GetNameBasedOnUrl(ctx.Args().Get(0))
		concurrency := ctx.Int("concurrency")
		if ctx.Args().Len() == 0 {
			fmt.Println("Please provide a URL to download something.")
			os.Exit(1)
		}

		if len(ctx.String("output")) != 0 && ctx.String("output") != util.GetNameBasedOnUrl(ctx.Args().Get(0)) {
			outputFileName = ctx.String("output")
		}

		doDownload(ctx.Args().Get(0), ctx.String("dir"), outputFileName, concurrency)
		return nil
	}
}
