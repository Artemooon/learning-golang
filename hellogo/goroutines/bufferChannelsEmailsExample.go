package main

import (
    "fmt"
    "sync"
    "time"
)

type EmailJob struct {
    CustomerEmail string
    OrderId int
}

func sendEmail(job EmailJob) {
    time.Sleep(2 * time.Second)
    fmt.Println("Sent order email:", job.OrderID, "to", job.CustomerEmail)
}

func sendEmail(job EmailJob) {
    // Real code would call an email provider here.
    time.Sleep(2 * time.Second)
    fmt.Println("Sent order email:", job.OrderID, "to", job.CustomerEmail)
}

func emailWorker(jobs <-chan EmailJob, wg *sync.WaitGroup) {
    defer wg.Done()

    for job := range jobs {
            sendEmail(job)
    }
}

func main() {
    // Queue up to 100 waiting email jobs.
    jobs := make(chan EmailJob, 100)

    var workers sync.WaitGroup

    // Only three email requests run concurrently.
    for i := 0; i < 3; i++ {
            workers.Add(1)
            go emailWorker(jobs, &workers)
    }

    // Imagine this loop is many incoming HTTP order requests.
    for orderID := 1; orderID <= 10; orderID++ {
            jobs <- EmailJob{
                    CustomerEmail: "customer@example.com",
                    OrderID:       orderID,
            }

            fmt.Println("Order accepted:", orderID)
    }

    close(jobs)
    workers.Wait()
}