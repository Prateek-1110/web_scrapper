#  GO-web-scraper

A fast, concurrent web scraper built in Go — my first real project while learning Go's concurrency model.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-green?style=flat)
![Status](https://img.shields.io/badge/status-active-brightgreen?style=flat)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-blue?style=flat)

---

## 🧠 Why I built this

I'm learning Go and wanted a project that forces me to understand the core of what makes Go special — **concurrency**. This scraper uses goroutines, channels, and worker pools to scrape hundreds of URLs in parallel, which would be painful to do in Python but feels natural in Go.

This is a learning project. The code is open for suggestions, improvements, and contributions.

---

## ✨ Features

- ⚡ **Concurrent scraping** — configurable worker pool (goroutines + channels)
- 🔁 **Retry logic** — exponential backoff on failures (1s → 2s → 4s)
- 🛑 **Graceful shutdown** — press `Ctrl+C` and it saves whatever it has scraped so far
- 🐢 **Rate limiting** — cap requests per second so you don't get IP banned
- 📌 **Resume flag** — skip already scraped URLs using a checkpoint file
- 📊 **Rich data** — scrapes title, meta description, status code, link count, image count
- 🎨 **Colored terminal output** — clean summary with success/failure breakdown
- 💾 **Dual output** — saves results as both JSON and CSV

---

## 📦 Install

> Requires Go 1.21+

```bash
go install github.com/Prateek-1110/web_scrapper@latest
```

Or clone and run locally:

```bash
git clone https://github.com/Prateek-1110/web_scrapper.git
cd web_scrapper
go run main.go urls.txt
```

---

## 🚀 Usage

Create a `urls.txt` file with one URL per line:

```
https://example.com
https://go.dev
https://github.com
```

Then run:

```bash
# Basic usage
go run main.go urls.txt

# With options
go run main.go -workers 20 -timeout 5s urls.txt

# Limit to 10 requests per second
go run main.go -rate 10 urls.txt

# Resume a previous scrape (skips already done URLs)
go run main.go -resume urls.txt

# Custom output files
go run main.go -json out.json -csv out.csv urls.txt
```

---

## ⚙️ Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-workers` | `10` | Number of concurrent workers |
| `-timeout` | `10s` | Per-request timeout |
| `-rate` | `0` | Max requests/sec (0 = unlimited) |
| `-resume` | `false` | Skip already scraped URLs |
| `-json` | `results.json` | JSON output filename |
| `-csv` | `results.csv` | CSV output filename |

---

## 📁 Project structure

```
web_scrapper/
├── main.go              ← CLI entry point, worker pool, orchestration
├── fetcher/
│   ├── fetcher.go       ← HTTP client, retry logic, HTML parsing
│   └── checkpoint.go    ← Resume flag checkpoint logic
├── output/
│   ├── writer.go        ← JSON and CSV writers
│   └── summary.go       ← Colored terminal summary
├── urls.txt             ← Your input URLs
└── go.mod
```

---

## 📄 Output

**results.json**
```json
[
  {
    "url": "https://go.dev",
    "title": "The Go Programming Language",
    "description": "Go is an open source programming language...",
    "status_code": 200,
    "link_count": 42,
    "image_count": 5,
    "attempts": 1
  }
]
```

**results.csv**
```
URL,Title,Error
https://go.dev,The Go Programming Language,
https://example.com,Example Domain,
```

---

## 🧩 What I learned building this

- **Goroutines** — lightweight threads managed by the Go runtime
- **Channels** — goroutine-safe pipes for passing data between workers
- **Worker pools** — producer/consumer pattern using buffered channels
- **WaitGroup** — synchronizing goroutine completion
- **context.Context** — timeout and cancellation propagation
- **Exponential backoff** — standard retry pattern for HTTP clients
- **Rate limiting** — using `time.Ticker` to cap request throughput

---

## 🤝 Contributing

This is a learning project and I'm actively improving it. If you spot something that could be done better — a more idiomatic Go pattern, a smarter approach, a bug — please open an issue or PR. All feedback is welcome.

```bash
git clone https://github.com/Prateek-1110/web_scrapper.git
cd web_scrapper
go mod tidy
go run main.go urls.txt
```

---

## 📋 Roadmap

- [ ] Add `--depth` flag for following links recursively
- [ ] Export to SQLite
- [ ] Docker support
- [ ] Prometheus metrics endpoint
- [ ] Web UI to visualize results

---

## 📜 License

MIT — do whatever you want with it.

---

> Built while learning Go. Feedback and stars are both appreciated ⭐