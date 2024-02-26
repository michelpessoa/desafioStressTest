package stresstester

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Report struct {
	URL                   string
	ConfiguredRequests    int
	ConfiguredConcurrency int
	StartTime             time.Time
	TotalRequest          int
	ResponseStatuses      map[int]int
	mutex                 sync.Mutex
}

func RunTester(url string, requests int, concurrency int) {
	log.Printf("Starting Stress Tester... \n\nParams:\n  url\t\t=> %s \n  requests\t=> %d \n  concurrency\t=> %d\n\n", url, requests, concurrency)

	report := newReport(url, requests, concurrency)

	if concurrency > requests {
		concurrency = requests
	}

	total_loops := requests / concurrency
	remainder := requests % concurrency

	if remainder > 0 {
		total_loops += 1
	}

	loops := make([]int, total_loops)

	for i := 0; i < total_loops; i++ {
		if i == total_loops-1 && remainder > 0 {
			loops[i] = remainder
			continue
		}
		loops[i] = concurrency
	}

	log.Printf("Executing")
	wg := &sync.WaitGroup{}

	for _, loop := range loops {
		wg.Add(loop)
		log.Printf("Dispatching %d\n", loop)
		for j := 0; j < loop; j++ {
			go executeRequest(url, wg, report)
		}
		wg.Wait()
		fmt.Println()
	}
	fmt.Printf(" OK\n")
	log.Println("Done")
	report.printResults()
}

var httpClient = http.Client{}

func executeRequest(url string, wg *sync.WaitGroup, report *Report) {
	fmt.Printf(".")
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		wg.Done()
		return
	}

	request.Close = true
	response, _ := httpClient.Do(request)
	defer response.Body.Close()
	report.registerResponseStatus(response.StatusCode)
	wg.Done()
}

func newReport(url string, requests int, concurrency int) *Report {
	report := &Report{}
	report.URL = url
	report.ConfiguredRequests = requests
	report.ConfiguredConcurrency = concurrency
	report.ResponseStatuses = map[int]int{}
	report.StartTime = time.Now()
	return report
}

func (r *Report) registerResponseStatus(statusCode int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	currentCount, exists := r.ResponseStatuses[statusCode]

	if !exists {
		currentCount = 0
	}
	r.ResponseStatuses[statusCode] = currentCount + 1
}

func (r *Report) getDuration() string {
	d := time.Since(r.StartTime)

	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.2f segundos", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%d minutos, %.2f segundos", int(d.Seconds())/60, d.Seconds()-float64(int(d.Seconds())/60)*60)
	}

	return fmt.Sprintf("%d horas, %d minutos, %.2f segundos", int(d.Hours()), int(d.Minutes())%60, d.Seconds()-float64(int(d.Seconds())/60)*60)
}

func (r *Report) getTotalRequests() int {
	t := 0

	for _, v := range r.ResponseStatuses {
		t += v
	}

	return t
}

func (r *Report) printResults() {
	requests_completed := r.getTotalRequests()
	fmt.Println()

	fmt.Printf("\nTotal time: \t\t%s", r.getDuration())
	fmt.Printf("\nFailed requests: \t%d", r.ConfiguredRequests-requests_completed)
	fmt.Printf("\nCompleted Requests: \t%d", requests_completed)

	fmt.Printf("\n\nHTTP STATUS\tTOTAL")
	for key, value := range r.ResponseStatuses {
		fmt.Printf("\n%d\t\t%d", key, value)
	}

	fmt.Println()
}
