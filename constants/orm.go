package constants

import "time"

const (
	// ConnMaxIdleTime is Connection Max Idle Time for access database, you can set by time.Second, time.Minute, time.Hour.
	ConnMaxIdleTime = time.Minute
	// ConnMaxLifeTime is Connection Max Life Time for access database, you can set by time.Second, time.Minute, time.Hour.
	ConnMaxLifeTime = time.Minute
	// MaxIdleConns is mean for maximum idle connection of database
	MaxIdleConns = 10
	// MaxOpenConns is mean for maximum open connection of database
	MaxOpenConns = 10
)
