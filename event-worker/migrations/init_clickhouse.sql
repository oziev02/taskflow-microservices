CREATE TABLE IF NOT EXISTS task_events (
                                           timestamp DateTime,
                                           topic String,
                                           payload String
) ENGINE = MergeTree()
    ORDER BY timestamp
