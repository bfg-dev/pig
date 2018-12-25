-- +pig Name: Insert into history_tracked table
-- +pig Requiremets: Create history_tracked table
-- +pig Up
INSERT INTO "history_tracked" ("id", "from", "to", "name", "native", "data")
VALUES (1, 'xmi', 'usd', 'XIAO.1810.AS', 'hkd', '[
  {
    "type": "column",
    "name": "Xiaomi Corp. Shares Daily Returns (USD)",
    "yAxis": "0",
    "yAxisTitle": "Xiaomi Corp. Shares Daily Returns",
    "negativeColor": "#b80000"
  },
  {
    "type": "spline",
    "name": "Xiaomi Corp. Shares Cumulative Daily Returns (USD)",
    "yAxis": "1",
    "yAxisTitle": "Xiaomi Corp. Shares Cumulative Daily Returns"
  }
]' :: JSONB);

-- +pig Down
DELETE FROM "history_tracked"
WHERE "id" = 1;
