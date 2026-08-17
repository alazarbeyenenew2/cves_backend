-- name: UpdateMetaData :exec 
INSERT INTO meta_datas (id, last_fetch, total_fetched, sync_status)
VALUES (1, $1,$2,$3)
ON CONFLICT (id)
DO UPDATE SET
    last_fetch = EXCLUDED.last_fetch,
    total_fetched = EXCLUDED.total_fetched,
    sync_status = EXCLUDED.sync_status;

-- name: UpdateMetaDataSyncStatus :exec 
UPDATE meta_datas set sync_status = $1; 

-- name: GetMetaData :one 
SELECT * FROM meta_datas;