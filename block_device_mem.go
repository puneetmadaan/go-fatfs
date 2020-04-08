package fatfs

import "fmt"

// MemBlockDevice is a block device implementation backed by a byte slice
type MemBlockDevice struct {
	memory     []byte
	blankBlock []byte
	blockSize  uint32
	blockCount uint32
}

var _ BlockDevice = (*MemBlockDevice)(nil)

func NewMemoryDevice(blockSize int, blockCount int) *MemBlockDevice {
	dev := &MemBlockDevice{
		memory:     make([]byte, blockSize*blockCount),
		blankBlock: make([]byte, blockSize),
		blockSize:  uint32(blockSize),
		blockCount: uint32(blockCount),
	}
	for i := range dev.blankBlock {
		dev.blankBlock[i] = 0xff
	}
	for i := uint32(0); i < uint32(blockCount); i++ {
		if err := dev.eraseBlock(i); err != nil {
			panic(fmt.Sprintf("could not initialize block %d: %s", i, err.Error()))
		}
	}
	return dev
}

func (bd *MemBlockDevice) ReadAt(buf []byte, off int64) (n int, err error) {
	return copy(buf, bd.memory[off:]), nil
}

func (bd *MemBlockDevice) WriteAt(buf []byte, off int64) (n int, err error) {
	return copy(bd.memory[off:], buf), nil
}

func (bd *MemBlockDevice) Size() int64 {
	return int64(bd.blockSize * bd.blockCount)
}

func (bd *MemBlockDevice) SectorSize() int64 {
	return SectorSize
}

func (bd *MemBlockDevice) EraseBlockSize() int64 {
	return int64(bd.blockSize)
}

func (bd *MemBlockDevice) EraseBlocks(start int64, len int64) error {
	for i := int64(0); i < len; i++ {
		if err := bd.eraseBlock(uint32(start + i)); err != nil {
			return err
		}
	}
	return nil
}

func (bd *MemBlockDevice) eraseBlock(block uint32) error {
	copy(bd.memory[bd.blockSize*block:], bd.blankBlock)
	return nil
}

func (bd *MemBlockDevice) Sync() error {
	return nil
}
