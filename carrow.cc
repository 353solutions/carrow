#include <arrow/api.h>
#include <iostream>

#include "carrow.h"

#ifdef __cplusplus
extern "C"
{
#endif

  const int INTEGER64_DTYPE = arrow::Type::INT64;
  const int FLOAT64_DTYPE = arrow::Type::DOUBLE;

  std::shared_ptr<arrow::DataType> data_type(int dtype)
  {
    switch (dtype)
    {
    case INTEGER64_DTYPE:
      return arrow::int64();
    case FLOAT64_DTYPE:
      return arrow::float64();
    }

    return NULL;
  }

  void *field_new(char *name, int dtype)
  {
    auto dt = data_type(dtype);
    return new arrow::Field(name, dt);
  }

  const char *field_name(void *vp)
  {
    auto field = (arrow::Field *)vp;
    return field->name().c_str();
  }

  int field_dtype(void *vp)
  {
    auto field = (arrow::Field *)vp;
    return field->type()->id();
  }

  void field_free(void *vp)
  {
    if (vp == NULL)
    {
      return;
    }
    auto field = (arrow::Field *)vp;
    delete field;
  }

  void *fields_new() { return new std::vector<std::shared_ptr<arrow::Field>>(); }

  void fields_append(void *vp, void *fp)
  {
    auto fields = (std::vector<std::shared_ptr<arrow::Field>> *)vp;
    std::shared_ptr<arrow::Field> field((arrow::Field *)fp);
    fields->push_back(field);
  }

  void fields_free(void *vp)
  {
    if (vp == NULL)
    {
      return;
    }
    auto fields = (std::vector<std::shared_ptr<arrow::Field>> *)vp;
    delete fields;
  }

  void *schema_new(void *vp)
  {
    auto fields = (std::vector<std::shared_ptr<arrow::Field>> *)vp;
    auto schema = new arrow::Schema(*fields);
    return (void *)schema;
  }

  void schema_free(void *vp)
  {
    if (vp == NULL)
    {
      return;
    }
    auto schema = (arrow::Schema *)vp;
    delete schema;
  }

  void *array_builder_new(int dtype)
  {
    switch (dtype)
    {
    case INTEGER64_DTYPE:
      return new arrow::Int64Builder();
    case FLOAT64_DTYPE:
      return new arrow::DoubleBuilder();
    }

    return NULL;
  }

  void array_builder_append_int(void *vp, long long value)
  {
    auto builder = (arrow::Int64Builder *)vp;
    builder->Append(value);
  }

  void array_builder_append_float(void *vp, double value)
  {
    auto builder = (arrow::DoubleBuilder *)vp;
    builder->Append(value);
  }

  finish_result_t array_builder_finish(void *vp)
  {
    auto builder = (arrow::ArrayBuilder *)vp;
    std::shared_ptr<arrow::Array> out;
    auto status = builder->Finish(&out);

    finish_result_t res = {NULL, NULL};
    if (!status.ok())
    {
      res.err = status.ToString().c_str();
    }
    else
    {
      res.arr = (void *)(out.get());
    }

    // TODO: Will out delete the underlying array?
    return res;
  }

  void array_free(void *vp)
  {
    if (vp == NULL)
    {
      return;
    }
    auto array = (arrow::Array *)vp;
    delete array;
  }

  void *column_new(void *fp, void *ap)
  {
    std::shared_ptr<arrow::Field> field((arrow::Field *)fp);
    std::shared_ptr<arrow::Array> array((arrow::Array *)ap);

    return new arrow::Column(field, array);
  }

  int column_dtype(void *vp)
  {
    auto column = (arrow::Column *)vp;
    return column->type()->id();
  }

  void column_free(void *vp)
  {
    if (vp == NULL)
    {
      return;
    }
    auto column = (arrow::Column *)vp;
    delete column;
  }

  void *column_field(void *vp)
  {
    auto column = (arrow::Column *)vp;
    return column->field().get();
  }

  void *columns_new()
  {
    return new std::vector<std::shared_ptr<arrow::Column>>();
  }

  void columns_append(void *vp, void *cp)
  {
    auto columns = (std::vector<std::shared_ptr<arrow::Column>> *)vp;
    std::shared_ptr<arrow::Column> column((arrow::Column *)cp);
    columns->push_back(column);
  }

  void columns_free(void *vp)
  {
    auto columns = (std::vector<std::shared_ptr<arrow::Column>> *)vp;
    delete columns;
  }

  void *table_new(void *sp, void *cp)
  {
    std::shared_ptr<arrow::Schema> schema((arrow::Schema *)sp);
    auto columns = (std::vector<std::shared_ptr<arrow::Column>> *)cp;

    auto table = arrow::Table::Make(schema, *columns);
    return new std::shared_ptr<arrow::Table>(table);
    ;
  }

  const char *table_validate(void *vp)
  {
    return NULL;
    /*
        auto table = (arrow::Table *)vp;
        // FIXME: arrow::Table::Validate is pure virtual
        auto status = table->Validate();
        if (status.ok()) {
        return NULL;
        }

        return status.ToString().c_str();
  */
  }

  long long table_num_cols(void *vp)
  {
    auto table = (std::shared_ptr<arrow::Table> *)vp;
    return table->get()->num_columns();
  }

  long long table_num_rows(void *vp)
  {
    auto table = (std::shared_ptr<arrow::Table> *)vp;
    return table->get()->num_rows();
  }

  void table_free(void *vp)
  {
    if (vp == NULL)
    {
      return;
    }
    auto table = (arrow::Table *)vp;
    delete table;
  }

#ifdef __cplusplus
} // extern "C"
#endif
