
Graceful Shutdown http server

EXPECTED:
* after calling Shutdown, live requests will get to finish processing


ACTUAL:
* After calling shutdown, http2 requests immediately send a GOAWAY frame to client
* 2 seconds later, server handler abrubtly stops executing, no panic, no error, etc.

WORKAROUND:
* Disable http2 by adding to the `http.Server` instantiation: `TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),`

(POSSIBLY) RELATED ISSUES:
* https://github.com/golang/go/issues/39776


REPRO:

1. Launch 3 terminals
1. In first terminal, run `make run`
1. In 2nd terminal run `tail -f server.log`
1. In 3rd terminal run `make curl`
1. back in first terminal, run `make stop`