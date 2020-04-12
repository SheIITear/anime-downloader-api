# anime-downloader-api

This is an api that downloads anime you send to it as post request to localhost:1337/download.
Needed parameters are name and ep. Additionally, you can specify a custom resolution if you dont want it to default to 480p.

To run this, first build the downloader:
>cd anime-cli && cargo build --release --no-default-features
Then move the file to /usr/bin/:
>sudo cp target/release/anime-cli /usr/bin
And lastly, go back to root of this repo and run the api:
>go run test-api.go
