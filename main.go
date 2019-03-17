package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	mr "math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type MessageSocket struct {
	Message string    `json:"message"`
	Status  WindWater `json:"status"`
}

type WindWater struct {
	Wind  int `json:"wind"`
	Water int `json:"water"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":4000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	// Open file json
	j, err := os.Open("angindanair.json")

	// Cek apakah ada error
	if err != nil {
		fmt.Println(err)
		return
	}

	// defer / close file sebelum return
	defer j.Close()

	// Read file reader untuk mendapatkan byte dari json tersebut
	byteJ, err := ioutil.ReadAll(j)

	// Cek apakah ada error
	if err != nil {
		fmt.Println(err)
		return
	}

	// deklarasi variable dari struct MessageSocket
	var data MessageSocket

	// Unmarshall atau parse dari json ke struct agar dapat di ganti valuenya
	// atau ditambahkan logic/business
	err = json.Unmarshal([]byte(byteJ), &data)

	// Cek apakah ada error
	if err != nil {
		fmt.Println(err)
		return
	}

	// Deklarasi variable websocket
	socket, err := upgrader.Upgrade(w, r, nil)

	// Cek apakah ada error
	if err != nil {
		log.Println(err)
		return
	}

	// Generate infinite loop agar koneksi websocket selalu alive
	for {

		// Melakukan always listen message websocket
		msgType, msg, err := socket.ReadMessage()

		// Cek apakah ada error
		if err != nil {
			log.Println(err)
			return
		}

		// Print messagetype dan message dari client
		fmt.Println(msgType, string(msg))

		// Generate infinite loop lagi agar mengirim message wind water secara
		// terus menerus
		for {

			// membuat variable pointer dari data(MessageSocket)
			dp := &data

			// Generate Random Number
			// deklarasi channel untuk goroutine
			chanInt := make(chan int)

			// Mengeksekusi goroutine dengan anonymous func
			// untuk mendapatkan random number
			go func() {

				rand.Seed(time.Now().UnixNano())
				min := 0
				max := 100

				for i := 0; i < 2; i++ {
					randNumber := mr.Intn(max - min)

					chanInt <- randNumber // send total to channelInt
				}
			}()
			rn1, rn2 := <-chanInt, <-chanInt // receive from channelInt

			// Mengisi message dan value dari water & wind
			dp.Message = "data baru"
			dp.Status.Water = rn1
			dp.Status.Wind = rn2

			// Parsing struct menjadi json kembali untuk mendapatkan []byte
			res, err := json.Marshal(data)

			// Cek apakah ada error
			if err != nil {
				fmt.Println(err)
				return
			}

			// Mengirim byte/json ke client
			err = socket.WriteMessage(msgType, res)

			// Cek apakah ada error
			if err != nil {
				fmt.Println(err)
				return
			}

			// Set interval loop per 15 detik
			time.Sleep(15 * time.Second)
		}

	}

}
