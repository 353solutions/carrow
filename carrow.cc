#include <arrow/api.h>
#include <iostream>

#include "carrow.h"

#ifdef __cplusplus
extern "C" {
#endif

const int INTEGER_DTYPE = arrow::Type::INT64;
const int FLOAT_DTYPE = arrow::Type::DOUBLE;

std::shared_ptr<arrow::DataType> data_type(int dtype) {
  switch (dtype) {
  case INTEGER_DTYPE:
    return arrow::int64();
  case FLOAT_DTYPE:
    return arrow::float64();
  }

  return NULL;
}

void *field_new(char *name, int dtype) {
  auto dt = data_type(dtype);
  return new arrow::Field(name, dt);
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

void schema_add_field(void *vp, void *fp) {
  auto field = (arrow::Field *)fp;
  auto ptr = std::make_shared<arrow::Field>(*field);
  /* TODO: Does not work
  auto schema = (arrow::Schema *)vp;
  schema->fields().push_back(ptr);
  */
}

void schema_free(void *vp) {
  if (vp == NULL) {
    return;
  }
  auto schema = (arrow::Schema *)vp;
  delete schema;
}

void *array_builder_new(int dtype) {
  switch (dtype) {
  case INTEGER_DTYPE:
    return new arrow::Int64Builder();
  case FLOAT_DTYPE:
    return new arrow::DoubleBuilder();
  }

  return NULL;
}

void array_builder_append_int(void *vp, long long value) {
  auto builder = (arrow::Int64Builder *)vp;
  builder->Append(value);
}

void array_builder_append_float(void *vp, double value) {
  auto builder = (arrow::DoubleBuilder *)vp;
  builder->Append(value);
}

finish_result array_builder_finish(void *vp) {
  auto builder = (arrow::ArrayBuilder *)vp;
  std::shared_ptr<arrow::Array> out;
  auto status = builder->Finish(&out);

  finish_result res = {NULL, NULL};
  if (!status.ok()) {
    res.err = status.ToString().c_str();
  } else {
    res.arr = (void *)(out.get());
  }

  return res;
}

#ifdef __cplusplus
} // extern "C"
#endif
