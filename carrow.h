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

void *schema_new();
void schema_add_field(void *sp, void *fp);
void schema_free(void *vp);

void *array_builder_new(int dtype);
void array_builder_append_int(void *vp, long long value);
void array_builder_append_float(void *vp, double value);
typedef struct {
  const char *err;
  void *arr;
} finish_result;
finish_result array_builder_finish(void *vp);

#ifdef __cplusplus
}
#endif // extern "C"

#endif
