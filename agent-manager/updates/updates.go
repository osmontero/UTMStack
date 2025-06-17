package updates

import (
	"crypto/tls"
	"github.com/gin-contrib/gzip"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/utmstack/UTMStack/agent-manager/config"
	"github.com/utmstack/UTMStack/agent-manager/utils"
)

func InitUpdatesManager() {
	ServeDependencies()
}

func ServeDependencies() {
	utils.ALogger.LogF(100, "Serving dependencies from %s", config.UpdatesDependenciesFolder)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(
		gin.Recovery(),
		gzip.Gzip(gzip.DefaultCompression),
	)

	r.NoRoute(notFound)

	group := r.Group("/private", HTTPAuthInterceptor())
	group.StaticFS("/dependencies", http.Dir(config.UpdatesDependenciesFolder))

	loadedCert, err := tls.LoadX509KeyPair(config.CertPath, config.CertKeyPath)
	if err != nil {
		utils.ALogger.Fatal("failed to load TLS credentials: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{loadedCert},
		MinVersion:   tls.VersionTLS13,
	}

	server := &http.Server{
		Addr:      ":8080",
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	utils.ALogger.Info("Starting HTTP server on port 8080")
	if err := server.ListenAndServeTLS("", ""); err != nil {
		utils.ALogger.ErrorF("error starting HTTP server: %v", err)
		return
	}
}

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
}
