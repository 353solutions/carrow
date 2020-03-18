#include <arrow/api.h>
#include <arrow/csv/api.h>
#include <arrow/io/api.h>

#include <memory>
#include <cstdint>

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

parse_options_t default_parse_options() {
	auto arrow_opts = arrow::csv::ParseOptions::Defaults();

	parse_options_t opts;
	opts.delimiter = arrow_opts.delimiter;
	opts.quoting = arrow_opts.quoting;
	opts.quote_char = arrow_opts.quote_char;
	opts.double_quote = arrow_opts.double_quote;
	opts.escaping = arrow_opts.escaping;
	opts.escape_char = arrow_opts.escape_char;
	opts.newlines_in_values = arrow_opts.newlines_in_values;
  opts.ignore_empty_lines = arrow_opts.ignore_empty_lines;
	return opts;
}

arrow::csv::ParseOptions popts_from_c(parse_options_t p) {
	auto opts = arrow::csv::ParseOptions::Defaults();
	opts.delimiter = p.delimiter;
	opts.quoting = p.quoting;
	opts.quote_char = p.quote_char;
	opts.double_quote = p.double_quote;
	opts.escaping = p.escaping;
	opts.escape_char = p.escape_char;
	opts.newlines_in_values = p.newlines_in_values;
  opts.ignore_empty_lines = p.ignore_empty_lines;

	return opts;
}


read_options_t default_read_options() {
	read_options_t opts;

	auto arrow_opts = arrow::csv::ReadOptions::Defaults();
	opts.use_threads = arrow_opts.use_threads;
	opts.block_size = arrow_opts.block_size;
	opts.skip_rows = arrow_opts.skip_rows;
	opts.column_names = nullptr;
	opts.column_name_count = 0;
	opts.autogenerate_column_names = arrow_opts.autogenerate_column_names;

	return opts;
}

arrow::csv::ReadOptions ropts_from_c(read_options_t p) {
	auto opts = arrow::csv::ReadOptions::Defaults();
	opts.use_threads = p.use_threads;
	opts.block_size = p.block_size;
	opts.skip_rows = p.skip_rows;

	if ((p.column_name_count > 0) && (p.column_names != nullptr)) {
		for (int i = 0; i < p.column_name_count; i++) {
			opts.column_names.push_back(p.column_names[i]);
		}
		// Allocated in Go with C.CString
		free(p.column_names);
	}

	opts.autogenerate_column_names = p.autogenerate_column_names;

	return opts;
}

read_res_t csv_read(
		long long id,
		read_options_t ro,
		parse_options_t po) {
	read_res_t res = {NULL, NULL};
	arrow::MemoryPool* pool = arrow::default_memory_pool();
	std::shared_ptr<arrow::io::InputStream> input = std::make_shared<GoStream>(id);

	// TODO: Allow user to pass options
	auto read_options = ropts_from_c(ro);
	auto parse_options = popts_from_c(po);
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
