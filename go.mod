module example.com/hello-go

go 1.16

require (
	example.com/greeting v0.0.0-00010101000000-000000000000
	golang.org/x/text v0.3.6 // indirect
	rsc.io/quote v1.5.2
	rsc.io/sampler v1.99.99 // indirect
)

replace example.com/greeting => ../greeting
