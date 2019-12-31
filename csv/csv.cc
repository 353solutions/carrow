#include <arrow/api.h>
#include <arrow/csv/api.h>
#include <arrow/io/api.h>

#include <memory>

#include "_cgo_export.h"
#include "csv.h"

// TODO: Unite with the one in carrow.cc
struct Table {
  std::shared_ptr<arrow::Table> table;
};


class GoStream: virtual public arrow::io::InputStream {
	long long id_;

	public:
	GoStream(long long id): id_(id) {}

	virtual arrow::Status Close() {
		return arrow::Status::OK();
	}

	arrow::Status Tell(int64_t *pos) const {
		auto res = istream_tell(id_);
		if (res.err != NULL) {
			auto err = std::string(res.err);
			return arrow::Status(arrow::StatusCode::UnknownError, err);
		}
		*pos = res.size;
		return arrow::Status::OK();
	}

	bool closed() const {
		auto res = istream_closed(id_);
		if (res.err != NULL) {
			return true;
		}

		return (res.size == 1) ? true : false;
	}

	arrow::Status Read(int64_t nbytes, int64_t* bytes_read, void* out) {
		auto res = istream_read(id_, nbytes);
		if (res.err != NULL) {
			auto err = std::string(res.err);
			return arrow::Status(arrow::StatusCode::UnknownError, err);
		}

		*bytes_read = res.size;
		memcpy(out, res.data, res.size);
		return arrow::Status::OK();
	}

	arrow::Status Read(int64_t nbytes, std::shared_ptr<arrow::Buffer>* out) {
		auto res = istream_read(id_, nbytes);
		if (res.err != NULL) {
			auto err = std::string(res.err);
			return arrow::Status(arrow::StatusCode::UnknownError, err);
		}

		auto data = (const uint8_t *)res.data;
		*out = std::make_shared<arrow::Buffer>(data, res.size);
		return arrow::Status::OK();
	}
};

read_res_t csv_read(long long id) {
	read_res_t res = {NULL, NULL};
	arrow::Status st;
	arrow::MemoryPool* pool = arrow::default_memory_pool();
	auto is = GoStream(id);
	// FIXME: How to convert GoStream to std::shared_ptr<arrow::io::InputStream>
	arrow::io::InputStream is = GoStream(id);
	auto input = std::make_shared<arrow::io::InputStream>(&is);

	auto read_options = arrow::csv::ReadOptions::Defaults();
	auto parse_options = arrow::csv::ParseOptions::Defaults();
	auto convert_options = arrow::csv::ConvertOptions::Defaults();

	std::shared_ptr<arrow::csv::TableReader> reader;
	st = arrow::csv::TableReader::Make(pool, input, read_options,
			parse_options, convert_options,
			&reader);
	if (!st.ok()) {
		res.err = st.message().c_str();
		return res;
	}

	std::shared_ptr<arrow::Table> table;
	st = reader->Read(&table);
	if (!st.ok()) {
		res.err = st.message().c_str();
		return res;
	}

	auto tp = new Table;
	tp->table = table;
	res.table = tp;
	return res;
}
