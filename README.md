# Denis DNS server
> [!NOTE]
> None of the code in this repo was AI Generated, everything is hand written. (necessary to mention)

*Using port 5354 beacuse MacOS uses 5353*

This is a custom DNS ([RFC 1035](https://datatracker.ietf.org/doc/html/rfc1035)) server written by me. At its **current** state, Denis is not being used as a replacement for my local DNS Server, but hopefully will be in the future. This project was just made out of interest of learning Go.

## Testing
Simply use the `dig` command as you would with any other DNS. *At it's current state, nothing happens*

```bash
dig @127.0.0.1 -p 5354 google.com
```

## Action Plan

### Phase 0: Set up project
Self explanatory.

### Phase 1: Listen and receive
Listen on UDP port 5354 (or 53) and prints out the raw bytes from incoming packets. Test with `dig`


### Phase 2: Parse the header
DNS header is the first 12 bytes of every packet. Parse the content, see [Header Section Format](https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.15). Just print for now.


### Phase 3: Parse the question 
After the header comes the question section. Parse out the QNAME (has length prefix), QTYPE (record type), and QCLASS (just `in` (one) for now). All done with bytes


### Phase 4: Build a response *(CURRENT)*
For now, hardcode a response for one domain. Take the transaction ID from the query, set the response flags, include the question, add an answer section with your IP. Send it back. Test with dig and see if you get your IP.


### Phase 5: Add a lookup table
Replace the hardcoded response with a map (either db (sqlite) or json). Look up the queried domain, return the IP if you have it.


### Phase 6: Forwarding
If you don't have the domain, forward the query to 8.8.8.8 (google) or 1.1.1.1 (cloudflare), get the response, send it back.

### Phase 7: Chaching 
Store responses from upstream to make queries faster. Respect the TTL value from the response. Caching won't be necessary for local because it's already very fast. Only do it if you use a DB.
