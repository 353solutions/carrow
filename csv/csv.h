#ifndef CARROW_CSV_H
#define CARROW_CSV_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
	void *data;
	unsigned long long size;
	char *err;
} csv_res_t;

typedef struct {
	void *table;
	const char *err;
} read_res_t;

typedef char cbool;

typedef struct {
  char delimiter;
  cbool quoting;
  char quote_char;
  cbool double_quote;
  cbool escaping;
  char escape_char;
  cbool newlines_in_values;
  cbool ignore_empty_lines;
} parse_options_t;

parse_options_t default_parse_options();

read_res_t csv_read(long long id, parse_options_t po);

#ifdef __cplusplus
}
#endif // extern "C"

#endif // CARROW_CSV_H
