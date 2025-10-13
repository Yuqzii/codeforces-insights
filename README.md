# Codeforces Insights
Website for displaying analytics about your Codeforces profile.

## Dev Environment
### Docker
This project uses Docker and Docker Compose for containerization and container orchestration.\
If you want to install for use with a CLI you need to follow these instructions:
- [Docker Engine](https://docs.docker.com/engine/install) (this also shows how to install the Docker CLI).
- [Docker Compose](https://docs.docker.com/compose/install/linux/)
If you prefer a GUI you can install Docker Desktop instead. Here are the [official instructions](https://docs.docker.com/desktop/setup/install/linux/).
> [!NOTE]
> Docker only runs natively on Linux. If you are on Windows I recommend using WSL2. [Installation docs](https://learn.microsoft.com/en-us/windows/wsl/install).\
> You can also use Docker Desktop, which should set this up automatically. (Haven't tested this).

### Environment Variables
To make the PostgreSQL database work, you need to create a `.env` file with these variables:
- `POSTGRES_DB`, the name of the database.
- `POSTGRES_USER`, the name of the database user.
- `POSTGRES_PASSWORD`, the password to the user.
It does not matter what these are for the dev environment, but they must be set.

### Running Services
For running the site locally use the provided `docker-compose.yml` with standard `docker compose` commands.\
It is intended to visit the site locally through Nginx on port 443, for HTTP/2 with support for request cancellation.\
Hot reload is also automatically configured with esbuild to watch for changes and rebuild the frontend,
while Browsersync automatically reloads the page.
To start and rebuild the Docker images normally run
```
docker compose up --build
```
> [!NOTE]
> The Nginx container uses self-signed SSL certificates which will likely cause your browser to warn you about entering the site.\
> You can ignore this warning by clicking "Advanced" and then "Continue"/"Proceed".
When the Nginx container has started visit `https://localhost`.

### Contest Fetcher
There is also another entrypoint for the Go backend called `fetcher`.\
This finds all available Codeforces contests that are not currently stored in the local database and fetches them.\
It drastically speeds up performance calculations, as the server does not have to wait on the Codeforces API.\
This is completely optional, and the server figures out whether to use stored contests or fetch live from Codeforces on its own.\
To run this do
```
docker compose run --build fetcher
```
