module github.com/kurakura967/tinygo-playground/namecard

go 1.25

replace github.com/kurakura967/wiodisplay => ../../wiodisplay

require (
	github.com/kurakura967/wiodisplay v0.0.0-00010101000000-000000000000
	tinygo.org/x/drivers v0.33.0
	tinygo.org/x/tinyfont v0.6.0
)

require (
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e // indirect
)
