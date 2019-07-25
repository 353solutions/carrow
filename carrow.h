#ifndef _CARROW_H_
#define _CARROW_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stddef.h>
#include <stdint.h>

extern const int BOOL_DTYPE;
extern const int FLOAT64_DTYPE;
extern const int INTEGER64_DTYPE;
extern const int STRING_DTYPE;
extern const int TIMESTAMP_DTYPE;

typedef struct {
  const char *err;
  void *ptr;
  int64_t i;
} result_t;

void *field_new(char *name, int type);
const char *field_name(void *field);
int field_dtype(void *vp);
void field_free(void *vp);

void *fields_new();
void fields_append(void *vp, void *fp);
void fields_free(void *vp);

void *schema_new(void *vp);
result_t schema_meta(void *vp);
result_t schema_set_meta(void *vp, void *meta);
void schema_free(void *vp);

result_t array_builder_new(int dtype);
result_t array_builder_append_bool(void *vp, int value);
result_t array_builder_append_float(void *vp, double value);
result_t array_builder_append_int(void *vp, int64_t value);
result_t array_builder_append_string(void *vp, char *value, size_t length);
result_t array_builder_append_timestamp(void *vp, long long value);
result_t array_builder_finish(void *vp);

int64_t array_length(void *vp);
int array_bool_at(void *vp, long long i);
int64_t array_int_at(void *vp, long long i);
double array_float_at(void *vp, long long i);
const char *array_str_at(void *vp, long long i);
int64_t array_timestamp_at(void *vp, long long i);

void array_free(void *vp);

void *column_new(void *field, void *array);
void *column_field(void *vp);
int column_dtype(void *vp);
void column_free(void *vp);

void *columns_new();
void columns_append(void *vp, void *cp);
void columns_free(void *vp);

void *table_new(void *sp, void *cp);
long long table_num_cols(void *vp);
long long table_num_rows(void *vp);
void table_free(void *vp);

void *meta_new();
result_t meta_set(void *vp, const char *key, const char *value);
result_t meta_size(void *vp);
result_t meta_key(void *vp, int64_t i);
result_t meta_value(void *vp, int64_t i);

result_t plasma_connect(char *path);
result_t plasma_write(void *cp, void *tp, char *oid);
result_t plasma_read(void *cp, char *oid, int64_t timeout_ms);
result_t plasma_release(void *cp, char *oid);
result_t plasma_disconnect(void *vp);

#ifdef __cplusplus
}
#endif // extern "C"

#endif // #ifdef _CARROW_H_
