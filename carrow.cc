#include <arrow/api.h>
#include <arrow/io/api.h>
#include <arrow/ipc/api.h>
#include <plasma/client.h>

#include <iostream>
#include <sstream>
#include <vector>

#include "carrow.h"

#ifdef __cplusplus
extern "C" {
#endif

const int BOOL_DTYPE = arrow::Type::BOOL;
const int FLOAT64_DTYPE = arrow::Type::DOUBLE;
const int INTEGER64_DTYPE = arrow::Type::INT64;
const int STRING_DTYPE = arrow::Type::STRING;
const int TIMESTAMP_DTYPE = arrow::Type::TIMESTAMP;

/*
static void debug_mark(std::string msg = "HERE") {
  std::cout << "\033[1;31m";
  std::cout << "<< " << msg << " >>\n";
  std::cout << "\033[0m";
  std::cout.flush();
}
*/

#define CARROW_RETURN_IF_ERROR(status)                                         \
  do {                                                                         \
    if (!status.ok()) {                                                        \
      return result_t{status.message().c_str(), nullptr};                      \
    }                                                                          \
  } while (false)

std::shared_ptr<arrow::DataType> data_type(int dtype) {
  switch (dtype) {
  case BOOL_DTYPE:
    return arrow::boolean();
  case FLOAT64_DTYPE:
    return arrow::float64();
  case INTEGER64_DTYPE:
    return arrow::int64();
  case STRING_DTYPE:
    return arrow::utf8();
  case TIMESTAMP_DTYPE:
    return arrow::timestamp(arrow::TimeUnit::NANO);
  }

  return nullptr;
}

/* TODO: Do it with template (currently not possible under extern "C")
so we can unite with Array

e.g.
template <class T>
struct Shared<T> {
  std::shared_ptr<T> ptr;
};
*/
struct Metadata {
  std::shared_ptr<arrow::KeyValueMetadata> ptr;
};

struct Schema {
  std::shared_ptr<arrow::Schema> ptr;
};

struct Array {
  std::shared_ptr<arrow::Array> ptr;
};

struct Table {
  std::shared_ptr<arrow::Table> ptr;
};

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
  if (vp == nullptr) {
    return;
  }
  auto field = (arrow::Field *)vp;
  delete field;
}

void *schema_new(void *vp, size_t count) {
  auto fields = (arrow::Field **)vp;
  auto vec = std::vector<std::shared_ptr<arrow::Field>>();
  for (size_t i = 0; i < count; i++) {
    vec.push_back(std::shared_ptr<arrow::Field>(fields[i]));
  }
  auto schema = new Schema;
  schema->ptr = std::make_shared<arrow::Schema>(vec);
  return schema;
}

result_t schema_set_meta(void *vp, void *mp) {
  result_t res = {nullptr, nullptr};
  auto schema = (Schema *)vp;
  if (schema == nullptr) {
    res.err = strdup("null schema");
    return res;
  }

  auto meta = (Metadata *)mp;
  if (meta == nullptr) {
    res.err = strdup("null meta");
    return res;
  }

  schema->ptr = schema->ptr->WithMetadata(meta->ptr);
  return res;
}

result_t schema_meta(void *vp) {
  result_t res = {nullptr, nullptr};
  auto schema = (Schema *)vp;
  if (schema == nullptr) {
    res.err = strdup("null schema");
    return res;
  }

  auto meta = new Metadata;
  meta->ptr = schema->ptr->metadata()->Copy();
  res.ptr = meta;
  return res;
}

void schema_free(void *vp) {
  if (vp == nullptr) {
    return;
  }
  auto schema = (Schema *)vp;
  delete schema;
}

result_t array_builder_new(int dtype) {
  result_t res = {nullptr, nullptr};
  switch (dtype) {
  case BOOL_DTYPE:
    res.ptr = new arrow::BooleanBuilder();
    break;
  case FLOAT64_DTYPE:
    res.ptr = new arrow::DoubleBuilder();
    break;
  case INTEGER64_DTYPE:
    res.ptr = new arrow::Int64Builder();
    break;
  case STRING_DTYPE:
    res.ptr = new arrow::StringBuilder();
    break;
  case TIMESTAMP_DTYPE:
    res.ptr = new arrow::TimestampBuilder(data_type(TIMESTAMP_DTYPE), nullptr);
    break;
  default:
    std::ostringstream oss;
    oss << "unknown dtype: " << dtype;
    res.err = oss.str().c_str();
  }

  return res;
}

// TODO: Check for nulls in all append
result_t array_builder_append_bool(void *vp, uint8_t value) {
  auto builder = (arrow::BooleanBuilder *)vp;
  auto status = builder->Append(value);
  CARROW_RETURN_IF_ERROR(status);
  return result_t{nullptr, nullptr};
}

