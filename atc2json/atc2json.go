package atc2json

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

var AtcFileSignature = [8]byte{'A', 'L', 'I', 'V', 'E', 0, 0, 0}

const ChecksumLength = 4

type AtcFileHeader struct {
	FileSignature [8]byte
	FileVersion   uint32
}

type BlockHeader struct {
	BlockId [4]byte
	Length  uint32
}

type FmtBlock struct {
	Format     byte
	Frequency  uint16
	Resolution uint16
	Flags      byte
	Reserved   uint16
}

// InfoBlock contains the ATC info block header
type InfoBlock struct {
	DateRecorded     [32]byte
	RecordingUUID    [40]byte
	PhoneUDID        [44]byte
	PhoneModel       [32]byte
	RecorderSoftware [32]byte
	RecorderHardware [32]byte
	Location         [52]byte
}

type EcgData struct {
	Frequency      float32    `json:"frequency"`
	MainsFrequency int        `json:"mainsFrequency"`
	Gain           float32    `json:"gain"`
	Samples        EcgSamples `json:"samples"`
	Info           *InfoBlock
}

type EcgSamples struct {
	LeadI   []int16 `json:"leadI"`
	LeadII  []int16 `json:"leadII,omitempty"`
	LeadIII []int16 `json:"leadIII,omitempty"`
	AVR     []int16 `json:"aVR,omitempty"`
	AVL     []int16 `json:"aVL,omitempty"`
	AVF     []int16 `json:"aVF,omitempty"`
}

