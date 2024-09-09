package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand" // 3. remove unused "math" is removed.
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
)

type rider struct {
	id        int64
	name      string
	entrance  string
	vipStatus bool
}

type rollercoaster struct {
	nextId      int64
	idMu        sync.Mutex
	rideQueue   []*rider
	rideQueueMu sync.Mutex
	ride        []*rider
	rideMu      sync.Mutex
}

const (
	rideDuration = 5000 // milliseconds
	numberOfCars = 8
	carCapacity  = 2
)

func main() {
	// Context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rc := rollercoaster{
		ride: make([]*rider, numberOfCars*carCapacity), // initialise slice
	}

	go rc.start(ctx)

	// start HTTP server
	err := http.ListenAndServe(":3000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d := &struct {
			Entrance string `json:"entrance"`
		}{}

		// Error handling for JSON
		if err := json.NewDecoder(r.Body).Decode(d); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		isVIP := false
		if rand.Float64() >= 0.7 {
			isVIP = true
		}

		// protect access to nextId using a mutex
		rc.idMu.Lock()
		rider := &rider{id: rc.nextId, entrance: d.Entrance, name: namesgenerator.GetRandomName(0), vipStatus: isVIP}
		rc.nextId++
		rc.idMu.Unlock()

		// 1. Thread Safety: protect access to rideQueue
		rc.rideQueueMu.Lock()
		if rider.vipStatus {
			rc.rideQueue = slices.Insert(rc.rideQueue, 0, rider)
		} else {
			rc.rideQueue = append(rc.rideQueue, rider)
		}
		rc.rideQueueMu.Unlock()
		log.Printf("Entrance %s: %s entered the queue. Size: %d\n", d.Entrance, rider.name, len(rc.rideQueue))
	}))
	if err != nil {
		panic(err)
	}
}

func (rc *rollercoaster) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done(): // handle cancellation
			log.Println("Context cancelled, stopping ride.")
			return
		default:
			rc.rideQueueMu.Lock() // lock ride to prevent concurrent modifiation
			rc.rideMu.Lock()      // lock rideQueue for whole operation

			i := 1

			// Ensure there is a rider in the queue before accessing it
			for i < len(rc.rideQueue) && i < numberOfCars*carCapacity {
				if len(rc.rideQueue) == 0 { // Exit loop if no riders in the queue
					break
				}

				r := rc.rideQueue[0]
				rc.rideQueue = rc.rideQueue[1:]

				car := i / carCapacity // (6) calculate the car seat
				carSeat := i % carCapacity
				rc.seatRider(r, car, carSeat)
				i++
			}
			rc.rideQueueMu.Unlock() // Unlock the queue
			rc.rideMu.Unlock()      // Unlock the ride

			if i > 0 {
				log.Println("Started: Ride")
				time.Sleep(rideDuration * time.Millisecond)
			}
		}
		log.Println("Finished: Ride")
	}
}

func (rc *rollercoaster) seatRider(r *rider, car, carSeat int) {
	//Log and seat the rider
	log.Printf("Ride: %s entering car %d in seat %d\n", r.name, car+1, carSeat)
	time.Sleep(1000 * time.Millisecond)
	// rc.rideMu.Lock()
	rc.ride = append(rc.ride, r) //  (2)Protect access to ride
	time.Sleep(1000 * time.Millisecond)
	// rc.rideMu.Unlock()
	log.Printf("Ride: %s entered car %d in seat %d\n", r.name, car+1, carSeat)
}
