# Virgil Benchmark — Detection Gap Report
#
# Generated: 2026-04-03 18:55:25
# Manifest entries:  42
# Virgil findings:   368
# Matched:           34
# Unmatched (gaps):  8
# Extra (new finds): 340
#
# Overall detection rate: 81.0%
# Precision: 9.1%
# Recall: 81.0%
# F1 Score: 16.3%

## Match Quality

| Match Type           |  Count | Percentage |
| -------------------- | ------ | ---------- |
| Pattern match        |     12 |      35.3% |
| Strict match         |     15 |      44.1% |
| Relaxed match        |      7 |      20.6% |

## Per-Category Summary

| Category                  |    TP |    FP |    FN |     Prec |   Recall |       F1 |
| ------------------------- | ----- | ----- | ----- | -------- | -------- | -------- |
| architecture              |     0 |    25 |     0 |     0.0% |     0.0% |     0.0% |
| documentation             |     0 |     1 |     0 |     0.0% |     0.0% |     0.0% |
| security                  |     0 |     8 |     0 |     0.0% |     0.0% |     0.0% |
| unknown                   |     0 |    36 |     0 |     0.0% |     0.0% |     0.0% |
| tech-debt                 |     2 |     0 |     2 |   100.0% |    50.0% |    66.7% |
| style                     |     6 |   240 |     3 |     2.4% |    66.7% |     4.7% |
| scalability               |    14 |    12 |     3 |    53.8% |    82.4% |    65.1% |
| code-quality              |    12 |    18 |     0 |    40.0% |   100.0% |    57.1% |

## Undetected Debt Instances

These manifest entries had NO corresponding virgil-cli finding.
Each represents a detection gap virgil-cli could be improved to cover.

### scalability (3 undetected / 17 total — 82.4% detected)

- `cmd/worker/main.go:23` **memory-leak-indicators** — go dispatcher.Start() and go scheduler.RunForever() -- goroutines spawned without context cancellation, leak on abnormal exit
- `internal/service/order.go:50` **n-plus-one-queries** — CreateOrder() decrements stock per item in a loop instead of batch update
- `internal/service/order.go:78` **n-plus-one-queries** — GetOrder() fetches product names one by one in a loop for each order item

### style (3 undetected / 9 total — 66.7% detected)

- `cmd/server/main.go:150` **dead-code** — Unreachable variable assignments after server exit (cacheClient, memQueue, userRepo, notificationSvc never used)
- `cmd/worker/main.go:38` **dead-code** — Unreachable cleanup code after infinite for loop -- fmt.Println and dispatcher.Stop() never execute
- `internal/api/middleware.go:129` **dead-code** — CORSMiddleware(), RecoveryMiddleware(), and RequestIDMiddleware() are defined but never registered in the router middleware chain

### tech-debt (2 undetected / 4 total — 50.0% detected)

- `go.mod:1` **outdated-dependency** — go.mod specifies older Go toolchain version; should target current stable Go release for security patches and improved runtime performance
- `go.mod:1` **version-drift** — Direct dependencies not pinned to specific patch versions in go.mod; go.sum ensures reproducibility but no upgrade policy prevents accumulation of patch-level drift

## Extra Findings (not in manifest)

These virgil-cli detections had no corresponding manifest entry.
They may be: (a) real debt not cataloged in the manifest, (b) false positives,
or (c) findings at a different granularity than the manifest entries.

### architecture (25 extra findings)

