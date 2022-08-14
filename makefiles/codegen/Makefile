
CORE_DIR=../../pkg/codegen/cookbook/core

TYPE=func
DIR=${CORE_DIR}/${TYPE}
FN=${DIR}/${TYPE}_template_gen
NEW_FN=${DIR}/${TYPE}

fn.gen:
	TYPE=func
	tojen gen ${FN}.go ${NEW_FN}.go

struct.gen:
	TYPE=struct
	tojen gen ${FN}.go ${NEW_FN}.go

gen.test:
	TYPE=test
	tojen gen ${FN}.go ${NEW_FN}.go