result_t array_builder_append_bools(void *vp, uint8_t *values, int64_t length) {
  auto builder = (arrow::BooleanBuilder *)vp;
  auto status = builder->AppendValues(values, length);
  CARROW_RETURN_IF_ERROR(status);
  return result_t{nullptr, nullptr};
}

result_t array_builder_append_float(void *vp, double value) {
  auto builder = (arrow::DoubleBuilder *)vp;
  auto status = builder->Append(value);
  CARROW_RETURN_IF_ERROR(status);
  return result_t{nullptr, nullptr};
}

result_t array_builder_append_floats(void *vp, double *values, int64_t length) {
  auto builder = (arrow::DoubleBuilder *)vp;
  auto status = builder->AppendValues(values, length);
  CARROW_RETURN_IF_ERROR(status);
  return result_t{nullptr, nullptr};
}

result_t array_builder_append_int(void *vp, int64_t value) {
  auto builder = (arrow::Int64Builder *)vp;
  auto status = builder->Append(value);
  CARROW_RETURN_IF_ERROR(status);
  return result_t{nullptr, nullptr};
}

result_t array_builder_append_ints(void *vp, int64_t *values, int64_t length) {
  auto builder = (arrow::Int64Builder *)vp;
  auto status = builder->AppendValues(values, length);
  CARROW_RETURN_IF_ERROR(status);
  return result_t{nullptr, nullptr};
}

result_t array_builder_append_string(void *vp, char *cp, size_t length) {
  auto builder = (arrow::StringBuilder *)vp;
  auto status = builder->Append(cp, length);
  CARROW_RETURN_IF_ERROR(status);
  return result_t{nullptr, nullptr};
}

result_t array_builder_append_strings(void *vp, char **cp, int64_t length) {
  auto builder = (arrow::StringBuilder *)vp;
  auto status = builder->AppendValues((const char **)cp, length);
  CARROW_RETURN_IF_ERROR(status);
  return result_t{nullptr, nullptr};
}

result_t array_builder_append_timestamp(void *vp, long value) {
  auto builder = (arrow::TimestampBuilder *)vp;
  auto status = builder->Append(value);
  CARROW_RETURN_IF_ERROR(status);
  return result_t{nullptr, nullptr};
}

result_t array_builder_append_timestamps(void *vp, long *values,
                                         int64_t length) {
  auto builder = (arrow::TimestampBuilder *)vp;
  auto status = builder->AppendValues(values, length);
  CARROW_RETURN_IF_ERROR(status);
  return result_t{nullptr, nullptr};
}

result_t array_builder_finish(void *vp) {
  auto builder = (arrow::ArrayBuilder *)vp;
  std::shared_ptr<arrow::Array> array;
  auto status = builder->Finish(&array);
  CARROW_RETURN_IF_ERROR(status);
  delete builder;

  auto wrapper = new Array;
  wrapper->ptr = array;
  return result_t{nullptr, wrapper};
}

int array_dtype(void *vp) {
  if (vp == nullptr) {
    return -1;
  }

  auto wrapper = (Array *)vp;
  return wrapper->ptr->type_id();
}

int64_t array_length(void *vp) {
  if (vp == nullptr) {
    return -1;
  }

  auto wrapper = (Array *)vp;
  return wrapper->ptr->length();
}

int array_bool_at(void *vp, long long i) {
  auto wrapper = (Array *)vp;
  if (wrapper == nullptr) {
    return -1;
  }

  if (wrapper->ptr->type_id() != BOOL_DTYPE) {
    return -1;
  }

  auto arr = (arrow::BooleanArray *)(wrapper->ptr.get());
  return arr->Value(i) ? 1 : 0;
}

double array_float_at(void *vp, long long i) {
  auto wrapper = (Array *)vp;
  if (wrapper == nullptr) {
    return -1;
  }

  if (wrapper->ptr->type_id() != FLOAT64_DTYPE) {
    return -1;
  }

  auto arr = (arrow::DoubleArray *)(wrapper->ptr.get());
  return arr->Value(i);
}

int64_t array_int_at(void *vp, long long i) {
  auto wrapper = (Array *)vp;
  if (wrapper == nullptr) {
    return -1;
  }

  if (wrapper->ptr->type_id() != INTEGER64_DTYPE) {
    return -1;
  }

  auto arr = (arrow::Int64Array *)(wrapper->ptr.get());
  return arr->Value(i);
}

const char *array_str_at(void *vp, long long i) {
  auto wrapper = (Array *)vp;
  if (wrapper == nullptr) {
    return nullptr;
  }

  if (wrapper->ptr->type_id() != STRING_DTYPE) {
    return nullptr;
  }

  auto arr = (arrow::StringArray *)(wrapper->ptr.get());
  auto str = arr->GetString(i);
  return strdup(str.c_str());
}

