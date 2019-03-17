#include <arrow/api.h>
#include <iostream>

#include "carrow.h"



#ifdef __cplusplus
extern "C" {
#endif

const int INTEGER = arrow::Type::INT64;
const int FLOAT = arrow::Type::DOUBLE;


void *field_new(char *name, int dtype) {
	std::shared_ptr<arrow::DataType> pt;

	switch (dtype) {
		case INTEGER:
			pt = arrow::int64();
			break;
		case FLOAT:
			pt = arrow::float64();
			break;

	}

	return new arrow::Field(name, pt);

	/*
	switch (dtype) {
		case INTEGER:
			return new arrow::Field(name, arrow::int64());
		case FLOAT:
			return new arrow::Field(name, arrow::float64());
	}

	return NULL;
	*/
}

const char *field_name(void *vp) {
	arrow::Field *field = (arrow::Field *)vp;
	return field->name().c_str();
}

int field_dtype(void *vp) {
	arrow::Field *field = (arrow::Field *)vp;
	std::cout << "TYPE: " << field->type() << "\n";

	return field->type()->id();
}

void field_free(void *vp) {
	arrow::Field *field = (arrow::Field *)vp;
	delete field;
}

#ifdef __cplusplus
} // extern "C"
#endif