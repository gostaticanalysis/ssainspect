module github.com/gostaticanalysis/ssainspect

go 1.23.4

require (
	github.com/gostaticanalysis/analysisutil v0.7.1
	github.com/gostaticanalysis/testutil v0.4.0
	github.com/tenntenn/golden v0.2.0
	golang.org/x/tools v0.28.0
)

require (
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/gostaticanalysis/comment v1.5.0 // indirect
	github.com/hashicorp/go-version v1.2.1 // indirect
	github.com/josharian/mapfs v0.0.0-20210615234106-095c008854e6 // indirect
	github.com/josharian/txtarfs v0.0.0-20210615234325-77aca6df5bca // indirect
	github.com/otiai10/copy v1.2.0 // indirect
	github.com/tenntenn/modver v1.0.1 // indirect
	github.com/tenntenn/text/transform v0.0.0-20200319021203-7eef512accb3 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/text v0.3.7 // indirect
)

retract (
	v0.2.0 // including a bug in Analyzer
)
