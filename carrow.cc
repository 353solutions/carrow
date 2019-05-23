#include <arrow/api.h>
#include <arrow/io/api.h>
#include <arrow/ipc/api.h>
#include <plasma/client.h>

#include <iostream>

#include "carrow.h"

#ifdef __cplusplus
extern "C" {
#endif

const int INTEGER64_DTYPE = arrow::Type::INT64;
const int FLOAT64_DTYPE = arrow::Type::DOUBLE;

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
  case INTEGER64_DTYPE:
    return arrow::int64();
  case FLOAT64_DTYPE:
    return arrow::float64();
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
  case INTEGER64_DTYPE:
    return new arrow::Int64Builder();
  case FLOAT64_DTYPE:
    return new arrow::DoubleBuilder();
  }

  return nullptr;
}

void array_builder_append_int(void *vp, long long value) {
  auto builder = (arrow::Int64Builder *)vp;
  builder->Append(value);
}

void array_builder_append_float(void *vp, double value) {
  auto builder = (arrow::DoubleBuilder *)vp;
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
  // FIXME: null checks
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

#ifdef __cplusplus
} // extern "C"
#endif
