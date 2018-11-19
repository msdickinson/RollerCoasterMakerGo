package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	var game = CreateGame()
	BuildStart(game.Builder, game.Coaster)
	BuildStright(game.Builder, game.Coaster)
	BuildStright(game.Builder, game.Coaster)
	BuildStright(game.Builder, game.Coaster)

	coasterPrint(game.Coaster)

}

type TrackType int
type TaskResults int

const (
	START_X                  float64 = 500
	START_Y                  float64 = -150
	START_Z                  float64 = 0
	START_YAW                float64 = 0
	START_PITCH              float64 = 0
	START_ROLL               float64 = 0
	BUILD_AREA_SIZE_X        float64 = 1000
	BUILD_AREA_SIZE_Y        float64 = 1000
	BUILD_AREA_SIZE_Z        float64 = 2500
	FINSH_AREA_X             int     = 260
	FINSH_AREA_Y             int     = 5
	FINSH_AREA_Z             int     = 150
	FINSH_AREA_X_RANGE       int     = 150
	FINSH_AREA_Y_RANGE       int     = 150
	FINSH_AREA_Z_RANGE       int     = 200
	TRACK_LENGTH             float64 = 7.7
	TRACK_HIGHT              float64 = 1.2
	TRACK_WIDTH              float64 = 3
	TRACK_LENGTH_2X          float64 = 7.7 * 2
	TRACK_HALF_LENGH         float64 = TRACK_LENGTH / 2
	TRACK_HALF_LENGH_SQUARED float64 = TRACK_HALF_LENGH * TRACK_HALF_LENGH
	MAX_TRACKS               int     = 20000
	CART_HEIGHT              float64 = 5
	CART_SCALE               float64 = .05
	STANDARD_ANGLE_CHANGE    float64 = 7.5
	CAMERA_HEIGHT            float64 = 7.0
	GRAVITY                  float64 = 9.8
	CART_MIN_SPEED           float64 = 2.0

	TRACK_TYPE_STRIGHT TrackType = 0
	TRACK_TYPE_UP      TrackType = 1
	TRACK_TYPE_DOWN    TrackType = 2
	TRACK_TYPE_LEFT    TrackType = 3
	TRACK_TYPE_RIGHT   TrackType = 4
	TRACK_TYPE_CUSTOM  TrackType = 5

	TASK_RESULTS_SUCCESSFUL              TaskResults = 0
	TASK_RESULTS_ANGLE                   TaskResults = 1
	TASK_RESULTS_COLLISON                TaskResults = 2
	TASK_RESULTS_MAX_X                   TaskResults = 3
	TASK_RESULTS_MAX_Y                   TaskResults = 4
	TASK_RESULTS_MIN_X                   TaskResults = 5
	TASK_RESULTS_MIN_Y                   TaskResults = 6
	TASK_RESULTS_MIN_Z                   TaskResults = 7
	TASK_RESULTS_MAX_Z                   TaskResults = 8
	TASK_RESULTS_OUT_OF_BOUNDS           TaskResults = 9
	TASK_RESULTS_SUPPORT_FUNCITON_FAILED TaskResults = 10
	TASK_RESULTS_REMOVE_STARTING_TRACKS  TaskResults = 11
	TASK_RESULTS_FAIL                    TaskResults = 12
	TASK_RESULTS_REMOVE_CUSTOM_TRACK     TaskResults = 13
)

type Track struct {
	X         float64
	Y         float64
	Z         float64
	Yaw       float64
	Pitch     float64
	TrackType TrackType
}
type Coaster struct {
	TrackCount         int
	LastNewTracks      int
	LastRemovedTracks  int
	ChunkCount         int
	TrackCountBuild    int
	NewTrackCount      int
	NewChunkCount      int
	LastBuildSucessful bool
	TracksStarted      bool
	TracksFinshed      bool
	TracksInFinshArea  bool
	Chunks             [MAX_TRACKS]int
	NewChunks          [MAX_TRACKS]int
	Tracks             [MAX_TRACKS]Track
	NewTracks          [MAX_TRACKS]Track
}
type BuildAction struct {
	RemoveTrack bool
	TrackType   TrackType
	YawOffset   float64
	PitchOffset float64
}

