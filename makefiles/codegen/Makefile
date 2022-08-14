

CORE_DIR=../../pkg/codegen/cookbook/core
DIR=${CORE_DIR}/func
FN=${DIR}/func/function_template_gen
NEW_FN=${DIR}/function

fn.gen:
	tojen gen ${FN}.go ${NEW_FN}.go

DIR=${CORE_DIR}/struct
FN=${DIR}/struct_template_gen
NEW_FN=${DIR}/struct

struct.gen:
	tojen gen ${FN}.go ${NEW_FN}.go

DIR=${CORE_DIR}/struct
FN=${DIR}/struct_template_gen
NEW_FN=${DIR}/struct

struct.gen:
	tojen gen ${FN}.go ${NEW_FN}.go


DIR=${CORE_DIR}/template_test
FN=${DIR}/template_test_suite
NEW_FN=${DIR}/_test_suite

gen.test:
	tojen gen ${FN}.go ${NEW_FN}.go