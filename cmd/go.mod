module github.com/tuan78/jsonconv/cmd

go 1.18

replace github.com/tuan78/jsonconv => ../

require (
	github.com/spf13/cobra v1.5.0
	github.com/spf13/pflag v1.0.5
	github.com/tuan78/jsonconv v1.0.1
)

require github.com/inconshreveable/mousetrap v1.0.0 // indirect