- `internal/api/cors.go:9` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `CORSConfig` has exported field `AllowedOrigins` — consider encapsulating with methods
- `internal/api/health.go:10` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `HealthStatus` has exported field `Status` — consider encapsulating with methods
- `internal/api/response.go:9` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `ErrorResponse` has exported field `Error` — consider encapsulating with methods
- `internal/config/config.go:8` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `Config` has exported field `Port` — consider encapsulating with methods
- `internal/config/validate.go:9` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `ValidationError` has exported field `Errors` — consider encapsulating with methods
- `internal/model/address.go:9` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `Address` has exported field `ID` — consider encapsulating with methods
- `internal/model/audit_log.go:18` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `AuditLog` has exported field `ID` — consider encapsulating with methods
- `internal/model/category.go:8` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `Category` has exported field `ID` — consider encapsulating with methods
- `internal/model/order.go:8` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `Order` has exported field `ID` — consider encapsulating with methods
- `internal/model/order.go:20` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `OrderItem` has exported field `ID` — consider encapsulating with methods
- `internal/model/product.go:9` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `Product` has exported field `ID` — consider encapsulating with methods
- `internal/model/user.go:8` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `User` has exported field `ID` — consider encapsulating with methods
- `internal/service/order.go:1` [api_surface_area/**excessive_public_api**] Module exports 9/11 symbols (82% exported, threshold: >80%)
- `internal/service/payment.go:14` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `PaymentResult` has exported field `TransactionID` — consider encapsulating with methods
- `internal/service/shipping.go:11` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `ShippingRate` has exported field `Method` — consider encapsulating with methods
- `internal/worker/dispatcher.go:1` [api_surface_area/**excessive_public_api**] Module exports 9/11 symbols (82% exported, threshold: >80%)
- `internal/worker/dispatcher.go:12` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `Job` has exported field `ID` — consider encapsulating with methods
- `internal/worker/dispatcher.go:21` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `JobResult` has exported field `JobID` — consider encapsulating with methods
- `internal/worker/dispatcher.go:28` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `DispatcherStats` has exported field `Processed` — consider encapsulating with methods
- `internal/worker/reporter.go:11` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `MetricSnapshot` has exported field `Timestamp` — consider encapsulating with methods
- `internal/worker/scheduler.go:12` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `ScheduledTask` has exported field `Name` — consider encapsulating with methods
- `pkg/cache/cache.go:16` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `CacheEntry` has exported field `Value` — consider encapsulating with methods
- `pkg/queue/priority.go:1` [api_surface_area/**excessive_public_api**] Module exports 12/13 symbols (92% exported, threshold: >80%)
- `pkg/queue/priority.go:9` [api_surface_area/**leaky_abstraction_boundary**] Exported struct `PriorityItem` has exported field `Value` — consider encapsulating with methods
- `pkg/queue/queue.go:1` [module_size_distribution/**anemic_module**] Module contains only 1 definition — consider merging into a related module

### code-quality (18 extra findings)

- `cmd/server/main.go:28` [function_length/**too_many_statements**] Function 'main' has 73 statements (threshold: 20)
- `cmd/server/main.go:106` [magic_numbers/**magic_number**] magic number `30` — consider extracting to a named constant for clarity
- `cmd/server/main.go:107` [magic_numbers/**magic_number**] magic number `120` — consider extracting to a named constant for clarity
- `internal/api/cors.go:22` [magic_numbers/**magic_number**] magic number `86400` — consider extracting to a named constant for clarity
- `internal/api/handler.go:44` [function_length/**too_many_statements**] Function 'CreateOrder' has 68 statements (threshold: 20)
- `internal/api/handler.go:44` [cognitive_complexity/**high_cognitive_complexity**] Cognitive complexity of 22 (threshold: 15) in function 'CreateOrder'
- `internal/api/middleware.go:32` [function_length/**too_many_statements**] Function 'AuthMiddleware' has 21 statements (threshold: 20)
- `internal/api/middleware.go:32` [cognitive_complexity/**high_cognitive_complexity**] Cognitive complexity of 17 (threshold: 15) in function 'AuthMiddleware'
- `internal/api/middleware.go:74` [function_length/**function_too_long**] Function 'RateLimitMiddleware' is 53 lines long (threshold: 50)
- `internal/api/middleware.go:74` [function_length/**too_many_statements**] Function 'RateLimitMiddleware' has 33 statements (threshold: 20)
- `internal/api/middleware.go:74` [cognitive_complexity/**high_cognitive_complexity**] Cognitive complexity of 22 (threshold: 15) in function 'RateLimitMiddleware'
- `internal/api/router.go:14` [function_length/**too_many_statements**] Function 'NewRouter' has 23 statements (threshold: 20)
- `internal/config/validate.go:24` [cyclomatic_complexity/**high_cyclomatic_complexity**] Cyclomatic complexity of 11 (threshold: 10) in function 'ValidateConfig'
- `internal/config/validate.go:24` [function_length/**too_many_statements**] Function 'ValidateConfig' has 21 statements (threshold: 20)
- `internal/config/validate.go:27` [magic_numbers/**magic_number**] magic number `65535` — consider extracting to a named constant for clarity
- `internal/service/shipping.go:48` [magic_numbers/**magic_number**] magic number `14` — consider extracting to a named constant for clarity
- `internal/worker/scheduler.go:37` [magic_numbers/**magic_number**] magic number `24` — consider extracting to a named constant for clarity
- `internal/worker/scheduler.go:47` [magic_numbers/**magic_number**] magic number `24` — consider extracting to a named constant for clarity

### documentation (1 extra findings)

- `cmd/migrate/main.go:1` [comment_to_code_ratio/**under_documented**] Comment-to-code ratio is 0.00 (0 comment lines, 76 code lines) — consider adding documentation

### scalability (12 extra findings)

- `cmd/migrate/main.go:77` [n_plus_one_queries/**db_query_in_loop**] Database/ORM call `db.Exec(migration)` inside loop — potential N+1 query problem
- `internal/api/middleware.go:82` [sync_blocking_in_async/**blocking_in_goroutine**] blocking call `time.Sleep` inside goroutine may stall concurrent execution
- `internal/repository/inventory_repo.go:57` [memory_leak_indicators/**unbounded_append_in_loop**] append() inside loop without apparent bound check — potential unbounded slice growth
- `internal/repository/inventory_repo.go:127` [memory_leak_indicators/**unbounded_append_in_loop**] append() inside loop without apparent bound check — potential unbounded slice growth
- `internal/repository/order_repo.go:40` [n_plus_one_queries/**db_query_in_loop**] Database/ORM call `tx.Exec(` inside loop — potential N+1 query problem
- `internal/repository/order_repo.go:92` [memory_leak_indicators/**unbounded_append_in_loop**] append() inside loop without apparent bound check — potential unbounded slice growth
- `internal/repository/order_repo.go:117` [memory_leak_indicators/**unbounded_append_in_loop**] append() inside loop without apparent bound check — potential unbounded slice growth
- `internal/repository/order_repo.go:152` [memory_leak_indicators/**unbounded_append_in_loop**] append() inside loop without apparent bound check — potential unbounded slice growth
- `internal/repository/user_repo.go:88` [memory_leak_indicators/**unbounded_append_in_loop**] append() inside loop without apparent bound check — potential unbounded slice growth
- `internal/worker/processor.go:90` [memory_leak_indicators/**unbounded_append_in_loop**] append() inside loop without apparent bound check — potential unbounded slice growth
- `pkg/queue/memory.go:45` [memory_leak_indicators/**unbounded_append_in_loop**] append() inside loop without apparent bound check — potential unbounded slice growth
- `pkg/queue/memory.go:75` [memory_leak_indicators/**unbounded_append_in_loop**] append() inside loop without apparent bound check — potential unbounded slice growth

### security (8 extra findings)

- `cmd/server/main.go:114` [resource_exhaustion/**unbounded_goroutine_spawn**] goroutine spawned in loop without bound may exhaust resources
- `internal/api/middleware.go:80` [resource_exhaustion/**unbounded_goroutine_spawn**] goroutine spawned in loop without bound may exhaust resources
- `internal/repository/inventory_repo.go:99` [sql_injection/**sql_string_concat**] SQL query via .Exec() uses string concatenation — use parameterized queries instead
- `internal/repository/inventory_repo.go:108` [sql_injection/**sql_string_concat**] SQL query via .Query() uses string concatenation — use parameterized queries instead
- `internal/service/payment.go:54` [resource_exhaustion/**unbounded_body_read**] io.ReadAll without MaxBytesReader may exhaust memory on large input
- `internal/service/payment.go:92` [resource_exhaustion/**unbounded_body_read**] io.ReadAll without MaxBytesReader may exhaust memory on large input
- `internal/worker/processor.go:82` [resource_exhaustion/**unbounded_goroutine_spawn**] goroutine spawned in loop without bound may exhaust resources
- `internal/worker/scheduler.go:63` [resource_exhaustion/**unbounded_goroutine_spawn**] goroutine spawned in loop without bound may exhaust resources

### style (240 extra findings)

- `cmd/server/main.go:3` [coupling/**excessive_imports**] file has 18 imports (threshold: 15) — consider splitting into smaller packages
- `cmd/server/main.go:10` [dead_code/**unused_import**] import `os` appears unused
- `internal/api/cors.go:8` [dead_exports/**dead_export**] Exported struct 'CORSConfig' is not referenced by any other file in the project
- `internal/api/cors.go:16` [dead_exports/**dead_export**] Exported function 'DefaultCORSConfig' is not referenced by any other file in the project
- `internal/api/cors.go:26` [dead_exports/**dead_export**] Exported function 'CORSMiddleware' is not referenced by any other file in the project
- `internal/api/handler.go:13` [dead_code/**unused_import**] import `model` appears unused
- `internal/api/handler.go:14` [dead_code/**unused_import**] import `service` appears unused
- `internal/api/handler.go:17` [dead_exports/**dead_export**] Exported struct 'Handler' is not referenced by any other file in the project
- `internal/api/health.go:9` [dead_exports/**dead_export**] Exported struct 'HealthStatus' is not referenced by any other file in the project
- `internal/api/health.go:17` [dead_exports/**dead_export**] Exported struct 'HealthHandler' is not referenced by any other file in the project
- `internal/api/health.go:23` [dead_exports/**dead_export**] Exported function 'NewHealthHandler' is not referenced by any other file in the project
- `internal/api/health.go:51` [dead_exports/**dead_export**] Exported method 'IsHealthy' is not referenced by any other file in the project
- `internal/api/health.go:52` [coupling/**low_cohesion**] method `IsHealthy` does not use its receiver `h` — consider making it a function
- `internal/api/middleware.go:8` [dead_code/**unused_import**] import `sync` appears unused
- `internal/api/middleware.go:11` [dead_code/**unused_import**] import `logger` appears unused
- `internal/api/ratelimit.go:5` [dead_code/**unused_import**] import `sync` appears unused
- `internal/api/ratelimit.go:9` [dead_exports/**dead_export**] Exported struct 'RateLimiter' is not referenced by any other file in the project
- `internal/api/ratelimit.go:24` [dead_exports/**dead_export**] Exported function 'NewRateLimiter' is not referenced by any other file in the project
- `internal/api/ratelimit.go:33` [dead_exports/**dead_export**] Exported method 'Allow' is not referenced by any other file in the project
- `internal/api/ratelimit.go:60` [dead_exports/**dead_export**] Exported method 'Middleware' is not referenced by any other file in the project
- `internal/api/ratelimit.go:72` [dead_exports/**dead_export**] Exported method 'Reset' is not referenced by any other file in the project
- `internal/api/response.go:8` [dead_exports/**dead_export**] Exported struct 'ErrorResponse' is not referenced by any other file in the project
- `internal/api/response.go:34` [dead_code/**unused_private_function**] unexported function `decodeJSON` appears unused within this file
- `internal/api/router.go:7` [dead_exports/**dead_export**] Exported struct 'Router' is not referenced by any other file in the project
- `internal/api/router.go:22` [duplicate_code/**duplicate_switch_case**] switch case at line 22 has a duplicate body (switch at line 21)
- `internal/api/router.go:41` [duplicate_code/**duplicate_switch_case**] switch case at line 41 has a duplicate body (switch at line 40)
- `internal/api/router.go:72` [dead_exports/**dead_export**] Exported method 'ListInventory' is not referenced by any other file in the project
- `internal/api/router.go:82` [dead_exports/**dead_export**] Exported method 'UpdateInventory' is not referenced by any other file in the project
- `internal/config/config.go:7` [dead_exports/**dead_export**] Exported struct 'Config' is not referenced by any other file in the project
- `internal/config/config.go:75` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'Error' exported from 3 files: internal/config/config.go, internal/config/validate.go, pkg/logger/logger.go
- `internal/config/defaults.go:8` [dead_exports/**dead_export**] Exported constant 'DefaultPort' is not referenced by any other file in the project
- `internal/config/defaults.go:9` [dead_exports/**dead_export**] Exported constant 'DefaultReadTimeout' is not referenced by any other file in the project
- `internal/config/defaults.go:10` [dead_exports/**dead_export**] Exported constant 'DefaultWriteTimeout' is not referenced by any other file in the project
- `internal/config/defaults.go:11` [dead_exports/**dead_export**] Exported constant 'DefaultMaxOpenConns' is not referenced by any other file in the project
- `internal/config/defaults.go:12` [dead_exports/**dead_export**] Exported constant 'DefaultMaxIdleConns' is not referenced by any other file in the project
- `internal/config/defaults.go:13` [dead_exports/**dead_export**] Exported constant 'DefaultLogLevel' is not referenced by any other file in the project
- `internal/config/defaults.go:14` [dead_exports/**dead_export**] Exported constant 'DefaultShutdownGrace' is not referenced by any other file in the project
- `internal/config/defaults.go:15` [dead_exports/**dead_export**] Exported constant 'DefaultWorkerCount' is not referenced by any other file in the project
- `internal/config/defaults.go:19` [dead_exports/**dead_export**] Exported function 'DefaultConfig' is not referenced by any other file in the project
- `internal/config/defaults.go:32` [dead_exports/**dead_export**] Exported method 'WithPort' is not referenced by any other file in the project
- `internal/config/defaults.go:38` [dead_exports/**dead_export**] Exported method 'WithDatabaseURL' is not referenced by any other file in the project
- `internal/config/defaults.go:44` [dead_exports/**dead_export**] Exported method 'WithLogLevel' is not referenced by any other file in the project
- `internal/config/defaults.go:50` [dead_exports/**dead_export**] Exported method 'WithWorkerCount' is not referenced by any other file in the project
- `internal/config/env.go:9` [dead_code/**unused_private_function**] unexported function `getEnvStr` appears unused within this file
- `internal/config/env.go:17` [dead_code/**unused_private_function**] unexported function `getEnvInt` appears unused within this file
- `internal/config/env.go:17` [duplicate_code/**duplicate_function_body**] function `getEnvInt` has a body identical to: `getEnvBool`
- `internal/config/env.go:30` [dead_code/**unused_private_function**] unexported function `getEnvBool` appears unused within this file
- `internal/config/env.go:30` [duplicate_code/**duplicate_function_body**] function `getEnvBool` has a body identical to: `getEnvInt`
- `internal/config/validate.go:8` [dead_exports/**dead_export**] Exported struct 'ValidationError' is not referenced by any other file in the project
- `internal/config/validate.go:13` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'Error' exported from 3 files: internal/config/config.go, internal/config/validate.go, pkg/logger/logger.go
- `internal/config/validate.go:18` [dead_exports/**dead_export**] Exported method 'HasErrors' is not referenced by any other file in the project
- `internal/config/validate.go:23` [dead_exports/**dead_export**] Exported function 'ValidateConfig' is not referenced by any other file in the project
- `internal/model/address.go:8` [dead_exports/**dead_export**] Exported struct 'Address' is not referenced by any other file in the project
- `internal/model/audit_log.go:8` [dead_exports/**dead_export**] Exported type_alias 'AuditAction' is not referenced by any other file in the project
- `internal/model/audit_log.go:11` [dead_exports/**dead_export**] Exported constant 'AuditActionCreate' is not referenced by any other file in the project
- `internal/model/audit_log.go:12` [dead_exports/**dead_export**] Exported constant 'AuditActionUpdate' is not referenced by any other file in the project
- `internal/model/audit_log.go:13` [dead_exports/**dead_export**] Exported constant 'AuditActionDelete' is not referenced by any other file in the project
- `internal/model/audit_log.go:17` [dead_exports/**dead_export**] Exported struct 'AuditLog' is not referenced by any other file in the project
- `internal/model/audit_log.go:28` [dead_exports/**dead_export**] Exported method 'Summary' is not referenced by any other file in the project
- `internal/model/audit_log.go:33` [dead_exports/**dead_export**] Exported method 'IsWrite' is not referenced by any other file in the project
- `internal/model/category.go:4` [dead_code/**unused_import**] import `time` appears unused
- `internal/model/category.go:7` [dead_exports/**dead_export**] Exported struct 'Category' is not referenced by any other file in the project
- `internal/model/category.go:18` [dead_exports/**dead_export**] Exported method 'IsRoot' is not referenced by any other file in the project
- `internal/model/category.go:23` [dead_exports/**dead_export**] Exported method 'IsChildOf' is not referenced by any other file in the project
- `internal/model/category.go:29` [dead_exports/**dead_export**] Exported function 'BuildBreadcrumb' is not referenced by any other file in the project
- `internal/model/order.go:4` [dead_code/**unused_import**] import `time` appears unused
- `internal/model/order.go:7` [dead_exports/**dead_export**] Exported struct 'Order' is not referenced by any other file in the project
- `internal/model/order.go:19` [dead_exports/**dead_export**] Exported struct 'OrderItem' is not referenced by any other file in the project
- `internal/model/order.go:30` [dead_exports/**dead_export**] Exported constant 'OrderStatusPending' is not referenced by any other file in the project
- `internal/model/order.go:31` [dead_exports/**dead_export**] Exported constant 'OrderStatusConfirmed' is not referenced by any other file in the project
- `internal/model/order.go:32` [dead_exports/**dead_export**] Exported constant 'OrderStatusProcessing' is not referenced by any other file in the project
- `internal/model/order.go:33` [dead_exports/**dead_export**] Exported constant 'OrderStatusShipped' is not referenced by any other file in the project
- `internal/model/order.go:34` [dead_exports/**dead_export**] Exported constant 'OrderStatusDelivered' is not referenced by any other file in the project
- `internal/model/order.go:35` [dead_exports/**dead_export**] Exported constant 'OrderStatusCancelled' is not referenced by any other file in the project
- `internal/model/order.go:36` [dead_exports/**dead_export**] Exported constant 'OrderStatusPaymentFailed' is not referenced by any other file in the project
- `internal/model/order.go:40` [dead_exports/**dead_export**] Exported method 'IsTerminal' is not referenced by any other file in the project
- `internal/model/order.go:47` [dead_exports/**dead_export**] Exported method 'ItemCount' is not referenced by any other file in the project
- `internal/model/product.go:5` [dead_code/**unused_import**] import `time` appears unused
- `internal/model/product.go:8` [dead_exports/**dead_export**] Exported struct 'Product' is not referenced by any other file in the project
- `internal/model/product.go:21` [dead_exports/**dead_export**] Exported method 'IsInStock' is not referenced by any other file in the project
- `internal/model/product.go:26` [dead_exports/**dead_export**] Exported method 'PriceFormatted' is not referenced by any other file in the project
- `internal/model/user.go:4` [dead_code/**unused_import**] import `time` appears unused
- `internal/model/user.go:7` [dead_exports/**dead_export**] Exported struct 'User' is not referenced by any other file in the project
- `internal/model/user.go:17` [dead_exports/**dead_export**] Exported method 'DisplayName' is not referenced by any other file in the project
- `internal/repository/audit_repo.go:4` [dead_code/**unused_import**] import `sync` appears unused
- `internal/repository/audit_repo.go:5` [dead_code/**unused_import**] import `time` appears unused
- `internal/repository/audit_repo.go:7` [dead_code/**unused_import**] import `model` appears unused
- `internal/repository/audit_repo.go:10` [dead_exports/**dead_export**] Exported struct 'AuditRepository' is not referenced by any other file in the project
- `internal/repository/audit_repo.go:17` [dead_exports/**dead_export**] Exported function 'NewAuditRepository' is not referenced by any other file in the project
- `internal/repository/audit_repo.go:36` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'ListAll' exported from 4 files: internal/repository/audit_repo.go, internal/repository/category_repo.go, internal/repository/inventory_repo.go, internal/repository/user_repo.go
- `internal/repository/category_repo.go:5` [dead_code/**unused_import**] import `sync` appears unused
- `internal/repository/category_repo.go:7` [dead_code/**unused_import**] import `model` appears unused
- `internal/repository/category_repo.go:10` [dead_exports/**dead_export**] Exported struct 'CategoryRepository' is not referenced by any other file in the project
- `internal/repository/category_repo.go:17` [dead_exports/**dead_export**] Exported function 'NewCategoryRepository' is not referenced by any other file in the project
- `internal/repository/category_repo.go:25` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'Create' exported from 3 files: internal/repository/category_repo.go, internal/repository/order_repo.go, internal/repository/user_repo.go
- `internal/repository/category_repo.go:48` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'ListAll' exported from 4 files: internal/repository/audit_repo.go, internal/repository/category_repo.go, internal/repository/inventory_repo.go, internal/repository/user_repo.go
- `internal/repository/category_repo.go:60` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'Delete' exported from 4 files: internal/repository/category_repo.go, internal/repository/user_repo.go, pkg/cache/lru.go, pkg/cache/redis.go
- `internal/repository/inventory_repo.go:4` [dead_code/**unused_import**] import `sql` appears unused
- `internal/repository/inventory_repo.go:7` [dead_code/**unused_import**] import `model` appears unused
- `internal/repository/inventory_repo.go:10` [dead_exports/**dead_export**] Exported struct 'InventoryRepository' is not referenced by any other file in the project
- `internal/repository/inventory_repo.go:37` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'ListAll' exported from 4 files: internal/repository/audit_repo.go, internal/repository/category_repo.go, internal/repository/inventory_repo.go, internal/repository/user_repo.go
- `internal/repository/order_repo.go:4` [dead_code/**unused_import**] import `sql` appears unused
- `internal/repository/order_repo.go:8` [dead_code/**unused_import**] import `model` appears unused
- `internal/repository/order_repo.go:11` [dead_exports/**dead_export**] Exported struct 'OrderRepository' is not referenced by any other file in the project
- `internal/repository/order_repo.go:21` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'Create' exported from 3 files: internal/repository/category_repo.go, internal/repository/order_repo.go, internal/repository/user_repo.go
- `internal/repository/user_repo.go:4` [dead_code/**unused_import**] import `sql` appears unused
- `internal/repository/user_repo.go:8` [dead_code/**unused_import**] import `model` appears unused
- `internal/repository/user_repo.go:11` [dead_exports/**dead_export**] Exported struct 'UserRepository' is not referenced by any other file in the project
- `internal/repository/user_repo.go:35` [dead_exports/**dead_export**] Exported method 'FindByEmail' is not referenced by any other file in the project
- `internal/repository/user_repo.go:49` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'Create' exported from 3 files: internal/repository/category_repo.go, internal/repository/order_repo.go, internal/repository/user_repo.go
- `internal/repository/user_repo.go:63` [dead_exports/**dead_export**] Exported method 'Update' is not referenced by any other file in the project
- `internal/repository/user_repo.go:72` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'ListAll' exported from 4 files: internal/repository/audit_repo.go, internal/repository/category_repo.go, internal/repository/inventory_repo.go, internal/repository/user_repo.go
- `internal/repository/user_repo.go:93` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'Delete' exported from 4 files: internal/repository/category_repo.go, internal/repository/user_repo.go, pkg/cache/lru.go, pkg/cache/redis.go
- `internal/service/audit.go:8` [dead_code/**unused_import**] import `repository` appears unused
- `internal/service/audit.go:11` [dead_exports/**dead_export**] Exported struct 'AuditService' is not referenced by any other file in the project
- `internal/service/audit.go:16` [dead_exports/**dead_export**] Exported function 'NewAuditService' is not referenced by any other file in the project
- `internal/service/audit.go:21` [dead_exports/**dead_export**] Exported method 'RecordAction' is not referenced by any other file in the project
- `internal/service/audit.go:30` [dead_exports/**dead_export**] Exported method 'GetUserHistory' is not referenced by any other file in the project
- `internal/service/audit.go:39` [dead_exports/**dead_export**] Exported method 'GetRecentActivity' is not referenced by any other file in the project
- `internal/service/audit.go:49` [dead_exports/**dead_export**] Exported method 'Count' is not referenced by any other file in the project
- `internal/service/category.go:7` [dead_code/**unused_import**] import `model` appears unused
- `internal/service/category.go:8` [dead_code/**unused_import**] import `repository` appears unused
- `internal/service/category.go:11` [dead_exports/**dead_export**] Exported struct 'CategoryService' is not referenced by any other file in the project
- `internal/service/category.go:16` [dead_exports/**dead_export**] Exported function 'NewCategoryService' is not referenced by any other file in the project
- `internal/service/category.go:21` [dead_exports/**dead_export**] Exported method 'CreateCategory' is not referenced by any other file in the project
- `internal/service/category.go:43` [dead_exports/**dead_export**] Exported method 'GetCategory' is not referenced by any other file in the project
- `internal/service/category.go:52` [dead_exports/**dead_export**] Exported method 'ListRootCategories' is not referenced by any other file in the project
- `internal/service/category.go:57` [dead_exports/**dead_export**] Exported method 'ListChildren' is not referenced by any other file in the project
- `internal/service/category.go:62` [dead_exports/**dead_export**] Exported method 'DeleteCategory' is not referenced by any other file in the project
- `internal/service/inventory.go:7` [dead_code/**unused_import**] import `model` appears unused
- `internal/service/inventory.go:8` [dead_code/**unused_import**] import `repository` appears unused
- `internal/service/inventory.go:11` [dead_exports/**dead_export**] Exported struct 'InventoryService' is not referenced by any other file in the project
- `internal/service/inventory.go:45` [dead_exports/**dead_export**] Exported method 'BulkUpdateStock' is not referenced by any other file in the project
- `internal/service/inventory.go:56` [dead_exports/**dead_export**] Exported method 'SearchProducts' is not referenced by any other file in the project
- `internal/service/inventory.go:61` [dead_exports/**dead_export**] Exported method 'GetLowStockProducts' is not referenced by any other file in the project
- `internal/service/notification.go:34` [dead_exports/**dead_export**] Exported method 'SendShippingUpdate' is not referenced by any other file in the project
- `internal/service/notification.go:80` [dead_exports/**dead_export**] Exported method 'ScheduleReminder' is not referenced by any other file in the project
- `internal/service/order.go:9` [dead_code/**unused_import**] import `model` appears unused
- `internal/service/order.go:10` [dead_code/**unused_import**] import `repository` appears unused
- `internal/service/order.go:174` [dead_exports/**dead_export**] Exported method 'ProcessExpiredOrders' is not referenced by any other file in the project
- `internal/service/payment.go:13` [dead_exports/**dead_export**] Exported struct 'PaymentResult' is not referenced by any other file in the project
- `internal/service/payment.go:20` [dead_exports/**dead_export**] Exported struct 'PaymentService' is not referenced by any other file in the project
- `internal/service/payment.go:100` [dead_exports/**dead_export**] Exported method 'ValidatePaymentMethod' is not referenced by any other file in the project
- `internal/service/shipping.go:7` [dead_code/**unused_import**] import `model` appears unused
- `internal/service/shipping.go:10` [dead_exports/**dead_export**] Exported struct 'ShippingRate' is not referenced by any other file in the project
- `internal/service/shipping.go:17` [dead_exports/**dead_export**] Exported struct 'ShippingService' is not referenced by any other file in the project
- `internal/service/shipping.go:24` [dead_exports/**dead_export**] Exported function 'NewShippingService' is not referenced by any other file in the project
- `internal/service/shipping.go:33` [dead_exports/**dead_export**] Exported method 'CalculateRate' is not referenced by any other file in the project
- `internal/service/shipping.go:60` [dead_exports/**dead_export**] Exported method 'ValidateAddress' is not referenced by any other file in the project
- `internal/service/shipping.go:61` [coupling/**low_cohesion**] method `ValidateAddress` does not use its receiver `s` — consider making it a function
- `internal/service/shipping.go:76` [dead_exports/**dead_export**] Exported method 'FormatLabel' is not referenced by any other file in the project
- `internal/service/shipping.go:77` [coupling/**low_cohesion**] method `FormatLabel` does not use its receiver `s` — consider making it a function
- `internal/worker/cleanup.go:4` [dead_code/**unused_import**] import `context` appears unused
- `internal/worker/cleanup.go:6` [dead_code/**unused_import**] import `sync` appears unused
- `internal/worker/cleanup.go:11` [dead_exports/**dead_export**] Exported type_alias 'CleanupFunc' is not referenced by any other file in the project
- `internal/worker/cleanup.go:14` [dead_exports/**dead_export**] Exported struct 'CleanupWorker' is not referenced by any other file in the project
- `internal/worker/cleanup.go:25` [dead_exports/**dead_export**] Exported function 'NewCleanupWorker' is not referenced by any other file in the project
- `internal/worker/cleanup.go:34` [dead_exports/**dead_export**] Exported method 'Run' is not referenced by any other file in the project
- `internal/worker/cleanup.go:69` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'Stats' exported from 3 files: internal/worker/cleanup.go, internal/worker/dispatcher.go, pkg/queue/memory.go
- `internal/worker/dispatcher.go:5` [dead_code/**unused_import**] import `sync` appears unused
- `internal/worker/dispatcher.go:11` [dead_exports/**dead_export**] Exported struct 'Job' is not referenced by any other file in the project
- `internal/worker/dispatcher.go:20` [dead_exports/**dead_export**] Exported struct 'JobResult' is not referenced by any other file in the project
- `internal/worker/dispatcher.go:27` [dead_exports/**dead_export**] Exported struct 'DispatcherStats' is not referenced by any other file in the project
- `internal/worker/dispatcher.go:34` [dead_exports/**dead_export**] Exported struct 'Dispatcher' is not referenced by any other file in the project
- `internal/worker/dispatcher.go:81` [dead_exports/**dead_export**] Exported method 'ProcessJobs' is not referenced by any other file in the project
- `internal/worker/processor.go:9` [dead_exports/**dead_export**] Exported struct 'Processor' is not referenced by any other file in the project
- `internal/worker/processor.go:15` [dead_exports/**dead_export**] Exported type_alias 'JobHandler' is not referenced by any other file in the project
- `internal/worker/processor.go:18` [dead_exports/**dead_export**] Exported function 'NewProcessor' is not referenced by any other file in the project
- `internal/worker/processor.go:34` [dead_exports/**dead_export**] Exported method 'Register' is not referenced by any other file in the project
- `internal/worker/processor.go:41` [dead_exports/**dead_export**] Exported method 'Process' is not referenced by any other file in the project
- `internal/worker/processor.go:77` [dead_exports/**dead_export**] Exported method 'ProcessBatch' is not referenced by any other file in the project
- `internal/worker/reporter.go:4` [dead_code/**unused_import**] import `context` appears unused
- `internal/worker/reporter.go:6` [dead_code/**unused_import**] import `sync` appears unused
- `internal/worker/reporter.go:10` [dead_exports/**dead_export**] Exported struct 'MetricSnapshot' is not referenced by any other file in the project
- `internal/worker/reporter.go:18` [dead_exports/**dead_export**] Exported struct 'Reporter' is not referenced by any other file in the project
- `internal/worker/reporter.go:29` [dead_exports/**dead_export**] Exported function 'NewReporter' is not referenced by any other file in the project
- `internal/worker/reporter.go:37` [dead_exports/**dead_export**] Exported method 'RecordRequest' is not referenced by any other file in the project
- `internal/worker/reporter.go:50` [dead_exports/**dead_export**] Exported method 'Run' is not referenced by any other file in the project
- `internal/worker/reporter.go:94` [dead_exports/**dead_export**] Exported method 'Snapshots' is not referenced by any other file in the project
- `internal/worker/scheduler.go:11` [dead_exports/**dead_export**] Exported struct 'ScheduledTask' is not referenced by any other file in the project
- `internal/worker/scheduler.go:19` [dead_exports/**dead_export**] Exported struct 'Scheduler' is not referenced by any other file in the project
- `internal/worker/scheduler.go:94` [dead_exports/**dead_export**] Exported method 'AddTask' is not referenced by any other file in the project
- `internal/worker/scheduler.go:100` [dead_exports/**dead_export**] Exported method 'RemoveTask' is not referenced by any other file in the project
- `internal/worker/scheduler.go:110` [dead_exports/**dead_export**] Exported method 'ListTasks' is not referenced by any other file in the project
- `pkg/cache/cache.go:7` [dead_exports/**dead_export**] Exported interface 'Cache' is not referenced by any other file in the project
- `pkg/cache/cache.go:15` [dead_exports/**dead_export**] Exported struct 'CacheEntry' is not referenced by any other file in the project
- `pkg/cache/lru.go:5` [dead_code/**unused_import**] import `sync` appears unused
- `pkg/cache/lru.go:16` [dead_exports/**dead_export**] Exported struct 'LRUCache' is not referenced by any other file in the project
- `pkg/cache/lru.go:24` [dead_exports/**dead_export**] Exported function 'NewLRUCache' is not referenced by any other file in the project
- `pkg/cache/lru.go:80` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'Delete' exported from 4 files: internal/repository/category_repo.go, internal/repository/user_repo.go, pkg/cache/lru.go, pkg/cache/redis.go
- `pkg/cache/lru.go:105` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'Len' exported from 3 files: pkg/cache/lru.go, pkg/queue/priority.go, pkg/queue/priority.go
- `pkg/cache/redis.go:10` [dead_exports/**dead_export**] Exported struct 'RedisCache' is not referenced by any other file in the project
- `pkg/cache/redis.go:55` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'Delete' exported from 4 files: internal/repository/category_repo.go, internal/repository/user_repo.go, pkg/cache/lru.go, pkg/cache/redis.go
- `pkg/cache/redis.go:76` [dead_exports/**dead_export**] Exported method 'Size' is not referenced by any other file in the project
- `pkg/cache/redis.go:82` [dead_exports/**dead_export**] Exported method 'Flush' is not referenced by any other file in the project
- `pkg/logger/formatter.go:6` [dead_code/**unused_import**] import `strings` appears unused
- `pkg/logger/formatter.go:14` [dead_exports/**dead_export**] Exported constant 'FormatText' is not referenced by any other file in the project
- `pkg/logger/formatter.go:16` [dead_exports/**dead_export**] Exported constant 'FormatJSON' is not referenced by any other file in the project
- `pkg/logger/formatter.go:20` [dead_exports/**dead_export**] Exported struct 'Formatter' is not referenced by any other file in the project
- `pkg/logger/formatter.go:26` [dead_exports/**dead_export**] Exported function 'NewFormatter' is not referenced by any other file in the project
- `pkg/logger/formatter.go:34` [dead_exports/**dead_export**] Exported method 'SetTimestamp' is not referenced by any other file in the project
- `pkg/logger/formatter.go:39` [dead_exports/**dead_export**] Exported method 'FormatLine' is not referenced by any other file in the project
- `pkg/logger/level.go:8` [dead_exports/**dead_export**] Exported type_alias 'Level' is not referenced by any other file in the project
- `pkg/logger/level.go:11` [dead_exports/**dead_export**] Exported constant 'LevelDebug' is not referenced by any other file in the project
- `pkg/logger/level.go:12` [dead_exports/**dead_export**] Exported constant 'LevelInfo' is not referenced by any other file in the project
- `pkg/logger/level.go:13` [dead_exports/**dead_export**] Exported constant 'LevelWarn' is not referenced by any other file in the project
- `pkg/logger/level.go:14` [dead_exports/**dead_export**] Exported constant 'LevelError' is not referenced by any other file in the project
- `pkg/logger/level.go:35` [dead_exports/**dead_export**] Exported function 'ParseLevel' is not referenced by any other file in the project
- `pkg/logger/level.go:38` [duplicate_code/**duplicate_switch_case**] switch case at line 38 has a duplicate body (switch at line 37)
- `pkg/logger/level.go:40` [duplicate_code/**duplicate_switch_case**] switch case at line 40 has a duplicate body (switch at line 37)
- `pkg/logger/level.go:42` [duplicate_code/**duplicate_switch_case**] switch case at line 42 has a duplicate body (switch at line 37)
- `pkg/logger/level.go:51` [dead_exports/**dead_export**] Exported function 'ShouldLog' is not referenced by any other file in the project
- `pkg/logger/level.go:56` [dead_exports/**dead_export**] Exported function 'AllLevels' is not referenced by any other file in the project
- `pkg/logger/logger.go:10` [dead_exports/**dead_export**] Exported struct 'Logger' is not referenced by any other file in the project
- `pkg/logger/logger.go:29` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'Error' exported from 3 files: internal/config/config.go, internal/config/validate.go, pkg/logger/logger.go
- `pkg/logger/logger.go:34` [dead_exports/**dead_export**] Exported method 'Debug' is not referenced by any other file in the project
- `pkg/logger/logger.go:42` [dead_exports/**dead_export**] Exported method 'Warn' is not referenced by any other file in the project
- `pkg/middleware/recovery.go:11` [dead_exports/**dead_export**] Exported function 'Recovery' is not referenced by any other file in the project
- `pkg/middleware/recovery.go:30` [dead_exports/**dead_export**] Exported function 'RecoveryWithHandler' is not referenced by any other file in the project
- `pkg/middleware/timeout.go:6` [dead_code/**unused_import**] import `time` appears unused
- `pkg/middleware/timeout.go:11` [dead_exports/**dead_export**] Exported function 'Timeout' is not referenced by any other file in the project
- `pkg/middleware/timeout.go:55` [dead_exports/**dead_export**] Exported method 'Write' is not referenced by any other file in the project
- `pkg/queue/memory.go:6` [dead_code/**unused_import**] import `sync` appears unused
- `pkg/queue/memory.go:9` [dead_exports/**dead_export**] Exported struct 'MemoryQueue' is not referenced by any other file in the project
- `pkg/queue/memory.go:56` [dead_exports/**dead_export**] Exported method 'Size' is not referenced by any other file in the project
- `pkg/queue/memory.go:61` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'Stats' exported from 3 files: internal/worker/cleanup.go, internal/worker/dispatcher.go, pkg/queue/memory.go
- `pkg/queue/memory.go:69` [dead_exports/**dead_export**] Exported method 'Drain' is not referenced by any other file in the project
- `pkg/queue/memory.go:82` [dead_exports/**dead_export**] Exported method 'EnqueueBatch' is not referenced by any other file in the project
- `pkg/queue/priority.go:5` [dead_code/**unused_import**] import `sync` appears unused
- `pkg/queue/priority.go:8` [dead_exports/**dead_export**] Exported struct 'PriorityItem' is not referenced by any other file in the project
- `pkg/queue/priority.go:17` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'Len' exported from 3 files: pkg/cache/lru.go, pkg/queue/priority.go, pkg/queue/priority.go
- `pkg/queue/priority.go:18` [dead_exports/**dead_export**] Exported method 'Less' is not referenced by any other file in the project
- `pkg/queue/priority.go:19` [dead_exports/**dead_export**] Exported method 'Swap' is not referenced by any other file in the project
- `pkg/queue/priority.go:26` [dead_exports/**dead_export**] Exported method 'Push' is not referenced by any other file in the project
- `pkg/queue/priority.go:33` [dead_exports/**dead_export**] Exported method 'Pop' is not referenced by any other file in the project
- `pkg/queue/priority.go:44` [dead_exports/**dead_export**] Exported struct 'PriorityQueue' is not referenced by any other file in the project
- `pkg/queue/priority.go:50` [dead_exports/**dead_export**] Exported function 'NewPriorityQueue' is not referenced by any other file in the project
- `pkg/queue/priority.go:85` [cross_file_duplicates/**cross_file_duplicate**] Cross-file duplicate: method 'Len' exported from 3 files: pkg/cache/lru.go, pkg/queue/priority.go, pkg/queue/priority.go
- `pkg/queue/priority.go:92` [dead_exports/**dead_export**] Exported method 'IsEmpty' is not referenced by any other file in the project
- `pkg/queue/queue.go:3` [dead_exports/**dead_export**] Exported interface 'Queue' is not referenced by any other file in the project

### unknown (36 extra findings)

- `cmd/migrate/main.go:77` [error_swallowing/**error_swallowed**] error return value discarded with blank identifier `_`
- `cmd/server/main.go:114` [goroutine_leak/**goroutine_missing_done_channel**] goroutine with for-loop but no select+ctx.Done() — may leak
- `internal/api/handler.go:146` [error_swallowing/**error_swallowed**] error return value discarded with blank identifier `_`
- `internal/api/handler.go:148` [error_swallowing/**error_swallowed**] error return value discarded with blank identifier `_`
- `internal/api/handler.go:152` [error_swallowing/**error_swallowed**] error return value discarded with blank identifier `_`
- `internal/api/handler.go:208` [error_swallowing/**error_swallowed**] error return value discarded with blank identifier `_`
- `internal/api/response.go:27` [naked_interface/**empty_interface_param**] empty interface (`interface{}` / `any`) loses type safety — consider a concrete interface
- `internal/api/response.go:34` [naked_interface/**empty_interface_param**] empty interface (`interface{}` / `any`) loses type safety — consider a concrete interface
- `internal/repository/inventory_repo.go:72` [error_swallowing/**error_swallowed**] error return value discarded with blank identifier `_`
- `internal/repository/inventory_repo.go:90` [error_swallowing/**error_swallowed**] error return value discarded with blank identifier `_`
- `internal/repository/inventory_repo.go:99` [error_swallowing/**error_swallowed**] error return value discarded with blank identifier `_`
- `internal/repository/order_repo.go:40` [error_swallowing/**error_swallowed**] error return value discarded with blank identifier `_`
- `internal/repository/order_repo.go:124` [error_swallowing/**error_swallowed**] error return value discarded with blank identifier `_`
- `internal/repository/user_repo.go:65` [error_swallowing/**error_swallowed**] error return value discarded with blank identifier `_`
- `internal/repository/user_repo.go:95` [error_swallowing/**error_swallowed**] error return value discarded with blank identifier `_`
- `internal/service/notification.go:85` [error_swallowing/**error_swallowed**] error return value discarded with blank identifier `_`
- `internal/service/order.go:159` [type_confusion/**unguarded_type_assertion**] type assertion without comma-ok check may panic at runtime
- `internal/service/order.go:162` [type_confusion/**unguarded_type_assertion**] type assertion without comma-ok check may panic at runtime
- `internal/service/order.go:164` [type_confusion/**unguarded_type_assertion**] type assertion without comma-ok check may panic at runtime
- `internal/service/order.go:166` [type_confusion/**unguarded_type_assertion**] type assertion without comma-ok check may panic at runtime
- `internal/service/order.go:182` [error_swallowing/**error_swallowed**] error return value discarded with blank identifier `_`
- `internal/service/payment.go:92` [error_swallowing/**error_swallowed**] error return value discarded with blank identifier `_`
- `internal/worker/dispatcher.go:83` [race_conditions/**loop_var_capture**] goroutine spawned inside for loop — loop variable may be captured by reference
- `internal/worker/processor.go:81` [race_conditions/**loop_var_capture**] goroutine spawned inside for loop — loop variable may be captured by reference
- `internal/worker/scheduler.go:61` [race_conditions/**loop_var_capture**] goroutine spawned inside for loop — loop variable may be captured by reference
- `pkg/cache/cache.go:9` [naked_interface/**empty_interface_param**] empty interface (`interface{}` / `any`) loses type safety — consider a concrete interface
- `pkg/cache/cache.go:10` [naked_interface/**empty_interface_param**] empty interface (`interface{}` / `any`) loses type safety — consider a concrete interface
- `pkg/cache/cache.go:17` [naked_interface/**empty_interface_field**] empty interface (`interface{}` / `any`) loses type safety — consider a concrete interface
- `pkg/cache/lru.go:12` [naked_interface/**empty_interface_field**] empty interface (`interface{}` / `any`) loses type safety — consider a concrete interface
- `pkg/cache/lru.go:35` [naked_interface/**empty_interface_param**] empty interface (`interface{}` / `any`) loses type safety — consider a concrete interface
- `pkg/cache/lru.go:55` [naked_interface/**empty_interface_param**] empty interface (`interface{}` / `any`) loses type safety — consider a concrete interface
- `pkg/middleware/recovery.go:31` [naked_interface/**empty_interface_param**] empty interface (`interface{}` / `any`) loses type safety — consider a concrete interface
- `pkg/queue/priority.go:10` [naked_interface/**empty_interface_field**] empty interface (`interface{}` / `any`) loses type safety — consider a concrete interface
- `pkg/queue/priority.go:27` [naked_interface/**empty_interface_param**] empty interface (`interface{}` / `any`) loses type safety — consider a concrete interface
- `pkg/queue/priority.go:60` [naked_interface/**empty_interface_param**] empty interface (`interface{}` / `any`) loses type safety — consider a concrete interface
- `pkg/queue/queue.go:5` [naked_interface/**empty_interface_param**] empty interface (`interface{}` / `any`) loses type safety — consider a concrete interface

## Detected Instances (matched)

These manifest entries were successfully matched to a virgil-cli finding.


**code-quality**
- `cmd/server/main.go:28` **god-functions** — main() is a 127-line god function handling server setup, routing, middleware, worker init, graceful shutdown, and health checks
- `cmd/server/main.go:39` **magic-numbers** — db.SetMaxOpenConns(25) and SetConnMaxLifetime(300) -- hardcoded connection pool parameters without named constants
- `cmd/server/main.go:105` **magic-numbers** — ReadTimeout/WriteTimeout/IdleTimeout use hardcoded 30/30/120 second values without named constants
- `internal/api/handler.go:44` **god-functions** — CreateOrder() is a 122-line handler doing auth extraction, input parsing, validation, business logic, payment, notification, and response formatting
- `internal/api/handler.go:44` **function-length** — CreateOrder() handler is 122 lines, far exceeding reasonable function length
- `internal/api/handler.go:116` **stringly-typed** — Order status uses string literals ('pending', 'payment_failed', 'confirmed') instead of typed constants defined in model
- `internal/api/handler.go:212` **magic-numbers** — limit := 20 -- hardcoded page size without named constant
- `internal/config/config.go:44` **magic-numbers** — Load() contains 6+ hardcoded default values (8080, 25, 5, 300, 10, 587) without named constants
- `internal/service/order.go:136` **deep-nesting** — _ = err error swallowed inside nested for-loop/if-block during stock restore, partial rollback on failure
- `internal/service/order.go:145` **magic-numbers** — FindByUserID(userID, 10000, 0) -- hardcoded 10000 limit to fetch all orders, no real pagination
- `internal/worker/processor.go:42` **deep-nesting** — Process() has 4 levels of nesting: for > if > if > if/else -- retry logic with deeply nested control flow
- `internal/worker/processor.go:50` **magic-numbers** — attempt < 3 and timeout: 30 -- hardcoded retry count and timeout without named constants

**scalability**
- `internal/api/middleware.go:80` **memory-leak-indicators** — Rate limit cleanup goroutine started with no cancellation mechanism -- leaks for the lifetime of the process
- `internal/service/notification.go:26` **sync-blocking-in-async** — SendOrderConfirmation() makes a blocking SMTP call in the request handler path
- `internal/service/notification.go:62` **sync-blocking-in-async** — SendBulkNotification() sends emails sequentially in a loop, blocking the caller for all SMTP round trips combined
- `internal/service/order.go:38` **n-plus-one-queries** — CreateOrder() fetches each product price individually in a loop instead of batch query
- `internal/service/order.go:100` **n-plus-one-queries** — ListOrders() fetches items separately for each order in a loop
- `internal/service/payment.go:44` **sync-blocking-in-async** — ProcessPayment() makes a blocking HTTP POST to external payment gateway in the request handler path, blocking goroutine for up to 30s
- `internal/service/payment.go:77` **sync-blocking-in-async** — RefundPayment() makes a blocking HTTP POST to external payment gateway, same blocking pattern as ProcessPayment
- `internal/worker/dispatcher.go:82` **resource-exhaustion** — ProcessJobs() spawns unbounded goroutines per batch with no concurrency limit, can exhaust memory with large batches
- `internal/worker/dispatcher.go:93` **memory-leak-indicators** — ProcessJobs() spawns one goroutine per job without context or WaitGroup -- goroutines leak if dispatcher stops
- `internal/worker/processor.go:79` **memory-leak-indicators** — ProcessBatch() creates done channel with hardcoded buffer of 10, causing goroutine leak when len(jobs) > 10
- `internal/worker/scheduler.go:60` **memory-leak-indicators** — RunForever() spawns goroutines per scheduled task with no WaitGroup or context -- goroutines leak indefinitely on shutdown
- `pkg/cache/redis.go:33` **memory-leak-indicators** — Get() reads from shared map without mutex -- concurrent map read causes panic under load
- `pkg/cache/redis.go:49` **memory-leak-indicators** — Set() writes to shared map without mutex -- concurrent map write causes panic under load
- `pkg/queue/memory.go:32` **sync-blocking-in-async** — Enqueue() blocks indefinitely on full channel with no timeout or context cancellation, causing goroutine pileup

**style**
- `internal/api/handler.go:158` **duplicate-code** — Response construction pattern (Content-Type + WriteHeader + json.NewEncoder) duplicated across CreateOrder, GetOrder, and ListOrders handlers
- `internal/service/notification.go:81` **dead-code** — ScheduleReminder() is defined but never called from any handler or worker
- `internal/service/order.go:175` **dead-code** — ProcessExpiredOrders() and calculateDiscount() are defined but never called from any handler or worker
- `internal/service/payment.go:101` **dead-code** — ValidatePaymentMethod() and retryPayment() are defined but never called from any service or handler
- `internal/worker/dispatcher.go:133` **dead-code** — collectResults() method is defined but never called from Start() or any external module
- `pkg/queue/memory.go:70` **dead-code** — Drain() and logQueueEvent() are defined but never called from any module

**tech-debt**
- `cmd/server/main.go:28` **legacy-pattern** — HTTP handler setup uses deprecated http.Get() without context propagation; should use http.NewRequestWithContext() to allow proper cancellation and deadline propagation
- `internal/api/handler.go:44` **deprecated-api-usage** — Uses ioutil.ReadAll() (deprecated since Go 1.16); should use io.ReadAll() from standard library -- deprecated package still imported across multiple handlers

---
_Report generated by virgil-benchmarks test harness_
