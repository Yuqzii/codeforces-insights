# Codeforces Insights
Website for displaying analytics and providing insights into your Codeforces profile.\
Visit it at https://cf-insights.org.


## Feature Highlights
### Analytics Dashboard
This is split into four parts:
- general user info,
- rating distribution of solved problems,
- tag distribution of solved problems,
- and a chart over rating history.


### Performance Calculation
Performance calculation for the entire history of the user, displayed over the rating history chart.\
This uses the standard elo rating system to calculate the probability of one rating achieving a higher rank in a contest than another.\
It also utilizes Goroutines in a master-worker pattern for concurrent computation.\
For more details see the [performance documentation](docs/performance.md).


### Problem Recommendation
Practice problem recommendations based on the first problem the user did not solve in recent contests.\
This is done by representing each problem as a vector of its tags.
Then we can take the vector sum of the recent unsolved problems, and find other problems with vectors similar to this.\
For more details see the [recommendation documentation](docs/recommend.md).


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
For running the site locally use the provided `Makefile`.\
To build and start the containers normally:
```
make
```
and to stop them:
```
make stop
```
It is intended to visit the site locally through Nginx on port 443, for HTTP/2 with support for request cancellation.\
Hot reload is also automatically configured with esbuild to watch for changes and rebuild the frontend,
while Browsersync automatically reloads the page.

> [!NOTE]
> The Nginx container uses self-signed SSL certificates which will likely cause your browser to warn you about entering the site.\
> You can ignore this warning by clicking "Advanced" and then "Continue"/"Proceed" (or something similar depending on what browser you use).

When the Nginx container has started visit `https://localhost`.


### Fetcher
The job of the fetcher is to fetch and store Codeforces data in a database to avoid waiting for the Codeforces API when requests are made.\
It is mainly used as a cron job on the server, but is also useful for fetching data for local development.\
To run this you can use the `Makefile` with arguments depending on what you want to fetch.\
To fetch all data run:
```
make fetch ARGS="-contests -problems"
```

#### Contests
- `-contests` to enable contest fetching.
- `-maxContestAge` contests older than this will be re-fetched. Example: `-maxContestAge=24h` to fetch all problems older than 24 hours.
- `-maxContestUpdates` to set maximum amount of contest fetches.\
This takes priority over `-maxContestAge`. Example: `-maxContestUpdates=100`.

The fetcher finds all available Codeforces contests that are not currently stored in the database and fetches them.\
This drastically speeds up performance calculations, as the server does not have to wait on the Codeforces API.\
However, this is not mandatory, and the server figures out whether to use stored contests or fetch live from Codeforces on its own.

#### Problems
- `-problems` to enable problem fetching.

The fetcher finds all Codeforces problems with at least one tag and stores them in the database.\
This is used for recommending problems, and is required for it to work, unlike the contests.

> [!NOTE]
> The fetcher only stores problems of contests that already exist in the DB.

