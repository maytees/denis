# DENIS
> [!NOTE]
> Code in this repo is not AI generated unless if specified, everything is hand written. If any parts were *aided* (not written) by AI, the file/block will mention.

> [!NOTE]
> Apple uses .local domains for their multicast dns, so adding a domain with .local won't work on MacOS; they use port `5353` for this.


## What is it?

DENIS right now is a working DNS server ([RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035)) meant to run on your local computer or local network. DENIS is just a hobby project for [me](https://maytees.net) learn golang.

<!-- The Documentation section below and the docs site (www/) were written by AI (Claude) -->

## Documentation

This README covers the basics. For the full walkthrough (configuration, records, testing, Docker), see the docs site: **[denis.matees.net/docs](https://denis.matees.net/docs)** — built with [Bloom](https://github.com/maytees/bloom), lives in [`www/`](./www).



https://github.com/user-attachments/assets/7efd7c2a-f245-4340-871b-511c58555be9



## Vision

The idea for DENIS came from wanting to route different local web servers to their own domains. For example, my [Memos](https://usememos.com/) notes server runs on `localhost:5230`, instead of going to that, I wanted to just go to `http://notes`. But there's two things to mention here:

1. Why not just use your local hosts file? (your computers DNS)
 
	- Why not build my own?

2. A DNS alone cant do this, beacuse all a DNS does is route you to an IP.
 
	- Thats what the second part of DENIS is for 👇👇

### More than a DNS

To make the *vision* work, I need to have a reverse proxy running on port 80 (http) wherever I have web apps, so the dns would route a domain to said IP, and then the reverse proxy goes to whatever port the web app is on. I could just use nginx or any other reverse proxy, but I'm choosing to abide by the *build my own* stance and make a custom http server to handle requests.

Moreover, to easily create, remove, and edit records and proxy routes, there will be a NextJS app running to do said actions, connected to an api service running in the DENIS binary

That being said, just note that this is the **vision**, at it's current state, DENIS is just a DNS server that works, but doesn't fully implement the [RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035) protocol fully *(yet)*

## Try it out

Here are the steps to getting DENIS up and running locally:

## Prerequisites 

Before you get started, you *should* install the following:

- [go](https://go.dev)
- [air](https://github.com/air-verse/air) - hot reloading (*optional for dev*)
- [just](https://github.com/casey/just) - task runner (*optional*)
- [gotestsum](https://github.com/gotestyourself/gotestsum) - pretty testing (*optional*)

1. Fork the repo

```bash
git clone https://github.com/maytees/denis
cd denis
```

2. Create config files

> [!NOTE]
> A proper CLI including args for the config folder is in the works. Detecting if a folder is in the the working dir is terrible design. If there isn't a config folder in your current dir, one will be created in your config dir depending on your OS ([where's that?](https://github.com/kirsle/configdir?tab=readme-ov-file#configdir-for-go)).

If you're running denis from the repo folder (if it detects a config folder in the current repo) just run:

```bash
# with just
just config-example

# without
cp ./config/config.example.toml ./config/config.toml
cp ./config/records.example.toml ./config/records.toml
```

3. Run it

```bash
# with just
just start

# without
go build -o dist/denist && sudo dist/denis
```

Or if you wish to run a dev server with air

```bash
# with just
just dev

# without
sudo air
```

4. Send query

Dig is a command line tool used to query DNS servers for records. If you haven't set DENIS as your default DNS, you need to pass the @addr (+ :port if you aren't running on :53) to tell dig to query DENIS.

Here's an example

```bash
dig @127.0.0.1 localhost
```

If DENIS is working correctly, it should grab the record from your records.toml, or fallback to your upstream dns server.

If it doesn't, dig should hang and ultimately error.

<!-- This section was written by AI (Claude) -->

## Run with Docker

Don't want to install Go? DENIS runs in a container — you only need [Docker](https://docs.docker.com/get-docker/).

```bash
# 1. create config files (see step 2 above)
just config-example

# 2. in config/config.toml, set:
#    host = '0.0.0.0'   (required in Docker, still loopback-only on your machine)

# 3. build & run
just docker-compose   # or: docker compose up -d --build

# 4. test it
just docker-check     # or: dig @127.0.0.1 localhost && dig @127.0.0.1 google.com
```

Edit `records.toml` anytime, then `just docker-restart` to apply. After changing code or pulling, rerun `just docker-compose`.

For the full guide (updating, using DENIS as your system DNS, the escape hatch when it's down), see **[denis.matees.net/docs/docker](https://denis.matees.net/docs/docker)**.

## Testing
Tests are WIP for DENIS, run them with the following

```bash
# just & with gotestsum
just prettytest

# or without gotestsum
just test

# without just
go test ./...
```

## Configuring
The config directory contains two files

### config.toml
This contains general settings and settings specific to the DNS.

Here's an example

```toml
[dns]
enabled = true # Toggle DNS
port = 53 # Default port for DNS
upstream = '8.8.8.8' # Where to route upstream DNS requesets. Google, Cloudflare, your router (common), etc.
host = '127.0.0.1' # Address to bind to. Use 0.0.0.0 when running in Docker
```
 
### records.toml
Stores all your DNS records

Here's an example

```toml
[[records]]
name = 'localhost' # Domain
type = 'A' # Record type (only A available)
value = '127.0.0.1' # Address
ttl = 300 # Cache time

[[records]]
name = 'my.google'
type = 'A'
value = '142.251.16.138'
ttl = 0
```

## How to route requests to DENIS automatically
This depends on your operating system. The idea is to route everything through DENIS, instead of your router's assigned dns.

If your OS isn't listed below, I reccomend you search up instructions for it.

### MacOS
Telling MacOS where to handle DNS requests is very simple. You need to go to your **System Settings** -> **Wifi** -> Click **Details** on the wifi you're connected to -> **DNS** -> Click **+** under DNS Servers, enter in where your server is running.

**Note: your DNS must be running on port 53**

### Browser not using DENIS properly?

If your browser seems to be ignoring DENIS even though you told your OS to use it, it's probably because your browser is using DoH (DNS over https).

For **Chrome based** browsers, go to your security settings and tell Chrome to use your OS' default dns.

For **Firefox**, go to the firefox Network Settings, and set DoH to off.

Alternatively, you could set exceptions for domains you want your browser to use DENIS for.

## Contributing

Anyone is free to contribute to DENIS. Only rule would be to not vibe code anything, or get AI to write any line of code by itself. If you want to contribute, please handwrite and don't give me slop. Your own handwritten slop is acceptable for review though.

## License

DENIS is under the MIT license, you are free to do whatever you want with the code.
