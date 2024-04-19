package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/go-fuego/fuego"
)

var (
	domainName string
	acmeEmail  string
)

// main starts a simple Fuego 🔥 server with ACME certificates from Let's Encrypt
// It also uses a simple HTTP server to solve the ACME HTTP challenges and redirect any HTTP request to HTTPS
// For this to run proper you need to use real domain name (and email) pointing to a server with ports 80 and 443
func main() {
	flag.StringVar(&domainName, "td", "", "domain name to use for TLS")
	flag.StringVar(&acmeEmail, "te", "", "email address for ACME Server")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	certmagic.DefaultACME.Agreed = true
	certmagic.DefaultACME.Email = acmeEmail

	// avoid  Let's Encrypt rate limits during testing
	certmagic.DefaultACME.CA = certmagic.LetsEncryptStagingCA // remove this for real certificates

	magic := certmagic.NewDefault()
	myACME := certmagic.NewACMEIssuer(magic, certmagic.DefaultACME)

	// create a simple HTTP server
	tlsHttpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:80", domainName),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       5 * time.Second,
		BaseContext:       func(listener net.Listener) context.Context { return ctx },
	}
	httpRedirectHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+domainName+r.RequestURI, http.StatusMovedPermanently)
	})
	tlsHttpServer.Handler = myACME.HTTPChallengeHandler(httpRedirectHandler)
	// async start HTTP listener to redirect to HTTPS and solve ACME HTTP challenges
	go func() {
		if err := tlsHttpServer.ListenAndServe(); err != nil {
			log.Println("http listener error: ", err)
		}
	}()

	// get or renew certificate from the ACME server
	err := magic.ManageSync(ctx, []string{domainName})
	if err != nil {
		log.Fatalln("error getting certs from ACME server: ", err)
	}

	httpsServer := fuego.NewServer(fuego.WithAddr(fmt.Sprintf("%s:443", domainName)))

	// use the updated TLS configuration that includes the ACME certificates
	httpsServer.Server.TLSConfig = magic.TLSConfig()
	httpsServer.Server.TLSConfig.NextProtos = append([]string{"h2", "http/1.1"}, httpsServer.Server.TLSConfig.NextProtos...)

	fuego.Get(httpsServer, "/", func(_ fuego.ContextNoBody) (string, error) {
		return "Hello, World!", nil
	})

	// async start Fuego 🔥 in TLS mode
	go func() {
		log.Printf("server listening on %s\n", httpsServer.Server.Addr)
		// no need for cert and key files, since those are in tls.Config from certmagic
		err := httpsServer.RunTLS("", "")
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("tls listener error: %v", err)
		}
	}()

	<-ctx.Done() // Wait for SIGINT or SIGTERM

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpsServer.Server.Shutdown(shutdownCtx); err != nil {
		log.Printf("http server shutodwn error: %v", err)
	}
	if err := tlsHttpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("tls server shutodwn error: %v", err)
	}
	log.Println("Server stopped")
}
