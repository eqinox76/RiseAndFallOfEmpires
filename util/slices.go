package util

func RemoveUint32(slice *[]uint32, elem uint32) []uint32{
	for i, id := range *slice{
		if id == elem{
			last := len(*slice) -1
			(*slice)[i] = (*slice)[last]
			*slice = (*slice)[:last]
			break
		}
	}
	return *slice
}

func RemoveUint64(slice *[]uint64, elem uint64) []uint64{
	for i, id := range *slice{
		if id == elem{
			last := len(*slice) -1
			(*slice)[i] = (*slice)[last]
			*slice = (*slice)[:last]
			return *slice
			break
		}
	}
	return *slice
}