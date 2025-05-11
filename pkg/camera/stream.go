package camera

// Streamer is an interface for a camera streamer.
type Streamer interface {
	Config() *Config
	Stop()
}
