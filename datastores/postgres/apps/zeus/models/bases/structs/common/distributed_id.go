package common

func SelectDistributedID() string {
	return "SELECT next_id()"
}
