package vm

func exactIsBranch(flags uint8) (uint8, bool) {
	return flags & 0b111, flags&0b1000 == 1
}