type Builder struct {
	LastBuildActionFail bool
	Results             TaskResults
	InitialTaskResults  TaskResults
	LastRuleIssueTrack  Track
	LastCollsionTrack   Track
	BuildActions        [MAX_TRACKS]BuildAction
	BuildActionCount    int
}
type Game struct {
	Coaster *Coaster
	Builder *Builder
}

//SUPPORT
func toRadians(degrees float64) float64 {
	return 0
	//	return Math.PI * degrees / 180.0
}
func keepBetween360Degrees(degrees float64) float64 {
	if degrees < 0 {
		return degrees + 360.0
	} else if degrees >= 360.0 {
		return degrees - 360.0
	} else {
		return degrees
	}
}

//COASTER
func createCoaster() *Coaster {
	var coaster Coaster

	coaster.TrackCount = 0
	coaster.LastNewTracks = 0
	coaster.LastRemovedTracks = 0
	coaster.ChunkCount = 0
	coaster.TrackCountBuild = 0
	coaster.NewTrackCount = 0
	coaster.NewChunkCount = 0
	coaster.LastBuildSucessful = false
	coaster.TracksStarted = false
	coaster.TracksFinshed = false
	coaster.TracksInFinshArea = false

	return &coaster
}
func coasterMerge(coaster *Coaster, startTracks bool) {
	coaster.LastBuildSucessful = true
	coaster.LastRemovedTracks = coaster.TrackCount - coaster.TrackCountBuild
	if coaster.TrackCount != coaster.TrackCountBuild {
		for coaster.TrackCount > coaster.TrackCountBuild && coaster.ChunkCount > 0 {
			coaster.Chunks[coaster.ChunkCount-1] = coaster.Chunks[coaster.ChunkCount-1] - 1
			if coaster.Chunks[coaster.ChunkCount-1] == 0 {
				coaster.ChunkCount = coaster.ChunkCount - 1
			}
			coaster.TrackCount--
		}
	}
	coaster.LastNewTracks = coaster.NewTrackCount
	for i := 0; i < coaster.NewTrackCount; i++ {
		coaster.Tracks[coaster.TrackCount] = coaster.NewTracks[i]
		coaster.TrackCount++
	}
	for i := 0; i < coaster.NewChunkCount; i++ {
		coaster.Chunks[coaster.ChunkCount] = coaster.NewChunks[i]
		coaster.ChunkCount++
	}
	coaster.TrackCountBuild = coaster.TrackCount
	coaster.NewTrackCount = 0
	coaster.NewChunkCount = 0
}
func coasterReset(coaster *Coaster) {
	coaster.TrackCountBuild = coaster.TrackCount
	coaster.NewTrackCount = 0
	coaster.NewChunkCount = 0
	coaster.LastBuildSucessful = false
}
func coasterLastTrack(coaster *Coaster) Track {
	if coaster.NewTrackCount == 0 {
		return coaster.Tracks[coaster.TrackCountBuild-1]
	} else {
		return coaster.NewTracks[coaster.NewTrackCount-1]
	}
}
func coasterPrint(coaster *Coaster) {

	for i := 0; i < coaster.TrackCount; i++ {
		fmt.Print("Track " + strconv.Itoa(i) + ": ")
		fmt.Println(coaster.Tracks[i])
	}
	fmt.Println(coaster.TrackCount)
}

//BUILDER
func createBuilder() *Builder {
	var builder Builder

	return &builder
}

