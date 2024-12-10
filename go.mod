module github.com/nikolalohinski/gonja/v2

go 1.22

toolchain go1.23.3

require (
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/dustin/go-humanize v1.0.1
	github.com/hexops/gotextdiff v1.0.3
	github.com/json-iterator/go v1.1.12
	github.com/onsi/ginkgo/v2 v2.20.1
	github.com/onsi/gomega v1.35.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.3
	github.com/yargevad/filepathx v1.0.0
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56
	golang.org/x/text v0.19.0
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20240827171923-fa2c70bbbfe5 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/tools v0.24.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Critical issue https://github.com/NikolaLohinski/gonja/pull/28
retract v2.3.2