int64_t array_timestamp_at(void *vp, long long i) {
  auto wrapper = (Array *)vp;
  if (wrapper == nullptr) {
    return -1;
  }

  if (wrapper->ptr->type_id() != TIMESTAMP_DTYPE) {
    return -1;
  }

  auto arr = (arrow::TimestampArray *)(wrapper->ptr.get());
  return arr->Value(i);
}

void array_free(void *vp) {
  if (vp == nullptr) {
    return;
  }

  delete (Array *)vp;
}


void *table_new(void *sp, void *ap, size_t ncols) {
  auto schema = (Schema *)sp;
  auto arrays = (Array**)ap;

  auto vec = std::vector<std::shared_ptr<arrow::Array>>();
  for (size_t i = 0; i < ncols; i++) {
    vec.push_back(arrays[i]->ptr);
  }

  auto table = arrow::Table::Make(schema->ptr, vec);
  if (table == nullptr) {
    return nullptr;
  }

  auto wrapper = new Table;
  wrapper->ptr = table;
  return wrapper;
}

long long table_num_cols(void *vp) {
  auto wrapper = (Table *)vp;
  return wrapper->ptr->num_columns();
}

long long table_num_rows(void *vp) {
  auto wrapper = (Table *)vp;
  return wrapper->ptr->num_rows();
}

void *table_schema(void *vp) {
  auto wrapper = (Table *)vp;
  auto ptr = wrapper->ptr->schema();
  if (ptr == NULL) {
    return NULL;
  }

  auto schema = new Schema;
  schema->ptr = ptr;
  return schema;
}

void *table_column(void *vp, int i) {
  auto wrapper = (Table *)vp;
  auto arr = wrapper->ptr->column(i);
  if (arr == NULL) {
    return NULL;
  }

  auto array = new Array;
  array->ptr = arr->chunk(0); // FIXME: Rethink Array
  return array;
}

void *table_field(void *vp, int i) {
  auto wrapper = (Table *)vp;
  auto field = wrapper->ptr->field(i);
  if (field == NULL) {
    return NULL;
  }

  return field.get();
}

void table_free(void *vp) {
  if (vp == nullptr) {
    return;
  }

  delete (Table *)vp;
}

void *table_slice(void *vp, int64_t offset, int64_t length) {
  auto wrapper = (Table *)vp;
  auto ptr = wrapper->ptr->Slice(offset, length);

  auto table = new Table;
  table->ptr = ptr;
  return table;
}

void *meta_new() {
  auto meta = new Metadata;
  meta->ptr = std::make_shared<arrow::KeyValueMetadata>();

  return meta;
}

result_t meta_set(void *vp, const char *key, const char *value) {
  result_t res = {nullptr, nullptr};
  auto meta = (Metadata *)vp;
  if (meta == nullptr) {
    res.err = strdup("null pointer");
    return res;
  }

  meta->ptr->Append(key, value);
  return res;
}

result_t meta_size(void *vp) {
  result_t res = {nullptr, nullptr};
  auto meta = (Metadata *)vp;
  if (meta == nullptr) {
    res.err = strdup("null pointer");
    return res;
  }
  res.i = meta->ptr->size();
  return res;
}

result_t meta_key(void *vp, int64_t i) {
  result_t res = {nullptr, nullptr};
  auto meta = (Metadata *)vp;
  if (meta == nullptr) {
    res.err = strdup("null pointer");
    return res;
  }

  res.ptr = strdup(meta->ptr->key(i).c_str());
  return res;
}

result_t meta_value(void *vp, int64_t i) {
  result_t res = {nullptr, nullptr};
  auto meta = (Metadata *)vp;
  if (meta == nullptr) {
    res.err = strdup("null pointer");
    return res;
  }

  res.ptr = strdup(meta->ptr->value(i).c_str());
  return res;
}

result_t plasma_connect(char *path) {
  plasma::PlasmaClient *client = new plasma::PlasmaClient();
  auto status = client->Connect(path, "", 0);
  if (!status.ok()) {
    client->Disconnect();
    delete client;
  }

  CARROW_RETURN_IF_ERROR(status);
  return result_t{nullptr, client};
}

arrow::Status
write_table(std::shared_ptr<arrow::Table> table,
            std::shared_ptr<arrow::ipc::RecordBatchWriter> writer) {
  arrow::TableBatchReader rdr(*table);

  while (true) {
    std::shared_ptr<arrow::RecordBatch> batch;
    auto status = rdr.ReadNext(&batch);
    if (!status.ok()) {
      return status;
    }

    if (batch == nullptr) {
      break;
    }

    status = writer->WriteRecordBatch(*batch);
    if (!status.ok()) {
      return status;
    }
  }

  return arrow::Status::OK();
}

