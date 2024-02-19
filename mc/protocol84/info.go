package protocol84

type Protocol84Info struct {
	SupportedVersions string
	PacketConverter   func()
}

var INFO = Protocol84Info{SupportedVersions: "0.15.10"}
