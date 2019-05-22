from argparse import ArgumentParser

from pyarrow import plasma
import pyarrow as pa

parser = ArgumentParser(description='read from plasma')
parser.add_argument('db_path', help='db path')
parser.add_argument('oid', help='object id')
args = parser.parse_args()

conn: plasma.PlasmaClient = plasma.connect(args.db_path, '', 0)
oid = args.oid.rjust(20, '0')
oid = plasma.ObjectID(oid.encode())
print(oid.binary())

if oid not in conn.list():
    raise SystemExit(f'error: unknown ID: {oid.binary()}')
buf, = conn.get_buffers([oid])
reader = pa.RecordBatchStreamReader(buf)
record_batch: pa.RecordBatch = reader.read_next_batch()
df = record_batch.to_pandas()
print(df)