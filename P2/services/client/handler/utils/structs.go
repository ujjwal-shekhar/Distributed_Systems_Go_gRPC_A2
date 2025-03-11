package utils

type ChunkMetadata struct {
	DocumentID 		string
	ChunkID 		int
	ChunkData 		[]byte
}

type ServerInfo struct {
	Address string
	Mapper 	bool
	Reducer bool
}