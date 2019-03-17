#include <arrow/api.h>
#include <iostream>

#include "carrow.h"



#ifdef __cplusplus
extern "C" {
#endif

const int INTEGER_DTYPE = arrow::Type::INT64;
const int FLOAT_DTYPE = arrow::Type::DOUBLE;


void *field_new(char *name, int dtype) {
	std::shared_ptr<arrow::DataType> pt;

	switch (dtype) {
		case INTEGER_DTYPE:
			pt = arrow::int64();
			break;
		case FLOAT_DTYPE:
			pt = arrow::float64();
			break;
		default:
			return NULL;

	}

	return new arrow::Field(name, pt);
}

const char *field_name(void *vp) {
	arrow::Field *field = (arrow::Field *)vp;
	return field->name().c_str();
}

int field_dtype(void *vp) {
	arrow::Field *field = (arrow::Field *)vp;
	return field->type()->id();
}

void field_free(void *vp) {
	arrow::Field *field = (arrow::Field *)vp;
	delete field;
}

#ifdef __cplusplus
} // extern "C"
#endif