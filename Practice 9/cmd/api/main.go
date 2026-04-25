package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"practice-9/internal/idempotency"
	"practice-9/internal/payment"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Request body for payment
type PaymentRequest struct {
	Amount float64 `json:"amount"`
	UserID string  `json:"user_id"`
}

// SimulatedPaymentServer handles payment requests with simulated failures
type SimulatedPaymentServer struct {
	failureCount int
	idempotency  idempotency.Store
}

func NewSimulatedPaymentServer() *SimulatedPaymentServer {
	return &SimulatedPaymentServer{
		failureCount: 3,
		idempotency:  idempotency.NewMemoryStore(),
	}
}

// SimultaneousPaymentRequests simulates a "double-click attack" with multiple goroutines
func SimultaneousPaymentRequests(numRequests int, idempotencyKey string) {
	fmt.Printf("[%s] [TEST] Sending %d simultaneous requests with same Idempotency-Key: '%s'\n",
		time.Now().Format("15:04:05.000"), numRequests, idempotencyKey)
	fmt.Printf("[%s] [TEST] Waiting 1 second before sending requests...\n", time.Now().Format("15:04:05.000"))

	time.Sleep(1 * time.Second)
	fmt.Printf("[%s] [TEST] Sending all requests now!\n", time.Now().Format("15:04:05.000"))

	var wg sync.WaitGroup
	client := &http.Client{Timeout: 15 * time.Second}

	for i := 1; i <= numRequests; i++ {
		wg.Add(1)
		go func(requestNum int) {
			defer wg.Done()

			req, _ := http.NewRequest("POST", "http://localhost:8080/execute-payment-idempotent", nil)
			req.Header.Set("Idempotency-Key", idempotencyKey)
			req.Header.Set("Content-Type", "application/json")

			fmt.Printf("[%s] [GOROUTINE #%d] Sending request...\n", time.Now().Format("15:04:05.000"), requestNum)

			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("[%s] [GOROUTINE #%d] ERROR: %v\n", time.Now().Format("15:04:05.000"), requestNum, err)
				return
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("[%s] [GOROUTINE #%d] Response - Status: %d, Body: %s\n",
				time.Now().Format("15:04:05.000"), requestNum, resp.StatusCode, string(body))
		}(i)
	}

	wg.Wait()
	fmt.Printf("[%s] [TEST] All requests completed\n", time.Now().Format("15:04:05.000"))
}

// HandlePayment handles incoming payment requests
func (s *SimulatedPaymentServer) HandlePayment(c *gin.Context) {
	s.failureCount--

	if s.failureCount > 0 {
		fmt.Printf("[%s] [Payment Server] Request failed (remaining failures: %d), returning 503\n",
			time.Now().Format("15:04:05.000"), s.failureCount)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Service temporarily unavailable",
		})
		return
	}

	fmt.Printf("[%s] [Payment Server] Request succeeded, returning 200\n", time.Now().Format("15:04:05.000"))
	c.JSON(http.StatusOK, gin.H{
		"status":         "success",
		"transaction_id": "txn_" + time.Now().Format("20060102150405"),
	})
}

// HandlePaymentWithIdempotency handles payments with idempotency key support
func (s *SimulatedPaymentServer) HandlePaymentWithIdempotency(c *gin.Context) {
	var req PaymentRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	key := c.GetHeader("Idempotency-Key")

	if cached, exists := s.idempotency.Get(key); exists {
		if cached.Completed {
			fmt.Printf("[%s] [Payment Server] Key '%s' - Returning cached response\n", time.Now().Format("15:04:05.000"), key)
			c.JSON(cached.StatusCode, gin.H{
				"status":         "paid",
				"amount":         1000,
				"transaction_id": "txn_cached_" + time.Now().Format("20060102150405"),
				"cached":         true,
			})
			return
		}
	}

	fmt.Printf("[%s] [Payment Server] Key '%s' - Processing started\n", time.Now().Format("15:04:05.000"), key)
	s.idempotency.StartProcessing(key)

	// Simulate heavy operation
	time.Sleep(2 * time.Second)

	fmt.Printf("[%s] [Payment Server] Key '%s' - Processing completed\n", time.Now().Format("15:04:05.000"), key)

	response := gin.H{
		"status":         "paid",
		"amount":         1000,
		"transaction_id": "txn_" + time.Now().Format("20060102150405"),
	}

	c.JSON(http.StatusOK, response)

	data, _ := c.GetRawData()
	s.idempotency.Finish(key, http.StatusOK, data)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	testServer := NewSimulatedPaymentServer()
	testRouter := gin.Default()
	testRouter.POST("/payment", testServer.HandlePayment)
	testRouter.POST("/payment-idempotent", testServer.HandlePaymentWithIdempotency)

	go func() {
		fmt.Println("[INFO] Test payment server listening on :9090")
		if err := testRouter.Run(":9090"); err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(500 * time.Millisecond)

	router := gin.Default()
	store := idempotency.NewMemoryStore()

	router.POST("/execute-payment", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cfg := payment.RetryConfig{
			MaxRetries: 5,
			BaseDelay:  100 * time.Millisecond,
			MaxDelay:   5 * time.Second,
		}

		fmt.Println("\n[CLIENT] Starting payment execution with retry logic...")
		statusCode, body, err := payment.ExecutePayment(ctx, cfg, "http://localhost:9090", nil)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		fmt.Printf("[CLIENT] Payment completed with status: %d\n", statusCode)
		c.JSON(statusCode, gin.H{
			"status": "completed",
			"result": string(body),
		})
	})

	router.POST("/execute-payment-idempotent", idempotency.Middleware(store), func(c *gin.Context) {
		var req PaymentRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		fmt.Println("\n[CLIENT] Processing payment with idempotency...")
		client := &http.Client{Timeout: 10 * time.Second}

		resp, err := client.Post("http://localhost:9090/payment-idempotent", "application/json", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		c.JSON(resp.StatusCode, gin.H{
			"status":         "success",
			"transaction_id": "txn_" + time.Now().Format("20060102150405"),
			"response":       string(body),
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Double-Click Attack Test Endpoint
	router.POST("/test-double-click", func(c *gin.Context) {
		type TestRequest struct {
			NumRequests    int    `json:"num_requests" binding:"required"`
			IdempotencyKey string `json:"idempotency_key" binding:"required"`
		}

		var req TestRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "num_requests and idempotency_key are required"})
			return
		}

		if req.NumRequests < 2 || req.NumRequests > 20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "num_requests must be between 2 and 20"})
			return
		}

		go SimultaneousPaymentRequests(req.NumRequests, req.IdempotencyKey)

		c.JSON(http.StatusAccepted, gin.H{
			"status":  "test_started",
			"message": fmt.Sprintf("Sending %d simultaneous requests with Idempotency-Key: %s", req.NumRequests, req.IdempotencyKey),
		})
	})

	fmt.Println("\n[INFO] Main server listening on :8080")
	fmt.Println("[INFO] Available endpoints:")
	fmt.Println("  - POST /execute-payment (with retry logic)")
	fmt.Println("  - POST /execute-payment-idempotent (with idempotency)")
	fmt.Println("  - POST /test-double-click (test with concurrent goroutines)")
	fmt.Println("  - GET /health")
	fmt.Println("\nTest server running on :9090")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
