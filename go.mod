module github.com/blablatov/syslogpars

go 1.20

require (
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
)

require (
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/blablatov/syslogpars/beeper => ./beeper

replace github.com/blablatov/syslogpars/syslog2mongo => ./syslog2mongo
