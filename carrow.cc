#include <arrow/api.h>
#include <arrow/io/api.h>
#include <arrow/ipc/api.h>
#include <plasma/client.h>

#include <iostream>
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

/* TODO: Remove these */
void warn(arrow::Status status) {
  if (status.ok()) {
    return;
  }
  std::cout << "CARROW:WARNING: " << status.message() << "\n";
}

void debug_mark(std::string msg = "HERE") {
  std::cout << "\033[1;31m";
  std::cout << "<< " <<  msg << " >>\n";
  std::cout << "\033[0m";
  std::cout.flush();
}

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

void *fields_new() { return new std::vector<std::shared_ptr<arrow::Field>>(); }

void fields_append(void *vp, void *fp) {
  auto fields = (std::vector<std::shared_ptr<arrow::Field>> *)vp;
  std::shared_ptr<arrow::Field> field((arrow::Field *)fp);
  fields->push_back(field);
}

void fields_free(void *vp) {
  if (vp == nullptr) {
    return;
  }
  delete (std::vector<std::shared_ptr<arrow::Field>> *)vp;
}

void *schema_new(void *vp) {
  auto fields = (std::vector<std::shared_ptr<arrow::Field>> *)vp;
  auto schema = new arrow::Schema(*fields);
  return (void *)schema;
}

void schema_free(void *vp) {
  if (vp == nullptr) {
    return;
  }
  auto schema = (arrow::Schema *)vp;
  delete schema;
}

void *array_builder_new(int dtype) {
  switch (dtype) {
  case BOOL_DTYPE:
    return new arrow::BooleanBuilder();
  case FLOAT64_DTYPE:
    return new arrow::DoubleBuilder();
  case INTEGER64_DTYPE:
    return new arrow::Int64Builder();
  case STRING_DTYPE:
    return new arrow::StringBuilder();
  case TIMESTAMP_DTYPE:
    return new arrow::TimestampBuilder(data_type(TIMESTAMP_DTYPE), nullptr);
  }

  return nullptr;
}

// TODO: Check for nulls in all append
void array_builder_append_bool(void *vp, int value) {
  auto builder = (arrow::BooleanBuilder *)vp;
  builder->Append(bool(value));
}

void array_builder_append_float(void *vp, double value) {
  auto builder = (arrow::DoubleBuilder *)vp;
  builder->Append(value);
}

void array_builder_append_int(void *vp, long long value) {
  auto builder = (arrow::Int64Builder *)vp;
  builder->Append(value);
}

void array_builder_append_string(void *vp, char *cp, size_t length) {
  auto builder = (arrow::StringBuilder *)vp;
  builder->Append(cp, length);
}

void array_builder_append_timestamp(void *vp, long long value) {
  auto builder = (arrow::TimestampBuilder *)vp;
  builder->Append(value);
}


// TODO: See comment in struct Table
struct Array {
  std::shared_ptr<arrow::Array> array;
};

result_t array_builder_finish(void *vp) {
  auto builder = (arrow::ArrayBuilder *)vp;
  std::shared_ptr<arrow::Array> array;
  auto status = builder->Finish(&array);
  warn(status);
  delete builder;


  result_t res = {nullptr, nullptr};
  if (!status.ok()) {
    res.err = status.ToString().c_str();
  } else {
    auto obj = new Array;
    obj->array = array;
    res.obj = obj;
  }

  // TODO: Will out delete the underlying array?
  return res;
}

void array_free(void *vp) {
  if (vp == nullptr) {
    return;
  }

  delete (Array *)vp;
}

void *column_new(void *fp, void *ap) {
  std::shared_ptr<arrow::Field> field((arrow::Field *)fp);
  auto wrapper = (Array *)ap;

  return new arrow::Column(field, wrapper->array);
}

int column_dtype(void *vp) {
  auto column = (arrow::Column *)vp;
  return column->type()->id();
}

void column_free(void *vp) {
  if (vp == nullptr) {
    return;
  }
  auto column = (arrow::Column *)vp;
  delete column;
}

void *column_field(void *vp) {
  auto column = (arrow::Column *)vp;
  return column->field().get();
}

void *columns_new() {
  return new std::vector<std::shared_ptr<arrow::Column>>();
}

void columns_append(void *vp, void *cp) {
  auto columns = (std::vector<std::shared_ptr<arrow::Column>> *)vp;
  std::shared_ptr<arrow::Column> column((arrow::Column *)cp);
  columns->push_back(column);
}

void columns_free(void *vp) {
  auto columns = (std::vector<std::shared_ptr<arrow::Column>> *)vp;
  delete columns;
}

/* TODO: Do it with template (currently not possible under extern "C")
so we can unite with Array

e.g.
template <class T>
struct Shared<T> {
  std::shared_ptr<T> ptr;
};
*/

struct Table {
  std::shared_ptr<arrow::Table> table;
};

void *table_new(void *sp, void *cp) {
  std::shared_ptr<arrow::Schema> schema((arrow::Schema *)sp);
  auto columns = (std::vector<std::shared_ptr<arrow::Column>> *)cp;

  auto wrapper = new Table;
  wrapper->table = arrow::Table::Make(schema, *columns);
  return wrapper;
}

const char *table_validate(void *vp) {
  return nullptr;
  /*
        auto wrapper = (Table *)vp;
        // FIXME: arrow::Table::Validate is pure virtual
        auto status = wrapper->table->Validate();
        if (status.ok()) {
        return nullptr;
        }

        return status.ToString().c_str();
  */
}

long long table_num_cols(void *vp) {
  auto wrapper = (Table *)vp;
  return wrapper->table->num_columns();
}