func BuildStright(builder *Builder, coaster *Coaster) bool {
	builder.BuildActionCount = 0
	for i := 0; i < 3; i++ {
		builder.BuildActionCount++
		builder.BuildActions[i].RemoveTrack = false
		builder.BuildActions[i].TrackType = TRACK_TYPE_STRIGHT
	}
	builder.Results = buildTracks(builder, coaster, true)
	return processBuildRequest(builder, coaster)
}
func BuildLeft(builder *Builder, coaster *Coaster) bool {
	builder.BuildActionCount = 0
	for i := 0; i < 3; i++ {
		builder.BuildActionCount++
		builder.BuildActions[i].RemoveTrack = false
		builder.BuildActions[i].TrackType = TRACK_TYPE_LEFT
	}
	builder.Results = buildTracks(builder, coaster, true)
	return processBuildRequest(builder, coaster)
}
func BuildRight(builder *Builder, coaster *Coaster) bool {
	builder.BuildActionCount = 0
	for i := 0; i < 3; i++ {
		builder.BuildActionCount++
		builder.BuildActions[i].RemoveTrack = false
		builder.BuildActions[i].TrackType = TRACK_TYPE_RIGHT
	}
	builder.Results = buildTracks(builder, coaster, true)
	return processBuildRequest(builder, coaster)
}
func BuildUp(builder *Builder, coaster *Coaster) bool {
	builder.BuildActionCount = 0
	for i := 0; i < 3; i++ {
		builder.BuildActionCount++
		builder.BuildActions[i].RemoveTrack = false
		builder.BuildActions[i].TrackType = TRACK_TYPE_UP
	}
	builder.Results = buildTracks(builder, coaster, true)
	return processBuildRequest(builder, coaster)
}
func BuildDown(builder *Builder, coaster *Coaster) bool {
	builder.BuildActionCount = 0
	for i := 0; i < 3; i++ {
		builder.BuildActionCount++
		builder.BuildActions[i].RemoveTrack = false
		builder.BuildActions[i].TrackType = TRACK_TYPE_DOWN
	}
	builder.Results = buildTracks(builder, coaster, true)
	return processBuildRequest(builder, coaster)
}
func BuildBack(builder *Builder, coaster *Coaster) bool {
	if coaster.ChunkCount == 1 {
		return false
	}

	builder.BuildActionCount = 0

	for i := 0; i < coaster.Chunks[coaster.ChunkCount-1]; i++ {
		builder.BuildActionCount++
		builder.BuildActions[i].RemoveTrack = true
	}
	builder.Results = buildTracks(builder, coaster, true)
	return processBuildRequest(builder, coaster)
}
func BuildLoop(builder *Builder, coaster *Coaster) bool {
	builder.BuildActionCount = 0

	for i := 0; i < 24; i++ {
		builder.BuildActionCount++
		builder.BuildActions[i].RemoveTrack = false
		builder.BuildActions[i].PitchOffset = STANDARD_ANGLE_CHANGE
		builder.BuildActions[i].YawOffset = .5
		builder.BuildActions[i].TrackType = TRACK_TYPE_CUSTOM
	}
	for i := 0; i < 24; i++ {
		builder.BuildActionCount++
		builder.BuildActions[i].RemoveTrack = false
		builder.BuildActions[i].PitchOffset = STANDARD_ANGLE_CHANGE
		builder.BuildActions[i].YawOffset = -.5
		builder.BuildActions[i].TrackType = TRACK_TYPE_CUSTOM
	}
	builder.Results = buildTracks(builder, coaster, true)
	return processBuildRequest(builder, coaster)
}
func BuildDownward(builder *Builder, coaster *Coaster) bool {
	builder.BuildActionCount = 0
	if coasterLastTrack(coaster).Pitch == 270.0 {
		for i := 0; i < 3; i++ {
			builder.BuildActionCount++
			builder.BuildActions[i].RemoveTrack = false
			builder.BuildActions[i].TrackType = TRACK_TYPE_STRIGHT
		}
		builder.Results = buildTracks(builder, coaster, true)

	} else {
		builder.Results = buildToPitch(builder, coaster, 270.0)
	}

	return processBuildRequest(builder, coaster)
}
func BuildUpward(builder *Builder, coaster *Coaster) bool {
	builder.BuildActionCount = 0
	if coasterLastTrack(coaster).Pitch == 90.0 {
		for i := 0; i < 3; i++ {
			builder.BuildActionCount++
			builder.BuildActions[i].RemoveTrack = false
			builder.BuildActions[i].TrackType = TRACK_TYPE_STRIGHT
		}
		builder.Results = buildTracks(builder, coaster, true)
	} else {
		builder.Results = buildToPitch(builder, coaster, 90.0)
	}

	return processBuildRequest(builder, coaster)
}
func BuildFlaten(builder *Builder, coaster *Coaster) bool {
	builder.BuildActionCount = 0
	if coasterLastTrack(coaster).Pitch == 0 {
		for i := 0; i < 3; i++ {
			builder.BuildActionCount++
			builder.BuildActions[i].RemoveTrack = false
			builder.BuildActions[i].TrackType = TRACK_TYPE_STRIGHT
		}
		builder.Results = buildTracks(builder, coaster, true)
	} else {
		builder.Results = buildToPitch(builder, coaster, 0.0)
	}
	return processBuildRequest(builder, coaster)
}
func BuildStart(builder *Builder, coaster *Coaster) bool {
	if coaster.TrackCount != 0 {
		return false
	}

	builder.BuildActionCount = 0
	for i := 0; i < 22; i++ {
		builder.BuildActions[builder.BuildActionCount].RemoveTrack = false
		builder.BuildActions[builder.BuildActionCount].TrackType = TRACK_TYPE_STRIGHT
		builder.BuildActionCount++
	}
	for i := 0; i < 12; i++ {

		builder.BuildActions[builder.BuildActionCount].RemoveTrack = false
		builder.BuildActions[builder.BuildActionCount].TrackType = TRACK_TYPE_LEFT
		builder.BuildActionCount++
	}
	for i := 0; i < 11; i++ {
		builder.BuildActions[builder.BuildActionCount].RemoveTrack = false
		builder.BuildActions[builder.BuildActionCount].TrackType = TRACK_TYPE_STRIGHT
		builder.BuildActionCount++
	}

	builder.Results = buildTracks(builder, coaster, false)
	coaster.NewChunks[coaster.NewChunkCount] = coaster.NewTrackCount
	coaster.NewChunkCount++
	coasterMerge(coaster, true)
	coaster.TracksStarted = true

	return true
}
func BuildFinsh() bool {
	return false
}
func buildToPitch(builder *Builder, coaster *Coaster, angle float64) TaskResults {
	var startAngle = coasterLastTrack(coaster).Pitch
	for i := 0; i < 15; i++ {
		builder.BuildActionCount = 0

		for j := 0; j < i; j++ {
			builder.BuildActionCount++
			builder.BuildActions[i].RemoveTrack = true
		}

		var differnce = angle - startAngle
		if differnce >= 180 {
			differnce -= 360
		} else if differnce < -180 {
			differnce += 360
		}
		var direction TrackType
		if differnce >= 0 {
			direction = TRACK_TYPE_UP
		} else {
			direction = TRACK_TYPE_DOWN
		}

		differnce = math.Abs(differnce)
		var tracks = int(math.Floor(differnce / STANDARD_ANGLE_CHANGE))

		for j := 0; j < tracks; j++ {
			builder.BuildActions[j].RemoveTrack = false
			builder.BuildActions[j].TrackType = direction
		}

		builder.Results = buildTracks(builder, coaster, false)
		if builder.Results == TASK_RESULTS_SUCCESSFUL {
			return builder.Results
		} else {
			coasterReset(coaster)
		}

	}
	return TASK_RESULTS_FAIL
}
func buildToYaw(builder *Builder, coaster *Coaster, angle float64) TaskResults {
	var startAngle = coasterLastTrack(coaster).Yaw
	for i := 0; i < 15; i++ {
		builder.BuildActionCount = 0

		for j := 0; j < i; j++ {
			builder.BuildActionCount++
			builder.BuildActions[i].RemoveTrack = true
		}

		var differnce = angle - startAngle
		if differnce > 180 {
			differnce -= 360
		} else if differnce < -180 {
			differnce += 360
		}
		var direction TrackType

		if differnce >= 0 {
			direction = TRACK_TYPE_LEFT
		} else {
			direction = TRACK_TYPE_RIGHT
		}

		differnce = math.Abs(differnce)
		var tracks = int(math.Floor(differnce / STANDARD_ANGLE_CHANGE))

		for j := 0; j < tracks; j++ {
			builder.BuildActions[j].RemoveTrack = false
			builder.BuildActions[j].TrackType = direction
		}

		var results = buildTracks(builder, coaster, false)
		if results == TASK_RESULTS_SUCCESSFUL {
			return results
		} else {
			coasterReset(coaster)
		}
	}
	return TASK_RESULTS_FAIL
}
func buildToRegion(builder *Builder, coaster *Coaster, x float64, y float64, z float64, xRange float64, yRange float64, ZRange float64) TaskResults {
	return TASK_RESULTS_FAIL
}
func maxX(x float64) bool {
	if x > BUILD_AREA_SIZE_X {
		return false
	} else {
		return true
	}
}
func maxY(y float64) bool {
	if y > BUILD_AREA_SIZE_Y {
		return false
	} else {
		return true
	}
}
func maxZ(z float64) bool {
	if z > BUILD_AREA_SIZE_Z {
		return false
	} else {
		return true
	}
}
func minX(x float64) bool {
	if x < 0 {
		return false
	} else {
		return true
	}
}
func minY(y float64) bool {
	if y < 0 {
		return false
	} else {
		return true
	}
}
func minZ(yaw float64, pitch float64, z float64) bool {
	if (pitch > 90 && pitch < 270) && (z < (0 + CART_HEIGHT*-1*math.Cos(toRadians(pitch)))) {
		return false
	} else if z < 0 {
		return false
	} else {
		return true
	}
}

