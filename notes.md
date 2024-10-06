## BEFORE:
```
EXPLAIN (ANALYZE, BUFFERS, SERIALIZE) SELECT name,
       namespace,
       UID,
       phase,
       startedat,
       finishedat,
       coalesce((workflow::JSON)->'metadata'->>'labels', '{}') AS labels,
       coalesce((workflow::JSON)->'metadata'->>'annotations', '{}') AS annotations,
       coalesce((workflow::JSON)->'status'->>'progress', '') AS progress,
       coalesce((workflow::JSON)->'metadata'->>'creationTimestamp', '') AS creationtimestamp,
       (workflow::JSON)->'spec'->>'suspend' AS suspend,
       coalesce((workflow::JSON)->'status'->>'message', '') AS message,
       coalesce((workflow::JSON)->'status'->>'estimatedDuration', '0') AS estimatedduration,
       coalesce((workflow::JSON)->'status'->>'resourcesDuration', '{}') AS resourcesduration
FROM "argo_archived_workflows"
WHERE (("clustername" = 'bmlk2r5q'
        AND "instanceid" = '')
       AND "namespace" = 'jk84shpm5nj'
       AND EXISTS
         (SELECT 1
          FROM argo_archived_workflows_labels
          WHERE clustername = argo_archived_workflows.clustername
            AND UID = argo_archived_workflows.uid
            AND name = 'workflows.argoproj.io/phase'
            AND value = 'Succeeded')
       AND EXISTS
         (SELECT 1
          FROM argo_archived_workflows_labels
          WHERE clustername = argo_archived_workflows.clustername
            AND UID = argo_archived_workflows.uid
            AND name = 'workflows.argoproj.io/completed'
            AND value = 'true'))
ORDER BY "startedat" DESC
LIMIT 1;
```

```
 Limit  (cost=1.14..94.49 rows=1 width=346) (actual time=0.109..0.109 rows=1 loops=1)
   ->  Nested Loop  (cost=1.14..69733.34 rows=747 width=346) (actual time=0.108..0.108 rows=1 loops=1)
         Join Filter: ((argo_archived_workflows_labels.uid)::text = (argo_archived_workflows_labels_1.uid)::text)
         ->  Nested Loop  (cost=0.71..64201.81 rows=2216 width=767) (actual time=0.043..0.043 rows=1 loops=1)
               ->  Index Scan Backward using argo_archived_workflows_i4 on argo_archived_workflows  (cost=0.29..13481.29 rows=20533 width=730) (actual time=0.024..0.024 rows=1 loops=1)
                     Filter: (((clustername)::text = 'bmlk2r5q'::text) AND ((instanceid)::text = ''::text) AND ((namespace)::text = 'jk84shpm5nj'::text))
                     Rows Removed by Filter: 1
               ->  Index Scan using argo_archived_workflows_labels_pkey on argo_archived_workflows_labels  (cost=0.42..2.47 rows=1 width=46) (actual time=0.017..0.017 rows=1 loops=1)
                     Index Cond: (((clustername)::text = 'bmlk2r5q'::text) AND ((uid)::text = (argo_archived_workflows.uid)::text) AND ((name)::text = 'workflows.argoproj.io/phase'::text))
                     Filter: ((value)::text = 'Succeeded'::text)
         ->  Index Scan using argo_archived_workflows_labels_pkey on argo_archived_workflows_labels argo_archived_workflows_labels_1  (cost=0.42..2.47 rows=1 width=46) (actual time=0.005..0.005 rows=1 loops=1)
               Index Cond: (((clustername)::text = 'bmlk2r5q'::text) AND ((uid)::text = (argo_archived_workflows.uid)::text) AND ((name)::text = 'workflows.argoproj.io/completed'::text))
               Filter: ((value)::text = 'true'::text)
 Planning Time: 0.384 ms
 Execution Time: 0.160 ms
(15 rows)
```

