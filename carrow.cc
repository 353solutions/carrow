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

// TODO: Find a better way, this is for debugging ATM
#define WARN(status) \
  do { \
    if (!status.ok()) { \
      std::cout << "CARROW:WARNING: " << status.message() << "\n"; \
    } \
  } while(false);

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
  auto fields = (std::vector<std::shared_ptr<arrow::Field>> *)vp;
  delete fields;
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

result_t array_builder_finish(void *vp) {
  auto builder = (arrow::ArrayBuilder *)vp;
  std::shared_ptr<arrow::Array> out;
  auto status = builder->Finish(&out);
  WARN(status);

  result_t res = {nullptr, nullptr};
  if (!status.ok()) {
    res.err = status.ToString().c_str();
  } else {
    res.obj = (void *)(out.get());
  }

  // TODO: Will out delete the underlying array?
  return res;
}

void array_free(void *vp) {
  if (vp == nullptr) {
    return;
  }
  auto array = (arrow::Array *)vp;
  delete array;
}

void *column_new(void *fp, void *ap) {
  std::shared_ptr<arrow::Field> field((arrow::Field *)fp);
  std::shared_ptr<arrow::Array> array((arrow::Array *)ap);

  return new arrow::Column(field, array);
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

void *table_new(void *sp, void *cp) {
  std::shared_ptr<arrow::Schema> schema((arrow::Schema *)sp);
  auto columns = (std::vector<std::shared_ptr<arrow::Column>> *)cp;

  auto table = arrow::Table::Make(schema, *columns);
  return table.get();
}

const char *table_validate(void *vp) {
  return nullptr;
  /*
        auto table = (arrow::Table *)vp;
        // FIXME: arrow::Table::Validate is pure virtual
        auto status = table->Validate();
        if (status.ok()) {
        return nullptr;
        }

        return status.ToString().c_str();
  */
}

long long table_num_cols(void *vp) {
  auto table = (arrow::Table *)vp;
  return table->num_columns();
}

long long table_num_rows(void *vp) {
  auto table = (arrow::Table *)vp;
  return table->num_rows();
}

void table_free(void *vp) {
  if (vp == nullptr) {
    return;
  }
  auto table = (arrow::Table *)vp;
  delete table;
}

void *plasma_connect(char *path) {
  plasma::PlasmaClient* client = new plasma::PlasmaClient();
  auto status = client->Connect(path);
  WARN(status);

  if (!status.ok()) {
    delete client;
    return nullptr; // TODO: Errors
  }

  return client;
}

int64_t table_size(arrow::Table *table) {
  arrow::TableBatchReader rdr(*table);
  std::shared_ptr<arrow::RecordBatch> batch;
  int64_t total_size = 0;

  while (true) {
    auto status = rdr.ReadNext(&batch);
    WARN(status);
    if (!status.ok()) {
      return -1;
    }

    if (batch == nullptr) {
      break;
    }

    int64_t size;
    status = arrow::ipc::GetRecordBatchSize(*batch, &size);
    WARN(status);
    if (!status.ok()) {
      return -1;
    }
    total_size += size;
  }

  return total_size;
}

bool write_table(arrow::Table *table, std::shared_ptr<arrow::ipc::RecordBatchWriter> wtr) {
  arrow::TableBatchReader rdr(*table);
  std::shared_ptr<arrow::RecordBatch> batch;

  while (true) {
    auto status = rdr.ReadNext(&batch);
    WARN(status);
    if (!status.ok()) {
      return false;
    }

    if (batch == nullptr) {
      break;
    }

    status = wtr->WriteRecordBatch(*batch, true);
    WARN(status);
    if (!status.ok()) {
      return false;
    }
  }

  return true;
}

int plasma_write(void *cp, char *oid, void *tp) {
  auto client = (plasma::PlasmaClient *)(cp);
  auto table = (arrow::Table *)(tp);
  auto size = table_size(table);

  plasma::ObjectID id = plasma::ObjectID::from_binary(oid);
  std::shared_ptr<arrow::Buffer> buf;
  // TODO: Check padding
  auto status = client->Create(id, size + 256, nullptr, 0, &buf);
  WARN(status);
  if (!status.ok()) {
    // TODO: Error
    return -1;
  }

  arrow::io::FixedSizeBufferWriter bw(buf);
  std::shared_ptr<arrow::ipc::RecordBatchWriter> wtr;
  status = arrow::ipc::RecordBatchStreamWriter::Open(&bw, table->schema(), &wtr);
  WARN(status);
  if (!status.ok()) {
    // TODO: Error
    return -1;
  }

  if (!write_table(table, wtr)) {
    // TODO: Error
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
  WARN(status);
  delete client;
}

#ifdef __cplusplus
} // extern "C"
#endif
