package client

import "time"

const (
	CLIENT_PREFIX = "/home/ujjwal-shake-her/distributed_systems/A2/P3/services/client/"
	RETRY_THRESHOLD = 5 * time.Second
	RETRY_FREQUENCY = 2 * time.Second
	MAX_RETRIES = 5
)