## AFTER:
```
EXPLAIN ANALYZE SELECT name,
       namespace,
       UID,
       phase,
       startedat,
       finishedat,
       coalesce((workflow::JSON)->'metadata'->>'labels', '{}') AS labels,
       coalesce((workflow::JSON)->'metadata'->>'annotations', '{}') AS annotations,
       coalesce((workflow::JSON)->'status'->>'progress', '') AS progress,
       coalesce((workflow::JSON)->'metadata'->>'creationTimestamp', '') AS creationtimestamp,
       (workflow::JSON)->'spec'->>'suspend' AS suspend,
       coalesce((workflow::JSON)->'status'->>'message', '') AS message,
       coalesce((workflow::JSON)->'status'->>'estimatedDuration', '0') AS estimatedduration,
       coalesce((workflow::JSON)->'status'->>'resourcesDuration', '{}') AS resourcesduration
FROM "argo_archived_workflows"
WHERE "clustername" = 'bmlk2r5q'
  AND UID IN
    (SELECT UID
     FROM "argo_archived_workflows"
     WHERE (("clustername" = 'bmlk2r5q'
             AND "instanceid" = '')
            AND "namespace" = 'jk84shpm5nj'
            AND EXISTS
              (SELECT 1
               FROM argo_archived_workflows_labels
               WHERE clustername = argo_archived_workflows.clustername
                 AND UID = argo_archived_workflows.uid
                 AND name = 'workflows.argoproj.io/phase'
                 AND value = 'Succeeded')
            AND EXISTS
              (SELECT 1
               FROM argo_archived_workflows_labels
               WHERE clustername = argo_archived_workflows.clustername
                 AND UID = argo_archived_workflows.uid
                 AND name = 'workflows.argoproj.io/completed'
                 AND value = 'true'))
     ORDER BY "startedat" DESC
     LIMIT 1);
```

EXPLAIN:
```
 Nested Loop  (cost=94.88..102.95 rows=1 width=346) (actual time=0.161..0.168 rows=1 loops=1)
   ->  HashAggregate  (cost=94.46..94.47 rows=1 width=37) (actual time=0.048..0.049 rows=1 loops=1)
         Group Key: ("ANY_subquery".uid)::text
         ->  Subquery Scan on "ANY_subquery"  (cost=1.14..94.46 rows=1 width=37) (actual time=0.044..0.045 rows=1 loops=1)
               ->  Limit  (cost=1.14..94.45 rows=1 width=45) (actual time=0.043..0.044 rows=1 loops=1)
                     ->  Nested Loop  (cost=1.14..69703.46 rows=747 width=45) (actual time=0.042..0.043 rows=1 loops=1)
                           Join Filter: ((argo_archived_workflows_labels.uid)::text = (argo_archived_workflows_labels_1.uid)::text)
                           ->  Nested Loop  (cost=0.71..64201.81 rows=2216 width=91) (actual time=0.036..0.037 rows=1 loops=1)
                                 ->  Index Scan Backward using argo_archived_workflows_i4 on argo_archived_workflows argo_archived_workflows_1  (cost=0.29..13481.29 rows=20533 width=54) (actual time=0.017..0.017 rows=1 loops=1)
                                       Filter: (((clustername)::text = 'bmlk2r5q'::text) AND ((instanceid)::text = ''::text) AND ((namespace)::text = 'jk84shpm5nj'::text))
                                       Rows Removed by Filter: 1
                                 ->  Index Scan using argo_archived_workflows_labels_pkey on argo_archived_workflows_labels  (cost=0.42..2.47 rows=1 width=46) (actual time=0.018..0.018 rows=1 loops=1)
                                       Index Cond: (((clustername)::text = 'bmlk2r5q'::text) AND ((uid)::text = (argo_archived_workflows_1.uid)::text) AND ((name)::text = 'workflows.argoproj.io/phase'::text))
                                       Filter: ((value)::text = 'Succeeded'::text)
                           ->  Index Scan using argo_archived_workflows_labels_pkey on argo_archived_workflows_labels argo_archived_workflows_labels_1  (cost=0.42..2.47 rows=1 width=46) (actual time=0.005..0.005 rows=1 loops=1)
                                 Index Cond: (((clustername)::text = 'bmlk2r5q'::text) AND ((uid)::text = (argo_archived_workflows_1.uid)::text) AND ((name)::text = 'workflows.argoproj.io/completed'::text))
                                 Filter: ((value)::text = 'true'::text)
   ->  Index Scan using argo_archived_workflows_pkey on argo_archived_workflows  (cost=0.42..8.44 rows=1 width=721) (actual time=0.046..0.046 rows=1 loops=1)
         Index Cond: (((clustername)::text = 'bmlk2r5q'::text) AND ((uid)::text = ("ANY_subquery".uid)::text))
 Planning Time: 0.595 ms
 Execution Time: 0.234 ms
(21 rows)
```


