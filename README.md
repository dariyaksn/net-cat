# TCP-Chat (NetCat Recreation)

A robust, concurrent Group Chat server written in Go that recreates the core behavioral mechanics of the `nc` (NetCat) utility using a Server-Client Architecture. The server listens for incoming TCP connections on a specified port, handles multiple clients concurrently, manages chat history, and broadcasts user actions in real-time.

---

## 🚀 Features

* **TCP Architecture:** 1-to-many connection architecture implemented via standard `net` sockets.
* **Concurrency Management:** Heavy use of Go-routines to handle multiple incoming client connections simultaneously without blocking the main server thread.
* **Thread Safety:** System states, user lists, and logs are protected using `sync.Mutex` to prevent race conditions.
* **Connection Control:** Strict restriction allowing a maximum of **10 concurrent connections** at any given time.
* **Interactive Onboarding:** Greets newly connected clients with an ASCII-art Linux penguin logo and enforces a non-empty username requirement.
* **Real-time Notifications:** Dynamically informs the group when a user joins or leaves the chat environment.
* **Message Formatting & History:** Every message is timestamped and structured as `[YYYY-MM-DD HH:MM:SS][client.name]:[client.message]`. New users automatically receive the entire chat history upon validation.
* **Input Validation:** Suppresses broadcast of empty or whitespace-only messages.
* **Stability:** If a client drops or exits, the remaining group chat components and clients remain completely unaffected.

---

## 🛠️ Project Structure

```text
net-cat/
├── cmd/
│   └── TCPChat/
│       └── main.go       # Entry point: handles CLI args and starts the server
├── netcat/
│   ├── server.go         # Core TCP server implementation, connection loops
│   ├── types.go          # Struct models (Server, User, Message)
│   └── welcome.txt       # Embedded ASCII Art welcome banner
├── .gitignore
├── go.mod
└── README.md
📥 Prerequisites
Go (version 1.16 or higher recommended)

NetCat (nc) client utility (built into Linux/macOS, available via WSL on Windows)

💻 Usage Instructions
Running the Server
By default, if no port is specified, the application falls back to port 8989.

Bash
# Run on default port (8989)
$ go run ./cmd/TCPChat

# Run on a custom port
$ go run ./cmd/TCPChat 2525
If too many arguments or invalid formats are provided, the system safely aborts with an instructional message:

Bash
$ go run ./cmd/TCPChat 2525 localhost
[USAGE]: ./TCPChat $port
Connecting as a Client
Open a separate terminal window and use nc to establish a connection with the server:

Bash
$ nc localhost 8989
💬 Chat Interface Preview
Upon a successful handshake, the client receives the following terminal workspace:

Plaintext
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]: Yenlik
[2026-06-22 18:03:43][Yenlik]: hello
[2026-06-22 18:03:46][Yenlik]: How are you?
[2026-06-22 18:04:10][Yenlik]: 
Lee has joined our chat...
[2026-06-22 18:04:32][Lee]: Hi everyone!
🧰 Allowed Standard Packages
This project strictly complies with academic constraints and implements network synchronization using only core packages:

net — Network socket primitives (TCP server)

sync — Synchronization tools (sync.Mutex for map control)

time — Explicit message timestamp generations

bufio / io — Stream data reading and output tracking

fmt / os / strings / log — System controls and formatting operations

