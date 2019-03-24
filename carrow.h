#ifndef _MY_PACKAGE_FOO_H_
#define _MY_PACKAGE_FOO_H_

#ifdef __cplusplus
extern "C" {
#endif

extern const int INTEGER_DTYPE;
extern const int FLOAT_DTYPE;

extern int NC_INTEGER_DTYPE;
extern int NC_FLOAT_DTYPE;

void *field_new(char *name, int type);
const char *field_name(void *field);
int field_dtype(void *vp);
void field_free(void *vp);

void *schema_new();
void schema_add_field(void *sp, void *fp);
void schema_free(void *vp);

#ifdef __cplusplus
}
#endif // extern "C"

#endif