// Parse will take atcData and return EcgData struct with error
func Parse(atcData []byte) (*EcgData, error) {

	dataLen := len(atcData)
	reader := bytes.NewReader(atcData)

	header := AtcFileHeader{}
	binary.Read(reader, binary.LittleEndian, &header)

	if header.FileSignature != AtcFileSignature {
		return nil, fmt.Errorf("Wrong file signature")
	}

	blockHeader := BlockHeader{}

	var leadISamples []int16
	var leadIISamples []int16
	var leadIIISamples []int16
	var aVRSamples []int16
	var aVLSamples []int16
	var aVFSamples []int16
	var fmtBlock *FmtBlock
	var infoBlock *InfoBlock

	for {
		blockStart := int64(dataLen - reader.Len())

		err := binary.Read(reader, binary.LittleEndian, &blockHeader)

		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("Error reading file: %s", err.Error())
		}

		blockType := string(blockHeader.BlockId[:])

		switch blockType {
		// Space after word is intended, per spec - cp 2019-2-19
		case "fmt ":
			fmtBlock = &FmtBlock{}
			err = binary.Read(reader, binary.LittleEndian, fmtBlock)
			if err != nil {
				return nil, fmt.Errorf("Error reading buffer: %s", err.Error())
			}
			err = verifyChecksum(atcData, blockStart, blockHeader.Length, reader)
			if err != nil {
				return nil, err
			}

		case "info":
			infoBlock = &InfoBlock{}
			err = binary.Read(reader, binary.LittleEndian, infoBlock)
			if err != nil {
				return nil, fmt.Errorf("Error reading buffer: %s", err.Error())
			}
			err = verifyChecksum(atcData, blockStart, blockHeader.Length, reader)
			if err != nil {
				return nil, err
			}

		// Space after word is intended, per spec - cp 2019-2-19
		case "ecg ":
			leadISamples = make([]int16, blockHeader.Length/2)
			err = binary.Read(reader, binary.LittleEndian, &leadISamples)
			if err != nil {
				return nil, fmt.Errorf("Error reading buffer: %s", err.Error())
			}

			err = verifyChecksum(atcData, blockStart, blockHeader.Length, reader)
			if err != nil {
				return nil, err
			}

		case "ecg2":
			leadIISamples = make([]int16, blockHeader.Length/2)
			err = binary.Read(reader, binary.LittleEndian, &leadIISamples)
			if err != nil {
				return nil, fmt.Errorf("Error reading buffer: %s", err.Error())
			}

			err = verifyChecksum(atcData, blockStart, blockHeader.Length, reader)
			if err != nil {
				return nil, err
			}

		case "ecg3":
			leadIIISamples = make([]int16, blockHeader.Length/2)
			err = binary.Read(reader, binary.LittleEndian, &leadIIISamples)
			if err != nil {
				return nil, fmt.Errorf("Error reading buffer: %s", err.Error())
			}

			err = verifyChecksum(atcData, blockStart, blockHeader.Length, reader)
			if err != nil {
				return nil, err
			}

		case "ecg4":
			aVRSamples = make([]int16, blockHeader.Length/2)
			err = binary.Read(reader, binary.LittleEndian, &aVRSamples)
			if err != nil {
				return nil, fmt.Errorf("Error reading buffer: %s", err.Error())
			}

			err = verifyChecksum(atcData, blockStart, blockHeader.Length, reader)
			if err != nil {
				return nil, err
			}

		case "ecg5":
			aVLSamples = make([]int16, blockHeader.Length/2)
			err = binary.Read(reader, binary.LittleEndian, &aVLSamples)
			if err != nil {
				return nil, fmt.Errorf("Error reading buffer: %s", err.Error())
			}

			err = verifyChecksum(atcData, blockStart, blockHeader.Length, reader)
			if err != nil {
				return nil, err
			}

		case "ecg6":
			aVFSamples = make([]int16, blockHeader.Length/2)
			err = binary.Read(reader, binary.LittleEndian, &aVFSamples)
			if err != nil {
				return nil, fmt.Errorf("Error reading buffer: %s", err.Error())
			}

			err = verifyChecksum(atcData, blockStart, blockHeader.Length, reader)
			if err != nil {
				return nil, err
			}
		default:
			discardBuf := make([]byte, blockHeader.Length+ChecksumLength)
			_, err = reader.Read(discardBuf)
			if err != nil {
				return nil, fmt.Errorf("Error reading input: %s", err.Error())
			}
		}
	}

	result := &EcgData{}

	result.Gain = 1e6 / float32(fmtBlock.Resolution)

	result.Frequency = float32(fmtBlock.Frequency)

	if fmtBlock.Flags&2 != 0 {
		result.MainsFrequency = 60
	} else {
		result.MainsFrequency = 50
	}

	if leadISamples != nil {
		result.Samples.LeadI = leadISamples
	}

	if leadIISamples != nil {
		result.Samples.LeadII = leadIISamples
	}

	if leadIIISamples != nil {
		result.Samples.LeadIII = leadIIISamples
	}

	if aVRSamples != nil {
		result.Samples.AVR = aVRSamples
	}

	if aVLSamples != nil {
		result.Samples.AVL = aVLSamples
	}

	if aVFSamples != nil {
		result.Samples.AVF = aVFSamples
	}

	result.Info = infoBlock

	return result, nil
}

// Convert marshals atcData to JSON string
func Convert(atcData []byte) (jsonStr string, err error) {
	ecgData, err := Parse(atcData)
	if err != nil {
		return "", err
	}

	output, err := json.Marshal(&ecgData)
	return string(output), err
}

func calcChecksum(data []byte) uint32 {
	var sum int32

	for _, b := range data {
		sum += int32(b)
	}

	return uint32(sum)
}

func verifyChecksum(data []byte, blockStart int64, blockLen uint32, reader io.Reader) (err error) {
	var checksum uint32
	binary.Read(reader, binary.LittleEndian, &checksum)

	sum := calcChecksum(data[blockStart : blockStart+8+int64(blockLen)])

	if checksum != sum {
		return fmt.Errorf("Checksum does not match. Expected: [%v] Calculated:[%v]", checksum, sum)
	}
	return nil
}

func calcMillivolts(data []int16, scale float32) []float32 {
	result := make([]float32, len(data))
	for i, sample := range data {
		result[i] = float32(sample) / scale
	}
	return result
}
