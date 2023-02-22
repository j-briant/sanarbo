package main

import (
	"embed"
	"fmt"
	"github.com/cristalhq/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/goserver"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/tools"
	"github.com/lao-tseu-is-alive/sanarbo/pkg/trees"
	"github.com/lao-tseu-is-alive/sanarbo/pkg/version"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

const (
	defaultPort            = 9999
	defaultDBPort          = 5432
	defaultDBIp            = "127.0.0.1"
	defaultDBSslMode       = "prefer"
	defaultWebRootDir      = "sanarboFront/dist/"
	defaultSecuredApi      = "/goapi/v1"
	defaultUsername        = "bill"
	defaultFakeStupidPass  = "board"
	charsetUTF8            = "charset=UTF-8"
	MIMEAppJSON            = "application/json"
	MIMEHtml               = "text/html"
	MIMEAppJSONCharsetUTF8 = MIMEAppJSON + "; " + charsetUTF8
	MIMEHtmlCharsetUTF8    = MIMEHtml + "; " + charsetUTF8
	HeaderContentType      = "Content-Type"
)

// content holds our static web server content.
//
//go:embed sanarboFront/dist/*
var content embed.FS

var dbConn database.DB

type ServiceExample struct {
	Log *log.Logger
	//Store       Storage
	dbConn      database.DB
	JwtSecret   []byte
	JwtDuration int
}

// login is just a trivial stupid example to test this server
// you should use the jwt token returned from LoginUser  in github.com/lao-tseu-is-alive/go-cloud-k8s-user-group'
// and share the same secret with the above component
func (s ServiceExample) login(ctx echo.Context) error {
	s.Log.Printf("trace: entering %v login()", ctx.Request().Method)
	username := ctx.FormValue("login")
	fakePassword := ctx.FormValue("pass")

	// Throws unauthorized error
	if username != defaultUsername || fakePassword != defaultFakeStupidPass {
		return ctx.JSON(http.StatusUnauthorized, "username not found or password invalid")
	}

	// Set custom claims
	claims := &goserver.JwtCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "",
			Audience:  nil,
			Issuer:    "",
			Subject:   "",
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(time.Minute * time.Duration(s.JwtDuration))},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			NotBefore: nil,
		},
		Id:       999,
		Name:     "Bill Whatever",
		Email:    "bill@whatever.com",
		Username: defaultUsername,
		IsAdmin:  false,
	}

	// Create token with claims
	signer, _ := jwt.NewSignerHS(jwt.HS512, s.JwtSecret)
	builder := jwt.NewBuilder(signer)
	token, err := builder.Build(claims)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("LoginUser(%s) succesfull login for user id (%d)", claims.Username, claims.Id)
	s.Log.Printf(msg)
	return ctx.JSON(http.StatusOK, echo.Map{
		"token": token.String(),
	})
}

func (s ServiceExample) restricted(ctx echo.Context) error {
	s.Log.Println("trace: entering restricted zone()")
	// get the current user from JWT TOKEN
	u := ctx.Get("jwtdata").(*jwt.Token)
	claims := goserver.JwtCustomClaims{}
	err := u.DecodeClaims(&claims)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	//callerUserId := claims.Id
	// you can check if the user is not active anymore and RETURN 401 Unauthorized
	//if !s.Store.IsUserActive(currentUserId) {
	//	return echo.NewHTTPError(http.StatusUnauthorized, "current calling user is not active anymore")
	//}
	return ctx.JSON(http.StatusCreated, claims)
}

func isDBAlive() bool {
	dbVer, err := dbConn.GetVersion()
	if err != nil {
		return false
	}
	if len(dbVer) < 2 {
		return false
	}
	return true
}

func checkReady(string) bool {
	// we decide what makes us ready, is a valid  connection to the database
	if !isDBAlive() {
		return false
	}
	return true
}

