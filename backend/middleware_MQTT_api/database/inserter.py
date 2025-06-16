from google.cloud import bigquery
from config import settings

client = bigquery.Client.from_service_account_json(settings.GOOGLE_APPLICATION_CREDENTIALS)

def insert_row_to_bigquery(data: dict):
    """
    Inserta una fila en BigQuery.
    :param data: Diccionario con los datos a insertar.
    """
    table_id = f"{settings.GCP_PROJECT_ID}.{settings.BQ_DATASET}.{settings.BQ_TABLE}"
    errors = client.insert_rows_json(table_id, [data])
    if errors:
        print(f"Error al insertar en BigQuery: {errors}")
    else:
        print("Datos insertados correctamente en BigQuery.")