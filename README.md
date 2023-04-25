# Emails Database

Emails Database is a database indexer written in Go that indexes a file-based MIME-emails database to a search engine called ZincSearch. It also includes a client-side interface built with Vue and Tailwind CSS for easy visualization of the emails data, and an API implemented using Chi router. Additionally, profiling is incorporated for performance analysis to optimize the indexing process.

## Indexer

To set up the indexer, follow these steps:

1. Download and install Go: [Go Installation Guide](https://golang.org/doc/install)
2. Download the emails database.
3. Download ZincObserve from 'docs.zinc.dev'.
4. Run the following commands in the terminal to set up ZincObserve:
	```
	set ZO_ROOT_USER_EMAIL=root@example.com
	set ZO_ROOT_USER_PASSWORD=Complexpass#123
	zincobserve.exe
	```
5. Open 'localhost:5080' in a web browser and sign in with the credentials provided above.
6. From the ingestion tab, copy the username and password token.
7. Create the indexer using the copied token.

## Profiling

Profiling is included in the indexer for performance analysis. Follow these steps to access and analyze the profile information:

1. Access the profile information during execution at 'http://localhost:8080/debug/pprof/'.
2. Download the profile using the command 'curl -o cpu.prof http://localhost:8080/debug/pprof/profile' or by visiting 'http://localhost:8080/debug/pprof/profile'.
3. Analyze the profile using the command 'go tool pprof cpu.prof' or by opening it on the localhost with 'go tool pprof -http=:8080 cpu.prof'.

## Usage

1. Start the Go indexer: `./indexer <database name>`
2. Start the Vue client: `cd client && npm run dev`
3. Open a web browser and go to `http://localhost:8080` to access the client-side interface.
4. Use the interface to search and visualize the emails data from the ZincSearch engine.
