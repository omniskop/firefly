package project

// Audio contains everything about the audio of a project
type Audio struct {
	Title  string     // title of the song
	Author string     // name of the interpret
	Genres []string   // genres of the song
	File   *AudioFile // the audio file
}

// AudioFile contains the (probably encoded) audio of the project
type AudioFile struct {
	Format string // the encoding format of the data
	Data   []byte // the audio data
}
