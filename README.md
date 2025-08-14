# Salt HTTP Server

Salt is a modern, minimal, and high-performance HTTP server written in Go, designed for serving static sites with simplicity and speed.

## Features
- Minimal configuration
- Fast and secure by default
- Serves static files from the `public` directory
- Easy to deploy on any Linux server or container

## Getting Started

1. **Clone the repository:**
   ```sh
   git clone https://github.com/josuesantos1/salt
   cd salt
   ```

2. **Build and run:**
   ```sh
   go run main.go
   ```
   The server will start on port `1112` by default.

3. **Access your site:**
   Place your static files (HTML, CSS, JS, images) in the `public` folder. Visit `http://localhost:1112` in your browser.


## Example
```
$ go run main.go
Salt HTTP server starting on :1112...
```
