package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
)

type Message struct {
    Username string `json:"username"`
    Content  string `json:"content"`
    Room     string `json:"room"`
}

var (
    rooms    map[string][]Message
    roomList []string
    mu       sync.Mutex
)

func main() {
    rooms = make(map[string][]Message)
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/", fs)
    http.HandleFunc("/send", sendMessage)
    http.HandleFunc("/receive", receiveMessages)
    http.HandleFunc("/create-room", createRoom)
    http.HandleFunc("/rooms", listRooms)
    fmt.Println("Server started at :8080")
    http.ListenAndServe(":8080", nil)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        var msg Message
        err := json.NewDecoder(r.Body).Decode(&msg)
        if err != nil {
            fmt.Println("Error decoding message:", err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        fmt.Printf("Received message: %+v\n", msg)
        mu.Lock()
        rooms[msg.Room] = append(rooms[msg.Room], msg)
        mu.Unlock()
        w.WriteHeader(http.StatusOK)
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

func receiveMessages(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        room := r.URL.Query().Get("room")
        mu.Lock()
        msgs := rooms[room]
        mu.Unlock()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(msgs)
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

func createRoom(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        var room struct {
            Room string `json:"room"`
        }
        err := json.NewDecoder(r.Body).Decode(&room)
        if err != nil {
            fmt.Println("Error decoding room:", err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        mu.Lock()
        if _, exists := rooms[room.Room]; !exists {
            rooms[room.Room] = []Message{}
            roomList = append(roomList, room.Room)
        }
        mu.Unlock()
        w.WriteHeader(http.StatusOK)
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}

func listRooms(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        mu.Lock()
        defer mu.Unlock()
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(roomList)
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }
}
