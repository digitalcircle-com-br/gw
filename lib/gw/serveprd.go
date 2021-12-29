package gw

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"path"

	"github.com/digitalcircle-com-br/caroot"
	"github.com/digitalcircle-com-br/gw/lib/base"
	"github.com/digitalcircle-com-br/gw/lib/util"
	"golang.org/x/crypto/acme/autocert"
)

func init() {
	Serve = ServePrd
	caroot.InitCA("caroot", func(ca string) {
		log.Printf("Initiating CA: %s", ca)
	})
}

var ServePrd = func() {
	log.Printf("Initiating GW")

	if !config.Acme {

		if config.SelfSigned && config.Https {
			log.Printf("Using https + self signed approach")
			tlscfg := &tls.Config{
				GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
					ca := caroot.GetOrGenFromRoot(info.ServerName)
					return ca, nil
				},
			}

			server := &http.Server{
				Addr:      config.Addr,
				Handler:   http.DefaultServeMux,
				TLSConfig: tlscfg,
			}
			go func() {
				err := server.ListenAndServeTLS("", "")
				if err != nil {
					log.Printf("Finishing server: %s", err.Error())
				}
			}()

		} else {

			if config.Addr == "" {
				config.Addr = ":8080"
			}

			log.Printf("No acme set - if required, set ENV VAR ACME")

			if config.Https {
				log.Printf("Using simple ssl approach")
				conn, err := net.Listen("tcp", config.Addr)
				util.ErrP(err)
				server := http.Server{}
				log.Printf("APIGW - Running HTTPS @ %v", config.Addr)
				go func() {
					err = server.ServeTLS(conn, config.Cert, config.Key)
					util.ErrP(err)
					//e := server.ListenAndServeTLS("", "")
				}()
			} else {
				log.Printf("Using simple approach")
				server := &http.Server{Addr: config.Addr}
				log.Printf("APIGW - Running HTTP @ %v", config.Addr)
				go func() {
					err := server.ListenAndServe()
					util.ErrP(err)
				}()
			}
		}

	} else {
		log.Printf("Using ACME approach")
		if config.Certs == "" {
			config.Certs = "./certs"
		}
		log.Printf("ACME Found!")
		certManager := autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache(path.Join("./certs")),
		}

		server := &http.Server{
			Addr: ":443",

			TLSConfig: &tls.Config{
				PreferServerCipherSuites: true,
				// Only use curves which have assembly implementations
				CurvePreferences: []tls.CurveID{
					tls.CurveP256,
					tls.X25519, // Go 1.8 only
				},
				MinVersion: tls.VersionTLS12,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

					// Best disabled, as they don't provide Forward Secrecy,
					// but might be necessary for some clients
					// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
					// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				},
				GetCertificate:     certManager.GetCertificate,
				InsecureSkipVerify: true,
			},
		}

		go func() {
			err := http.ListenAndServe(":80", certManager.HTTPHandler(nil))
			util.ErrP(err)
		}()
		go func() {
			err := server.ListenAndServeTLS("", "")
			util.ErrP(err)
		}()
	}
	log.Printf("Server Running")
	base.LockGoRoutine()
}
