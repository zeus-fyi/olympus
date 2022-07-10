package strings

const farFutureEpochToInt64MAX = int64(9223372036854775807)
const farFutureEpochToUINT64Postgres = uint64(9223372036854775807)
const farFutureEpochString = "18446744073709551615"

func FarFutureEpoch(epochString string) int64 {

	if epochString == farFutureEpochString {
		return farFutureEpochToInt64MAX
	}
	epoch := Uint64StringParser(epochString)
	if epoch >= farFutureEpochToUINT64Postgres {
		return farFutureEpochToInt64MAX
	}
	return int64(epoch)
}
