from pyarrow import plasma
from argparse import ArgumentParser

parser = ArgumentParser(description='read from plasma')
parser.add_argument('db_path', help='db path')
parser.add_argument('oid', help='object id')
args = parser.parse_args()

conn = plasma.PlasmaClient(args.db_path)
oid = args.oid.rjust(20, '0')
oid = plasma.ObjectID(oid.encode())
buf, = conn.get([oid])