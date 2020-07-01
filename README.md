# anime-downloader-api

This is an api that downloads and can compress the anime you send to it as post request to localhost:1337/download.
Needed parameters are name and ep. Additionally, you can specify a custom resolution if you dont want it to default to 480p.
Also, unless your username is shelltear, change it to yours in anime-cli/src/main.rs:112 before compiling.


To run this, first build the downloader:
```
# sudo apt install openssl pkg-config
# cd anime-cli && cargo build --release --no-default-features
```

Then move the file to /usr/bin/ and create the dir to download anime to:
```
# sudo cp target/release/anime-cli /usr/bin
# mkdir ~/AnimeDownloads
```
And lastly, go back to root of this repo and run the api:
```
# go run api.go
```
Examples
```
Without custom resolution:
# curl -X POST localhost:1337/download -d "name=sword art online 0&ep=1"

With custom resolution (1080p):
# curl -X POST localhost:1337/download -d "name=sword art online 0&ep=1&reso=1080"

With compression:
# curl -X POST localhost:1337/download -d "name=sword art online 0&ep=1&compress=true"

Search:
# curl -X POST localhost:1337/search -d "name=charlotte&ep=1&reso=720p"
Note: on search if you want custom resolution, you need to add "p" to the resolution. 720 -> 720p. Also compression uses ffmpeg.
```
## Disclaimer
When downloading anime, users are subject to country-specific software distribution laws. This is not designed to enable illegal activity. We do not promote piracy nor do we allow it under any circumstances. You should own an original copy of every content downloaded through this tool. Please take the time to review copyright and video distribution laws and/or policies for your country before proceeding.


