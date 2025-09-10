package main

import (
	"context"
	"fmt"
	"log"
	"maps"
	"net/url"
	// "reflect"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/fatih/color"
	"strings"
	"time"
)

type PageCalls struct {
	URL      string
	anaCalls map[string][]string
}

func runBrowser() []PageCalls {
	calls := make([]PageCalls, 0)
	// Create a new context
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))
	ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, _ = chromedp.NewContext(ctx)
	// defer cancel()

	// Create a timeout to prevent the program from running indefinitely
	// ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	// defer cancel()

	// Create a channel to signal when we are done
	// done := make(chan bool)

	// Set up the listener for network events
	chromedp.ListenTarget(ctx, func(ev any) {
		// Check if the event is the one we're interested in
		if req, ok := ev.(*network.EventRequestWillBeSent); ok {
			// Filter for specific requests, e.g., POST requests to an analytics endpoint
			isPost := req.Request.Method == "POST"
			// isAnalytics := strings.Contains(req.Request.URL, "analytics")

			if isPost {
				log.Println("✅ Call detected")
				currentCall := PageCalls{
					URL:      req.Request.URL,
					anaCalls: make(map[string][]string, 0),
				}
				switch {
				case strings.Contains(req.Request.URL, "https://www.google-analytics.com"):
					log.Println("Google Analytics Detected")
				}
				params, err := url.ParseQuery(req.Request.URL)
				currentCall.URL = req.Request.URL
				if err != nil {
					log.Fatalf("Failed to prase query: %s", err)
				}
				maps.Copy(currentCall.anaCalls, params)
				calls = append(calls, currentCall)
				// log.Printf("   Request Body: %s", req.Request.PostData)
			}
		}

		// You can also listen for when the page is fully loaded
		// if _, ok := ev.(*page.EventLoadEventFired); ok {
		// 	close(done) // Signal that the page has loaded
		// }
	})

	// // Define the tasks to run
	if err := chromedp.Run(ctx,
		// 1. Enable network event listening
		network.Enable(),
		// 2. Navigate to the target page
		chromedp.Navigate(`http://www.concordusa.com/cja`), // Example page
		// 3. Wait for a moment to allow async analytics calls to fire
		chromedp.Sleep(5*time.Second),
	); err != nil {
		log.Fatalf("Failed to run chromedp tasks: %v", err)
	}

	log.Println("Finished capturing network requests.")
	return calls

}

func runTests(testConfig TestConfig, calls []PageCalls) {
	red := color.New(color.FgRed).PrintlnFunc()
	green := color.New(color.FgGreen).PrintfFunc()
	log.Println("Running tests...")
	for i, test := range testConfig.Tests {
		fmt.Printf("--- Test Case %d ---\n", i+1)
		fmt.Printf("Name: %s\n", test.Name)
		fmt.Printf("URL: %s\n", test.URL)
		fmt.Println("Dependencies:")
		for _, dep := range test.Deps {
			fmt.Printf("  - Key: %s, Value: %s\n", dep.Key, dep.Value)
		}

		fmt.Println("Comparisons:")
		for _, comp := range test.Compare {
			fmt.Printf("  - Key: %s, Value: %s\n", comp.Key, comp.Value)
		}
		for _, call := range calls {
			if strings.Contains(call.URL, test.URL) {
				log.Println("Test Match")
				log.Println("Checking Deps...")
				// Assume the current call is a match until proven otherwise.
				allDepsMatch := true
				// 2. Inner loop: Check every required dependency for the current call.
				for _, dep := range test.Deps {
					val, ok := call.anaCalls[dep.Key]
					// Check if the property ekxists and if its value matches the dependency's value.
					if !ok || val[0] != dep.Value {
						// If a dependency is missing or the value is wrong, it's not a match.
						allDepsMatch = false
						break // Exit the inner loop; no need to check other deps for this call.
					}
				}
				if allDepsMatch {
					log.Println("Deps OKAY")
					log.Println("Testing...")
					allCompsMatch := true
					for _, comp := range test.Compare {
						val, ok := call.anaCalls[comp.Key]
						// Check if the property ekxists and if its value matches the dependency's value.
						if !ok || val[0] != comp.Value {
							log.Printf("No Match - Call Value: %s, Test Value: %s", comp.Value, val[0])
							// If a dependency is missing or the value is wrong, it's not a match.
							allCompsMatch = false
							break // Exit the inner loop; no need to check other deps for this call.
						} else {
							log.Printf("Match - Call Value: %s, Test Value: %s", comp.Value, val[0])
						}
					}
					if allCompsMatch {
						green("PASS! ")
						fmt.Printf("%s\n", test.Name)
					} else {
						red("FAIL! ")
						fmt.Printf("%s\n", test.Name)
					}
				}
			}
		}

	}

}

func main() {
	testConfig := parseTests("my_test.yaml")
	calls := runBrowser()
	runTests(testConfig, calls)
}
