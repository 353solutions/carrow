#ifndef _CARROW_H_
#define _CARROW_H_

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

extern const int BOOL_DTYPE;
extern const int FLOAT64_DTYPE;
extern const int INTEGER64_DTYPE;
extern const int STRING_DTYPE;
extern const int TIMESTAMP_DTYPE;

void *field_new(char *name, int type);
const char *field_name(void *field);
int field_dtype(void *vp);
void field_free(void *vp);

void *fields_new();
void fields_append(void *vp, void *fp);
void fields_free(void *vp);

void *schema_new(void *vp);
void schema_free(void *vp);

typedef struct {
  const char *err;
  void *obj;
} result_t;

void *array_builder_new(int dtype);
void array_builder_append_bool(void *vp, int value);
void array_builder_append_float(void *vp, double value);
void array_builder_append_int(void *vp, long long value);
void array_builder_append_string(void *vp, char *value, size_t length);
void array_builder_append_timestamp(void *vp, long long value);
result_t array_builder_finish(void *vp);

void array_free(void *vp);

void *column_new(void *field, void *array);
void *column_field(void *vp);
int  column_dtype(void *vp);
void column_free(void *vp);

void *columns_new();
void columns_append(void *vp, void *cp);
void columns_free(void *vp);

void *table_new(void *sp, void *cp);
const char *table_validate(void *vp);
long long table_num_cols(void *vp);
long long table_num_rows(void *vp);
void table_free(void *vp);

void *plasma_connect(char *path);
int plasma_write(void *cp, void *tp, char *oid);
void *plasma_read(void *cp, char *oid, int64_t timeout_ms);
int plasma_release(void *cp, char *oid);
void plasma_disconnect(void *vp);

#ifdef __cplusplus
}
#endif // extern "C"

#endif // #ifdef _CARROW_H_
