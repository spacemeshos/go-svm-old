PROJ_DIR := $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
BIN_DIR ?= $(PROJ_DIR)build/
export GOOS
export GOARCH
export GOARM
export BIN_DIR

ifeq ($(OS),Windows_NT)
  HOST_OS := windows
else
  HOST_OS := $(shell uname | tr "[A-Z]" "[a-z]")
endif

ifeq ($(GOOS),)
  GOOS := $(HOST_OS)
endif

GOTAGS = # -tags ...
export CGO_LDFLAGS := -L$(BIN_DIR)

CGO_TEST_LDFLAGS := $(CGO_LDFLAGS) -Wl,-rpath,$(BIN_DIR)

ifeq ($(GOOS),windows)
  platform := windows
  export PATH := $(PATH):$(PROJ_DIR)build
  EXE := .exe
else
  TEMP := /tmp
  ifeq ($(GOOS),darwin)
    platform := macos
    CGO_LDFLAGS := $(CGO_LDFLAGS) -Wl,-rpath,@loader_path
    ULIMIT := ulimit -n 9999;
  else
    platform := linux
    CGO_LDFLAGS := $(CGO_LDFLAGS) -Wl,-rpath,$$ORIGIN
  endif
endif

ifneq ($(VERBOSE),)
  $(info "OS: $(OS), HOST_OS: $(HOST_OS), GOOS: $(GOOS), GOARCH: $(GOARH), BIN_DIR: $(BIN_DIR), platform: $(platform)")
endif

ifeq ($(platform), windows)
  SVM_LIBS = svm.dll libsvm.dll.a
else
  ifeq ($(platform), macos)
  	SVM_LIBS = libsvm.dylib
  else
  	SVM_LIBS = libsvm.so
  endif
endif

BINDIR_SVM_LIBS = $(foreach X,$(SVM_LIBS),$(BIN_DIR)$(X))

SUBDIRS_LVL1 := $(foreach X,$(wildcard $(dir $(PROJ_DIR))*/.), $(lastword $(subst /, ,$(dir $(X)))))
SUBDIRS_ONLY := $(sort \
	$(foreach X,$(SUBDIRS_LVL1),\
		$(foreach Y, $(wildcard $(X)/*_test.go), test-$(X)) ))
SUBDIRS_ALL := $(sort \
	$(foreach X,$(SUBDIRS_LVL1),\
		$(foreach Y, $(wildcard $(X)/*/*_test.go), test-$(X)-all ) \
		$(foreach Y, $(wildcard $(X)/*_test.go), test-$(X)-all ) ))
SUBDIRS_LVL2 := $(sort \
	$(foreach X,$(SUBDIRS_LVL1),\
		$(foreach Y, $(wildcard $(X)/*/*_test.go), test-$(X)-$(lastword $(subst /, ,$(dir $(Y))) ))))

# defines test-* targets from subdirs level 1 and 2
$(SUBDIRS_ONLY): $(BINDIR_SVM_LIBS)
	$(ULIMIT) CGO_LDFLAGS="$(CGO_TEST_LDFLAGS)" go test -v ./$(@:test-%=%)
$(SUBDIRS_LVL2): $(BINDIR_SVM_LIBS)
	$(ULIMIT) CGO_LDFLAGS="$(CGO_TEST_LDFLAGS)" go test -v ./$(subst -,/,$(@:test-%=%))/...
$(SUBDIRS_ALL): $(BINDIR_SVM_LIBS)
	$(ULIMIT) CGO_LDFLAGS="$(CGO_TEST_LDFLAGS)" go test -v ./$(@:test-%-all=%)/...

test-all: $(BINDIR_SVM_LIBS)
	$(ULIMIT) CGO_LDFLAGS="$(CGO_TEST_LDFLAGS)" go test -v ./...

print-test-targets:
	@for i in test-all test-tidy test-fmt \
		$(sort \
			$(SUBDIRS_LVL2) \
			$(foreach X, $(SUBDIRS_ALL), $(if $(filter $(SUBDIRS_ONLY),$(X:%-all=%)), \
				$(X:%-all=%)[-all],\
				$(X)))); \
	do echo $$i; \
	done
.PHONY: $(SUBDIRS_ALL) $(SUBDIRS_LVL2) $(SUBDIRS_ONLY) print-test-targets

$(info $(word 1,$@))

go-env: $(BINDIR_SVM_LIBS)
	go env -w CGO_LDFLAGS="$(CGO_LDFLAGS)"
.PHONY: go-env

go-env-test: $(BINDIR_SVM_LIBS)
	go env -w CGO_LDFLAGS="$(CGO_TEST_LDFLAGS)"
.PHONY: go-env-test

print-ldflags: $(BINDIR_SVM_LIBS)
	@echo '$(CGO_LDFLAGS)'
.PHONY: print-ldflags

print-test-ldflags: $(BINDIR_SVM_LIBS)
	@echo '$(CGO_TEST_LDFLAGS)'
.PHONY: print-test-ldflags