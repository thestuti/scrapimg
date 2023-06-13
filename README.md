```
███████  ██████ ██████   █████  ██████  ██ ███    ███  ██████  
██      ██      ██   ██ ██   ██ ██   ██ ██ ████  ████ ██       
███████ ██      ██████  ███████ ██████  ██ ██ ████ ██ ██   ███ 
     ██ ██      ██   ██ ██   ██ ██      ██ ██  ██  ██ ██    ██ 
███████  ██████ ██   ██ ██   ██ ██      ██ ██      ██  ██████  
```


-----


<br>

A CLI tool written in Go for extracting images from a webpage. Given a URL, it fetches the HTML, parses it, and downloads all the images found on the page.

---


<h2> Installation </h2> 


``` 
git clone github.com/thestuti/scrapimg
go run src/main.go

```
---

<h2> Usage </h2> 


| Flag | Description                           | Example                              |
| ---- | ------------------------------------- | ------------------------------------ |
| -h   | Display help related to usage         | go run src/main.go -h                |
| -u   | Extract images from the url           | go run src/main.go -u https://example.example.com |
| -l   | Path where you want images to be stored | go run src/main.go -l  /path/to/save/location    |
| -f  | File where different urls are stored | go run src/main.go -f  urls.txt   |
