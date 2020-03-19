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

	arrow::Result<int64_t> Tell() const {
		auto res = istream_tell(id_);
		if (res.err != NULL) {
			auto err = std::string(res.err);
			return arrow::Status::IOError(err);
		}
		return res.size;
	}

	bool closed() const {
		auto res = istream_closed(id_);
		if (res.err != NULL) {
			return true;
		}

		return (res.size == 1) ? true : false;
	}

	arrow::Result<int64_t> Read(int64_t nbytes, void* out) {
		auto res = istream_read(id_, nbytes);
		if (res.err != NULL) {
			auto err = std::string(res.err);
			return arrow::Status::IOError(err);
		}

		memcpy(out, res.data, res.size);
		return res.size;
	}

	arrow::Result<std::shared_ptr<arrow::Buffer>> Read(int64_t nbytes) {
		auto res = istream_read(id_, nbytes);
		if (res.err != NULL) {
			auto err = std::string(res.err);
			return arrow::Status(arrow::StatusCode::UnknownError, err);
		}

		auto data = (const uint8_t *)res.data;
		return std::make_shared<arrow::Buffer>(data, res.size);
	}
};

read_res_t csv_read(long long id) {
	read_res_t res = {NULL, NULL};
	arrow::MemoryPool* pool = arrow::default_memory_pool();
	std::shared_ptr<arrow::io::InputStream> input = std::make_shared<GoStream>(id);

	// TODO: Allow user to pass options
	auto read_options = arrow::csv::ReadOptions::Defaults();
	auto parse_options = arrow::csv::ParseOptions::Defaults();
	auto convert_options = arrow::csv::ConvertOptions::Defaults();
	
	auto ptr = arrow::csv::TableReader::Make(pool, input, read_options,
			parse_options, convert_options);
	if (!ptr.ok()) {
		res.err = ptr.status().message().c_str();
		return res;
	}
	
	std::shared_ptr<arrow::csv::TableReader> reader = ptr.ValueOrDie();
	auto rptr = reader->Read();
	if (!rptr.ok()) {
		res.err = rptr.status().message().c_str();
		return res;
	}

	auto tp = new Table;
	tp->table = rptr.ValueOrDie();
	res.table = tp;
	return res;
}
