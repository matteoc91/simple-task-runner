# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).


## [Unreleased]
### Added
- Added read operation
- Added operation management
- Added command line argument validation

### Changed
- Changed **taskmanager.Read** return arguments, provided slices of tasks
- Changed **taskmanager.Create** return arguments, provided only error

## [0.0.0] 2020-07-26
### Added
- Added **taskmanager.Create** method
- Added **simpletask.Task** and **simpletask.Comment** definitions
- Inited project **Simple Task Runner**