package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type Downloader struct {
	Url         string
	Dir         string
	FileName    string
	Concurrency int
}

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

func (d *Downloader) doDownload() {
	res, err := http.Get(d.Url)
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

	file, err := os.Create(path.Join(d.Dir, d.FileName))
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

	chunk := fileSize / int64(d.Concurrency)

	var wg sync.WaitGroup
	wg.Add(d.Concurrency)

	// Start the goroutines for concurrent downloading
	for i := 0; i < d.Concurrency; i++ {
		start := int64(i) * chunk
		end := start + chunk - 1

		// For the last goroutine, download the remaining bytes
		if i == d.Concurrency-1 {
			end = fileSize - 1
		}

		go downloadRange(client, d.Url, file, start, end, progress, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Printf("\n%s downloaded successfully.\n", d.FileName)
}

func StartDownload(url string, fullPath string, concurrency int) {
	d := Downloader{
		Url:         url,
		Dir:         path.Dir(fullPath),
		FileName:    path.Base(fullPath),
		Concurrency: concurrency,
	}

	d.doDownload()
}