func checkHealthy(string) bool {
	// you decide what makes you ready, may be it is the connection to the database
	//if !isDBAlive() {
	//	return false
	//}
	return true
}

func main() {
	l := log.New(os.Stdout, fmt.Sprintf("%s ", version.APP), log.Ldate|log.Ltime|log.Lshortfile)
	l.Printf("INFO: 'Starting %s v:%s  rev:%s  build: %s'", version.APP, version.VERSION, version.REVISION, version.BuildStamp)
	l.Printf("INFO: 'Repository url: https://%s'", version.REPOSITORY)
	secret, err := config.GetJwtSecretFromEnv()
	if err != nil {
		l.Fatalf("💥💥 ERROR: 'in NewGoHttpServer config.GetJwtSecretFromEnv() got error: %v'\n", err)
	}
	tokenDuration, err := config.GetJwtDurationFromEnv(60)
	if err != nil {
		l.Fatalf("💥💥 ERROR: 'in NewGoHttpServer config.GetJwtDurationFromEnv() got error: %v'\n", err)
	}
	dbDsn, err := config.GetPgDbDsnUrlFromEnv(defaultDBIp, defaultDBPort,
		tools.ToSnakeCase(version.APP), version.AppSnake, defaultDBSslMode)
	if err != nil {
		l.Fatalf("💥💥 error doing config.GetPgDbDsnUrlFromEnv. error: %v\n", err)
	}
	dbConn, err = database.GetInstance("pgx", dbDsn, runtime.NumCPU(), l)
	if err != nil {
		l.Fatalf("💥💥 error doing users.GetPgxConn(postgres, dbDsn  : %v\n", err)
	}
	defer dbConn.Close()

	yourService := ServiceExample{
		Log:         l,
		dbConn:      dbConn,
		JwtSecret:   []byte(secret),
		JwtDuration: tokenDuration,
	}

	listenAddr, err := config.GetPortFromEnv(defaultPort)
	if err != nil {
		l.Fatalf("💥💥 ERROR: 'calling GetPortFromEnv got error: %v'\n", err)
	}
	l.Printf("INFO: 'Will start HTTP server listening on port %s'", listenAddr)
	server := goserver.NewGoHttpServer(listenAddr, l, defaultWebRootDir, content, defaultSecuredApi)
	e := server.GetEcho()
	e.GET("/readiness", server.GetReadinessHandler(checkReady, "Connection to DB"))
	e.GET("/health", server.GetHealthHandler(checkHealthy, "Connection to DB"))
	// Login route
	e.POST("/login", yourService.login)
	r := server.GetRestrictedGroup()
	r.GET("/secret", yourService.restricted)

	objStore, err := trees.GetStorageInstance("pgx", dbConn, l)
	if err != nil {
		l.Fatalf("💥💥 ERROR: 'calling objects.GetStorageInstance got error: %v'\n", err)
	}
	// now with restricted group reference you can register your secured handlers defined in OpenApi objects.yaml
	objService := trees.Service{
		Log:         l,
		Store:       objStore,
		JwtSecret:   []byte(secret),
		JwtDuration: tokenDuration,
	}
	trees.RegisterHandlers(r, &objService)

	loginExample := fmt.Sprintf("curl -v -X POST -d 'login=%s' -d 'pass=%s' http://localhost%s/login", defaultUsername, defaultFakeStupidPass, listenAddr)
	getSecretExample := fmt.Sprintf(" curl -v  -H \"Authorization: Bearer ${TOKEN}\" http://localhost%s%s/secret |jq\n", listenAddr, defaultSecuredApi)
	l.Printf("INFO: from another terminal just try :\n %s", loginExample)
	l.Printf("INFO: then type export TOKEN=your_token_above_goes_here   \n %s", getSecretExample)

	err = server.StartServer()
	if err != nil {
		l.Fatalf("💥💥 ERROR: 'calling echo.Start(%s) got error: %v'\n", listenAddr, err)
	}

}
