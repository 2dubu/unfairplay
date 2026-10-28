package recovery

type sidbSparseTerm struct {
	target byte
	linear frameByteLinear
}

// sidbSparseAffine is an affine 128-bit layer whose byte-level dependency
// graph is sparse. It keeps the producer's actual byte transforms without the
// separate input and output encodings used by the protected implementation.
type sidbSparseAffine struct {
	constant [16]byte
	terms    [16][]sidbSparseTerm
}

func (value sidbSparseAffine) apply(input [16]byte) (output [16]byte) {
	output = value.constant
	for source, byteValue := range input {
		for _, term := range value.terms[source] {
			output[term.target] ^= term.linear.apply(byteValue)
		}
	}
	return output
}

func sidbSparseAffineToFull(value sidbSparseAffine) (result frameAffine128) {
	result.constant = value.constant
	for source, terms := range value.terms {
		for _, term := range terms {
			for bit, mapped := range term.linear {
				result.columns[source*8+bit][term.target] = mapped
			}
		}
	}
	return result
}
