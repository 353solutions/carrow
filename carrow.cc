#include <arrow/api.h>
#include <iostream>

#include "carrow.h"



#ifdef __cplusplus
extern "C" {
#endif

const int INTEGER_DTYPE = arrow::Type::INT64;
const int FLOAT_DTYPE = arrow::Type::DOUBLE;
int NC_INTEGER_DTYPE = INTEGER_DTYPE;
int NC_FLOAT_DTYPE = FLOAT_DTYPE;


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
	auto field = (arrow::Field *)vp;
	return field->name().c_str();
}

int field_dtype(void *vp) {
	auto field = (arrow::Field *)vp;
	return field->type()->id();
}

void field_free(void *vp) {
	auto field = (arrow::Field *)vp;
	delete field;
}

void *schema_new() {
	std::vector<std::shared_ptr<arrow::Field>> fields;
	auto schema = new arrow::Schema(fields);
	return (void *)schema;
}

void schema_add_field(void *sp, void *fp) {
	auto schema = (arrow::Schema *)sp;
	auto field = (arrow::Field *)fp;
	auto ptr = std::make_shared<arrow::Field>(*field);
	/* TODO: Does not work
	schema->fields().push_back(ptr);
	*/
}

void schema_free(void *vp) {
	auto schema = (arrow::Schema *)vp;
	delete schema;
}


#ifdef __cplusplus
} // extern "C"
#endif