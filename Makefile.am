
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUFFIXES = .go .c
OBJEXT= none


sbin_PROGRAMS = dummyapi

dummyapi_SOURCES = dummyapi.go

GOFLAGS = -ldflags="-s -w"

EXTRA_dummyapi_SOURCES = \
	dummyapi.go \
	server/config/config.go \
	server/controller/hello/hellocont.go \
	server/controller/tools/gin.go \
	server/daemon/daemon.go \
	server/middleware/requestlog.go \
	server/middleware/responselog.go \
	server/server.go \
	tools/crypto.go \
	tools/crypto_test.go \
	tools/file.go 

dummyapi$(EXEEXT): $(dummyapi_SOURCES) $(EXTRA_dummyapi_SOURCES)
	$(GO) build $(GOFLAGS) -o dummyapi$(EXEEXT) $(dummyapi_SOURCES)

EXTRA_DIST = \
	README.md \
	go.mod \
	go.sum

clean-local:
	rm -f builder$(EXEEXT)
	rm -rf autom4te.cache

install-data-local:
	test -z $(DESTDIR)$(APP_LIBDIR) || $(MKDIR_P) $(DESTDIR)$(APP_LIBDIR)
#	cp -r public $(DESTDIR)$(APP_LIBDIR)
	test -z $(DESTDIR)$(APP_CONFDIR) || $(MKDIR_P) $(DESTDIR)$(APP_CONFDIR)
	test -z $(DESTDIR)$(APP_LOGDIR) || $(MKDIR_P) $(DESTDIR)$(APP_LOGDIR)
	test -z $(DESTDIR)$(APP_RUNDIR) || $(MKDIR_P) $(DESTDIR)$(APP_RUNDIR)
	test -z $(DESTDIR)$(APP_DATABASEDIR) || $(MKDIR_P) $(DESTDIR)$(APP_DATABASEDIR)
#EOF
