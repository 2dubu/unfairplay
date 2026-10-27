package recovery

import "math/bits"

type outerDigestRoute struct {
	rotate int8
	mask   uint32
}

type outerDigestByteMap struct {
	x27      [3]outerDigestRoute
	residual [2]outerDigestRoute
	x4       [2]outerDigestRoute
	x10      [2]outerDigestRoute
	x28      [2]outerDigestRoute
}

var outerDigestHardwareMaps = [6]outerDigestByteMap{
	{x27: [3]outerDigestRoute{{18, 0x03fc0000}}, residual: [2]outerDigestRoute{{3, 0x000007f8}, {28, 0xf000000f}}, x10: [2]outerDigestRoute{{28, 0xf000000f}}, x28: [2]outerDigestRoute{{11, 0x0007f800}}},
	{x27: [3]outerDigestRoute{{18, 0x03fc0000}}, residual: [2]outerDigestRoute{{12, 0x000ff000}}, x28: [2]outerDigestRoute{{1, 0x000001fe}, {16, 0x00ff0000}}},
	{x27: [3]outerDigestRoute{{27, 0xf8000007}, {31, 0x8000007f}}, residual: [2]outerDigestRoute{{7, 0x00003f80}, {30, 0xc000003f}}, x4: [2]outerDigestRoute{{7, 0x00007f80}}, x28: [2]outerDigestRoute{{14, 0x003fc000}}},
	{x27: [3]outerDigestRoute{{2, 0x000003fc}, {15, 0x007f8000}, {29, 0xe000001f}}, residual: [2]outerDigestRoute{{9, 0x0001be00}}, x4: [2]outerDigestRoute{{9, 0x0001fe00}}, x28: [2]outerDigestRoute{{30, 0xc000003f}}},
	{x27: [3]outerDigestRoute{{28, 0xf000000f}}, residual: [2]outerDigestRoute{{6, 0x00003dc0}, {25, 0xfe000001}}, x4: [2]outerDigestRoute{{25, 0xfe000001}}, x10: [2]outerDigestRoute{{6, 0x00003fc0}}},
	{x27: [3]outerDigestRoute{{19, 0x07f80000}, {30, 0xc000003f}}, residual: [2]outerDigestRoute{{8, 0x0000fd00}}, x10: [2]outerDigestRoute{{8, 0x0000ff00}}, x28: [2]outerDigestRoute{{2, 0x000003fc}}},
}

var outerDigestCodecMaps = [13]outerDigestByteMap{
	{x27: [3]outerDigestRoute{{21, 0x1fe00000}}, residual: [2]outerDigestRoute{{10, 0x0003fc00}}},
	{x27: [3]outerDigestRoute{{9, 0x0001fe00}}, x28: [2]outerDigestRoute{{24, 0xff000000}}},
	{x27: [3]outerDigestRoute{{11, 0x0007f800}}},
	{x27: [3]outerDigestRoute{{4, 0x00000ff0}, {16, 0x00ff0000}}, residual: [2]outerDigestRoute{{26, 0xfc000003}}, x10: [2]outerDigestRoute{{26, 0xfc000003}}, x28: [2]outerDigestRoute{{19, 0x07f80000}}},
	{},
	{x27: [3]outerDigestRoute{{4, 0x00000ff0}, {23, 0x7f800000}}, residual: [2]outerDigestRoute{{14, 0x003fc000}}, x10: [2]outerDigestRoute{{14, 0x003fc000}}, x28: [2]outerDigestRoute{{6, 0x00003fc0}}},
	{x27: [3]outerDigestRoute{{9, 0x0001fe00}}, residual: [2]outerDigestRoute{{21, 0x1fe00000}}, x4: [2]outerDigestRoute{{21, 0x1fe00000}}},
	{x27: [3]outerDigestRoute{{31, 0x8000007f}}, residual: [2]outerDigestRoute{{9, 0x0001fc00}}, x10: [2]outerDigestRoute{{9, 0x0001fe00}}},
	{x27: [3]outerDigestRoute{{2, 0x000003fc}}, residual: [2]outerDigestRoute{{14, 0x00378000}}, x4: [2]outerDigestRoute{{14, 0x003fc000}}},
	{residual: [2]outerDigestRoute{{23, 0x7f800000}}},
	{x27: [3]outerDigestRoute{{28, 0xf000000f}}, residual: [2]outerDigestRoute{{8, 0x0000bf00}}, x4: [2]outerDigestRoute{{8, 0x0000ff00}}},
	{x27: [3]outerDigestRoute{{23, 0x7f800000}}, residual: [2]outerDigestRoute{{18, 0x03fc0000}}},
	{},
}

func outerDigestRouteByte(value byte, routes []outerDigestRoute) uint32 {
	var output uint32
	for _, route := range routes {
		if route.mask != 0 {
			output ^= bits.RotateLeft32(uint32(value), int(route.rotate)) & route.mask
		}
	}
	return output
}

func applyOuterDigestStep(
	state *outerDigestState, hardware *[6]byte, codec *[33]byte,
	hardwareBase, codecPhase uint32,
) bool {
	if state == nil || hardware == nil || codec == nil || hardwareBase >= 6 || codecPhase >= 13 {
		return false
	}
	exit := state.Counter == outerDigestStop

	core := state.X0 ^
		bits.RotateLeft32(state.X4, 20) ^
		bits.RotateLeft32(state.X10, 22) ^
		bits.RotateLeft32(state.X28, 17)
	nextX27 := core ^ uint32(0x4ff57b2b)
	residual := uint32(0xe11c7ff6) ^
		(state.X4 & 0xfff7bfff) ^
		(state.X10 & 0xfffffdff)
	nextX4, nextX10, nextX28 := state.X4, state.X10, state.X28

	for index, mapping := range outerDigestHardwareMaps {
		value := hardware[(int(hardwareBase)+index)%len(hardware)]
		nextX27 ^= outerDigestRouteByte(value, mapping.x27[:])
		residual ^= outerDigestRouteByte(value, mapping.residual[:])
		nextX4 ^= outerDigestRouteByte(value, mapping.x4[:])
		nextX10 ^= outerDigestRouteByte(value, mapping.x10[:])
		nextX28 ^= outerDigestRouteByte(value, mapping.x28[:])
	}
	for index, mapping := range outerDigestCodecMaps {
		value := codec[4+(int(codecPhase)+index)%len(outerDigestCodecMaps)]
		nextX27 ^= outerDigestRouteByte(value, mapping.x27[:])
		residual ^= outerDigestRouteByte(value, mapping.residual[:])
		nextX4 ^= outerDigestRouteByte(value, mapping.x4[:])
		nextX10 ^= outerDigestRouteByte(value, mapping.x10[:])
		nextX28 ^= outerDigestRouteByte(value, mapping.x28[:])
	}

	state.X0 = state.X27 ^ (nextX27 & 0xc7a9ffef) ^ residual
	state.X4 = nextX4
	state.X10 = nextX10
	state.X27 = nextX27
	state.X28 = nextX28
	state.Counter--
	return exit
}
