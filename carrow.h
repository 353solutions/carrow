#ifndef _MY_PACKAGE_FOO_H_
#define _MY_PACKAGE_FOO_H_

#ifdef __cplusplus
extern "C" {
#endif

extern const int INTEGER_DTYPE;
extern const int FLOAT_DTYPE;

void *field_new(char *name, int type);
const char *field_name(void *field);
int field_dtype(void *vp);
void field_free(void *vp);

void *fields_new();
void fields_append(void *vp, void *fp);
void fields_free(void *vp);

void *schema_new(void *vp);
void schema_free(void *vp);

void *array_builder_new(int dtype);
void array_builder_append_int(void *vp, long long value);
void array_builder_append_float(void *vp, double value);
typedef struct {
  const char *err;
  void *arr;
} finish_result_t;
finish_result_t array_builder_finish(void *vp);

void array_free(void *vp);

void *column_new(void *field, void *array);
void column_free(void *vp);

#ifdef __cplusplus
}
#endif // extern "C"

#endif