func collison(builder *Builder, coaster *Coaster, x float64, y float64, z float64) bool {
	var raidus float64 = TRACK_HALF_LENGH
	var count int = 0
	var j float64 = 0
	var q float64 = 0
	var d float64 = 0
	for i := 0; i < coaster.TrackCountBuild; i++ {
		count++
		if count > coaster.TrackCountBuild-1 {
			return true
		}
		j = x - coaster.Tracks[i].X
		q = y - coaster.Tracks[i].Y
		d = z - coaster.Tracks[i].Z
		if ((j * j) + (q * q) + (d * d)) <= raidus*raidus {
			builder.LastCollsionTrack = coaster.Tracks[i]
			return false
		}
	}
	return true
}
func fixX(builder *Builder, coaster *Coaster) TaskResults {
	var results TaskResults = buildToYaw(builder, coaster, 90)
	if results == TASK_RESULTS_SUCCESSFUL {
		return results
	}

	results = buildToYaw(builder, coaster, 270)
	if results == TASK_RESULTS_SUCCESSFUL {
		return results
	}

	results = buildToYaw(builder, coaster, keepBetween360Degrees(builder.LastRuleIssueTrack.Yaw+180))
	if results == TASK_RESULTS_SUCCESSFUL {
		return results
	}

	return TASK_RESULTS_FAIL
}
func fixY(builder *Builder, coaster *Coaster) TaskResults {
	var results TaskResults = buildToYaw(builder, coaster, 180)
	if results == TASK_RESULTS_SUCCESSFUL {
		return results
	}

	results = buildToYaw(builder, coaster, 0)
	if results == TASK_RESULTS_SUCCESSFUL {
		return results
	}

	results = buildToYaw(builder, coaster, keepBetween360Degrees(builder.LastRuleIssueTrack.Yaw+180))
	if results == TASK_RESULTS_SUCCESSFUL {
		return results
	}

	return TASK_RESULTS_FAIL
}

