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


read_res_t csv_read(long long id);

#ifdef __cplusplus
}
#endif // extern "C"

#endif // CARROW_CSV_H
