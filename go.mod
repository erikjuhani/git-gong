module github.com/erikjuhani/git-gong

go 1.13

require (
	github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6 // indirect
	github.com/libgit2/git2go/v31 v31.4.7
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.10.1 // indirect
	github.com/ugorji/go v1.1.4 // indirect
	github.com/xordataexchange/crypt v0.0.3-0.20170626215501-b2862e3d0a77 // indirect
)

replace github.com/libgit2/git2go/v31 => ./vendor/git2go
