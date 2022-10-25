package structs

func SelectDistributedID() string {
	return "SELECT next_id()"
}