result_t table_size(std::shared_ptr<arrow::Table> table) {
  arrow::TableBatchReader rdr(*table);
  std::shared_ptr<arrow::RecordBatch> batch;
  arrow::io::MockOutputStream stream;

  std::shared_ptr<arrow::ipc::RecordBatchWriter> writer;
  auto status = arrow::ipc::RecordBatchStreamWriter::Open(
      &stream, table->schema(), &writer);
  CARROW_RETURN_IF_ERROR(status);
  status = write_table(table, writer);
  CARROW_RETURN_IF_ERROR(status);
  status = writer->Close();
  CARROW_RETURN_IF_ERROR(status);

  auto num_written = stream.GetExtentBytesWritten();
  return result_t{nullptr, (void *)num_written};
}

result_t plasma_write(void *cp, void *tp, char *oid) {
  if ((cp == nullptr) || (tp == nullptr) || (oid == nullptr)) {
    return result_t{"null pointer", nullptr};
  }

  auto client = (plasma::PlasmaClient *)(cp);
  auto wrapper = (Table *)(tp);
  auto table = wrapper->ptr;

  auto res = table_size(table);
  if (res.err != nullptr) {
    return res;
  }

  auto size = int64_t(res.ptr);

  plasma::ObjectID id = plasma::ObjectID::from_binary(oid);
  std::shared_ptr<arrow::Buffer> buf;
  // TODO: Check padding
  auto status = client->Create(id, size, nullptr, 0, &buf);
  CARROW_RETURN_IF_ERROR(status);

  arrow::io::FixedSizeBufferWriter bw(buf);
  std::shared_ptr<arrow::ipc::RecordBatchWriter> writer;
  status =
      arrow::ipc::RecordBatchStreamWriter::Open(&bw, table->schema(), &writer);
  CARROW_RETURN_IF_ERROR(status);

  status = write_table(table, writer);
  CARROW_RETURN_IF_ERROR(status);
  status = client->Seal(id);
  CARROW_RETURN_IF_ERROR(status);

  return result_t{nullptr, (void *)size};
}

result_t plasma_disconnect(void *vp) {
  if (vp == nullptr) {
    return result_t{nullptr, nullptr};
  }

  auto client = (plasma::PlasmaClient *)(vp);
  auto status = client->Disconnect();
  CARROW_RETURN_IF_ERROR(status);
  delete client;
  return result_t{nullptr, nullptr};
}

// TODO: Do we want allowing multiple IDs? (like the client Get)
result_t plasma_read(void *cp, char *oid, int64_t timeout_ms) {
  if ((cp == nullptr) || (oid == nullptr)) {
    return result_t{"null pointer", nullptr};
  }

  auto client = (plasma::PlasmaClient *)(cp);

  plasma::ObjectID id = plasma::ObjectID::from_binary(oid);
  std::vector<plasma::ObjectID> ids;
  ids.push_back(id);
  std::vector<plasma::ObjectBuffer> buffers;

  auto status = client->Get(ids, timeout_ms, &buffers);
  CARROW_RETURN_IF_ERROR(status);

  // TODO: Support multiple buffers
  if (buffers.size() != 1) {
    std::ostringstream oss;
    oss << "more than one buffer for " << oid;
    return result_t{oss.str().c_str(), nullptr};
  }

  auto buf_reader = std::make_shared<arrow::io::BufferReader>(buffers[0].data);
  std::shared_ptr<arrow::ipc::RecordBatchReader> reader;
  status = arrow::ipc::RecordBatchStreamReader::Open(buf_reader, &reader);
  CARROW_RETURN_IF_ERROR(status);

  std::vector<std::shared_ptr<arrow::RecordBatch>> batches;
  while (true) {
    std::shared_ptr<arrow::RecordBatch> batch;
    status = reader->ReadNext(&batch);
    CARROW_RETURN_IF_ERROR(status);
    if (batch == nullptr) {
      break;
    }
    batches.push_back(batch);
  }

  std::shared_ptr<arrow::Table> table;
  status = arrow::Table::FromRecordBatches(batches, &table);
  CARROW_RETURN_IF_ERROR(status);

  auto wrapper = new Table;
  wrapper->ptr = table;
  return result_t{nullptr, wrapper};
}

result_t plasma_release(void *cp, char *oid) {
  if ((cp == nullptr) || (oid == nullptr)) {
    return result_t{"null pointer", nullptr};
  }

  auto client = (plasma::PlasmaClient *)(cp);
  plasma::ObjectID id = plasma::ObjectID::from_binary(oid);
  auto status = client->Release(id);
  CARROW_RETURN_IF_ERROR(status);
  return result_t{nullptr, nullptr};
}

#ifdef __cplusplus
} // extern "C"
#endif
