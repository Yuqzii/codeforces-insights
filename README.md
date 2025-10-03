# Codeforces Stats
Website for displaying stats about your Codeforces profile.

## Running dev environment
### Setting up SSL certs
These are necessary for nginx to run HTTPS and therefore HTTP/2.\
There are many ways to generate these, but I use
```
mkcrt localhost
```
Be sure to place these inside `/certs` from the project root.\
You can also optionally install these certificates with `mkcrt -install` to avoid a warning from the browser.

### Running services
For running the site locally use the provided `docker-compose.yml` with standard `docker compose` commands.\
It is intended to visit the site locally through nginx on port 443, although the frontend will still work when visiting through Browsersync on port 3000.
This is because nginx reverse-proxies requests with HTTPS and HTTP/2 to support request cancellation.\
Hot reload is also automatically configured with esbuild to watch for changes and rebuild the frontend,
while Browsersync automatically reloads the page.
To start the services normally run
```
docker compose up --build
```

### Contest Fetcher
There is also another entrypoint for the Go backend called `fetcher`.\
This finds all available Codeforces contests that are not currently stored in the local database and fetches them.\
It drastically speeds up performance calculations, as the server does not have to wait on the Codeforces API.\
This is completely optional, and the server figures out whether to use stored contests or fetch live from Codeforces on its own.\
To run this do
```
docker compose run --build fetcher
```
