# Changelog

## [1.0.1](https://github.com/samuel032khoury/gopherfeed/compare/v1.0.0...v1.0.1) (2026-01-16)


### Bug Fixes

* correct version extraction regex in update-api-version workflow ([964b1a3](https://github.com/samuel032khoury/gopherfeed/commit/964b1a3705ebc948376ad50774d6a7225ab41b78))
* set version to 1.0.0 in main.go ([af31dd6](https://github.com/samuel032khoury/gopherfeed/commit/af31dd637709d51942c59a52c3036b88a719964d))

## 1.0.0 (2026-01-16)


### Features

* add audience (aud) field to JWT configuration and update authenticator ([e0da229](https://github.com/samuel032khoury/gopherfeed/commit/e0da22967c3a2637ee48c17b0843d2922117ac08))
* add automation workflow ([027aa14](https://github.com/samuel032khoury/gopherfeed/commit/027aa147aaac860219bea3cd548b6f816ff94505))
* add expvar support for application statistics and health checks ([7d01607](https://github.com/samuel032khoury/gopherfeed/commit/7d01607457f74008b55f14430eeef2a33f0ff902))
* add GitHub Actions workflow for automated releases ([8944bb0](https://github.com/samuel032khoury/gopherfeed/commit/8944bb068addfcbbcd7f3b0cca520c6045797147))
* add LICENSE and README files with project details and setup instructions ([20d3b43](https://github.com/samuel032khoury/gopherfeed/commit/20d3b4315a058dafdaf58089c73ce502dd245f9b))
* add logout endpoint to authentication API with JWT cookie clearing ([9657679](https://github.com/samuel032khoury/gopherfeed/commit/96576792ad6e64e4c2b639ac49b9db257cb38ca3))
* add Redis caching for user data and update user role handling ([dee6d78](https://github.com/samuel032khoury/gopherfeed/commit/dee6d783fb0e4ea252e99851fe8e68ac1838fefc))
* add unit tests for feed retrieval and implement mock stores for testing ([9a1e6e4](https://github.com/samuel032khoury/gopherfeed/commit/9a1e6e4ef4b957c00f3a1ea9c11b647b180bf50d))
* implement graceful shutdown for the API server and adjust kill delay in configuration ([6497b52](https://github.com/samuel032khoury/gopherfeed/commit/6497b525d7a355abca6e35e6739965bb5845d467))
* implement JWT authentication and basic auth middleware ([fc350dc](https://github.com/samuel032khoury/gopherfeed/commit/fc350dca813ebdf495ba132614d42aa438bf6a87))
* implement logout functionality and update authentication flow with cookie management ([ba7cb42](https://github.com/samuel032khoury/gopherfeed/commit/ba7cb423760a297f988f4151de91d27659372db6))
* implement rate limiting middleware and error handling for rate limit exceeded ([383b69a](https://github.com/samuel032khoury/gopherfeed/commit/383b69aed87f6b419825d688ddf9cce123d5ccbe))
* implement role-based access control and user roles management ([fbf0356](https://github.com/samuel032khoury/gopherfeed/commit/fbf03561824dd61e515fc128afaa59eaedd2048e))
* implement structured logging across the application and refactor email publishing and consumption ([90220ce](https://github.com/samuel032khoury/gopherfeed/commit/90220cea5655c49dc67f7952a3f1879e40137e61))
* initialize web application with React, TypeScript, and Vite ([af75458](https://github.com/samuel032khoury/gopherfeed/commit/af754587eafdba695fe1f8b7afc23f7cafa3aed6))
* prevent users from following or unfollowing themselves ([9f0e930](https://github.com/samuel032khoury/gopherfeed/commit/9f0e930e134b1ab30c923c9b7d3d67c1bebe706c))
* replace custom logger with zap.SugaredLogger across email and mq packages ([d3f0280](https://github.com/samuel032khoury/gopherfeed/commit/d3f0280a5e071d77d8a1ee33a107a3e8cb3aa240))
* update CI workflow to use ubuntu-latest for improved compatibility ([0775b8e](https://github.com/samuel032khoury/gopherfeed/commit/0775b8e9c05e191b818640c14f068586ad04980c))


### Bug Fixes

* correct error message casing for Mailtrap credentials ([8ddefd5](https://github.com/samuel032khoury/gopherfeed/commit/8ddefd56fa86246cf958975b4218c161ea139b37))
