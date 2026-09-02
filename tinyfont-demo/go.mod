module github.com/kurakura967/tinygo-playground/tinyfont-demo

go 1.25

replace github.com/kurakura967/wiodisplay => ../../wiodisplay

require (
	github.com/kurakura967/wiodisplay v0.0.0-00010101000000-000000000000
	github.com/sago35/koebiten v0.6.0
	tinygo.org/x/drivers v0.35.0
	tinygo.org/x/tinyfont v0.6.0
)

require (
	github.com/chewxy/math32 v1.11.1 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	tinygo.org/x/tinydraw v0.4.0 // indirect
)
