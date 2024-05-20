# Changelog

## [0.0.5](https://github.com/krystal/zadara-exporter/compare/v0.0.4...v0.0.5) (2024-05-20)


### Features

* **config:** add health_path config option, defaulting to "/healthz" ([#13](https://github.com/krystal/zadara-exporter/issues/13)) ([75be2e9](https://github.com/krystal/zadara-exporter/commit/75be2e9983b2cdee409bcff1c626a99f2918429a))
* **config:** add metric namespace config option, defaulting to "zadara" ([#12](https://github.com/krystal/zadara-exporter/issues/12)) ([75ae588](https://github.com/krystal/zadara-exporter/commit/75ae5884d900df13f32feb0948cf88a2d1bb0e93))
* **health:** switch health endpoint from Info to Debug log calls on success ([#14](https://github.com/krystal/zadara-exporter/issues/14)) ([fcbd3d3](https://github.com/krystal/zadara-exporter/commit/fcbd3d35ec6aec331139d2e9276f9108e1494bc8))
* **labels:** add "store" label with format "{store_name}@{cloud_name}" ([#9](https://github.com/krystal/zadara-exporter/issues/9)) ([0b6d42d](https://github.com/krystal/zadara-exporter/commit/0b6d42ded7555aebe502961e2f2f195fe9f4954b))


### Bug Fixes

* **labels:** add "store" label to policy metrics ([#15](https://github.com/krystal/zadara-exporter/issues/15)) ([b079d5b](https://github.com/krystal/zadara-exporter/commit/b079d5bf74f5dced39264144d3d9a379c469bd2f))
* **metrics:** hide OpenTelemetry specific labels and metrics ([#10](https://github.com/krystal/zadara-exporter/issues/10)) ([c8e1a42](https://github.com/krystal/zadara-exporter/commit/c8e1a42239d9c9254182d6ddd59a2dea003814e3))

## [0.0.4](https://github.com/krystal/zadara-exporter/compare/v0.0.3...v0.0.4) (2024-05-17)


### Bug Fixes

* **chart/servicemonitor:** ensure correct labels are set ([#7](https://github.com/krystal/zadara-exporter/issues/7)) ([6cf11db](https://github.com/krystal/zadara-exporter/commit/6cf11dbf14315f72b670b35051f6b15e97391653))

## [0.0.3](https://github.com/krystal/zadara-exporter/compare/v0.0.2...v0.0.3) (2024-05-17)


### Bug Fixes

* **health:** bail when an error is encountered ([af831bc](https://github.com/krystal/zadara-exporter/commit/af831bc457b1bf5f494d8c8c49be6da399d27852))

## [0.0.2](https://github.com/krystal/zadara-exporter/compare/v0.0.1...v0.0.2) (2024-05-10)


### Features

* add zadara exporter initial implementation ([#1](https://github.com/krystal/zadara-exporter/issues/1)) ([b49b2c1](https://github.com/krystal/zadara-exporter/commit/b49b2c17f70fe228bac62d8a45308400661f71fd))