func fixZ(builder *Builder, coaster *Coaster) TaskResults {
	var results TaskResults = buildToPitch(builder, coaster, 0)
	if results == TASK_RESULTS_SUCCESSFUL {
		return results
	}

	results = buildToPitch(builder, coaster, 180)
	if results == TASK_RESULTS_SUCCESSFUL {
		return results
	}

	return TASK_RESULTS_FAIL
}
func fixTrackCollison(builder *Builder, coaster *Coaster) TaskResults {
	var results TaskResults = buildToYaw(builder, coaster, builder.LastCollsionTrack.Yaw)
	if results == TASK_RESULTS_SUCCESSFUL {
		return results
	}

	results = buildToYaw(builder, coaster, builder.LastCollsionTrack.Yaw+180.0)
	if results == TASK_RESULTS_SUCCESSFUL {
		return results
	}

	results = buildToYaw(builder, coaster, keepBetween360Degrees(builder.LastRuleIssueTrack.Yaw+180))
	if results == TASK_RESULTS_SUCCESSFUL {
		return results
	}

	return TASK_RESULTS_FAIL
}
func processBuildRequest(builder *Builder, coaster *Coaster) bool {
	if builder.Results != TASK_RESULTS_SUCCESSFUL {
		coasterReset(coaster)
	}
	switch builder.Results {
	case TASK_RESULTS_MAX_X:
		builder.Results = fixX(builder, coaster)
		break
	case TASK_RESULTS_MAX_Y:
		builder.Results = fixY(builder, coaster)
		break
	case TASK_RESULTS_MIN_X:
		builder.Results = fixX(builder, coaster)
		break
	case TASK_RESULTS_MIN_Y:
		builder.Results = fixY(builder, coaster)
		break
	case TASK_RESULTS_MIN_Z:
		builder.Results = fixZ(builder, coaster)
		break
	case TASK_RESULTS_COLLISON:
		builder.Results = fixTrackCollison(builder, coaster)
		break
	}

	if builder.Results == TASK_RESULTS_SUCCESSFUL {
		//Chunk anything that has not already been Chunked
		var totalNewTracksChunked = 0
		for i := 0; i < coaster.NewChunkCount; i++ {
			totalNewTracksChunked += coaster.NewChunks[i]
		}

		var tracksWihtNoChunk = coaster.NewTrackCount - totalNewTracksChunked
		if tracksWihtNoChunk > 0 {
			coaster.NewChunks[coaster.NewChunkCount] = tracksWihtNoChunk
			coaster.NewChunkCount++
		}

		builder.LastBuildActionFail = false

		coasterMerge(coaster, false)
		return true
	} else {
		builder.LastBuildActionFail = true
		coasterReset(coaster)
		return false
	}
}
func buildTracks(builder *Builder, coaster *Coaster, removeChunk bool) TaskResults {
	var result TaskResults = TASK_RESULTS_FAIL
	for i := 0; i < builder.BuildActionCount; i++ {
		if builder.BuildActions[i].RemoveTrack {
			result = removeTrack(coaster, removeChunk)
		} else {
			result = buildTrack(builder, coaster, builder.BuildActions[i])
		}
		if result != TASK_RESULTS_SUCCESSFUL {
			break
		}
	}
	return result
}
func buildTrack(builder *Builder, coaster *Coaster, action BuildAction) TaskResults {
	//Check If Coater Finshed
	var yaw float64 = 0
	var pitch float64 = 0
	var x float64 = 0
	var y float64 = 0
	var z float64 = 0
	var lastTrack Track
	//Determine Starting Position
	if coaster.TrackCountBuild == 0 && coaster.NewTrackCount == 0 {
		yaw = START_YAW
		pitch = START_PITCH
		x = START_X
		y = START_Y
		z = START_Z
	} else {
		if coaster.NewTrackCount > 0 {
			lastTrack = coaster.NewTracks[coaster.NewTrackCount-1]
		} else {
			lastTrack = coasterLastTrack(coaster)
		}

		yaw = lastTrack.Yaw
		pitch = lastTrack.Pitch
		x = lastTrack.X
		y = lastTrack.Y
		z = lastTrack.Z

		//Determine Yaw And Pitch
		switch action.TrackType {
		case TRACK_TYPE_STRIGHT:
			break
		case TRACK_TYPE_LEFT:
			yaw = yaw + STANDARD_ANGLE_CHANGE
			break
		case TRACK_TYPE_RIGHT:
			yaw = yaw - STANDARD_ANGLE_CHANGE
			break
		case TRACK_TYPE_UP:
			pitch = pitch + STANDARD_ANGLE_CHANGE
			break
		case TRACK_TYPE_DOWN:
			pitch = pitch - STANDARD_ANGLE_CHANGE
			break
		case TRACK_TYPE_CUSTOM:
			yaw = yaw + action.YawOffset
			pitch = pitch + action.PitchOffset
			break
		}

		//IF X out of 360
		if yaw < 0 {
			yaw += 360
		} else if yaw >= 360 {
			yaw += -360
		}

		if pitch < 0 {
			pitch += 360
		} else if pitch >= 360 {
			pitch += -360
		}

		//Determine X, Y, And Z
		x = lastTrack.X +
			math.Cos(toRadians(lastTrack.Yaw))*(math.Cos(toRadians(lastTrack.Pitch))*TRACK_HALF_LENGH) +
			math.Cos(toRadians(yaw))*(math.Cos(toRadians(pitch))*TRACK_HALF_LENGH)
		y = lastTrack.Y +
			math.Sin(toRadians(lastTrack.Yaw))*(math.Cos(toRadians(lastTrack.Pitch))*TRACK_HALF_LENGH) +
			math.Sin(toRadians(yaw))*(math.Cos(toRadians(pitch))*TRACK_HALF_LENGH)
		z = lastTrack.Z +
			math.Sin(toRadians(lastTrack.Pitch))*TRACK_HALF_LENGH +
			math.Sin(toRadians(pitch))*TRACK_HALF_LENGH

	}

	//Check Rules
	var result TaskResults = checkRules(builder, coaster, x, y, z, yaw, pitch, action.TrackType)

	//Add Track
	if TASK_RESULTS_SUCCESSFUL == result {
		coaster.NewTracks[coaster.NewTrackCount].X = x
		coaster.NewTracks[coaster.NewTrackCount].Y = y
		coaster.NewTracks[coaster.NewTrackCount].Z = z
		coaster.NewTracks[coaster.NewTrackCount].Pitch = pitch
		coaster.NewTracks[coaster.NewTrackCount].Yaw = yaw
		coaster.NewTracks[coaster.NewTrackCount].TrackType = action.TrackType
		coaster.NewTrackCount++
	}

	return result
}
func removeTrack(coaster *Coaster, removeChunk bool) TaskResults {
	if coaster.NewTrackCount == 0 && coaster.TrackCountBuild == 45 {
		return TASK_RESULTS_OUT_OF_BOUNDS
	} else {
		if coaster.NewTrackCount > 0 {
			coaster.NewTrackCount--
		} else if removeChunk == true || coasterLastTrack(coaster).TrackType != TRACK_TYPE_CUSTOM {
			coaster.TrackCountBuild--
		} else {
			return TASK_RESULTS_REMOVE_CUSTOM_TRACK
		}
	}
	return TASK_RESULTS_SUCCESSFUL
}
func checkRules(builder *Builder, coaster *Coaster, x float64, y float64, z float64, yaw float64, pitch float64, trackType TrackType) TaskResults {
	var result TaskResults = TASK_RESULTS_SUCCESSFUL

	if coaster.TracksStarted == false || coaster.TracksInFinshArea == true {
		return result
	}

	if !maxX(x) {
		result = TASK_RESULTS_MAX_X
	} else if !maxY(y) {
		result = TASK_RESULTS_MAX_Y
	} else if !minX(x) {
		result = TASK_RESULTS_MIN_X
	} else if !minY(y) {
		result = TASK_RESULTS_MIN_Y
	} else if !minZ(yaw, pitch, z) {
		result = TASK_RESULTS_MIN_Z
	} else if !maxZ(z) {
		result = TASK_RESULTS_MAX_Z
	} else if !collison(builder, coaster, x, y, z) {
		result = TASK_RESULTS_COLLISON
	}

	if result != TASK_RESULTS_SUCCESSFUL {
		builder.LastRuleIssueTrack.X = x
		builder.LastRuleIssueTrack.Y = y
		builder.LastRuleIssueTrack.Z = z
		builder.LastRuleIssueTrack.Pitch = pitch
		builder.LastRuleIssueTrack.Yaw = yaw
		builder.LastRuleIssueTrack.TrackType = trackType
	}
	return result
}

//GAME
func CreateGame() Game {
	var game Game
	game.Builder = createBuilder()
	game.Coaster = createCoaster()
	return game
}
