# Changelog

## [1.0.1](https://github.com/samuel032khoury/gopherfeed/compare/v1.0.0...v1.0.1) (2026-01-16)


### Bug Fixes

* correct version extraction regex in update-api-version workflow ([964b1a3](https://github.com/samuel032khoury/gopherfeed/commit/964b1a3705ebc948376ad50774d6a7225ab41b78))
* set version to 1.0.0 in main.go ([af31dd6](https://github.com/samuel032khoury/gopherfeed/commit/af31dd637709d51942c59a52c3036b88a719964d))

## 1.0.0 (2026-01-16)


### Features

* implement basic API server with health check endpoint ([84ea2e9](https://github.com/samuel032khoury/gopherfeed/commit/84ea2e9ac1e4259bc032f07113113a5c1ebb8ada))
* add database connection and storage implementations ([7d52d11](https://github.com/samuel032khoury/gopherfeed/commit/7d52d1197f3e05084dff3c9f9d49e78d2d8889b7))
* add database configuration and connection handling ([212085c](https://github.com/samuel032khoury/gopherfeed/commit/212085c4aad505ead65c4f4e2969f0e96879dbbe))
* add Docker Compose configuration and database initialization ([78cc07d](https://github.com/samuel032khoury/gopherfeed/commit/78cc07d9a2e2bd4174e8cc73e67a27c5eed64b2e))
* add database migration scripts for users and posts tables ([b7cdcf3](https://github.com/samuel032khoury/gopherfeed/commit/b7cdcf3248297e1524cb243fccd3ec30ca5a944c))
* add create and get handlers ([f9c790a](https://github.com/samuel032khoury/gopherfeed/commit/f9c790a8a9bce488fbe6c597c86f88540a0156fe))
* add comments functionality with CRUD operations ([e583576](https://github.com/samuel032khoury/gopherfeed/commit/e5835760f1f8ffe2699465770608a2b3ace2a8ee))
* add user retrieval functionality and database seeding script ([b0d8b7d](https://github.com/samuel032khoury/gopherfeed/commit/b0d8b7d0f281565529f0cb99788455f618cc8944))
* implement user follow/unfollow functionality with middleware ([7f4fdd5](https://github.com/samuel032khoury/gopherfeed/commit/7f4fdd5608b633b1913ffdfcc1d3ac6e2ad9147a))
* add indexes for comments, posts, and users ([6474d85](https://github.com/samuel032khoury/gopherfeed/commit/6474d850883607f22bc2fa1a8a2bfb61b0a8a955))
* add feed retrieval functionality with handler and storage method ([064c8a3](https://github.com/samuel032khoury/gopherfeed/commit/064c8a3d6b92957fb13d354784c699edc3fbd90d))
* implement pagination for feed retrieval with validation ([f3e8125](https://github.com/samuel032khoury/gopherfeed/commit/f3e8125be53e4dc7c27ce1f247376f8e0d42a929))
* add swagger documentation ([704e102](https://github.com/samuel032khoury/gopherfeed/commit/704e102df82bef3d2024db658972cd0c4b79e8a6))
* add structured logging ([808635d](https://github.com/samuel032khoury/gopherfeed/commit/808635d8e3aa5dce4c01d61b24896674ae34da3f))
* add user authentication endpoints and invitation functionality ([4850aa6](https://github.com/samuel032khoury/gopherfeed/commit/4850aa635f57acce27029bb2056cd8c5d12b0ed2))
* add account activation functionality and documentation ([a8c5225](https://github.com/samuel032khoury/gopherfeed/commit/a8c5225de5e202450f39e99ab657fe78a061d109))
* add email invitation functionality for user registration ([0c3cc7c](https://github.com/samuel032khoury/gopherfeed/commit/0c3cc7c5d582062b5ceaf9e99988b328d198044e))
* implement asynchronous email processing with RabbitMQ integration ([c578328](https://github.com/samuel032khoury/gopherfeed/commit/c5783285d9b20eef27346fe7ccaee214f678fccb))
* initialize web application with React, TypeScript, and Vite ([af75458](https://github.com/samuel032khoury/gopherfeed/commit/af754587eafdba695fe1f8b7afc23f7cafa3aed6))
* implement JWT authentication and basic auth middleware ([fc350dc](https://github.com/samuel032khoury/gopherfeed/commit/fc350dca813ebdf495ba132614d42aa438bf6a87))
* implement role-based access control and user roles management ([fbf0356](https://github.com/samuel032khoury/gopherfeed/commit/fbf03561824dd61e515fc128afaa59eaedd2048e))
* add Redis caching for user data and update user role handling ([dee6d78](https://github.com/samuel032khoury/gopherfeed/commit/dee6d783fb0e4ea252e99851fe8e68ac1838fefc))
* implement structured logging across the application and refactor email publishing and consumption ([90220ce](https://github.com/samuel032khoury/gopherfeed/commit/90220cea5655c49dc67f7952a3f1879e40137e61))
* implement graceful shutdown for the API server and adjust kill delay in configuration ([6497b52](https://github.com/samuel032khoury/gopherfeed/commit/6497b525d7a355abca6e35e6739965bb5845d467))
* implement rate limiting middleware and error handling for rate limit exceeded ([383b69a](https://github.com/samuel032khoury/gopherfeed/commit/383b69aed87f6b419825d688ddf9cce123d5ccbe))
* add logout endpoint to authentication API with JWT cookie clearing ([9657679](https://github.com/samuel032khoury/gopherfeed/commit/96576792ad6e64e4c2b639ac49b9db257cb38ca3))
* implement logout functionality and update authentication flow with cookie management ([ba7cb42](https://github.com/samuel032khoury/gopherfeed/commit/ba7cb423760a297f988f4151de91d27659372db6))
* add audience (aud) field to JWT configuration and update authenticator ([e0da229](https://github.com/samuel032khoury/gopherfeed/commit/e0da22967c3a2637ee48c17b0843d2922117ac08))
* prevent users from following or unfollowing themselves ([9f0e930](https://github.com/samuel032khoury/gopherfeed/commit/9f0e930e134b1ab30c923c9b7d3d67c1bebe706c))
* replace custom logger with zap.SugaredLogger across email and mq packages ([d3f0280](https://github.com/samuel032khoury/gopherfeed/commit/d3f0280a5e071d77d8a1ee33a107a3e8cb3aa240))
* add unit tests for feed retrieval and implement mock stores for testing ([9a1e6e4](https://github.com/samuel032khoury/gopherfeed/commit/9a1e6e4ef4b957c00f3a1ea9c11b647b180bf50d))
* add expvar support for application statistics and health checks ([7d01607](https://github.com/samuel032khoury/gopherfeed/commit/7d01607457f74008b55f14430eeef2a33f0ff902))
* add automation workflow ([027aa14](https://github.com/samuel032khoury/gopherfeed/commit/027aa147aaac860219bea3cd548b6f816ff94505))
* update CI workflow to use ubuntu-latest for improved compatibility ([0775b8e](https://github.com/samuel032khoury/gopherfeed/commit/0775b8e9c05e191b818640c14f068586ad04980c))
* add LICENSE and README files with project details and setup instructions ([20d3b43](https://github.com/samuel032khoury/gopherfeed/commit/20d3b4315a058dafdaf58089c73ce502dd245f9b))
* add GitHub Actions workflow for automated releases ([8944bb0](https://github.com/samuel032khoury/gopherfeed/commit/8944bb068addfcbbcd7f3b0cca520c6045797147))


### Bug Fixes

* correct error message casing for Mailtrap credentials ([8ddefd5](https://github.com/samuel032khoury/gopherfeed/commit/8ddefd56fa86246cf958975b4218c161ea139b37))
* set version to 1.0.0 in main.go ([af31dd6](https://github.com/samuel032khoury/gopherfeed/commit/af31dd637709d51942c59a52c3036b88a719964d))
* correct version extraction regex in update-api-version workflow ([964b1a3](https://github.com/samuel032khoury/gopherfeed/commit/964b1a3705ebc948376ad50774d6a7225ab41b78))


### Refactoring

* refactor health check handler to return structured JSON ([3e51875](https://github.com/samuel032khoury/gopherfeed/commit/3e5187528b3f8e48502304270a3bbb813876bc96))
* refactor error handling in API handlers with centralized functions ([6360aef](https://github.com/samuel032khoury/gopherfeed/commit/6360aef4596356e14f13e54bc5ed836ad565c2d2))
* rename follower_id to followee_id in followers table ([5377f3b](https://github.com/samuel032khoury/gopherfeed/commit/5377f3b03aa9982a155015f729a76b8476d03a13))
* refactor database context handling with common timeout function ([cc40912](https://github.com/samuel032khoury/gopherfeed/commit/cc40912a93dc74be1c5e7d0a224a71218e885ffe))
* refactor mailer package and add async mailer with RabbitMQ ([6a231f4](https://github.com/samuel032khoury/gopherfeed/commit/6a231f4f1071fbbe3e416c9f39e214a1eaa52332))
* refactor email handling to use publisher model with RabbitMQ ([f5492b8](https://github.com/samuel032khoury/gopherfeed/commit/f5492b8692006131e42f5e52bbc23c79f51dbbfb))
* streamline configuration loading and enhance code organization ([a53d5d5](https://github.com/samuel032khoury/gopherfeed/commit/a53d5d587356534143833e30c057a6b9a0e92041))
* simplify middleware function signatures and improve context handling ([9655d30](https://github.com/samuel032khoury/gopherfeed/commit/9655d30495adf79ec1a28c572001be997cd52948))


### Chore

* update kill_delay and send_interrupt settings ([151cf47](https://github.com/samuel032khoury/gopherfeed/commit/151cf4765e5d7269226c73ee2595c622689cb2f1))
* enable main log only output in configuration files ([4885e64](https://github.com/samuel032khoury/gopherfeed/commit/4885e64a98faf99c01c0b309f86e0b72bb9dcfa9))
* add GitHub Actions workflow for automatic version updates ([c4fe8a2](https://github.com/samuel032khoury/gopherfeed/commit/c4fe8a28676eae1e80c97059918d8c8da71643bd))
