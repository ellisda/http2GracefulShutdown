package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log"
)

var testCACrt = []byte(`-----BEGIN CERTIFICATE-----
MIICxDCCAaygAwIBAgIBATANBgkqhkiG9w0BAQsFADASMRAwDgYDVQQDEwd0ZXN0
LWNhMCAXDTE3MDcyMDIxMTc1MloYDzIxMTcwNjI2MjExNzUzWjASMRAwDgYDVQQD
Ewd0ZXN0LWNhMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuv/sT2xH
VS1/uXVNAEIwvEb2yTMbXwP6FD38LWkc37Ri7YMB9xiXEDBrbr6K1JThsqyitBxU
22QIl53LUm6I7c/vej1tdYtE2rDVuviiiRgy6omR8imVSv9vU024rgDe+nC9zTT1
3aNKR03olCG6fkygdcZOghzlyQLhyh8LG75XdnLNksnakum2dNxQ5QIFmBKAuev3
A069oRMNjudot+t/nFP9UDZ8dL80PNTNPF22bPsnxiau7KLZ4I0Lf7gt6yHlNcue
Fd5sqzqsw/LUFJR5Xuo1+0e7NV3SwCH5CymG6hkboM4Rf5S3EDDyXTxPbXzbQHf1
7ksW6gjAxh4x/wIDAQABoyMwITAOBgNVHQ8BAf8EBAMCAqQwDwYDVR0TAQH/BAUw
AwEB/zANBgkqhkiG9w0BAQsFAAOCAQEATgmDrW1BjFp+Vmw6T+ojVK4lJuIoerGw
TCCqabHs6O1iWkNi5KsY6vV86tofBIEXsf6S3mV2jcBn87+CIbNHlHFKrXwmcydA
WOc0LWVqqoeqIvEcMNoWQskzmOOUDTanX9mXkirm8d8BljC351TH17rSjLGzFuNh
Cy48xyKFM7kPauNZGfCyaZsGbNJP3Keeu35dOLZMDdBJw7ZvYEUqX7MLOO+d7vlO
JGNA5jsU2uBnSo6qsjxfsbGemk2uRO0nLhplWurw+4qzA79D0lKNLtH9yTn12KZb
/kUpsOSCtLomjWwp67lQyA/yFXf897pSKMXbnIfZfIlDg51CI3U2Sw==
-----END CERTIFICATE-----`)

// valid for hostname test-service.test-ns.svc
// signed by testCACrt
var svcCrt = []byte(`-----BEGIN CERTIFICATE-----
MIIDDDCCAfSgAwIBAgIBBDANBgkqhkiG9w0BAQsFADASMRAwDgYDVQQDEwd0ZXN0
LWNhMCAXDTE3MDcyMDIxMjAzN1oYDzIxMTcwNjI2MjEyMDM4WjAjMSEwHwYDVQQD
Exh0ZXN0LXNlcnZpY2UudGVzdC1ucy5zdmMwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQDOKgoTmlVeDhImiBLBccxdniKkS+FZSaoAEtoTvJG1wjk0ewzF
vKhjbHolJ+/qEANiQ6CpTz4hU3m/Iad6IrnmKd1jnkh9yKEaU32B2Xbh6VaV7Sca
Hv4cKWTe50sBvufZinTT8hlFcGufFlJIOLXya5t6HH1Ld7Xf2qwNqusHdmFlJko7
0By8jhTtD7+2OAJsIPQDWfAsXxFa6LeQ/lqS2DCFnp45DirTNetXoIH8ZJvTBjak
bQuAAA3H+61gRm1blIu8/JjHYTDOcUe5pFyrFLFPgA+eIcpIbzTD61UTNhVlusV2
eRrBr5BlRM13Zj6ZMcWp0Iiw5QI/W9QU7O4jAgMBAAGjWjBYMA4GA1UdDwEB/wQE
AwIFoDATBgNVHSUEDDAKBggrBgEFBQcDATAMBgNVHRMBAf8EAjAAMCMGA1UdEQQc
MBqCGHRlc3Qtc2VydmljZS50ZXN0LW5zLnN2YzANBgkqhkiG9w0BAQsFAAOCAQEA
kpULlml6Ct0cjOuHgDKUnTboFTUm2FHJY27p4NXUvXUMoipg/eSxk0r5JynzXrPa
jaJfY2bC45ixLjJv9irp9ER/lGYUcBQ8OHouXy+uKtA5YKHI3wYf8ITZrQwzRpsK
C5v7qDW+iyb9dn4T6qgpZvylKtkH5cH31hQooNgjZd5aEq3JSFnCM12IVqs/5kjL
NnbPXzeBQ9CHbC+mx7Mm6eSQVtGcQOy4yXFrh2/vrIB2t4gNeWaI1b+7l4MaJjV/
kRrOirhZaJ90ow/PdYrILtEAdpeC/2Varpr3l4rIKhkkop4gfPwaFeWhG38shH3E
eG5PW2waPpxJzEnGBoAWxQ==
-----END CERTIFICATE-----`)

var svcKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAzioKE5pVXg4SJogSwXHMXZ4ipEvhWUmqABLaE7yRtcI5NHsM
xbyoY2x6JSfv6hADYkOgqU8+IVN5vyGneiK55indY55IfcihGlN9gdl24elWle0n
Gh7+HClk3udLAb7n2Yp00/IZRXBrnxZSSDi18mubehx9S3e139qsDarrB3ZhZSZK
O9AcvI4U7Q+/tjgCbCD0A1nwLF8RWui3kP5aktgwhZ6eOQ4q0zXrV6CB/GSb0wY2
pG0LgAANx/utYEZtW5SLvPyYx2EwznFHuaRcqxSxT4APniHKSG80w+tVEzYVZbrF
dnkawa+QZUTNd2Y+mTHFqdCIsOUCP1vUFOzuIwIDAQABAoIBABiX9z/DZ2+i6hNi
pCojcyev154V1zoZiYgct5snIZK3Kq/SBgIIsWW66Q9Jplsbseuk+aN46oZ7OMjO
MPZm8ho84EYj+a3XozBKyWwWDxKADW4xLjr1e4bMgVX97Xq11V6kH6+w78bS1GPT
+9jVuw7CO3fjsiawjye3JFM1Enh/NeRLEpT/oaQoWIV8b0IQB0VyqrdxWOO0rQhd
xA5w39tAZPDQ79MbMQyNWtPgBy0FuulP0GB12PrEbE+SXxsFhWViEwdB5Qx6Gqsx
KGn9vB1oaeSuuKIAjyBV0rXszrGektorDchsOY9UQi1mQsPSvvRFTM9T3qqSFIpu
oPNQLvECgYEA3ox3WJGjEve6VI4RMRt0l6ZFswNbNaHcTMPVsayqsl9KfebG+uyn
Z7TyyoCRzZZQa+3Z9jjW3hAGM9e7MG8jkeHbZpJpZv9X7eB3dgq3eZ1Zt5dyoDrU
PTdIPA2efFAf6V1ejyqH9h6RPQMeAb4uFU9nbI4rPagMxRdp5qIveIUCgYEA7Scb
0zWplDit4EUo+Fq80wzItwJZv64my8KIkEPpW3Fu6UPQvY74qyhE2fCSCwHqRpYJ
jVylyE0GIMx42kjwBgOpi4yEg8M3uMTal+Iy9SgrxZ5cPetaFpEF3Wk7/tz6ppr+
wnZQTO2WH3YLzv7JIWVrOKuBNVfNEbguVFWw4IcCgYB54mp2uoSancySBKDLyWKo
r6raqQrqK7TQ4iyGO6/dMy1EGQF/ad8hgEu8tn+kHh/7jG/kVyruwc3z1MIze5r6
ib00xxktDMnmgRpMLwBffdsmHq7rrGyS/lT0du0G3ocrszRXqo5+MC2RQcTMZZEt
oKhfHtn10bT0uKcKZmcjVQKBgEls2WWccMOuhM8yOowic+IYTDC1bpo1Tle6BFQ+
YoroZQGd+IwoLv+3ORINNPppfmKaY5y7+aw5hNM025oiCQajraPCPukY0TDI6jEq
XMKgzGSkMkUNkFf6UMmLooK3Yneg94232gbnbJqTDvbo1dccMoVaPGgKpjh9QQLl
gR0TAoGACFOvhl8txfbkwLeuNeunyOPL7J4nIccthgd2ioFOr3HTou6wzN++vYTa
a3OF9jH5Z7m6X1rrwn6J1+Gw9sBme38/GeGXHigsBI/8WaTvyuppyVIXOVPoTvVf
VYsTwo5YgV1HzDkV+BNmBCw1GYcGXAElhJI+dCsgQuuU6TKzgl8=
-----END RSA PRIVATE KEY-----`)

func main() { //nolint:funlen //DELETEME after test code removal

	log.Println("Starting service")

	cert, err := tls.X509KeyPair(svcCrt, svcKey)
	if err != nil {
		log.Fatalf("backend: invalid x509/key pair: %v", err)
	}

	server := simpleServer(":8443", &cert)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("listening for incoming requests...")
	go func() {
		if err := server.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("listenAndServe", err)
		}
	}()

	<-signalChan
	log.Println("shutdown signal received")
	log.Println("shutting down webhook server gracefully")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatal(err, "failed to shut down webhook server gracefully")
	}

	log.Println("shutdown has returned")
}

func simpleServer(addr string, cert *tls.Certificate) *http.Server {

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("received request")
			for i := 0; i < 1e3; i++ {
				time.Sleep(30 * time.Millisecond)
				log.Println(fmt.Sprintf("sleeping %d", i))
			}
			log.Println(" -- done sleeping")

			defer func() { log.Println("defer complete") }()

			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered in f", r)
				}
			}()
			_, _ = w.Write([]byte("All good\n"))
			log.Println("wrote response")
		},
	))

	return &http.Server{
		Addr:    addr,
		Handler: mux,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*cert},
			MinVersion:   tls.VersionTLS12,
		},
		// TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
}