long long table_num_rows(void *vp) {
  auto wrapper = (Table *)vp;
  return wrapper->table->num_rows();
}

void table_free(void *vp) {
  if (vp == nullptr) {
    return;
  }

  delete (Table *)vp;
}

void *plasma_connect(char *path) {
  plasma::PlasmaClient* client = new plasma::PlasmaClient();
  auto status = client->Connect(path, "", 0);
  warn(status);

  if (!status.ok()) {
    delete client;
    return nullptr; // TODO: Errors
  }

  return client;
}

bool write_table(std::shared_ptr<arrow::Table> table, std::shared_ptr<arrow::ipc::RecordBatchWriter> writer) {
  arrow::TableBatchReader rdr(*table);

  while (true) {
    std::shared_ptr<arrow::RecordBatch> batch;
    auto status = rdr.ReadNext(&batch);
    warn(status);
    if (!status.ok()) {
      return false;
    }

    if (batch == nullptr) {
      break;
    }

    status = writer->WriteRecordBatch(*batch, true);
    warn(status);
    if (!status.ok()) {
      return false;
    }
  }

  return true;
}


int64_t table_size(std::shared_ptr<arrow::Table> table) {
  arrow::TableBatchReader rdr(*table);
  std::shared_ptr<arrow::RecordBatch> batch;
  arrow::io::MockOutputStream stream;

  std::shared_ptr<arrow::ipc::RecordBatchWriter> writer;
  auto status = arrow::ipc::RecordBatchStreamWriter::Open(&stream, table->schema(), &writer);
  warn(status);
  if (!status.ok()) {
    return -1;
  }

  write_table(table, writer);

  status = writer->Close();
  warn(status);
  if (!status.ok()) {
    return -1;
  }

  return stream.GetExtentBytesWritten();
}

int plasma_write(void *cp, void *tp, char *oid) {
  // TODO: Log
  if ((cp == nullptr) || (tp == nullptr) || (oid == nullptr)) {
    return -1;
  }

  auto client = (plasma::PlasmaClient *)(cp);
  auto ptr = (Table *)(tp);
  auto table = ptr->table;

  auto size = table_size(table);

  plasma::ObjectID id = plasma::ObjectID::from_binary(oid);
  std::shared_ptr<arrow::Buffer> buf;
  // TODO: Check padding
  auto status = client->Create(id, size, nullptr, 0, &buf);
  warn(status);
  if (!status.ok()) {
    // TODO: Error
    return -1;
  }

  arrow::io::FixedSizeBufferWriter bw(buf);
  std::shared_ptr<arrow::ipc::RecordBatchWriter> writer;
  status = arrow::ipc::RecordBatchStreamWriter::Open(&bw, table->schema(), &writer);
  warn(status);
  if (!status.ok()) {
    // TODO: Error
    return -1;
  }

  if (!write_table(table, writer)) {
    // TODO: Error
    return -1;
  }

  status = client->Seal(id);
  warn(status);
  if (!status.ok()) {
    return -1;
  }

  return int(size);
}

void plasma_disconnect(void *vp) {
  if (vp == nullptr) {
    return;
  }

  auto client = (plasma::PlasmaClient*)(vp);
  auto status = client->Disconnect();
  warn(status);
  delete client;
}

// TODO: Do we want allowing multiple IDs? (like the client Get)
void *plasma_read(void *cp, char *oid, int64_t timeout_ms) {
  // TODO: Log
  if ((cp == nullptr) || (oid == nullptr)) {
    return nullptr;
  }

  auto client = (plasma::PlasmaClient *)(cp);

  plasma::ObjectID id = plasma::ObjectID::from_binary(oid);
  std::vector<plasma::ObjectID> ids;
  ids.push_back(id);
  std::vector<plasma::ObjectBuffer> buffers;

  auto status = client->Get(ids, timeout_ms, &buffers);
  warn(status);
  if (!status.ok()) {
    return nullptr;
  }

  // TODO: Support multiple buffers
  if (buffers.size() != 1) {
    std::cout << "CARROW:WARNING: more than one buffer for " << oid << "\n";
    return nullptr;
  }

  auto buf_reader = std::make_shared<arrow::io::BufferReader>(buffers[0].data);
  std::shared_ptr<arrow::ipc::RecordBatchReader> reader;
  status = arrow::ipc::RecordBatchStreamReader::Open(buf_reader, &reader);
  warn(status);
  if (!status.ok()) {
    return nullptr;
  }

  std::vector<std::shared_ptr<arrow::RecordBatch>> batches;
  while (true)
  {
    std::shared_ptr<arrow::RecordBatch> batch;
    status = reader->ReadNext(&batch);
    warn(status);
    if (!status.ok()) {
      return nullptr;
    }
    if (batch == nullptr) {
      break;
    }
    batches.push_back(batch);
  }

  std::shared_ptr<arrow::Table> table;
  status = arrow::Table::FromRecordBatches(batches, &table);
  warn(status);
  if (!status.ok()) {
    return nullptr;
  }

  auto ptr = new Table;
  ptr->table = table;
  return ptr;
}

int plasma_release(void *cp, char *oid) {
  // TODO: Log
  if ((cp == nullptr) || (oid == nullptr)) {
    return -1;
  }

  auto client = (plasma::PlasmaClient *)(cp);
  plasma::ObjectID id = plasma::ObjectID::from_binary(oid);
  auto status = client->Release(id);

  warn(status);
  if (!status.ok()) {
    return -1;
  }

  return 0;
}

#ifdef __cplusplus
} // extern "C"
#endif
