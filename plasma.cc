#include <plasma/client.h>
#include <arrow/api.h>
#include <arrow/io/api.h>
#include <arrow/ipc/api.h>
#include <iostream>
#include <vector>

std::shared_ptr<arrow::RecordBatch> build_batch() {
  arrow::Int64Builder builder;
  for (int64_t i = 0; i < 10; i++) {
    builder.Append(i);
  }
  std::shared_ptr<arrow::Array> array;
  auto status = builder.Finish(&array);
  if (!status.ok()) {
      std::cerr << "error: can't create array" << status.message() << "\n";
      return NULL;
  }
  std::vector<std::shared_ptr<arrow::Array>> arrays;
  arrays.push_back(array);

  std::shared_ptr<arrow::Field> field(new arrow::Field("i", arrow::int64()));

  std::vector<std::shared_ptr<arrow::Field>> fields;
  fields.push_back(field);
  std::shared_ptr<arrow::Schema> schema(new arrow::Schema(fields));

  return arrow::RecordBatch::Make(schema, array->length(), arrays);
}

int main(int argc, char** argv) {
  // Start up and connect a Plasma client.
  plasma::PlasmaClient client;
  auto status = client.Connect("/tmp/plasma", "");
  if (!status.ok()) {
      std::cerr << "error: can't connect" << status.message() << "\n";
      std::exit(1);
  }

  auto batch = build_batch();
  if (batch == NULL) {
    std::cerr << "error: build\n";
    std::exit(1);
  }

  int64_t size;
  status = arrow::ipc::GetRecordBatchSize(*batch, &size);
  if (!status.ok()) {
    std::cerr << "error: batch size: " << status.message() << "\n";
    std::exit(1);
  }

  std::cout << "batch size = " << size << "\n";

  plasma::ObjectID id = plasma::ObjectID::from_binary("00000000000000000007");
  std::shared_ptr<arrow::Buffer> buf;
  status = client.Create(id, size, NULL, 0, &buf);
  if (!status.ok()) {
    std::cerr << "error: create obj: " << status.message() << "\n";
    std::exit(1);
  }

  std::cout << "buf size = " << buf->size() << "\n";
  arrow::io::FixedSizeBufferWriter wb(buf);
  std::shared_ptr<arrow::ipc::RecordBatchWriter> wtr;
  status = arrow::ipc::RecordBatchStreamWriter::Open(&wb, batch->schema(), &wtr);
  if (!status.ok()) {
    std::cerr << "error: create writer: " << status.message() << "\n";
    std::exit(1);
  }

  status = wtr->WriteRecordBatch(*batch, true);
  if (!status.ok()) {
    std::cerr << "error: write: " << status.message() << "\n";
    std::exit(1);
  }
  client.Seal(id);


  std::cout << "OK\n";
}