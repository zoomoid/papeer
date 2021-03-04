from minio import Minio
from minio.error import S3Error
import argparse
import logging
import os

logger = logging.getLogger("papeer-cli")
parser = argparse.ArgumentParser(
    description="Minimal python CLI to interact with the papeer backend, e.g. publishing releases.")

parser.add_argument("--file", dest="local_file")
parser.add_argument("--title", dest="title")
parser.add_argument("--description", dest="description")
parser.add_argument("--commit-file", dest="commit_file")
parser.add_argument("--no-file", dest="no_file")
parser.add_argument("--commit", dest="commmit")

def make_release(title, description, commit, file):


def bootstrap_client() -> Minio:
    """
    Bootstraps a Minio S3-compliant client for uploading assets to a bucket
    """
    endpoint = os.environ['MINIO_ENDPOINT']
    if endpoint:
        access_key = os.environ["MINIO_ACCESS_KEY"]
        secret_key = os.environ["MINIO_SECRET_KEY"]
        if access_key and secret_key:
            client = Minio(endpoint, access_key=access_key, secret_key=secret_key)
            bucket = os.environ["MINIO_BUCKET"]
            if bucket:
                found = client.bucket_exists(bucket)
                if found:
                    return client, { bucket: bucket, access_key: access_key, secret_key: secret_key }
                else:
                    logger.error("Given Minio bucket could not be found, aborting...")
            else:
                logger.error("Lacking Minio bucket, aborting...")
        else:
            logger.error("Lacking Minio Access Key or Secret Key, aborting...")
    else:
        logger.error("Lacking Minio endpoint, aborting...")

def main():
    args = parser.parse_args()
    if args.no_file:
        return
    else:
        client, config = bootstrap_client()
        if not client:
            exit(1)
        logger.info("Initialized minio client, uploading file...")
        if not os.path.exists(args.file):
            logger.error(f"File at {args.file} does not exist, aborting...")
            exit(1)
        path_prefix = os.environ["BUCKET_PATH_PREFIX"]
        basename = os.path.basename(arg.file)
        try
            result = client.fput_object(config["bucket"], path_prefix + basename, args.file)
        except S3Error as err:
            logger.error(f"Failed to upload {basename} to bucket: {config}, {path_prefix}, aborting...")
            exit(1)
        logger.info(f"Uploaded {basename} to {config["bucket"]}. Continuing with making a release...")
        try:
            release = make_release()
        except:
            logger.error("Failed to create release, aborting...")
            exit(1)
        

if __name__ == "main":
    try:
        main()
    except S3Error as err:
        logger.critical(f"S3 error occured! {err}",)
        exit(1)
